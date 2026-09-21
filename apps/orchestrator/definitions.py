from dagster import Definitions

from apps.orchestrator.jobs.foundation_model import foundation_model_job


defs = Definitions(jobs=[foundation_model_job])
