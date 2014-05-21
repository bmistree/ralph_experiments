package read_only

import "path/filepath"

const STACKED_OUTPUT_NAME = "stacked_output.csv"

func stackedReadOnlyTests(jarDir,outputFolder string) {
    fqJar := filepath.Join(jarDir,STACKED_READ_ONLY_JAR_NAME)
    params := createDefaultParameter()

    var results [] * LoggedReadOnlyResult
    for i := 0; i < NUM_TIMES_TO_RUN_EACH_EXPERIMENT; i++ {
        result := loggedReadOnlyJar(fqJar,params)
        results = append(results,result)
    }
    loggedResultsToFile(
        results, filepath.Join(outputFolder,STACKED_OUTPUT_NAME))
}
