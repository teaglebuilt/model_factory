import pytest

from model_factory.contracts.artifacts import CheckpointRef, ModelRef
from model_factory.evaluation.evaluator import EvaluationResult
from model_factory.manifests.model import ModelManifest
from model_factory.registry.promotion import approve_candidate, promote_to_production
from model_factory.registry.registry import InMemoryModelRegistry


@pytest.mark.asyncio
async def test_registry_and_promotion_flow() -> None:
    model = ModelRef(
        uri="s3://model-factory/models/coder/0.0.1",
        model_id="coder",
        version="0.0.1",
    )
    checkpoint = CheckpointRef(
        uri="s3://model-factory/checkpoints/run-1/step-100",
        run_id="run-1",
        step=100,
    )

    manifest = ModelManifest(
        model=model,
        run_id="run-1",
        experiment_name="smoke-10m-v001",
        source_checkpoint=checkpoint,
    )

    registry = InMemoryModelRegistry()
    await registry.register(manifest)

    stored = await registry.get("coder", "0.0.1")
    assert stored.model.status == "candidate"

    approved = approve_candidate(
        stored.model,
        EvaluationResult(suite="smoke", passed=True),
    )
    production = promote_to_production(approved)
    assert production.status == "production"
