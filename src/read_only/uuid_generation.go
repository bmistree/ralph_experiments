package read_only

import "path/filepath"
import "fmt"
import "strconv"

const UUID_GENERATION_OUTPUT_NAME = "uuid_generation.csv"
const UUID_GENERATION_TEST_NUM_READS = 100000
const UUID_GENERATION_OP_TYPE = READ_NUM 

func uuidGenerationTests(
    readOnly* ReadOnly, jarDir,outputFolder string) {
    
    fqJar := filepath.Join(jarDir,READ_ONLY_JAR_NAME)
    fmt.Println("Running uuid generation tests: ")
    
    var results [] * ReadOnlyResult
    for i := 0; i < NUM_TIMES_TO_RUN_EACH_EXPERIMENT; i++ {
        fmt.Println(
            "\t" + strconv.Itoa(i+1) + " of " +
            strconv.Itoa(NUM_TIMES_TO_RUN_EACH_EXPERIMENT));

        // with standard uuid generation
        result := readOnly.perfReadOnlyJar(
            fqJar,UUID_GENERATION_TEST_NUM_READS,1,
            UUID_GENERATION_OP_TYPE,0,0,false,false,false)
        results = append(results,result)

        // with atomic number uuid generation
        result = readOnly.perfReadOnlyJar(
            fqJar,UUID_GENERATION_TEST_NUM_READS,1,
            UUID_GENERATION_OP_TYPE,0,0,true,false,false)
        results = append(results,result)
    }
    // write results to file
    resultsToFile(
        results,filepath.Join(outputFolder,UUID_GENERATION_OUTPUT_NAME))
}
