# ADR 0004: immutable dataset contracts

Status: accepted

Experiments use logical dataset URIs so model configuration is not coupled to physical storage.

Example:

```text
platform-data://foundation-models/code-pretrain@v1
```

Before execution, the logical reference must resolve to a physical URI plus immutable revision.
Mutable revisions such as `main`, `master`, `latest`, or `HEAD` are rejected by the resolved
dataset contract.

Run manifests persist the resolved reference so every model can be traced back to the exact
dataset bytes that produced it.
