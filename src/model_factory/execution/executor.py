from __future__ import annotations

from dataclasses import dataclass
from typing import Protocol


@dataclass(frozen=True)
class JobResources:
    gpu: int = 0
    cpu: int = 1
    memory: str = "2Gi"


@dataclass(frozen=True)
class JobSpec:
    name: str
    image: str
    command: list[str]
    env: dict[str, str]
    resources: JobResources


class Executor(Protocol):
    async def submit(self, job: JobSpec) -> str: ...
    async def status(self, job_id: str) -> str: ...
