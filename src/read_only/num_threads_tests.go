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
    commonNumThreadsTests(
        readOnly,fqJar,outputFilename,false,true,false,false,
        ALL_OPERATION_TYPES,"numThreadsTest")
}

func perfNumThreadsTests(readOnly* ReadOnly,jarDir,outputFolder string) {
    outputFilename := filepath.Join(outputFolder,PERF_NUM_THREADS_OUTPUT_NAME)
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    commonNumThreadsTests(
        readOnly,fqJar,outputFilename,true,true,false,false,
        ALL_OPERATION_TYPES,"perfNumThreadsTest")
}

func perfWoundWaitNumThreadsTests(
    readOnly* ReadOnly,jarDir,outputFolder string) {

    outputFilename :=
        filepath.Join(outputFolder,WOUND_WAIT_PERF_NUM_THREADS_OUTPUT_NAME)
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    commonNumThreadsTests(
        readOnly,fqJar,outputFilename,true,false,true,false,
        []operationType{READ_ATOM_NUM},"perfWoundWaitNumThreadsTest")
}

func perfReadsOnDifferentObjectsNumThreadsTests(
    readOnly* ReadOnly,jarDir,outputFolder string) {

    outputFilename :=
        filepath.Join(outputFolder,WOUND_WAIT_PERF_NUM_THREADS_OUTPUT_NAME)
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    commonNumThreadsTests(
        readOnly,fqJar,outputFilename,true,false,true,true,
        []operationType{READ_ATOM_NUM},
        "perfReadsOnDifferentObjectsNumThreadsTest")
}


func perfGCOffNumThreadsTests(
    readOnly* ReadOnly,jarDir,outputFolder string) {

    outputFilename :=
        filepath.Join(outputFolder,GC_OFF_PERF_NUM_THREADS_OUTPUT_NAME)
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    commonNumThreadsTests(
        readOnly,fqJar,outputFilename,true,false,false,false,
        []operationType{READ_ATOM_NUM},"perfGCOffNumThreadsTest")
}

func perfLocksOffNumThreadsTests(
    readOnly* ReadOnly,jarDir,outputFolder string) {

    outputFilename :=
        filepath.Join(outputFolder,LOCKS_OFF_PERF_NUM_THREADS_OUTPUT_NAME)
    fqJar := filepath.Join(jarDir,LOCKS_OFF_READ_ONLY_JAR_NAME)
    commonNumThreadsTests(
        readOnly,fqJar,outputFilename,true,false,false,false,
        []operationType{READ_ATOM_NUM},"perfLocksOffNumThreadsTest")
}

func commonNumThreadsTests(
    readOnly* ReadOnly,fqJar,outputFilename string, perfTest,
    gcOn, woundWaitOn, readsOnOtherAtomNum bool, opsToRun []operationType,
    testDescription string) {

    // describe test that's running to the user
    fmt.Println(testDescription)
    
    var results [] * ReadOnlyResult
    for i := 0; i < NUM_TIMES_TO_RUN_EACH_EXPERIMENT; i++ {
        fmt.Println(
            "\t" + strconv.Itoa(i+1) + " of " +
            strconv.Itoa(NUM_TIMES_TO_RUN_EACH_EXPERIMENT));
        
        for _,numThreads := range NUM_THREADS_TEST_NUM_THREADS {
            // try thread size tests across all operation types
            for _, opType := range opsToRun {

                var result *ReadOnlyResult
                if perfTest {

                    if gcOn {
                        result = readOnly.perfReadOnlyJar(
                            fqJar,NUM_THREADS_TEST_NUM_READS,numThreads,
                            opType,0,0,false,woundWaitOn,readsOnOtherAtomNum)
                    } else {
                        result = readOnly.perfReadOnlyJarGCOff(
                            fqJar,NUM_THREADS_TEST_NUM_READS,numThreads,
                            opType,0,0,false,woundWaitOn,readsOnOtherAtomNum)
                    }
                } else {
                    result = readOnly.readOnlyJar(
                        fqJar,NUM_THREADS_TEST_NUM_READS,numThreads,
                        opType,0,0,false,woundWaitOn,readsOnOtherAtomNum)
                }
                results = append(results,result)
            }
        }
    }
    
    // write results to file
    resultsToFile(results,outputFilename)
}