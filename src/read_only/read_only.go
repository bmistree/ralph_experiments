package read_only

import "strconv"
import "os/exec"
import "os"
import "bytes"
import "strings"
import "common"
import "io/ioutil"


/**
 Runs all read only performance tests for ralph
 */

const READ_ONLY_JAR_NAME = "read_perf.jar"
const LOCKS_OFF_READ_ONLY_JAR_NAME = "locks_off_read_perf.jar"
const STACKED_READ_ONLY_JAR_NAME = "logging_read_perf.jar"

type operationType uint32
const (
    READ_ATOM_NUM operationType = iota
    READ_ATOM_MAP
    READ_NUM
    READ_MAP
)
var ALL_OPERATION_TYPES =
    []operationType {READ_ATOM_NUM,READ_NUM,READ_ATOM_MAP,READ_MAP}


// For gc off experiments, need to be able to run with such a large
// heap that garbage collection will be turned off.
const GC_OFF_MIN_HEAP_SIZE_FLAG = "-Xms1000m"
const GC_OFF_MAX_HEAP_SIZE_FLAG = "-Xmx1000m"

// Corresponding command line args for operation types
const READ_ATOM_NUM_ARG = "-an"
const READ_ATOM_MAP_ARG = "-am"
const READ_NUM_ARG = "-nan"
const READ_MAP_ARG = "-nam"

const PERF_STAT_OUTPUT_FILENAME = "perf_stats.txt"
const NUM_TIMES_TO_RUN_EACH_EXPERIMENT = 10


func RunAll(jarDir,outputFolder string) {
    singleThreadWarmTests(jarDir,outputFolder)
    numThreadsTests(jarDir,outputFolder)
    perfNumThreadsTests(jarDir,outputFolder)
    perfGCOffNumThreadsTests(jarDir,outputFolder)
    perfLocksOffNumThreadsTests(jarDir,outputFolder)
    perfWoundWaitNumThreadsTests(jarDir,outputFolder)
    perfReadsOnDifferentObjectsNumThreadsTests(jarDir,outputFolder)
    threadPoolSizeTests(jarDir,outputFolder)
    uuidGenerationTests(jarDir,outputFolder)
    memLeakTests(jarDir,outputFolder)
}


/**
 @returns --- Slice containing all of arguments to pass into exec.
 */
func argBuilder (fqJar string, params * Parameter) [] string {
    var argSlice [] string
    
    if params.perfOn {
        argSlice = append(
            argSlice,
            []string{"perf","stat","-o",PERF_STAT_OUTPUT_FILENAME,"java"}...)
    } else {
        argSlice = append(argSlice,"java")
    }
    
    if !params.gcOn {
        argSlice = append(argSlice,GC_OFF_MIN_HEAP_SIZE_FLAG)
        argSlice = append(argSlice,GC_OFF_MAX_HEAP_SIZE_FLAG)
    }
    
    argSlice = append(argSlice,[]string{"-jar",fqJar}...)
    argSlice = append(argSlice,addReadsPerThread(params.numReads)...)
    argSlice = append(argSlice,addNumThreads(params.numThreads)...)
    argSlice = append(argSlice,addOperationType(params.opType)...)
    if params.persistentThreadPoolSize != 0 {
        threadPoolSizesArgs := addThreadPoolSizes(
            params.persistentThreadPoolSize,params.maxThreadPoolSize)
        argSlice = append(argSlice,threadPoolSizesArgs...)
    }
    if params.atomIntUUIDGeneration {
        argSlice = append(argSlice,"-a")
    }
    if params.woundWaitOn {
        argSlice = append(argSlice,"-w")
    }

    if params.readsOnOtherAtomNum {
        argSlice = append(argSlice,"-oan")
    }
    return argSlice
}

func baseReadOnlyJar(fqJar string, params * Parameter) (*common.PerfOutput,string) {
    argSlice := argBuilder(fqJar,params)
    var stdOut bytes.Buffer
    cmd := exec.Command(argSlice[0], argSlice[1:]...)
    cmd.Stdout = &stdOut
	err := cmd.Run()
	if err != nil {
        panic(err)
	}
    // get number of ops from stdout
    outputString := stdOut.String()

    // get results from perf
    var perfOutput * common.PerfOutput = nil
    if params.perfOn {
        perfStatsByteData, err := ioutil.ReadFile(PERF_STAT_OUTPUT_FILENAME)
        if err != nil {
            panic(err)
        }
        perfStatsString := string(perfStatsByteData[:])
        os.Remove(PERF_STAT_OUTPUT_FILENAME)

        perfOutput = common.ParsePerfOutput(perfStatsString)
    }

    return perfOutput, outputString
}


func commonReadOnlyJar(fqJar string, params * Parameter) *ReadOnlyResult {

    perfOutput, outputString := baseReadOnlyJar(fqJar,params)
    // returns read only resutls
    return testRunOutputToResults(
        outputString,params.numReads,params.numThreads,params.opType,
        perfOutput)
}

func loggedReadOnlyJar(fqJar string, params * Parameter) *LoggedReadOnlyResult {

    perfOutput, outputString := baseReadOnlyJar(fqJar,params)
    return loggedTestRunOutputToResults(
        outputString, params.numReads,params.numThreads,params.opType,
        perfOutput)
}

func testRunOutputToResults(
    cmdResult string,numReads,numThreads uint32, opType operationType,
    perfOutput * common.PerfOutput) * ReadOnlyResult {
    splitResults := strings.Split(cmdResult," ")
    if len(splitResults) != 2 {
        panic("Expected 2 results when splitting")
    }

    opsPerSecond, err :=
        strconv.ParseFloat(strings.TrimSpace(splitResults[1]),64)
    if err != nil {
        panic(err)
    }
    toReturn := &ReadOnlyResult{
        numReads: numReads,
        numThreads: numThreads,
        opType: opType,
        opsPerSecond: opsPerSecond,
        perfOutput: perfOutput,
    }
    return toReturn
}

func addThreadPoolSizes(
    persistentThreadPoolSize,maxThreadPoolSize uint32) [] string {
    
    return [] string {
        "-p",strconv.FormatUint(uint64(persistentThreadPoolSize),10),
        "-m",strconv.FormatUint(uint64(maxThreadPoolSize),10),
    }
}

func addReadsPerThread(numOps uint32) []string {
    
    return [] string {"-r",strconv.FormatUint(uint64(numOps),10)}
}

func addNumThreads(numThreads uint32) []string {
    return []string{"-t", strconv.FormatUint(uint64(numThreads),10)}
}

func addOperationType(op operationType) []string {

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
