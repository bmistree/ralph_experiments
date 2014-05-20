package read_only

import "path/filepath"
import "fmt"
import "strconv"

const THREAD_POOL_SIZE_OUTPUT_NAME = "thread_pool_size.csv"
const THREAD_POOL_SIZE_TEST_NUM_READS = 100000

var THREAD_POOL_SIZES [6]uint32 = [6]uint32{5,10,20,40,60,80}
var NUM_THREADS [2] uint32 = [2] uint32 {1,2}

const THREAD_POOL_SIZE_OP_TYPE = READ_ATOM_NUM 


func threadPoolSizeTests(jarDir,outputFolder string) {
    
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    fmt.Println("Running thread pool size experiment: ")

    params := createDefaultParameter()
    params.opType = THREAD_POOL_SIZE_OP_TYPE
    params.numReads = NUM_THREADS_TEST_NUM_READS
    params.perfOn = true
    
    var results [] * ReadOnlyResult
    for i := 0; i < NUM_TIMES_TO_RUN_EACH_EXPERIMENT; i++ {
        fmt.Println(
            "\t" + strconv.Itoa(i+1) + " of " +
            strconv.Itoa(NUM_TIMES_TO_RUN_EACH_EXPERIMENT));
        
        for _,numThreads := range NUM_THREADS {
            for _, threadPoolSize := range THREAD_POOL_SIZES {

                params.numThreads = numThreads
                params.persistentThreadPoolSize = threadPoolSize
                params.maxThreadPoolSize = threadPoolSize
                
                result := commonReadOnlyJar(fqJar,params)
                results = append(results,result)
            }
        }
    }
    // write results to file
    resultsToFile(
        results,filepath.Join(outputFolder,THREAD_POOL_SIZE_OUTPUT_NAME))
}
