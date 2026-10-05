from __future__ import annotations

from typing import Protocol

from model_factory.contracts.artifacts import ModelRef
from model_factory.manifests.model import ModelManifest


class ModelRegistry(Protocol):
    async def register(self, manifest: ModelManifest) -> ModelRef:
        ...

    async def get(self, model_id: str, version: str) -> ModelManifest:
        ...


class InMemoryModelRegistry:
    """Small reference implementation used by tests and local development."""

    def __init__(self) -> None:
        self._models: dict[tuple[str, str], ModelManifest] = {}

    async def register(self, manifest: ModelManifest) -> ModelRef:
        key = (manifest.model.model_id, manifest.model.version)
        self._models[key] = manifest
        return manifest.model

    async def get(self, model_id: str, version: str) -> ModelManifest:
        return self._models[(model_id, version)]
