# GCN-Informer

本目录包含 GCN-Informer 的训练/验证/测试代码，以及配套的 CPU / Memory / Response Time 数据集样例。

## 目录结构

- `train_gcn_informer.py`：训练/验证/测试入口（默认使用 `data/myData/2024-05-5_train_constant.csv`）
- `data/myData/`：数据集 CSV
- `models/`：模型组件（Informer + GCN）
- `utils/`：mask 与指标计算
- `checkpoints/`：默认保存 `checkpoints/gcn_informer/best_model.pt`
- `results/`：默认输出 `results/gcn_informer/test_predictions.csv`

## 安装依赖

```bash
pip install -r requirements.txt
```

## 运行

默认训练（结束后会跑 test 并导出预测 CSV）：

```bash
python train_gcn_informer.py
```

指定数据文件：

```bash
python train_gcn_informer.py --data_path data/myData/2024-07-08_train_non-constant.csv
```
