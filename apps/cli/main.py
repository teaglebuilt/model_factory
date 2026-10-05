"""CLI shim — prefer `model-factory` / `model_factory.cli`."""

from model_factory.cli import app

__all__ = ["app"]

if __name__ == "__main__":
    app()
