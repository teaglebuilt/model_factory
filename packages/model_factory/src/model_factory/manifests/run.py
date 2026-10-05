from __future__ import annotations

from datetime import datetime, timezone

from pydantic import BaseModel, Field

from model_factory.contracts.datasets import ResolvedDatasetRef


class RunManifest(BaseModel):
    run_id: str
    experiment_name: str
    experiment_hash: str
    git_sha: str | None = None
    container_digest: str | None = None
    datasets: dict[str, ResolvedDatasetRef]
    tokenizer: str
    architecture: str
    seed: int
    training_backend: str
    gpu_type: str | None = None
    cuda_version: str | None = None
    torch_version: str | None = None
    started_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    finished_at: datetime | None = None
    parent_run_id: str | None = None
