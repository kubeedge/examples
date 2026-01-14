import torch
import torch.nn as nn


def num_param(net: nn.Module) -> int:
    return sum(p.numel() for p in net.parameters() if p.requires_grad)


def trace_net(net: nn.Module, inp: torch.Tensor, save_pth: str = "traced_model.pt") -> None:
    traced_script_module = torch.jit.trace(net, inp, strict=True)
    traced_script_module.save(save_pth)


def channel_shuffle(x: torch.Tensor, groups: int) -> torch.Tensor:
    n, c, h, w = x.shape
    if groups <= 1 or c % groups != 0:
        return x
    x = x.view(n, groups, c // groups, h, w)
    x = x.transpose(1, 2).contiguous()
    return x.view(n, c, h, w)


class MultiBatchNorm2d(nn.Module):
    def __init__(self, n1: int, n2: int):
        super().__init__()
        self.b1 = nn.BatchNorm2d(n1)
        self.b2 = nn.BatchNorm2d(n2)

    def forward(self, x):
        x1, x2 = x
        return (self.b1(x1), self.b2(x2))


class Concat2d(nn.Module):
    def __init__(self, shuffle: bool = False):
        super().__init__()
        self.shuffle = shuffle

    def forward(self, x):
        if self.shuffle:
            b, _, h, w = x[0].shape
            x = [_x.unsqueeze(1) for _x in x]
            out = torch.cat(x, 1)
            out = out.transpose(1, 2)
            return torch.reshape(out, (b, -1, h, w))
        return torch.cat(x, 1)


class SE1(nn.Module):
    def __init__(self, c_in: int, c_out: int, g: int = 1, reduction: int = 16, ver: int = 1):
        super().__init__()
        self.ver = ver
        self.avg_pool = nn.AdaptiveAvgPool2d(1)
        hidden = max(1, c_out // reduction)
        self.fc = nn.Sequential(
            nn.Conv2d(c_in, hidden, kernel_size=1, groups=g, bias=False),
            nn.ReLU(inplace=True),
            nn.Conv2d(hidden, c_out, kernel_size=1, groups=g, bias=False),
            nn.Sigmoid(),
        )

    def forward(self, x):
        y = self.fc(self.avg_pool(x))
        if self.ver == 2:
            y = 2 * y
        return y

