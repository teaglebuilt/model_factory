from __future__ import annotations

from typing import Protocol

from pydantic import BaseModel, Field

from model_factory.contracts.artifacts import ModelRef
from model_factory.contracts.datasets import ResolvedDatasetRef


class EvaluationRequest(BaseModel):
    model: ModelRef
    dataset: ResolvedDatasetRef | None = None
    suite: str = "smoke"


class EvaluationResult(BaseModel):
    suite: str
    passed: bool
    metrics: dict[str, float] = Field(default_factory=dict)
    checks: dict[str, bool] = Field(default_factory=dict)


class Evaluator(Protocol):
    async def evaluate(self, request: EvaluationRequest) -> EvaluationResult:
        ...
