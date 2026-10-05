from __future__ import annotations

from typing import Literal

from pydantic import BaseModel, ConfigDict


class ArtifactRef(BaseModel):
    model_config = ConfigDict(frozen=True)

    uri: str
    digest: str | None = None


class CheckpointRef(ArtifactRef):
    run_id: str
    step: int


class ModelRef(ArtifactRef):
    model_id: str
    version: str
    status: Literal["candidate", "approved", "rejected", "production"] = "candidate"
