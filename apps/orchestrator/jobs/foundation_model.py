from dagster import job, op


@op
def resolve_dataset() -> str:
    return "dataset-snapshot"


@op
def prepare_dataset(snapshot: str) -> str:
    return f"prepared:{snapshot}"


@op
def pretrain(prepared: str) -> str:
    return f"checkpoint:{prepared}"


@op
def fine_tune(checkpoint: str) -> str:
    return f"model:{checkpoint}"


@op
def register_model(model: str) -> str:
    return f"registered:{model}"


@job
def foundation_model_job():
    register_model(fine_tune(pretrain(prepare_dataset(resolve_dataset()))))
