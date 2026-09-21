FROM python:3.12-slim
WORKDIR /app
COPY . /app
RUN pip install uv && uv sync --extra training
ENTRYPOINT ["uv", "run", "model-factory"]
