from __future__ import annotations

from typing import Protocol

from pydantic import BaseModel

from model_factory.contracts.artifacts import CheckpointRef
from model_factory.contracts.datasets import ResolvedDatasetRef
from model_factory.specs.experiment import ExperimentSpec


class TrainingRequest(BaseModel):
    run_id: str
    experiment: ExperimentSpec
    dataset: ResolvedDatasetRef


class TrainingResult(BaseModel):
    run_id: str
    checkpoint: CheckpointRef
    metrics: dict[str, float] = {}


class TrainerBackend(Protocol):
    """Training engine boundary. Dagster orchestrates; backends train."""

    async def train(self, request: TrainingRequest) -> TrainingResult:
        ...

    async def resume(
        self,
        request: TrainingRequest,
        checkpoint: CheckpointRef,
    ) -> TrainingResult:
        ...
