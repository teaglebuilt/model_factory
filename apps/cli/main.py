from pathlib import Path

import typer
from rich import print

from model_factory.specs.experiment import ExperimentSpec

app = typer.Typer(no_args_is_help=True)


@app.command()
def validate(path: Path) -> None:
    spec = ExperimentSpec.from_yaml(path)
    print(f"[green]valid[/green] {spec.metadata.name}")


@app.command()
def show(path: Path) -> None:
    spec = ExperimentSpec.from_yaml(path)
    print(spec.model_dump_json(indent=2))


if __name__ == "__main__":
    app()
