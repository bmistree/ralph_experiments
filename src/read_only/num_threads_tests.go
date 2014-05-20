package read_only

import "fmt"
import "path/filepath"
import "strconv"

const NUM_THREADS_OUTPUT_NAME = "num_threads.csv"
const PERF_NUM_THREADS_OUTPUT_NAME = "perf_num_threads.csv"
const GC_OFF_PERF_NUM_THREADS_OUTPUT_NAME = "gc_off_perf_num_threads.csv"
const LOCKS_OFF_PERF_NUM_THREADS_OUTPUT_NAME = "locks_off_perf_num_threads.csv"
const WOUND_WAIT_PERF_NUM_THREADS_OUTPUT_NAME =
    "wound_wait_perf_num_threads.csv"
const READS_ON_DIFFERENT_OBJECTS_PERF_NUM_THREADS_OUTPUT_NAME =
    "reads_on_different_objects.csv"
const NUM_THREADS_TEST_NUM_READS = 100000
var NUM_THREADS_TEST_NUM_THREADS [4]uint32 = [4]uint32{1,2,5,10}
// var NUM_THREADS_TEST_NUM_THREADS [1]uint32 = [1]uint32{1}
// var NUM_THREADS_TEST_NUM_THREADS [2]uint32 = [2]uint32{1,2}

func numThreadsTests(readOnly* ReadOnly,jarDir,outputFolder string) {
    outputFilename := filepath.Join(outputFolder,NUM_THREADS_OUTPUT_NAME)
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)

    params := createDefaultParameter()
    commonNumThreadsTests(
        readOnly,fqJar,outputFilename,ALL_OPERATION_TYPES,params,
        "numThreadsTest")
}

func perfNumThreadsTests(readOnly* ReadOnly,jarDir,outputFolder string) {
    params := createDefaultParameter()
    params.perfOn = true
    outputFilename := filepath.Join(outputFolder,PERF_NUM_THREADS_OUTPUT_NAME)
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    commonNumThreadsTests(
        readOnly,fqJar,outputFilename,ALL_OPERATION_TYPES,params,
        "perfNumThreadsTest")
}

func perfWoundWaitNumThreadsTests(
    readOnly* ReadOnly,jarDir,outputFolder string) {
    params := createDefaultParameter()
    params.perfOn = true
    params.woundWaitOn = true
    outputFilename :=
        filepath.Join(outputFolder,WOUND_WAIT_PERF_NUM_THREADS_OUTPUT_NAME)
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    commonNumThreadsTests(
        readOnly,fqJar,outputFilename, []operationType{READ_ATOM_NUM},
        params,"perfWoundWaitNumThreadsTest")
}

func perfReadsOnDifferentObjectsNumThreadsTests(
    readOnly* ReadOnly,jarDir,outputFolder string) {

    params := createDefaultParameter()
    params.perfOn = true
    params.readsOnOtherAtomNum = true
    
    outputFilename :=
        filepath.Join(outputFolder,WOUND_WAIT_PERF_NUM_THREADS_OUTPUT_NAME)
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    commonNumThreadsTests(
        readOnly,fqJar,outputFilename,[]operationType{READ_ATOM_NUM},
        params,"perfReadsOnDifferentObjectsNumThreadsTest")
}

func perfGCOffNumThreadsTests(
    readOnly* ReadOnly,jarDir,outputFolder string) {

    params := createDefaultParameter()
    params.perfOn = true
    params.gcOn = false

    outputFilename :=
        filepath.Join(outputFolder,GC_OFF_PERF_NUM_THREADS_OUTPUT_NAME)
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    commonNumThreadsTests(
        readOnly,fqJar,outputFilename,[]operationType{READ_ATOM_NUM},
        params,"perfGCOffNumThreadsTest")
}

func perfLocksOffNumThreadsTests(
    readOnly* ReadOnly,jarDir,outputFolder string) {

    params := createDefaultParameter()
    params.perfOn = true
    
    outputFilename :=
        filepath.Join(outputFolder,LOCKS_OFF_PERF_NUM_THREADS_OUTPUT_NAME)
    fqJar := filepath.Join(jarDir,LOCKS_OFF_READ_ONLY_JAR_NAME)
    commonNumThreadsTests(
        readOnly,fqJar,outputFilename,[]operationType{READ_ATOM_NUM},
        params,"perfLocksOffNumThreadsTest")
}

func commonNumThreadsTests(
    readOnly* ReadOnly,fqJar,outputFilename string,
    opsToRun []operationType, params * Parameter, testDescription string) {


    // describe test that's running to the user
    fmt.Println(testDescription)

    params.numReads = NUM_THREADS_TEST_NUM_READS
        
    var results [] * ReadOnlyResult
    for i := 0; i < NUM_TIMES_TO_RUN_EACH_EXPERIMENT; i++ {
        fmt.Println(
            "\t" + strconv.Itoa(i+1) + " of " +
            strconv.Itoa(NUM_TIMES_TO_RUN_EACH_EXPERIMENT));
        
        for _,numThreads := range NUM_THREADS_TEST_NUM_THREADS {
            params.numThreads = numThreads
            // try thread size tests across all operation types
            for _, opType := range opsToRun {
                params.opType = opType
                result := readOnly.commonReadOnlyJar(fqJar,params)
                results = append(results,result)
            }
        }
    }
    
    // write results to file
    resultsToFile(results,outputFilename)
}