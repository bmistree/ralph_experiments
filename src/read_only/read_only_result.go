package read_only
import "strconv"

type ReadOnlyResult struct {
    numReads uint32
    numThreads uint32
    opType operationType
    opsPerSecond float64
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