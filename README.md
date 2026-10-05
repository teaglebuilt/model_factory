# Model Factory

Application framework for building, evaluating, versioning, and promoting custom language models on the teaglebuilt AI platform.

## Repository boundary

Model Factory owns model-building behavior:

```text
DatasetRef
   ↓
ExperimentSpec
   ↓
preparation / tokenization / packing
   ↓
pretraining
   ↓
checkpoint
   ↓
fine-tuning
   ↓
evaluation
   ↓
model registry
   ↓
promotion
```

It deliberately does **not** own the Kubernetes cluster, GPU runtime, lakeFS, Spark, Iceberg,
SeaweedFS, DVC, or source-data ingestion.

- `platform-data` produces curated, versioned datasets and exposes them to this package through
  immutable dataset references.
- Dagster orchestrates Model Factory runs.
- Kubernetes is the execution substrate for Dagster run/step jobs.
- Training backends such as TorchTitan and TRL remain independently executable.
- Serving is downstream of the registry/promotion boundary.

## Dataset contract

Experiments refer to datasets using logical platform-data URIs:

```yaml
datasets:
  pretraining:
    uri: platform-data://foundation-models/code-pretrain@v1
```

Before training, a resolver must turn that logical reference into a `ResolvedDatasetRef` containing
an immutable revision and physical URI. Training code consumes only resolved references.

Direct lakeFS references are supported as an escape hatch for very large immutable snapshots, but
lakeFS lifecycle management remains owned by `platform-data`.

## Dagster + Kubernetes

Production runs are launched by the shared Dagster deployment on Kubernetes. The Dagster
`K8sRunLauncher` isolates each Dagster run in a Kubernetes Job; `k8s_job_executor` can isolate
individual steps in their own Jobs and apply per-step CPU/GPU resource requirements.

Model Factory does not contain a second Kubernetes scheduler.

## Local development

```bash
uv sync --all-extras
uv run model-factory validate experiments/iteration-01/smoke-10m.yaml
uv run pytest
uv run dagster dev -m apps.orchestrator.definitions
```
