# Architecture

The Model Factory is the application layer for turning immutable dataset references into evaluated,
promotable model artifacts.

## Ownership boundaries

### homelab / platform infrastructure

Owns Kubernetes/Talos, NVIDIA runtime, GPU device plugin, storage classes, cluster networking,
secrets, shared observability, and the shared Dagster deployment.

### platform-data

Owns source ingestion, Spark/Iceberg/lakeFS/SeaweedFS/DVC, data-quality checks, and publication of
curated training/evaluation datasets.

### model_factory

Owns experiment specifications, training-time preparation, training backends, checkpoint lineage,
evaluation, model registration, and promotion.

## Data contract

Model Factory experiments reference logical datasets:

```text
platform-data://foundation-models/code-pretrain@v1
```

A platform-data integration resolves the logical URI to a `ResolvedDatasetRef`. The resolved
reference must contain an immutable revision. The training package does not own lakeFS branch
lifecycle or DVC publication.

## Execution model

```text
shared Dagster deployment
        |
        | K8sRunLauncher
        v
Dagster run Job
        |
        | k8s_job_executor
        +-------------------+
        |                   |
        v                   v
CPU preparation Job     GPU training Job
                            |
                            v
                        checkpoint
                            |
                            v
                       evaluation
                            |
                            v
                         registry
```

Dagster decides dependency order, retries, schedules, and step configuration. Kubernetes schedules
and executes pods. TorchTitan/TRL perform model computation.

There is intentionally no custom Kubernetes executor inside the Model Factory package.

## Primary flow

1. Parse and validate an `ExperimentSpec`.
2. Resolve each logical dataset reference to an immutable `ResolvedDatasetRef`.
3. Apply training-specific tokenization/packing/mixing.
4. Run pretraining.
5. Persist checkpoints and a `RunManifest`.
6. Optionally run supervised fine-tuning.
7. Evaluate the resulting candidate.
8. Register a `ModelManifest`.
9. Promote only candidates that pass the configured evaluation gate.
