from __future__ import annotations

from datetime import datetime, timezone

from pydantic import BaseModel, Field

from model_factory.contracts.artifacts import CheckpointRef


class CheckpointManifest(BaseModel):
    checkpoint: CheckpointRef
    experiment_hash: str
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    metrics: dict[str, float] = Field(default_factory=dict)
