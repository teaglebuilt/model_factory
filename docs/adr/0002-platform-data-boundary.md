# ADR 0002: platform-data owns dataset production and versioning

Status: accepted

## Context

The platform-data repository already owns Spark, Iceberg, lakeFS, SeaweedFS, DVC, data quality,
and Dagster-managed Write-Audit-Publish workflows.

Duplicating lakeFS or source-ingestion logic inside Model Factory creates two owners for dataset
lifecycle.

## Decision

Model Factory consumes logical `DatasetRef` values and requires them to be resolved to immutable
`ResolvedDatasetRef` values before training begins.

The default path is a curated DVC dataset produced by platform-data. Direct immutable lakeFS
snapshots may be supported for large datasets where copying into DVC would be wasteful, but
Model Factory never manages lakeFS branches or mutable refs.
