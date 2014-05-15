package read_only

import "strconv"
import "fmt"

/**
 Runs all read only performance tests for ralph
 */

const NUM_WARM_OPS = 600000

type operationType uint32
const (
    READ_ATOM_NUM operationType = iota
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
}

func(readOnly ReadOnly) RunAll(jar,outputFolder string) {
    readOnly.singleThreadWarmTests(jar,outputFolder)
}

func (readOnly ReadOnly) singleThreadWarmTests(jar,outputFolder string) {
    execString := "java -jar " + jar

    execString = readOnly.addReadsPerThread(execString, 5000)
    execString = readOnly.addNumThreads(execString, 1)
    execString = readOnly.addOperationType(execString,READ_ATOM_NUM)

    fmt.Println("\n\n")
    fmt.Println(execString)
    fmt.Println("\n\n")
}

func (readOnly ReadOnly) addReadsPerThread(execString string, numOps uint32) string {
    return execString + " " + "-r " + strconv.Itoa(numOps)
}

func (readOnly ReadOnly) addNumThreads(execString string, numThreads uint32) string {
    return execString + " " + "-t " + strconv.Itoa(numThreads)
}

func (readOnly ReadOnly) addOperationType(execString string, op operationType) string {

    if op == READ_ATOM_NUM {
        return execString + " " + READ_ATOM_NUM_ARG
    }
    
    if op == READ_ATOM_MAP {
        return execString + " " + READ_ATOM_MAP_ARG
    }

    if op == READ_MAP {
        return execString + " " + READ_MAP_ARG
    }
    
    if op == READ_NUM {
        return execString + " " + READ_NUM_ARG
    }

    panic("Unknown op in read_only.go")
}
