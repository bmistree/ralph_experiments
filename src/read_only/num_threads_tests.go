package read_only

import "fmt"
import "path/filepath"
import "strconv"

const NUM_THREADS_OUTPUT_NAME = "num_threads.csv"
const PERF_NUM_THREADS_OUTPUT_NAME = "perf_num_threads.csv"
const GC_OFF_PERF_NUM_THREADS_OUTPUT_NAME = "gc_off_perf_num_threads.csv"
const NUM_THREADS_TEST_NUM_READS = 100000
// var NUM_THREADS_TEST_NUM_THREADS [4]uint32 = [4]uint32{1,2,5,10}
var NUM_THREADS_TEST_NUM_THREADS [1]uint32 = [1]uint32{1}

func numThreadsTests(readOnly* ReadOnly,jarDir,outputFolder string) {
    outputFilename := filepath.Join(outputFolder,NUM_THREADS_OUTPUT_NAME)
    commonNumThreadsTests(readOnly,jarDir,outputFilename,false,true)
}

func perfNumThreadsTests(readOnly* ReadOnly,jarDir,outputFolder string) {
    outputFilename := filepath.Join(outputFolder,PERF_NUM_THREADS_OUTPUT_NAME)
    commonNumThreadsTests(readOnly,jarDir,outputFilename,true,true)
}

func perfGCOffNumThreadsTests(readOnly* ReadOnly,jarDir,outputFolder string) {
    outputFilename := filepath.Join(outputFolder,GC_OFF_PERF_NUM_THREADS_OUTPUT_NAME)
    commonNumThreadsTests(readOnly,jarDir,outputFilename,true,false)
}

func commonNumThreadsTests(
    readOnly* ReadOnly,jarDir,outputFilename string, perfTest, gcOn bool) {
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    if perfTest {
        fmt.Println("Running perf num threads experiment: ")
    } else {
        fmt.Println("Running num threads experiment: ")
    }
    
    var results [] * ReadOnlyResult
    for i := 0; i < NUM_TIMES_TO_RUN_EACH_EXPERIMENT; i++ {
        fmt.Println(
            "\t" + strconv.Itoa(i+1) + " of " +
            strconv.Itoa(NUM_TIMES_TO_RUN_EACH_EXPERIMENT));
        
        for _,numThreads := range NUM_THREADS_TEST_NUM_THREADS {
            // try thread size tests across all operation types
            for _, opType := range ALL_OPERATION_TYPES {

                var result *ReadOnlyResult
                if perfTest {
                    result = readOnly.perfReadOnlyJar(
                        fqJar,NUM_THREADS_TEST_NUM_READS,numThreads,
                        opType,0,0,false)
                } else {
                    result = readOnly.readOnlyJar(
                        fqJar,NUM_THREADS_TEST_NUM_READS,numThreads,
                        opType,0,0,false)
                }
                results = append(results,result)
            }
        }
    }
    
    // write results to file
    resultsToFile(results,outputFilename)
}