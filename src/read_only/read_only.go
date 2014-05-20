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
const NUM_TIMES_TO_RUN_EACH_EXPERIMENT = 1


type ReadOnly struct {
}

func(readOnly* ReadOnly) RunAll(jarDir,outputFolder string) {
    // singleThreadWarmTests(readOnly,jarDir,outputFolder)
    // numThreadsTests(readOnly,jarDir,outputFolder)
    perfNumThreadsTests(readOnly,jarDir,outputFolder)
    // perfGCOffNumThreadsTests(readOnly,jarDir,outputFolder)
    // perfLocksOffNumThreadsTests(readOnly,jarDir,outputFolder)
    // perfWoundWaitNumThreadsTests(readOnly,jarDir,outputFolder)
    // perfReadsOnDifferentObjectsNumThreadsTests(readOnly,jarDir,outputFolder)
    // threadPoolSizeTests(readOnly, jarDir,outputFolder)
    // uuidGenerationTests(readOnly, jarDir,outputFolder)
    // memLeakTests(readOnly,jarDir,outputFolder)
}

func (readOnly* ReadOnly) commonReadOnlyJar(
    perfOn bool, fqJar string, numReads, numThreads uint32,
    opType operationType, persistentThreadPoolSize,
    maxThreadPoolSize uint32, atomIntUUIDGeneration, gcOn, woundWaitOn,
    readsOnOtherAtomNum bool) *ReadOnlyResult {

    var argSlice [] string
    
    if perfOn {
        argSlice = append(
            argSlice,
            []string{"stat","-o",PERF_STAT_OUTPUT_FILENAME,"java"}...)
    }
    
    if !gcOn {
        argSlice = append(argSlice,GC_OFF_MIN_HEAP_SIZE_FLAG)
        argSlice = append(argSlice,GC_OFF_MAX_HEAP_SIZE_FLAG)
    }
    
    argSlice = append(argSlice,[]string{"-jar",fqJar}...)
    argSlice = append(argSlice,readOnly.addReadsPerThread(numReads)...)
    argSlice = append(argSlice,readOnly.addNumThreads(numThreads)...)
    argSlice = append(argSlice,readOnly.addOperationType(opType)...)
    if persistentThreadPoolSize != 0 {
        threadPoolSizesArgs := readOnly.addThreadPoolSizes(
            persistentThreadPoolSize,maxThreadPoolSize)
        argSlice = append(argSlice,threadPoolSizesArgs...)
    }
    if atomIntUUIDGeneration {
        argSlice = append(argSlice,"-a")
    }
    if woundWaitOn {
        argSlice = append(argSlice,"-w")
    }

    if readsOnOtherAtomNum {
        argSlice = append(argSlice,"-oan")
    }
    
    var stdOut bytes.Buffer
    var cmd * exec.Cmd = nil

    if perfOn {
        cmd = exec.Command("perf", argSlice...)
    } else {
        cmd = exec.Command("java", argSlice...)
    }
    cmd.Stdout = &stdOut
	err := cmd.Run()
	if err != nil {
        panic(err)
	}
    // get number of ops from stdout
    outputString := stdOut.String()

    // get results from perf
    var perfOutput * common.PerfOutput = nil
    if perfOn {
        perfStatsByteData, err := ioutil.ReadFile(PERF_STAT_OUTPUT_FILENAME)
        if err != nil {
            panic(err)
        }
        perfStatsString := string(perfStatsByteData[:])
        os.Remove(PERF_STAT_OUTPUT_FILENAME)

        perfOutput = common.ParsePerfOutput(perfStatsString)
    }

    // returns read only resutls
    return testRunOutputToResults(
        outputString,numReads,numThreads,opType,perfOutput)
}

func (readOnly * ReadOnly) perfReadOnlyJar(
    fqJar string, numReads, numThreads uint32, opType operationType,
    persistentThreadPoolSize,maxThreadPoolSize uint32,
    atomIntUUIDGeneration, woundWaitOn,
    readsOnOtherAtomNum bool) * ReadOnlyResult {

    return readOnly.commonReadOnlyJar(
        true, fqJar, numReads, numThreads, opType,persistentThreadPoolSize,
        maxThreadPoolSize,atomIntUUIDGeneration,true,woundWaitOn,
        readsOnOtherAtomNum)
}

/**
Cannot really run gc off, but can set maximum heap size to be so
high that gc is unlikely to run.
*/
func (readOnly * ReadOnly) perfReadOnlyJarGCOff(
    fqJar string, numReads, numThreads uint32, opType operationType,
    persistentThreadPoolSize,maxThreadPoolSize uint32,
    atomIntUUIDGeneration, woundWaitOn,
    readsOnOtherAtomNum  bool) * ReadOnlyResult {

    return readOnly.commonReadOnlyJar(
        true, fqJar, numReads, numThreads, opType,persistentThreadPoolSize,
        maxThreadPoolSize,atomIntUUIDGeneration,false,woundWaitOn,
        readsOnOtherAtomNum)
}

func (readOnly ReadOnly) readOnlyJar (
    fqJar string, numReads, numThreads uint32, opType operationType,
    persistentThreadPoolSize,maxThreadPoolSize uint32,
    atomIntUUIDGeneration, woundWaitOn,
    readsOnOtherAtomNum bool) * ReadOnlyResult {

    return readOnly.commonReadOnlyJar(
        false, fqJar, numReads, numThreads, opType,persistentThreadPoolSize,
        maxThreadPoolSize,atomIntUUIDGeneration,true,woundWaitOn,
        readsOnOtherAtomNum)
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

func (readOnly ReadOnly) addThreadPoolSizes(
    persistentThreadPoolSize,maxThreadPoolSize uint32) [] string {
    
    return [] string {
        "-p",strconv.FormatUint(uint64(persistentThreadPoolSize),10),
        "-m",strconv.FormatUint(uint64(maxThreadPoolSize),10),
    }
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
