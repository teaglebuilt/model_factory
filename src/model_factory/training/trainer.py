from __future__ import annotations

from typing import Protocol

from model_factory.specs.experiment import ExperimentSpec


class TrainingResult(Protocol):
    checkpoint_uri: str


class TrainerBackend(Protocol):
    async def train(self, experiment: ExperimentSpec, run_id: str) -> TrainingResult: ...
    async def resume(self, experiment: ExperimentSpec, run_id: str, checkpoint_uri: str) -> TrainingResult: ...
