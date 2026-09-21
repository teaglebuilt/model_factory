from model_factory.specs.experiment import ExperimentSpec


def test_iteration_one_spec_loads():
    spec = ExperimentSpec.from_yaml("experiments/iteration-01/coder-150m.yaml")
    assert spec.metadata.name == "coder-150m-v001"
    assert spec.runtime.resources.gpu == 1
    assert spec.pretraining.backend == "torchtitan"
