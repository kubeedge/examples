import argparse
import os
import random
import sys
from dataclasses import dataclass

import numpy as np
import pandas as pd
import torch
import torch.nn as nn
from sklearn.preprocessing import StandardScaler
from torch.utils.data import DataLoader, Dataset
from tqdm import tqdm

from models.gcn import edge_index_to_normalized_adjacency
from models.gcn_informer import GCNInformer
from utils.metrics import metric


PROJECT_ROOT = os.path.dirname(os.path.abspath(__file__))
if PROJECT_ROOT not in sys.path:
    sys.path.insert(0, PROJECT_ROOT)


def set_seed(seed: int) -> None:
    random.seed(seed)
    np.random.seed(seed)
    torch.manual_seed(seed)
    torch.cuda.manual_seed_all(seed)


@dataclass(frozen=True)
class Split:
    train_ratio: float
    val_ratio: float

    def __post_init__(self):
        if not (0.0 < self.train_ratio < 1.0):
            raise ValueError("train_ratio must be in (0, 1)")
        if not (0.0 <= self.val_ratio < 1.0):
            raise ValueError("val_ratio must be in [0, 1)")
        if self.train_ratio + self.val_ratio >= 1.0:
            raise ValueError("train_ratio + val_ratio must be < 1.0")


def create_time_features(data_with_time: np.ndarray) -> np.ndarray:
    """
    data_with_time: (L, 1+F) where first column is datetime-like.
    returns: (L, 5) float32
    """
    # Fast-path for the provided CSV format; fall back to generic parsing if needed.
    time = pd.to_datetime(data_with_time[:, 0], format="%Y-%m-%d %H:%M:%S", errors="coerce")
    if time.isna().any():
        time = pd.to_datetime(data_with_time[:, 0])
    start_time = pd.to_datetime(data_with_time[0, 0], format="%Y-%m-%d %H:%M:%S", errors="coerce")
    if pd.isna(start_time):
        start_time = pd.to_datetime(data_with_time[0, 0])
    total_seconds = (time - start_time).total_seconds().values.astype(np.int32)
    max_seconds = (pd.to_datetime(data_with_time[-1, 0]) - start_time).total_seconds()
    if max_seconds <= 0:
        max_seconds = 1.0
    absolute_time = (total_seconds / max_seconds).astype(np.float32)[:, None]

    local_seconds_60 = (total_seconds % 60).astype(np.float32)[:, None]
    local_sec_sin_60 = np.sin(2 * np.pi * local_seconds_60 / 60).astype(np.float32)
    local_sec_cos_60 = np.cos(2 * np.pi * local_seconds_60 / 60).astype(np.float32)

    local_seconds_100 = (total_seconds % 100).astype(np.float32)[:, None]
    local_sec_sin_100 = np.sin(2 * np.pi * local_seconds_100 / 100).astype(np.float32)
    local_sec_cos_100 = np.cos(2 * np.pi * local_seconds_100 / 100).astype(np.float32)

    return np.concatenate([absolute_time, local_sec_sin_60, local_sec_cos_60, local_sec_sin_100, local_sec_cos_100], axis=-1)


def load_csv_as_numpy(file_path: str) -> np.ndarray:
    df = pd.read_csv(file_path)
    # Keep column order; assume first column is time, remaining are numeric features.
    return df.values


def make_windows(data: np.ndarray, seq_len: int, label_len: int, pred_len: int, stride: int) -> tuple[np.ndarray, np.ndarray]:
    """
    data: (L, 1+F), first column is time, remaining are scaled features.
    Returns:
      X: (N, seq_len, 1+F)
      Y: (N, label_len+pred_len, 1+F)
    """
    x_list = []
    y_list = []
    end = len(data) - seq_len - pred_len + 1
    for i in range(0, end, stride):
        seq_data = data[i : i + seq_len]
        x_list.append(seq_data)

        target = data[i + seq_len - label_len : i + seq_len + pred_len]
        y_list.append(target)
    return np.array(x_list), np.array(y_list)


def split_and_scale(
    data_raw: np.ndarray,
    *,
    split: Split,
    seq_len: int,
    label_len: int,
    pred_len: int,
    stride: int,
) -> tuple[np.ndarray, np.ndarray, np.ndarray, np.ndarray, np.ndarray, np.ndarray, StandardScaler]:
    n = len(data_raw)
    split_train = int(n * split.train_ratio)
    split_val = int(n * (split.train_ratio + split.val_ratio))

    train_raw = data_raw[:split_train]
    val_raw = data_raw[max(0, split_train - seq_len) : split_val]
    test_raw = data_raw[max(0, split_val - seq_len) :]

    scaler = StandardScaler()
    scaler.fit(train_raw[:, 1:].astype(np.float32))

    def transform(raw: np.ndarray) -> np.ndarray:
        time_col = raw[:, 0].reshape(-1, 1)
        feats = raw[:, 1:].astype(np.float32)
        feats_scaled = scaler.transform(feats)
        return np.hstack((time_col, feats_scaled))

    train = transform(train_raw)
    val = transform(val_raw)
    test = transform(test_raw)

    x_train, y_train = make_windows(train, seq_len, label_len, pred_len, stride)
    x_val, y_val = make_windows(val, seq_len, label_len, pred_len, stride)
    x_test, y_test = make_windows(test, seq_len, label_len, pred_len, stride=1)
    return x_train, x_val, x_test, y_train, y_val, y_test, scaler


class TimeSeriesDataset(Dataset):
    def __init__(self, x: np.ndarray, y: np.ndarray):
        self.x = x
        self.y = y

    def __len__(self) -> int:
        return len(self.x)

    def __getitem__(self, idx: int):
        x = self.x[idx]
        y = self.y[idx]

        x_time = create_time_features(x)
        y_time = create_time_features(y)

        x_value = x[:, 1:].astype(np.float32)
        y_value = y[:, 1:].astype(np.float32)

        return (
            torch.from_numpy(x_value),
            torch.from_numpy(y_value),
            torch.from_numpy(x_time.astype(np.float32)),
            torch.from_numpy(y_time.astype(np.float32)),
        )


def create_edge_index_and_weights(device: torch.device) -> tuple[torch.Tensor, torch.Tensor]:
    # 3 nodes: [CPU, MEM, ResponseTime]
    edge_index = torch.tensor(
        [
            [0, 2],
            [2, 0],
            [0, 1],
            [1, 0],
            [1, 2],
            [2, 1],
        ],
        dtype=torch.long,
        device=device,
    ).t().contiguous()

    edge_weight = torch.tensor([0.56, 0.56, 0.66, 0.66, 0.52, 0.52], dtype=torch.float32, device=device)
    return edge_index, edge_weight


def weighted_mse(pred: torch.Tensor, true: torch.Tensor, weights: torch.Tensor) -> torch.Tensor:
    # pred/true: (B, L, 3)
    return (weights * (pred - true) ** 2).mean()


def save_predictions_csv(
    *,
    trues: np.ndarray,
    preds: np.ndarray,
    out_dir: str,
    prefix: str,
) -> None:
    os.makedirs(out_dir, exist_ok=True)
    cols = ["cpu", "mem", "response_time"]
    df_true = pd.DataFrame(trues, columns=[f"{c}_true" for c in cols])
    df_pred = pd.DataFrame(preds, columns=[f"{c}_pred" for c in cols])
    df = pd.concat([df_true, df_pred], axis=1)
    df.to_csv(os.path.join(out_dir, f"{prefix}_predictions.csv"), index=False)


def main() -> int:
    parser = argparse.ArgumentParser(description="Train GCN-Informer on custom 3-metric time series (CPU/MEM/RT).")
    parser.add_argument("--data_path", type=str, default=os.path.join("data", "myData", "2024-05-5_train_constant.csv"))
    parser.add_argument("--seq_len", type=int, default=10)
    parser.add_argument("--label_len", type=int, default=10)
    parser.add_argument("--pred_len", type=int, default=5)
    parser.add_argument("--stride", type=int, default=5)
    parser.add_argument("--epochs", type=int, default=100)
    parser.add_argument("--batch_size", type=int, default=32)
    parser.add_argument("--lr", type=float, default=1e-4)
    parser.add_argument("--train_ratio", type=float, default=0.8)
    parser.add_argument("--val_ratio", type=float, default=0.1)
    parser.add_argument("--patience", type=int, default=10)
    parser.add_argument("--min_delta", type=float, default=1e-4)
    parser.add_argument("--seed", type=int, default=42)
    parser.add_argument("--device", type=str, default="auto", choices=["auto", "cpu", "cuda"])
    parser.add_argument("--save_dir", type=str, default=os.path.join("checkpoints", "gcn_informer"))
    parser.add_argument("--results_dir", type=str, default=os.path.join("results", "gcn_informer"))
    args = parser.parse_args()

    os.chdir(PROJECT_ROOT)
    set_seed(args.seed)

    if args.device == "auto":
        device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    elif args.device == "cuda":
        device = torch.device("cuda")
    else:
        device = torch.device("cpu")

    data_raw = load_csv_as_numpy(args.data_path)
    split = Split(train_ratio=args.train_ratio, val_ratio=args.val_ratio)

    x_train, x_val, x_test, y_train, y_val, y_test, scaler = split_and_scale(
        data_raw,
        split=split,
        seq_len=args.seq_len,
        label_len=args.label_len,
        pred_len=args.pred_len,
        stride=args.stride,
    )

    train_loader = DataLoader(TimeSeriesDataset(x_train, y_train), batch_size=args.batch_size, shuffle=True)
    val_loader = DataLoader(TimeSeriesDataset(x_val, y_val), batch_size=args.batch_size, shuffle=False)
    test_loader = DataLoader(TimeSeriesDataset(x_test, y_test), batch_size=1, shuffle=False)

    edge_index, edge_weight = create_edge_index_and_weights(device)
    adj_norm = edge_index_to_normalized_adjacency(edge_index, edge_weight, num_nodes=3)

    model = GCNInformer(
        gcn_in_channels=args.seq_len,
        gcn_hidden_channels=3,
        gcn_out_channels=args.seq_len,
        num_gcn_layers=3,
        informer_enc_in=3,
        informer_dec_in=3,
        informer_c_out=3,
        informer_out_len=args.pred_len,
    ).to(device)

    optimizer = torch.optim.Adam(model.parameters(), lr=args.lr)
    weights = torch.tensor([1.0, 1.0, 1.2], dtype=torch.float32, device=device).view(1, 1, 3)

    os.makedirs(args.save_dir, exist_ok=True)
    best_path = os.path.join(args.save_dir, "best_model.pt")
    best_val = float("inf")
    patience_counter = 0

    for epoch in range(args.epochs):
        model.train()
        train_losses = []
        for x_value, y_value, x_time, y_time in tqdm(train_loader, desc=f"Epoch {epoch+1}/{args.epochs}", leave=False):
            x_value = x_value.to(device)  # (B, seq_len, 3)
            y_value = y_value.to(device)  # (B, label_len+pred_len, 3)
            x_time = x_time.to(device)  # (B, seq_len, 5)
            y_time = y_time.to(device)  # (B, label_len+pred_len, 5)

            x_nodes = x_value.transpose(1, 2)  # (B, 3, seq_len)
            dec_mask = torch.zeros_like(y_value[:, args.label_len :], device=device)
            x_dec = torch.cat([y_value[:, : args.label_len], dec_mask], dim=1)

            pred = model(x_nodes, adj_norm, x_time, x_dec, y_time)
            true = y_value[:, args.label_len :]
            loss = weighted_mse(pred, true, weights)

            optimizer.zero_grad()
            loss.backward()
            optimizer.step()

            train_losses.append(loss.item())

        # Validation
        model.eval()
        val_losses = []
        with torch.no_grad():
            for x_value, y_value, x_time, y_time in val_loader:
                x_value = x_value.to(device)
                y_value = y_value.to(device)
                x_time = x_time.to(device)
                y_time = y_time.to(device)

                x_nodes = x_value.transpose(1, 2)
                dec_mask = torch.zeros_like(y_value[:, args.label_len :], device=device)
                x_dec = torch.cat([y_value[:, : args.label_len], dec_mask], dim=1)

                pred = model(x_nodes, adj_norm, x_time, x_dec, y_time)
                true = y_value[:, args.label_len :]
                val_losses.append(weighted_mse(pred, true, weights).item())

        train_loss = float(np.mean(train_losses)) if train_losses else float("nan")
        val_loss = float(np.mean(val_losses)) if val_losses else float("nan")
        print(f"Epoch {epoch+1}: train_loss={train_loss:.6f} val_loss={val_loss:.6f}")

        if val_loss < best_val - args.min_delta:
            best_val = val_loss
            patience_counter = 0
            # Save weights only (compatible with PyTorch 2.6+ default weights_only loader).
            torch.save(model.state_dict(), best_path)
        else:
            patience_counter += 1
            if patience_counter >= args.patience:
                print(f"Early stopping: no val improvement for {args.patience} epochs. Best val={best_val:.6f}")
                break

    # Test (load best if exists)
    if os.path.exists(best_path):
        state = torch.load(best_path, map_location=device)
        model.load_state_dict(state)

    model.eval()
    preds = []
    trues = []
    with torch.no_grad():
        for x_value, y_value, x_time, y_time in test_loader:
            x_value = x_value.to(device)
            y_value = y_value.to(device)
            x_time = x_time.to(device)
            y_time = y_time.to(device)

            x_nodes = x_value.transpose(1, 2)
            dec_mask = torch.zeros_like(y_value[:, args.label_len :], device=device)
            x_dec = torch.cat([y_value[:, : args.label_len], dec_mask], dim=1)

            pred = model(x_nodes, adj_norm, x_time, x_dec, y_time)  # (1, pred_len, 3)
            true = y_value[:, args.label_len :]  # (1, pred_len, 3)

            preds.append(pred.squeeze(0).cpu().numpy())
            trues.append(true.squeeze(0).cpu().numpy())

    preds = np.array(preds).reshape(-1, 3)
    trues = np.array(trues).reshape(-1, 3)

    # Inverse scale
    preds_inv = scaler.inverse_transform(preds)
    trues_inv = scaler.inverse_transform(trues)

    mae, mse, rmse, mape, mspe = metric(preds_inv, trues_inv)
    print(f"Test: mse={mse:.6f} mae={mae:.6f} rmse={rmse:.6f} mape={mape:.6f} mspe={mspe:.6f}")

    save_predictions_csv(trues=trues_inv, preds=preds_inv, out_dir=args.results_dir, prefix="test")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
