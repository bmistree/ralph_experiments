package read_only

import "path/filepath"
import "fmt"
import "strconv"

const MEM_LEAK_OUTPUT_NAME = "mem_leak.csv"
const MEM_LEAK_SINGLE_THREAD_NUM_READS = 100000
const MEM_LEAK_OP_TYPE = READ_ATOM_NUM 
var MEM_LEAK_MULTIPLE_DURATIONS_TO_RUN [5] uint32 =
    [5]uint32{1,2,4,8,16}

func memLeakTests(readOnly* ReadOnly, jarDir,outputFolder string) {
    
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    fmt.Println("Running mem leak experiment: ")
    
    var results [] * ReadOnlyResult
    
    // baseline
    baselineResult := readOnly.perfReadOnlyJar(
        fqJar,MEM_LEAK_SINGLE_THREAD_NUM_READS,2,
        MEM_LEAK_OP_TYPE,0,0,false,false,false)
    
    results = append(results,baselineResult)
    
    // comparison to baseline
    for i := 0; i < NUM_TIMES_TO_RUN_EACH_EXPERIMENT; i++ {
        fmt.Println(
            "\t" + strconv.Itoa(i+1) + " of " +
            strconv.Itoa(NUM_TIMES_TO_RUN_EACH_EXPERIMENT));
        
        for _, numOpsMultiplier := range MEM_LEAK_MULTIPLE_DURATIONS_TO_RUN {
            totalNumOps := numOpsMultiplier * MEM_LEAK_SINGLE_THREAD_NUM_READS
            result := readOnly.perfReadOnlyJar(
                fqJar,totalNumOps,1,THREAD_POOL_SIZE_OP_TYPE,0,0,
                false,false,false)
            results = append(results,result)
        }
    }

    // write results to file
    resultsToFile(
        results,filepath.Join(outputFolder,MEM_LEAK_OUTPUT_NAME))
}
