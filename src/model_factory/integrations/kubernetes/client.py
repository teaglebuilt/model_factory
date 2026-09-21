from __future__ import annotations

from kubernetes import client, config

from model_factory.execution.executor import JobSpec


class KubernetesExecutor:
    def __init__(self, namespace: str = "ai") -> None:
        self.namespace = namespace
        try:
            config.load_incluster_config()
        except config.ConfigException:
            config.load_kube_config()
        self.batch = client.BatchV1Api()

    async def submit(self, job: JobSpec) -> str:
        container = client.V1Container(
            name="workload",
            image=job.image,
            command=job.command,
            env=[client.V1EnvVar(name=k, value=v) for k, v in job.env.items()],
            resources=client.V1ResourceRequirements(
                requests={"cpu": str(job.resources.cpu), "memory": job.resources.memory},
                limits={
                    "cpu": str(job.resources.cpu),
                    "memory": job.resources.memory,
                    **({"nvidia.com/gpu": str(job.resources.gpu)} if job.resources.gpu else {}),
                },
            ),
        )
        body = client.V1Job(
            metadata=client.V1ObjectMeta(name=job.name),
            spec=client.V1JobSpec(
                backoff_limit=0,
                template=client.V1PodTemplateSpec(
                    metadata=client.V1ObjectMeta(labels={"app.kubernetes.io/part-of": "model-factory"}),
                    spec=client.V1PodSpec(restart_policy="Never", containers=[container]),
                ),
            ),
        )
        self.batch.create_namespaced_job(namespace=self.namespace, body=body)
        return job.name

    async def status(self, job_id: str) -> str:
        job = self.batch.read_namespaced_job(job_id, self.namespace)
        if job.status.succeeded:
            return "succeeded"
        if job.status.failed:
            return "failed"
        if job.status.active:
            return "running"
        return "pending"
