package read_only
import "strconv"
import "io/ioutil"
import "common"

type ReadOnlyResult struct {
    numReads uint32
    numThreads uint32
    opType operationType
    opsPerSecond float64
    perfOutput * common.PerfOutput
}

func resultsToFile(results []*ReadOnlyResult, filename string) {
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

func (readOnlyResult* ReadOnlyResult) toJSONString() string {
    toReturn := "{"
    toReturn += "\"num_reads\": " + strconv.FormatUint(uint64(readOnlyResult.numReads),10) + "," +
        "\"num_threads\": " + strconv.FormatUint(uint64(readOnlyResult.numThreads),10) + "," +
        "\"op_type\": " + strconv.FormatUint(uint64(readOnlyResult.opType),10) + "," +
        "\"ops_per_second\": " + strconv.FormatFloat(readOnlyResult.opsPerSecond,'f',2,64) + ","

    perfOutput := readOnlyResult.perfOutput
    toReturn += "\"perf_output\": "
    if perfOutput != nil {
        toReturn += perfOutput.ToJSONString()
    } else {
        toReturn += "null"
    }
    return toReturn
}

// output will be num reads
func (readOnlyResult* ReadOnlyResult) toCSVString () string {
    toReturn := ""
    toReturn += strconv.FormatUint(uint64(readOnlyResult.numReads),10) + "," +
        strconv.FormatUint(uint64(readOnlyResult.numThreads),10) + "," +
        strconv.FormatUint(uint64(readOnlyResult.opType),10) + "," +
        strconv.FormatFloat(readOnlyResult.opsPerSecond,'f',2,64)

    perfOutput := readOnlyResult.perfOutput
    if perfOutput != nil {
        toReturn += "," +
            strconv.FormatFloat(perfOutput.CpusUtilized,'f',2,64) + "," +
            strconv.FormatUint(uint64(perfOutput.ContextSwitches),10) + "," +
            strconv.FormatUint(uint64(perfOutput.CpuMigrations),10) + "," +
            strconv.FormatUint(uint64(perfOutput.PageFaults),10) + "," +
            strconv.FormatUint(uint64(perfOutput.Cycles),10) + "," +
            strconv.FormatUint(uint64(perfOutput.StalledCyclesFrontend),10) + "," +
            strconv.FormatFloat(perfOutput.FrontendCyclesIdle,'f',2,64) + "," +
            strconv.FormatUint(uint64(perfOutput.StalledCyclesBackend),10) + "," +
            strconv.FormatFloat(perfOutput.BackendCyclesIdle,'f',2,64) + "," +
            strconv.FormatUint(uint64(perfOutput.Instructions),10) + "," +
            strconv.FormatUint(uint64(perfOutput.Branches),10) + "," +
            strconv.FormatUint(uint64(perfOutput.BranchMisses),10)
    }
    return toReturn
}