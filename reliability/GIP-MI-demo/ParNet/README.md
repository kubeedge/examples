# ParNet

本目录包含 ParNet 的二分类训练与评估代码，输入为 CPU / Memory / Response Time 标注数据。

## 数据

`data/` 下为标注数据 CSV，列为：`cpu, mem, response_time, label(0/1)`。

## 安装依赖

```bash
pip install -r requirements.txt
```

## 运行

默认使用 `data/2024-07-08_label_non-constant.csv`，滑动窗口大小 `window_size=21`：

```bash
python parnet.py
```

指定数据文件：

```bash
python parnet.py --data_path data/2024-05-5_label_constant.csv
```

输出：

- 最佳模型：`checkpoints/parnet/best_model_params.pth`
- 指标：`results/parnet/metrics.txt`
