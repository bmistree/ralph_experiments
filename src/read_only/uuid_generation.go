package read_only

import "path/filepath"
import "fmt"
import "strconv"

const UUID_GENERATION_OUTPUT_NAME = "uuid_generation.csv"
const UUID_GENERATION_TEST_NUM_READS = 100000
const UUID_GENERATION_OP_TYPE = READ_NUM 

func uuidGenerationTests(jarDir,outputFolder string) {
    
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    fmt.Println("Running uuid generation tests: ")

    params := createDefaultParameter()
    params.numReads = UUID_GENERATION_TEST_NUM_READS
    params.opType = UUID_GENERATION_OP_TYPE
    params.perfOn = true
    
    var results [] * ReadOnlyResult
    for i := 0; i < NUM_TIMES_TO_RUN_EACH_EXPERIMENT; i++ {
        fmt.Println(
            "\t" + strconv.Itoa(i+1) + " of " +
            strconv.Itoa(NUM_TIMES_TO_RUN_EACH_EXPERIMENT));

        // with standard uuid generation
        params.atomIntUUIDGeneration = false
        result := commonReadOnlyJar(fqJar,params)
        results = append(results,result)

        // with atomic number uuid generation
        params.atomIntUUIDGeneration = true
        result = commonReadOnlyJar(fqJar,params)
        results = append(results,result)
    }
    // write results to file
    resultsToFile(
        results,filepath.Join(outputFolder,UUID_GENERATION_OUTPUT_NAME))
}
