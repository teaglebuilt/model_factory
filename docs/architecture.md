# Architecture

The Model Factory is the application/control plane for model-building workflows. It deliberately does not own Talos, Kubernetes cluster lifecycle, storage classes, networking, GPU drivers, AI Gateway, or shared observability.

## Dependency rule

Domain/application modules depend on protocols. Integrations implement those protocols.

- training does not import Kubernetes, boto3, lakeFS, Dagster, or Prometheus
- execution adapters translate application jobs into runtime-specific workloads
- dataset adapters resolve immutable dataset versions
- model registration is independent of serving

## Primary flow

1. Parse and validate an ExperimentSpec.
2. Resolve immutable dataset reference.
3. Materialize prepared/tokenized/packed dataset assets.
4. Submit pretraining workload.
5. Track checkpoints and lineage.
6. Optionally submit SFT workload.
7. Evaluate and register the artifact.
8. Promote by policy or explicit approval.
