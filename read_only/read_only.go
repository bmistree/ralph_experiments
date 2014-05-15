package read_only;

/**
 Runs all read only performance tests for ralph
 */

const NUM_WARM_OPS = 600000

type operationType uint32
const (
    READ_ATOM_NUM operation_type = iota
    READ_ATOM_MAP
    READ_NUM
    READ_MAP
)

// Corresponding command line args for operation types
const READ_ATOM_NUM_ARG = "-an"
const READ_ATOM_MAP_ARG = "-am"
const READ_NUM_ARG = "-nan"
const READ_MAP_ARG = "-nam"


type ReadOnly struct {
    addNumOps(cmdString string, numOps uint32) string
    addNumThreads(cmdString string, numThreads uint32) string
    addOperationType(cmdString string, op operationType) string
}

func(readOnly ReadOnly) RunAll(jar,outputFolder string) {
    readOnly.singleThreadWarmTests(jar,outputFolder)
}

func (readOnly ReadOnly) singleThreadWarmTests(jar,outputFolder string) {
    execString := "java -jar " + jar

    execString = readOnly.addNumOps(execString, 5000)
    execString = readOnly.addNumThreads(execString, 1)
    execString = addOperationType(execString,READ_ATOM_NUM_ARG)

    fmt.Println("\n\n")
    fmt.Println(execString)
    fmt.Println("\n\n")
}

