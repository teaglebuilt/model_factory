from __future__ import annotations

from typing import Literal
from pydantic import BaseModel


class ModelArtifact(BaseModel):
    model_id: str
    run_id: str
    experiment_name: str
    artifact_uri: str
    dataset_ref: str
    architecture: str
    tokenizer: str
    status: Literal["candidate", "accepted", "rejected", "production"] = "candidate"


class ModelRegistry:
    def __init__(self) -> None:
        self._models: dict[str, ModelArtifact] = {}

    def register(self, artifact: ModelArtifact) -> ModelArtifact:
        self._models[artifact.model_id] = artifact
        return artifact

    def get(self, model_id: str) -> ModelArtifact:
        return self._models[model_id]
