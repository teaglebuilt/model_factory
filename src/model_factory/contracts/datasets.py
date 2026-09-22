from __future__ import annotations

from typing import Literal

from pydantic import BaseModel, ConfigDict, Field, model_validator


class DatasetRef(BaseModel):
    """Logical dataset reference supplied by an experiment."""

    model_config = ConfigDict(frozen=True)

    uri: str = Field(
        description=(
            "Logical dataset URI, normally platform-data://<project>/<dataset>@<version>."
        )
    )


class DatasetProvenance(BaseModel):
    """Optional upstream lineage preserved by platform-data."""

    model_config = ConfigDict(frozen=True)

    producer: str = "platform-data"
    source_repository: str | None = None
    source_revision: str | None = None
    lakefs_repository: str | None = None
    lakefs_commit: str | None = None
    lakefs_path: str | None = None


class ResolvedDatasetRef(BaseModel):
    """Immutable physical dataset reference used by a training run."""

    model_config = ConfigDict(frozen=True)

    logical_uri: str
    source_type: Literal["dvc", "lakefs", "huggingface", "object-store"]
    physical_uri: str
    revision: str = Field(min_length=7)
    digest: str | None = None
    provenance: DatasetProvenance = Field(default_factory=DatasetProvenance)

    @model_validator(mode="after")
    def reject_mutable_revisions(self) -> "ResolvedDatasetRef":
        if self.revision.lower() in {"main", "master", "latest", "head"}:
            raise ValueError("resolved dataset revisions must be immutable")
        return self
