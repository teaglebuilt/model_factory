from __future__ import annotations

from typing import Literal

from pydantic import BaseModel


class PackingSpec(BaseModel):
    strategy: Literal["none", "concat", "best-fit"] = "concat"
    sequence_length: int = 4096


class PreparationSpec(BaseModel):
    packing: PackingSpec = PackingSpec()
    shuffle: bool = True
