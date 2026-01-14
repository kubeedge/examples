import torch
import torch.nn as nn

from utils import SE1, channel_shuffle
from typing import Optional


def conv_bn(in_channels: int, out_channels: int, kernel_size: int, stride: int, padding: int, groups: int = 1) -> nn.Sequential:
    return nn.Sequential(
        nn.Conv2d(
            in_channels=in_channels,
            out_channels=out_channels,
            kernel_size=kernel_size,
            stride=stride,
            padding=padding,
            groups=groups,
            bias=False,
        ),
        nn.BatchNorm2d(out_channels),
    )


class RepVGGBlock(nn.Module):
    def __init__(
        self,
        in_channels: int,
        out_channels: int,
        kernel_size: int,
        *,
        stride: int = 1,
        padding: int = 1,
        groups: int = 1,
        avg_pool: bool = False,
        se_block: bool = False,
        activation: Optional[nn.Module] = None,
    ):
        super().__init__()
        if activation is None:
            activation = nn.ReLU()

        self.groups = groups
        self.stride = stride
        self.kernel_size = kernel_size
        self.in_channels = in_channels
        self.out_channels = out_channels
        self.padding = padding
        self.se_block = se_block
        self.fused = False

        assert padding == 1
        padding_11 = padding - 3 // 2

        self.nonlinearity = activation
        self.rbr_identity = nn.BatchNorm2d(in_channels) if (out_channels == in_channels and stride == 1) else None
        self.rbr_dense = conv_bn(in_channels, out_channels, kernel_size=3, stride=stride, padding=padding, groups=groups) if (kernel_size != 1) else None
        self.rbr_1x1 = conv_bn(in_channels, out_channels, kernel_size=1, stride=stride, padding=padding_11, groups=groups)
        if stride == 2 and avg_pool:
            self.rbr_1x1 = nn.Sequential(nn.AvgPool2d(2, 2), conv_bn(in_channels, out_channels, kernel_size=1, stride=1, padding=0, groups=groups))

        self.channel_shuffle_enabled = groups > 1
        if self.se_block:
            self.se = SE1(in_channels, out_channels, g=groups, ver=2 if (out_channels != in_channels or stride != 1) else 1)

    def forward(self, inputs: torch.Tensor) -> torch.Tensor:
        id_out = 0 if self.rbr_identity is None else self.rbr_identity(inputs)
        rbr_1x1_output = 0 if self.fused else self.rbr_1x1(inputs)
        dense_output = 0 if self.rbr_dense is None else self.rbr_dense(inputs)

        if self.se_block and self.rbr_identity is not None:
            id_out = id_out * self.se(id_out)

        out = dense_output + rbr_1x1_output + id_out
        if self.se_block and (self.rbr_identity is None):
            out = out * self.se(inputs)

        out = self.nonlinearity(out)
        if self.channel_shuffle_enabled:
            out = channel_shuffle(out, self.groups)
        return out
