import torch
import torch.nn as nn
import torch.nn.functional as F

import math


class PositionalEmbedding(nn.Module):
    def __init__(self, d_model, max_len=5000):
        super(PositionalEmbedding, self).__init__()
        # Compute the positional encodings once in log space.
        pe = torch.zeros(max_len, d_model).float()
        pe.require_grad = False

        position = torch.arange(0, max_len).float().unsqueeze(1)
        div_term = (torch.arange(0, d_model, 2).float() * -(math.log(10000.0) / d_model)).exp()

        pe[:, 0::2] = torch.sin(position * div_term)
        pe[:, 1::2] = torch.cos(position * div_term)

        pe = pe.unsqueeze(0)
        self.register_buffer('pe', pe)

    def forward(self, x):
        return self.pe[:, :x.size(1)]


class TokenEmbedding(nn.Module):
    def __init__(self, c_in, d_model):
        super(TokenEmbedding, self).__init__()
        # torch.__version__ can be like "2.1.0+cu121"; avoid lexicographic compare (e.g. "1.10" vs "1.5").
        version_core = torch.__version__.split("+", 1)[0]
        major_minor = version_core.split(".", 2)[:2]
        try:
            major, minor = (int(major_minor[0]), int(major_minor[1]))
        except Exception:
            major, minor = (0, 0)
        padding = 1 if (major, minor) >= (1, 5) else 2
        self.tokenConv = nn.Conv1d(in_channels=c_in, out_channels=d_model,
                                   kernel_size=3, padding=padding, padding_mode='circular')
        for m in self.modules():
            if isinstance(m, nn.Conv1d):
                nn.init.kaiming_normal_(m.weight, mode='fan_in', nonlinearity='leaky_relu')

    def forward(self, x):
        x = self.tokenConv(x.permute(0, 2, 1)).transpose(1, 2)
        return x


class FixedEmbedding(nn.Module):
    def __init__(self, c_in, d_model):
        super(FixedEmbedding, self).__init__()

        w = torch.zeros(c_in, d_model).float()
        w.require_grad = False

        position = torch.arange(0, c_in).float().unsqueeze(1)
        div_term = (torch.arange(0, d_model, 2).float() * -(math.log(10000.0) / d_model)).exp()

        w[:, 0::2] = torch.sin(position * div_term)
        w[:, 1::2] = torch.cos(position * div_term)

        self.emb = nn.Embedding(c_in, d_model)
        self.emb.weight = nn.Parameter(w, requires_grad=False)

    def forward(self, x):
        return self.emb(x).detach()


class TemporalEmbedding(nn.Module):
    def __init__(self, d_model, embed_type='fixed', freq='t'):
        super(TemporalEmbedding, self).__init__()
        self.freq = freq
        # 根据频率定义输入特征维度（与 create_time 的输出维度匹配）
        self.input_features = 5  # second_sin, second_cos, minute_sin, minute_cos
        # 使用线性层将时间特征映射到 d_model
        self.linear = nn.Linear(self.input_features, d_model)

    def forward(self, x):
        # x 的形状为 (batch_size, seq_len, input_features)
        return self.linear(x)


# class TemporalEmbedding(nn.Module):
#     def __init__(self, d_model, embed_type='fixed', freq='t'):
#         super(TemporalEmbedding, self).__init__()
#
#         minute_size = 60;
#         hour_size = 24;
#         second_size = 60
#
#         Embed = FixedEmbedding if embed_type == 'fixed' else nn.Embedding
#         self.minute_embed = Embed(minute_size, d_model)
#         self.hour_embed = Embed(hour_size, d_model)
#         self.second_embed = Embed(second_size, d_model)
#
#     def forward(self, x):
#         x = x.long()
#         second_x = self.second_embed(x[:, :, 0])
#         minute_x = self.minute_embed(x[:, :, 1]) if hasattr(self, 'minute_embed') else 0.
#         hour_x = self.hour_embed(x[:, :, 2])
#         return second_x + minute_x + hour_x

class TimeFeatureEmbedding(nn.Module):
    def __init__(self, d_model, embed_type='timeF', freq='h'):
        super(TimeFeatureEmbedding, self).__init__()

        freq_map = {'h': 4, 't': 5, 's': 6, 'm': 1, 'a': 1, 'w': 2, 'd': 3, 'b': 3}
        d_inp = freq_map[freq]
        self.embed = nn.Linear(d_inp, d_model)

    def forward(self, x):
        return self.embed(x)


class DataEmbedding(nn.Module):
    def __init__(self, c_in, d_model, embed_type='fixed', freq='h', dropout=0.1):
        super(DataEmbedding, self).__init__()
        self.value_embedding = TokenEmbedding(c_in=c_in, d_model=d_model)
        self.position_embedding = PositionalEmbedding(d_model=d_model)
        # 统一使用 TemporalEmbedding
        self.temporal_embedding = TemporalEmbedding(d_model=d_model, freq=freq)
        self.dropout = nn.Dropout(p=dropout)

    def forward(self, x, x_mark):
        x1 = self.value_embedding(x)
        x2 = self.position_embedding(x)
        x3 = self.temporal_embedding(x_mark)  # x_mark 来自 create_time 的输出
        x = x1 + x2 + x3
        x = x1 + x2 + x3
        return self.dropout(x)
