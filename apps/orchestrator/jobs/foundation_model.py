from __future__ import annotations

from dagster import Field, String, job, op
from dagster_k8s import k8s_job_executor

from model_factory.contracts.datasets import DatasetRef


@op(
    config_schema={"dataset_uri": Field(String)},
)
def resolve_dataset(context) -> str:
    """Resolve a logical platform-data dataset URI.

    The real resolver is supplied as a resource/integration by the shared platform.
    For now the op preserves the contract boundary and returns the logical URI.
    """
    ref = DatasetRef(uri=context.op_config["dataset_uri"])
    context.log.info("Resolving dataset %s", ref.uri)
    return ref.uri


@op
def prepare_dataset(dataset_uri: str) -> str:
    context = f"prepared:{dataset_uri}"
    return context


@op(
    tags={
        "dagster-k8s/config": {
            "container_config": {
                "resources": {
                    "requests": {
                        "cpu": "4",
                        "memory": "16Gi",
                        "nvidia.com/gpu": "1",
                    },
                    "limits": {
                        "nvidia.com/gpu": "1",
                    },
                }
            },
            "pod_spec_config": {
                "runtimeClassName": "nvidia",
                "nodeSelector": {"nvidia.com/gpu.present": "true"},
            },
        }
    }
)
def pretrain(prepared_dataset: str) -> str:
    return f"checkpoint:{prepared_dataset}"


@op
def evaluate(checkpoint: str) -> str:
    return f"evaluation:{checkpoint}"


@op
def register_model(evaluation: str) -> str:
    return f"registered:{evaluation}"


@job(executor_def=k8s_job_executor)
def foundation_model_job():
    register_model(evaluate(pretrain(prepare_dataset(resolve_dataset()))))
