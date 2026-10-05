from pathlib import Path

import pytest
from pydantic import ValidationError

from model_factory.contracts.datasets import ResolvedDatasetRef
from model_factory.specs.experiment import ExperimentSpec


def test_smoke_experiment_loads() -> None:
    spec = ExperimentSpec.from_yaml(
        Path("experiments/iteration-01/smoke-10m.yaml")
    )

    assert spec.metadata.name == "smoke-10m-v001"
    assert spec.datasets.pretraining.uri.startswith("platform-data://")
    assert spec.pretraining.max_steps == 100
    assert spec.runtime.resources.gpu == 1


def test_resolved_dataset_rejects_mutable_revision() -> None:
    with pytest.raises(ValidationError):
        ResolvedDatasetRef(
            logical_uri="platform-data://foundation-models/code-pretrain@v1",
            source_type="lakefs",
            physical_uri="s3://lakehouse/main/code-pretrain",
            revision="main",
        )
