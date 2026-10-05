from __future__ import annotations

from model_factory.contracts.artifacts import ModelRef
from model_factory.evaluation.evaluator import EvaluationResult


def approve_candidate(model: ModelRef, evaluation: EvaluationResult) -> ModelRef:
    if not evaluation.passed:
        raise ValueError("model cannot be approved when evaluation has failed")
    return model.model_copy(update={"status": "approved"})


def promote_to_production(model: ModelRef) -> ModelRef:
    if model.status != "approved":
        raise ValueError("only approved models can be promoted to production")
    return model.model_copy(update={"status": "production"})
