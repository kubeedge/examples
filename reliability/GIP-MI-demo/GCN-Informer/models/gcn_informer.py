import torch.nn as nn

from models.gcn import GCN
from models.informer import Informer


class GCNInformer(nn.Module):
    def __init__(
        self,
        *,
        gcn_in_channels: int,
        gcn_hidden_channels: int,
        gcn_out_channels: int,
        num_gcn_layers: int,
        informer_enc_in: int,
        informer_dec_in: int,
        informer_c_out: int,
        informer_out_len: int,
        gcn_dropout: float = 0.5,
    ):
        super().__init__()
        self.gcn = GCN(
            in_channels=gcn_in_channels,
            hidden_channels=gcn_hidden_channels,
            out_channels=gcn_out_channels,
            num_layers=num_gcn_layers,
            dropout=gcn_dropout,
        )

        # Map node dimension (e.g. CPU/MEM/RT = 3) to Informer input channels.
        self.fc = nn.Linear(3, informer_enc_in)
        self.informer = Informer(
            enc_in=informer_enc_in,
            dec_in=informer_dec_in,
            c_out=informer_c_out,
            out_len=informer_out_len,
        )

    def forward(self, x_nodes, adj_norm, x_time_enc, x_dec, x_time_dec):
        """
        x_nodes: (B, 3, seq_len)    3 nodes = [CPU, MEM, RT], feature dim = seq_len
        adj_norm: (3, 3)
        x_time_enc: (B, seq_len, 5)
        x_dec: (B, label_len + pred_len, 3)
        x_time_dec: (B, label_len + pred_len, 5)
        """
        gcn_out = self.gcn(x_nodes, adj_norm)  # (B, 3, seq_len)
        gcn_out = gcn_out.permute(0, 2, 1)  # (B, seq_len, 3)
        x_enc = self.fc(gcn_out)  # (B, seq_len, informer_enc_in)
        return self.informer(x_enc, x_time_enc, x_dec, x_time_dec)

