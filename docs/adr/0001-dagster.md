# ADR 0001: Dagster for orchestration

Status: accepted

## Decision

Use the shared Dagster deployment as the orchestration control plane for Model Factory.

Production runs execute on Kubernetes. The Dagster `K8sRunLauncher` provides run isolation and
`k8s_job_executor` provides step-level Kubernetes Jobs and per-step resource configuration.

Model Factory does not implement a second Kubernetes scheduler or job state machine.

Training engines remain independently executable and do not import Dagster.
