package read_only

import "strconv"
import "fmt"
import "os/exec"
import "log"
import "bytes"
import "path/filepath"

/**
 Runs all read only performance tests for ralph
 */

const NUM_WARM_OPS = 600000
const READ_ONLY_JAR_NAME = "read_perf.jar"

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
var WARM_TEST_NUM_OPS [6]uint32 =
    [6]uint32{1000,5000,10000,50000,100000,150000}

type ReadOnly struct {
}

func(readOnly ReadOnly) RunAll(jar,outputFolder string) {
    readOnly.singleThreadWarmTests(jar,outputFolder)
}

func (readOnly ReadOnly) singleThreadWarmTests(jar_dir,outputFolder string) {
    fqJar := filepath.Join(jar_dir,READ_ONLY_JAR_NAME)

    for _,numReads := range WARM_TEST_NUM_OPS {
        argSlice := []string {"-jar",fqJar}
        argSlice = append(argSlice,readOnly.addReadsPerThread(numReads)...)
        argSlice = append(argSlice,readOnly.addNumThreads(1)...)
        argSlice = append(argSlice,readOnly.addOperationType(READ_NUM)...)
        readOnly.readOnlyJar(argSlice)
    }
}

func (readOnly ReadOnly) readOnlyJar (argSlice []string) {
    var out bytes.Buffer    
    cmd := exec.Command("java", argSlice...)
    cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
        log.Fatal(err)
	}
    fmt.Println("\n\n")    
    fmt.Println(out.String())
    fmt.Println("\n\n")    
}


func (readOnly ReadOnly) addReadsPerThread(numOps uint32) []string {
    
    return [] string {"-r",strconv.FormatUint(uint64(numOps),10)}
}

func (readOnly ReadOnly) addNumThreads(numThreads uint32) []string {
    return []string{"-t", strconv.FormatUint(uint64(numThreads),10)}
}

func (readOnly ReadOnly) addOperationType(op operationType) []string {

    if op == READ_ATOM_NUM {
        return [] string {READ_ATOM_NUM_ARG}
    }
    
    if op == READ_ATOM_MAP {
        return [] string {READ_ATOM_MAP_ARG}
    }

    if op == READ_MAP {
        return [] string {READ_MAP_ARG}
    }
    
    if op == READ_NUM {
        return [] string {READ_NUM_ARG}
    }

    panic("Unknown op in read_only.go")
}
