# ADR 0003: Pluggable training backends

Status: accepted

Pretraining and fine-tuning implementations sit behind protocols. The initial direction is TorchTitan for pretraining and TRL for SFT, without coupling orchestration to either implementation.
