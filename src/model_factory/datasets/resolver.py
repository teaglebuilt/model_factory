from __future__ import annotations

from typing import Protocol

from model_factory.contracts.datasets import DatasetRef, ResolvedDatasetRef


class DatasetResolver(Protocol):
    """Boundary implemented by the platform-data integration."""

    async def resolve(self, ref: DatasetRef) -> ResolvedDatasetRef:
        """Resolve a logical dataset URI to an immutable physical dataset reference."""
        ...
