from __future__ import annotations

from dataclasses import dataclass
from typing import Protocol


@dataclass(frozen=True)
class DatasetSnapshot:
    repository: str
    commit: str
    path: str
    uri: str


class DatasetStore(Protocol):
    async def resolve(self, repository: str, ref: str, path: str) -> DatasetSnapshot: ...
