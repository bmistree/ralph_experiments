package read_only
import "strconv"
import "io/ioutil"

type ReadOnlyResult struct {
    numReads uint32
    numThreads uint32
    opType operationType
    opsPerSecond float64
}

func resultsToFile(results []ReadOnlyResult, filename string) {
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

// output will be num reads
func (readOnlyResult ReadOnlyResult) toCSVString () string {
    toReturn := ""
    toReturn += strconv.FormatUint(uint64(readOnlyResult.numReads),10) + "," +
        strconv.FormatUint(uint64(readOnlyResult.numThreads),10) + "," +
        strconv.FormatUint(uint64(readOnlyResult.opType),10) + "," +
        strconv.FormatFloat(readOnlyResult.opsPerSecond,'f',2,64)
    return toReturn
}