import torch
import torch.nn as nn
import torch.nn.functional as F


def edge_index_to_normalized_adjacency(edge_index: torch.Tensor, edge_weight: torch.Tensor, num_nodes: int) -> torch.Tensor:
    """
    Build a symmetrically-normalized adjacency with self-loops:  D^{-1/2} (A + I) D^{-1/2}
    edge_index: (2, E)
    edge_weight: (E,)
    returns: (N, N)
    """
    device = edge_index.device
    dtype = edge_weight.dtype
    adj = torch.zeros((num_nodes, num_nodes), device=device, dtype=dtype)
    adj[edge_index[0], edge_index[1]] = edge_weight
    adj = adj + torch.eye(num_nodes, device=device, dtype=dtype)

    deg = adj.sum(dim=1)
    deg_inv_sqrt = torch.pow(deg, -0.5)
    deg_inv_sqrt[torch.isinf(deg_inv_sqrt)] = 0.0
    d = torch.diag(deg_inv_sqrt)
    return d @ adj @ d


class GraphConvolution(nn.Module):
    """
    A simple batched GCN layer for small fixed graphs.

    x: (B, N, Fin)
    adj_norm: (N, N)
    """

    def __init__(self, in_features: int, out_features: int, bias: bool = True):
        super().__init__()
        self.weight = nn.Parameter(torch.empty(in_features, out_features))
        self.bias = nn.Parameter(torch.empty(out_features)) if bias else None
        self.reset_parameters()

    def reset_parameters(self) -> None:
        nn.init.xavier_uniform_(self.weight)
        if self.bias is not None:
            nn.init.zeros_(self.bias)

    def forward(self, x: torch.Tensor, adj_norm: torch.Tensor) -> torch.Tensor:
        support = x @ self.weight
        out = torch.einsum("ij,bjf->bif", adj_norm, support)
        if self.bias is not None:
            out = out + self.bias
        return out


class GCN(nn.Module):
    def __init__(self, in_channels: int, hidden_channels: int, out_channels: int, num_layers: int, dropout: float = 0.5):
        super().__init__()
        if num_layers < 1:
            raise ValueError("num_layers must be >= 1")
        self.dropout = dropout

        layers: list[nn.Module] = []
        for i in range(num_layers):
            if i == 0 and num_layers == 1:
                layers.append(GraphConvolution(in_channels, out_channels))
            elif i == 0:
                layers.append(GraphConvolution(in_channels, hidden_channels))
            elif i == num_layers - 1:
                layers.append(GraphConvolution(hidden_channels, out_channels))
            else:
                layers.append(GraphConvolution(hidden_channels, hidden_channels))
        self.layers = nn.ModuleList(layers)

    def forward(self, x: torch.Tensor, adj_norm: torch.Tensor) -> torch.Tensor:
        for i, layer in enumerate(self.layers):
            x = layer(x, adj_norm)
            x = F.relu(x)
            if i != len(self.layers) - 1:
                x = F.dropout(x, p=self.dropout, training=self.training)
        return x

