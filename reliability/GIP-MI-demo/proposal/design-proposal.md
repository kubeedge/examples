---
title: Software Aging Prediction and Rejuvenation for Edge Nodes
authors:
  - "@niuchaoyu"
owning-sig: sig-node
participating-sigs:
  - sig-ai
status: provisional
creation-date: '2025-01-06'
---

# Software Aging Prediction and Rejuvenation for Edge Nodes

## Summary

This proposal introduces a reliability enhancement framework (GIP-MI) for KubeEdge to proactively detect software aging and perform rejuvenation. It integrates a **GCN-Informer** model for metric prediction and a **ParNet** model for aging state determination. Upon detection of aging, a **MOEA/D-IFM** based strategy triggers task offloading and node rejuvenation to ensure service continuity.

## Motivation

### Goals

* Enable KubeEdge to predict resource exhaustion trends (CPU/Memory) on edge nodes.
* Provide a mechanism to identify "aging" states before a hard failure (Crash/NotReady) occurs.
* Implement a rejuvenation strategy that migrates workloads before restarting components, minimizing downtime.

### Non-Goals

* This proposal does not replace the existing node health check mechanism (NodeStatus) but enhances it with predictive capabilities.

## Proposal

### User Stories

* **As a cluster administrator**, I want to receive alerts when an edge node is predicted to exhaust its memory in the next hour, so I can intervene early.
* **As an edge application developer**, I want my services to be gracefully migrated to other nodes if the current node needs a restart due to software aging.

### Design Details

#### 1. Architecture
The framework consists of three main components interacting with KubeEdge:

* **Data Collection**: Scrapes metrics (CPU, Memory, Response Time) from EdgeCore via Prometheus.
* **Decision Engine**:
    * **Prediction**: Runs GCN-Informer to forecast future metric trends.
    * **Determination**: Runs ParNet to classify the node state (Normal vs. Aging).
* **Execution**:
    * If "Aging" is detected, the MOEA/D-IFM algorithm calculates the optimal offloading plan.
    * Commands are sent to KubeEdge (via Kubernetes API) to evict pods and restart the `edgecore` service or the node.

#### 2. Algorithms
* **GCN-Informer**: Captures spatial-temporal dependencies in system metrics for accurate long-term forecasting.
* **ParNet**: A parallel network model used to identify aging patterns from time-series data slices.
* **MOEA/D-IFM**: A multi-objective evolutionary algorithm to optimize task offloading decisions, balancing latency and resource usage.
