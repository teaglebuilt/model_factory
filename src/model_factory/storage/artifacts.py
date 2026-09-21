from __future__ import annotations

from pathlib import Path
from typing import Protocol


class ArtifactStore(Protocol):
    async def put_file(self, source: Path, destination: str) -> str: ...
    async def put_text(self, content: str, destination: str) -> str: ...
