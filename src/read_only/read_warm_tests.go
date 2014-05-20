package read_only

import "fmt"
import "path/filepath"
import "strconv"

const READ_WARM_TEST_OUTPUT_NAME = "read_warm.csv"
var WARM_TEST_NUM_OPS [6]uint32 =
    [6]uint32{1000,5000,10000,50000,100000,150000}

func singleThreadWarmTests(jarDir,outputFolder string) {
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    fmt.Println("Running warm experiment: ");

    params := createDefaultParameter()
    params.opType = READ_NUM
    
    var results []*ReadOnlyResult
    for i := 0; i < NUM_TIMES_TO_RUN_EACH_EXPERIMENT; i++ {
        fmt.Println(
            "\t" + strconv.Itoa(i+1) + " of " +
            strconv.Itoa(NUM_TIMES_TO_RUN_EACH_EXPERIMENT));
        for _,numReads := range WARM_TEST_NUM_OPS {
            params.numReads = numReads
            result := commonReadOnlyJar(fqJar,params)
            results = append(results,result)
        }
    }
    
    // write results to file
    resultsToFile(
        results,filepath.Join(outputFolder,READ_WARM_TEST_OUTPUT_NAME))
}

