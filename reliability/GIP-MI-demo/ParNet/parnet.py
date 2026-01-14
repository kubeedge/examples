import argparse
import copy
import os
import sys

import numpy as np
import pandas as pd
import torch
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim
from sklearn.metrics import accuracy_score, confusion_matrix, f1_score, precision_score, recall_score
from sklearn.model_selection import train_test_split
from torch.utils.data import DataLoader, TensorDataset
from tqdm import tqdm

from repvgg import RepVGGBlock
from utils import Concat2d, MultiBatchNorm2d


PROJECT_ROOT = os.path.dirname(os.path.abspath(__file__))
if PROJECT_ROOT not in sys.path:
    sys.path.insert(0, PROJECT_ROOT)


class ClassifierHead(nn.Module):
    def __init__(self, in_chs: int, num_classes: int, pool_type: str = "avg", drop_rate: float = 0.0):
        super().__init__()
        if pool_type == "avg":
            self.pool = nn.AdaptiveAvgPool2d(1)
        elif pool_type == "max":
            self.pool = nn.AdaptiveMaxPool2d(1)
        else:
            raise ValueError(f"Unsupported pool_type: {pool_type}")
        self.drop = nn.Dropout(p=drop_rate) if drop_rate and drop_rate > 0 else nn.Identity()
        self.fc = nn.Linear(in_chs, num_classes)

    def forward(self, x: torch.Tensor) -> torch.Tensor:
        x = self.pool(x).flatten(1)
        x = self.drop(x)
        return self.fc(x)


class RepVGGOur(nn.Module):
    expansion: int = 1

    def __init__(self, inplanes, planes, stride=1, groups=1, kernel_size=3, se_block=True, additional_branches=None):
        super().__init__()
        additional_branches = additional_branches or []

        activation = nn.ReLU()
        if "swish" in additional_branches:
            activation = nn.SiLU()

        self.block = RepVGGBlock(
            inplanes,
            planes,
            kernel_size,
            stride=stride,
            padding=1,
            groups=groups,
            avg_pool=True,
            se_block=se_block,
            activation=activation,
        )

    def forward(self, x):
        return self.block(x)


class SimpNet(nn.Module):
    def __init__(self, planes, num_blocks, dropout_lin, additional_branches=None):
        super().__init__()
        additional_branches = additional_branches or []
        block = RepVGGOur

        last_planes = planes[-1]
        planes = planes[0:-1]
        strides = [2] * len(planes)
        assert num_blocks[-1] == 1
        assert num_blocks[-2] != 1
        num_blocks = num_blocks[0:-1]

        self.inits = nn.ModuleList()
        in_planes = min(64, planes[0])
        self.inits.append(
            nn.Sequential(
                block(3, in_planes, stride=2, additional_branches=additional_branches),
                block(in_planes, planes[0], stride=2, additional_branches=additional_branches),
            )
        )

        for stride, in_plane, out_plane in zip(strides[1:], planes[0:-1], planes[1:]):
            self.inits.append(
                block(
                    in_plane * block.expansion,
                    out_plane * block.expansion,
                    stride,
                    additional_branches=additional_branches,
                )
            )

        self.streams = nn.ModuleList()

        def stream_block(plane):
            return block(plane, plane, kernel_size=3, stride=1, se_block=True, additional_branches=additional_branches)

        for num_block, plane in zip(num_blocks, planes):
            stream = nn.ModuleList()
            for _ in range(num_block - 1):
                stream.append(stream_block(plane * block.expansion))
            self.streams.append(nn.Sequential(*stream))

        self.downsamples_2 = nn.ModuleList()
        in_planes_list = planes[0:-1]
        out_planes_list = planes[1:]
        for i, (stride, in_plane, out_plane) in enumerate(zip(strides[1:], in_planes_list, out_planes_list)):
            if i == 0:
                self.downsamples_2.append(
                    block(
                        in_plane * block.expansion,
                        out_plane * block.expansion,
                        stride,
                        kernel_size=3,
                        additional_branches=additional_branches,
                    )
                )
            else:
                self.downsamples_2.append(
                    nn.Sequential(
                        MultiBatchNorm2d(in_plane * block.expansion, in_plane * block.expansion),
                        Concat2d(shuffle=True),
                        block(
                            2 * in_plane * block.expansion,
                            out_plane * block.expansion,
                            stride=2,
                            groups=2,
                            kernel_size=3,
                            additional_branches=additional_branches,
                        ),
                    )
                )

        in_planes_combine = planes[-1]
        self.combine = nn.Sequential(
            MultiBatchNorm2d(in_planes_combine * block.expansion, in_planes_combine * block.expansion),
            Concat2d(shuffle=True),
            block(
                2 * in_planes_combine * block.expansion,
                in_planes_combine * block.expansion,
                stride=1,
                groups=2,
                additional_branches=additional_branches,
            ),
            block(planes[-1], last_planes, stride=2, additional_branches=additional_branches),
        )

        self.head = ClassifierHead(last_planes * block.expansion, 1000, pool_type="avg", drop_rate=dropout_lin)
        self.num_features = last_planes * block.expansion

    def forward(self, img):
        x = img
        x_list = []
        for init in self.inits:
            x = init(x)
            x_list.append(x)

        y_old = None
        for i, (x_s, stream) in enumerate(zip(x_list, self.streams)):
            y = stream(x_s)
            if y_old is None:
                y_old = self.downsamples_2[i](y)
            elif i < len(self.downsamples_2):
                y_old = self.downsamples_2[i]((y, y_old))
            else:
                y_old = (y, y_old)

        out = self.combine(y_old)
        return self.head(out)


def slice_data(data_path: str, window_size: int):
    df = pd.read_csv(data_path)

    sequences = []
    labels = []
    for i in range(len(df) - window_size + 1):
        sequence = df.iloc[i : i + window_size, :3]
        label_1_count = (df.iloc[i : i + window_size, 3] == 0).sum()
        label_2_count = (df.iloc[i : i + window_size, 3] == 1).sum()
        label = 0 if label_1_count > label_2_count else 1
        sequences.append(sequence)
        labels.append(label)

    return np.array(sequences), np.array(labels)


def divide_data(sequences, labels, image_size: int):
    x_train, x_temp, y_train, y_temp = train_test_split(sequences, labels, test_size=0.2, stratify=labels, random_state=None)
    x_test, x_val, y_test, y_val = train_test_split(x_temp, y_temp, test_size=0.3, stratify=y_temp, random_state=None)

    x_train = np.expand_dims(np.array(x_train), axis=1)
    x_test = np.expand_dims(np.array(x_test), axis=1)
    x_val = np.expand_dims(np.array(x_val), axis=1)

    x_train = np.transpose(x_train, (0, 3, 2, 1))
    x_test = np.transpose(x_test, (0, 3, 2, 1))
    x_val = np.transpose(x_val, (0, 3, 2, 1))

    x_train_tensor = torch.tensor(x_train, dtype=torch.float32)
    y_train_tensor = torch.tensor(y_train, dtype=torch.long)
    x_test_tensor = torch.tensor(x_test, dtype=torch.float32)
    y_test_tensor = torch.tensor(y_test, dtype=torch.long)
    x_val_tensor = torch.tensor(x_val, dtype=torch.float32)
    y_val_tensor = torch.tensor(y_val, dtype=torch.long)

    x_train_tensor = F.interpolate(x_train_tensor, size=(image_size, image_size), mode="bilinear", align_corners=False)
    x_test_tensor = F.interpolate(x_test_tensor, size=(image_size, image_size), mode="bilinear", align_corners=False)
    x_val_tensor = F.interpolate(x_val_tensor, size=(image_size, image_size), mode="bilinear", align_corners=False)

    return x_train_tensor, y_train_tensor, x_test_tensor, y_test_tensor, x_val_tensor, y_val_tensor


def train_with_parnet(
    x_train_tensor,
    y_train_tensor,
    x_test_tensor,
    y_test_tensor,
    x_val_tensor,
    y_val_tensor,
    *,
    epochs: int,
    batch_size: int,
    lr: float,
    save_path: str,
):
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

    train_loader = DataLoader(TensorDataset(x_train_tensor, y_train_tensor), batch_size=batch_size, shuffle=True)
    test_loader = DataLoader(TensorDataset(x_test_tensor, y_test_tensor), batch_size=batch_size, shuffle=False)
    val_loader = DataLoader(TensorDataset(x_val_tensor, y_val_tensor), batch_size=batch_size, shuffle=False)

    model = SimpNet(planes=[92, 192, 384, 1280], num_blocks=[5, 6, 6, 1], dropout_lin=0.0).to(device)
    criterion = nn.CrossEntropyLoss()
    optimizer = optim.Adam(model.parameters(), lr=lr)

    best_model_params = None
    best_accuracy = 0.0

    model.train()
    for epoch in range(epochs):
        running_loss = 0.0
        for inputs, labels in tqdm(train_loader, desc=f"Epoch {epoch + 1}", leave=False):
            inputs, labels = inputs.to(device), labels.to(device)
            optimizer.zero_grad()
            outputs = model(inputs)
            loss = criterion(outputs, labels)
            loss.backward()
            optimizer.step()
            running_loss += loss.item() * inputs.size(0)

        epoch_loss = running_loss / len(train_loader.dataset)
        print(f"Epoch {epoch + 1}, Loss: {epoch_loss:.4f}")

        model.eval()
        correct = 0
        total = 0
        with torch.no_grad():
            for inputs, labels in val_loader:
                inputs, labels = inputs.to(device), labels.to(device)
                outputs = model(inputs)
                _, predicted = torch.max(outputs, 1)
                total += labels.size(0)
                correct += (predicted == labels).sum().item()

        accuracy = correct / total
        print(f"Validation Accuracy after Epoch {epoch + 1}: {accuracy:.4f}")

        if accuracy > best_accuracy:
            best_accuracy = accuracy
            best_model_params = copy.deepcopy(model.state_dict())

    os.makedirs(os.path.dirname(save_path), exist_ok=True)
    torch.save(best_model_params, save_path)
    model.load_state_dict(best_model_params)

    model.eval()
    predictions = []
    true_labels = []
    with torch.no_grad():
        for inputs, labels in tqdm(test_loader, desc="Testing", leave=False):
            inputs, labels = inputs.to(device), labels.to(device)
            outputs = model(inputs)
            _, predicted = torch.max(outputs, 1)
            predictions.extend(predicted.cpu().numpy())
            true_labels.extend(labels.cpu().numpy())

    conf_matrix = confusion_matrix(true_labels, predictions)
    print("Confusion matrix:")
    print(conf_matrix)

    accuracy = accuracy_score(true_labels, predictions)
    precision = precision_score(true_labels, predictions, average="macro", zero_division=0)
    recall = recall_score(true_labels, predictions, average="macro", zero_division=0)
    f1 = f1_score(true_labels, predictions, average="macro", zero_division=0)

    print(f"Accuracy: {accuracy:.4f}")
    print(f"Precision: {precision:.4f}")
    print(f"Recall: {recall:.4f}")
    print(f"F1: {f1:.4f}")

    return accuracy, precision, recall, f1, conf_matrix


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--data_path", type=str, default=os.path.join("data", "2024-07-08_label_non-constant.csv"))
    parser.add_argument("--window_size", type=int, default=21)
    parser.add_argument("--image_size", type=int, default=128)
    parser.add_argument("--epochs", type=int, default=10)
    parser.add_argument("--batch_size", type=int, default=10)
    parser.add_argument("--lr", type=float, default=0.001)
    args = parser.parse_args()

    os.chdir(PROJECT_ROOT)

    sequences, labels = slice_data(args.data_path, args.window_size)
    x_train, y_train, x_test, y_test, x_val, y_val = divide_data(sequences, labels, args.image_size)

    save_path = os.path.join("checkpoints", "parnet", "best_model_params.pth")
    accuracy, precision, recall, f1, cm = train_with_parnet(
        x_train,
        y_train,
        x_test,
        y_test,
        x_val,
        y_val,
        epochs=args.epochs,
        batch_size=args.batch_size,
        lr=args.lr,
        save_path=save_path,
    )

    os.makedirs(os.path.join("results", "parnet"), exist_ok=True)
    with open(os.path.join("results", "parnet", "metrics.txt"), "w", encoding="utf-8") as f:
        f.write(f"accuracy\t{accuracy}\n")
        f.write(f"precision\t{precision}\n")
        f.write(f"recall\t{recall}\n")
        f.write(f"f1\t{f1}\n")
        f.write("confusion_matrix\n")
        f.write(np.array2string(cm))
        f.write("\n")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
