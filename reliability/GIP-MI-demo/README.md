# GIP-MI: Software Aging Prediction and Rejuvenation for KubeEdge

## Introduction
This example demonstrates the **GIP-MI** framework, a reliability enhancement solution designed for KubeEdge edge nodes. It aims to proactively detect software aging trends and perform zero-downtime rejuvenation.

The framework consists of **three core modules** working in a pipeline:

1.  **Aging Prediction (GCN-Informer)**:
    * Captures spatial-temporal dependencies of system metrics (CPU, Memory, etc.).
    * Predicts future metric trends to provide early warning signals.

2.  **Aging State Determination (ParNet)**:
    * Analyzes the predicted data slices using a Parallel Network (ParNet).
    * Accurately classifies the node's state (e.g., *Normal* vs. *Aging*).

3.  **Rejuvenation Strategy (MOEA/D-IFM)**:
    * Triggered when an "Aging" state is confirmed.
    * Uses Multi-Objective Evolutionary Algorithm (MOEA/D-IFM) to generate an optimal task offloading and node restart strategy, ensuring service continuity.

## Directory Structure

This example includes the implementation of the three algorithms:

* `GCN-Informer/`: Code for the aging prediction model.
* `ParNet/`: Code for the aging state determination model.
* `MOEAD-IFM/`: Code for the rejuvenation strategy generation algorithm.

## Prerequisites

* Python 3.8+
* PyTorch (for deep learning models)
* NumPy, Pandas, Scikit-learn

## Usage Guide

### Step 1: Aging Prediction
Navigate to the prediction module to forecast system metrics:

```bash
cd GCN-Informer
Run the prediction script
python main.py  (Note: Replace with your actual entry script)
```

### Step 2: Aging Determination
Use the predicted data to determine if the node is entering an aging state:
```bash
cd ../ParNet
Run the state determination model
python main.py  (Note: Replace with your actual entry script)
```

### Step 3: Strategy Generation
If aging is detected, generate the offloading and restart strategy:
```bash
cd ../MOEAD-IFM
Calculate the optimal rejuvenation strategy
python main.py  (Note: Replace with your actual entry script)
```

## Reference
For more design details and the background of this research, please refer to the Feature Request Issue:  https://github.com/kubeedge/kubeedge/issues/6607

The full research prototype and datasets are also available in the author's repository: https://github.com/niuchaoyu/Software-Aging-State-Determination-and-Rejuvenation-Strategy-Generation-for-the-KubeEdge

Signed-off-by: niuchaoyu <411654689@qq.com>
