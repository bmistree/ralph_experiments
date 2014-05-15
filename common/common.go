package common

type ExperimentModuler interface {
    Description() string
    RunAll(jar, outputFolder string)
}
