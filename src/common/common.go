package common

type ExperimentModuler interface {
    Description() string
    RunAll(jar_dir, outputFolder string)
}
