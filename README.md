# Model Factory

Application framework for building, evaluating, versioning, and promoting custom language models on the teaglebuilt AI platform.

This repository owns **model-building behavior**. The homelab repository owns Kubernetes/Talos, GPU runtime, networking, storage classes, gateways, and shared observability.

## v0.1 pipeline

`ExperimentSpec -> dataset snapshot -> preparation/tokenization -> pretraining -> checkpoint -> SFT -> registration -> serving handoff`

The framework keeps infrastructure-specific code behind interfaces so training logic can run locally or on Kubernetes.

## Local development

```bash
uv sync --all-extras
uv run model-factory validate experiments/iteration-01/coder-150m.yaml
uv run pytest
uv run dagster dev -m apps.orchestrator.definitions
```
