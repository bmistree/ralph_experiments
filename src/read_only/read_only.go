package read_only

import "strconv"
import "os/exec"
import "log"
import "bytes"
import "strings"
import "fmt"


/**
 Runs all read only performance tests for ralph
 */

const READ_ONLY_JAR_NAME = "read_perf.jar"


type operationType uint32
const (
    READ_ATOM_NUM operationType = iota
    READ_ATOM_MAP
    READ_NUM
    READ_MAP
)
var ALL_OPERATION_TYPES =
    [4]operationType {READ_ATOM_NUM,READ_NUM,READ_ATOM_MAP,READ_MAP}


// Corresponding command line args for operation types
const READ_ATOM_NUM_ARG = "-an"
const READ_ATOM_MAP_ARG = "-am"
const READ_NUM_ARG = "-nan"
const READ_MAP_ARG = "-nam"

const NUM_TIMES_TO_RUN_EACH_EXPERIMENT = 2

type ReadOnly struct {
}

func(readOnly* ReadOnly) RunAll(jarDir,outputFolder string) {
    singleThreadWarmTests(readOnly,jarDir,outputFolder)
    numThreadsTests(readOnly,jarDir,outputFolder)
    perfNumThreadsTests(readOnly,jarDir,outputFolder)
}

func (readOnly ReadOnly) perfReadOnlyJar(
    fqJar string, numReads, numThreads uint32, opType operationType) ReadOnlyResult {
    
    argSlice := []string {"stat","java","-jar",fqJar}
    argSlice = append(argSlice,readOnly.addReadsPerThread(numReads)...)
    argSlice = append(argSlice,readOnly.addNumThreads(numThreads)...)
    argSlice = append(argSlice,readOnly.addOperationType(opType)...)
    
    var stdOut bytes.Buffer
    var stdErr bytes.Buffer
    cmd := exec.Command("perf", argSlice...)
    cmd.Stdout = &stdOut
    cmd.Stderr = &stdErr
	err := cmd.Run()
	if err != nil {
        log.Fatal(err)
	}

    outputString := stdOut.String()
    perfStatsString := stdErr.String()

    fmt.Println("\n\n")
    fmt.Println(perfStatsString)
    fmt.Println("\n\n")

    return testRunOutputToResults(
        outputString,numReads,numThreads,opType)
}

func (readOnly ReadOnly) readOnlyJar (
    fqJar string, numReads, numThreads uint32, opType operationType) ReadOnlyResult {

    argSlice := []string {"-jar",fqJar}
    argSlice = append(argSlice,readOnly.addReadsPerThread(numReads)...)
    argSlice = append(argSlice,readOnly.addNumThreads(numThreads)...)
    argSlice = append(argSlice,readOnly.addOperationType(opType)...)
    
    var out bytes.Buffer    
    cmd := exec.Command("java", argSlice...)
    cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
        log.Fatal(err)
	}
    return testRunOutputToResults(
        out.String(),numReads,numThreads,opType)
}

func testRunOutputToResults(
    cmdResult string,numReads,numThreads uint32, opType operationType) ReadOnlyResult {
    splitResults := strings.Split(cmdResult," ")
    if len(splitResults) != 2 {
        log.Fatal("Expected 2 results when splitting")
    }

    opsPerSecond, err :=
        strconv.ParseFloat(strings.TrimSpace(splitResults[1]),64)
    if err != nil {
        log.Fatal(err)
    }
    toReturn := ReadOnlyResult{
        numReads: numReads,
        numThreads: numThreads,
        opType: opType,
        opsPerSecond: opsPerSecond,
    }
    return toReturn
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
