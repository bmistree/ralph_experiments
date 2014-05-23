package read_only

import "path/filepath"

const STACKED_OUTPUT_NAME = "stacked_output.js"
const ACTIVE_EVENT_MAP_OUTPUT_NAME = "active_event_map_output.js"
var STACKED_TEST_NUM_THREADS [3]uint32 = [3]uint32{1,2,4}


func stackedReadOnlyTests(jarDir,outputFolder string) {
    fqJar := filepath.Join(jarDir,STACKED_READ_ONLY_JAR_NAME)
    params := createDefaultParameter()

    var results [] * LoggedReadOnlyResult
    toSplitOn := "begin_first_phase_commit bottom"
    for i := 0; i < NUM_TIMES_TO_RUN_EACH_EXPERIMENT; i++ {

        for _,numThreads := range STACKED_TEST_NUM_THREADS {
            params.numThreads = numThreads
            result := loggedReadOnlyJar(fqJar,params,toSplitOn)
            results = append(results,result)
        }
    }

    loggedResultsToFile(
        results, filepath.Join(outputFolder,STACKED_OUTPUT_NAME))
}


func loggedActiveEventMapTests(jarDir,outputFolder string) {
    fqJar := filepath.Join(jarDir,ACTIVE_EVENT_MAP_JAR_NAME)
    params := createDefaultParameter()

    var results [] * LoggedReadOnlyResult
    toSplitOn := "end_sentinel"
    for i := 0; i < NUM_TIMES_TO_RUN_EACH_EXPERIMENT; i++ {

        for _,numThreads := range STACKED_TEST_NUM_THREADS {
            params.numThreads = numThreads
            result := loggedReadOnlyJar(fqJar,params,toSplitOn)
            results = append(results,result)
        }
    }
    loggedResultsToFile(
        results, filepath.Join(outputFolder,ACTIVE_EVENT_MAP_OUTPUT_NAME))
}