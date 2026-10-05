from __future__ import annotations

from pathlib import Path
from typing import Literal

import yaml
from pydantic import BaseModel, Field

from model_factory.contracts.datasets import DatasetRef
from model_factory.preparation.specs import PreparationSpec


class ExperimentDatasets(BaseModel):
    pretraining: DatasetRef
    instruction: DatasetRef | None = None
    evaluation: DatasetRef | None = None


class TokenizerSpec(BaseModel):
    name: str
    vocab_size: int = 32000


class ModelParameters(BaseModel):
    hidden_size: int
    num_layers: int
    num_attention_heads: int
    max_sequence_length: int = 4096


class ModelSpec(BaseModel):
    family: str
    architecture: str
    parameters: ModelParameters


class OptimizerSpec(BaseModel):
    type: str = "adamw"
    learning_rate: float = 3e-4
    weight_decay: float = 0.1


class SchedulerSpec(BaseModel):
    type: str = "cosine"
    warmup_steps: int = 1000


class CheckpointSpec(BaseModel):
    interval: int = 1000


class TrainingSpec(BaseModel):
    backend: str = "torchtitan"
    seed: int = 42
    max_steps: int | None = None
    batch_size: int = 8
    gradient_accumulation_steps: int = 8
    optimizer: OptimizerSpec = Field(default_factory=OptimizerSpec)
    scheduler: SchedulerSpec = Field(default_factory=SchedulerSpec)
    checkpoints: CheckpointSpec = Field(default_factory=CheckpointSpec)


class FineTuneSpec(BaseModel):
    enabled: bool = False
    backend: str = "trl"


class ResourceSpec(BaseModel):
    gpu: int = 1
    cpu: int = 8
    memory: str = "32Gi"


class RuntimeSpec(BaseModel):
    queue: str = "gpu"
    resources: ResourceSpec = Field(default_factory=ResourceSpec)


class OutputSpec(BaseModel):
    checkpoint_store: str
    model_store: str


class Metadata(BaseModel):
    name: str
    description: str | None = None


class ExperimentSpec(BaseModel):
    schema_version: Literal["v1"] = "v1"
    metadata: Metadata
    datasets: ExperimentDatasets
    tokenizer: TokenizerSpec
    preparation: PreparationSpec = Field(default_factory=PreparationSpec)
    model: ModelSpec
    pretraining: TrainingSpec
    finetuning: FineTuneSpec = Field(default_factory=FineTuneSpec)
    runtime: RuntimeSpec = Field(default_factory=RuntimeSpec)
    outputs: OutputSpec

    @classmethod
    def from_yaml(cls, path: str | Path) -> "ExperimentSpec":
        with Path(path).open("r", encoding="utf-8") as handle:
            return cls.model_validate(yaml.safe_load(handle))
