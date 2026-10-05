from __future__ import annotations

from datetime import datetime, timezone

from pydantic import BaseModel, Field

from model_factory.contracts.artifacts import CheckpointRef, ModelRef


class ModelManifest(BaseModel):
    model: ModelRef
    run_id: str
    experiment_name: str
    source_checkpoint: CheckpointRef
    evaluation_suite: str | None = None
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
