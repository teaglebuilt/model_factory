FROM python:3.12-slim
WORKDIR /app
COPY . /app
RUN pip install uv && uv sync
CMD ["uv", "run", "dagster", "api", "grpc", "-m", "apps.orchestrator.definitions"]
