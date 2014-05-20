package read_only
import "common"
import "io/ioutil"

type LoggedReadOnlyResult struct {
    numReads uint32
    numThreads uint32
    opType operationType
    opsPerSecond float64
    perfOutput * common.PerfOutput
}

func (logged * LoggedReadOnlyResult) toCSVString() string {
    return "FIXME: must fill in toCSVString"
}

func loggedTestRunOutputToResults(
    outputString string, numReads,numThreads uint32,opType operationType,
    perfOutput * common.PerfOutput) * LoggedReadOnlyResult {

    // FIXME: actually parse output to get runtime for each logged
    // entry.
    toReturn := LoggedReadOnlyResult {
        numReads: numReads,
        numThreads: numThreads,
        opType: opType,
        opsPerSecond: -1,
        perfOutput: perfOutput,
    }

    return &toReturn
}

func loggedResultsToFile(results []*LoggedReadOnlyResult, filename string) {
    fileOutputString := ""
    for _, result := range results {
        fileOutputString += result.toCSVString() + "\n"
    }

    // automatically creates file if doesn't exist
    err := ioutil.WriteFile(filename, []byte(fileOutputString), 0777)
    if err != nil {
        panic(err)
    }
}
