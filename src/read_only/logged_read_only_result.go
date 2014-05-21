package read_only
import "common"
import "io/ioutil"
import "regexp"
import "strings"
import "strconv"

type LoggedReadOnlyResult struct {
    readOnlyResult * ReadOnlyResult
    allTraces [] * Trace
}

type Trace struct {
    timestampedEvents [] *TimestampedEvent
}
type TimestampedEvent struct {
    timestamp uint64
    eventString string
}

func (logged * LoggedReadOnlyResult) toJSONString() string {
    return "FIXME: must fill in toCSVString"
}


/**
 @param outputString --- Format is something like this:

529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743160027| Creation
529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743210253| get_val top
529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743254263| Add touched top
529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743287048| Add touched bottom
529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743334574| get_val bottom
529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743355640| begin_first_phase_commit top
529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743407517| obj first_phase_commit top
529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743433792| obj complete_commit top
529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743479146| obj complete_commit bottom
529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743535767| second_phase_commit top
529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743546713| obj complete_commit top
529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743564614| obj complete_commit bottom
529f429c-ff62-4e3a-8e76-f317b50d500b: 952598743593402| second_phase_commit bottom

...

476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744168775| Creation
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744192025| get_val top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744210586| Add touched top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744221727| Add touched bottom
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744234575| get_val bottom
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744245623| begin_first_phase_commit top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744265168| obj first_phase_commit top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744275005| obj complete_commit top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744291870| obj complete_commit bottom
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744317010| second_phase_commit top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744328428| obj complete_commit top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744344167| obj complete_commit bottom
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744366975| second_phase_commit bottom
ATOMIC_NUMBER_READ	 448.55989676932245

 */
func loggedTestRunOutputToResults(
    outputString string, numReads,numThreads uint32,opType operationType,
    perfOutput * common.PerfOutput) * LoggedReadOnlyResult {

    // each element of slice contains a text trace for a single
    // program, except the last, which just contains a string like
    // this: ATOMIC_NUMBER_READ 448.55989676932245
    stringTraces := strings.SplitAfter(outputString,"second_phase_commit bottom")

    // read overall throughput results
    overallThroughputString := stringTraces[len(stringTraces)-1]
    // remove last element
    stringTraces = stringTraces[0:len(stringTraces) -2]
    readOnlyResult := testRunOutputToResults(
        overallThroughputString,numReads,numThreads,opType,perfOutput)

    
    var allTraces [] * Trace
    for _, singleStringTrace := range stringTraces {
        trace := Trace {
            timestampedEvents: createTimestampsFromString(singleStringTrace),
        }
        allTraces = append(allTraces,&trace)
    }
    
    // FIXME: actually parse output to get runtime for each logged
    // entry.
    toReturn := LoggedReadOnlyResult {
        readOnlyResult: readOnlyResult,
        allTraces: allTraces,
    }

    return &toReturn
}

/**
 @param singleStringTrace --- Should have format:

476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744168775| Creation
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744192025| get_val top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744210586| Add touched top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744221727| Add touched bottom
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744234575| get_val bottom
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744245623| begin_first_phase_commit top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744265168| obj first_phase_commit top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744275005| obj complete_commit top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744291870| obj complete_commit bottom
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744317010| second_phase_commit top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744328428| obj complete_commit top
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744344167| obj complete_commit bottom
476c92b8-330a-4414-ba90-ac62e5f3d21c: 952598744366975| second_phase_commit bottom
*/
func createTimestampsFromString(singleStringTrace string) [] * TimestampedEvent {
    var toReturn []*TimestampedEvent
    
    singleStringTrace = strings.TrimSpace(singleStringTrace)
    individualStringEvents := strings.Split(singleStringTrace,"\n")

    timestampStringRegex := regexp.MustCompile(": ([0-9]+)\\|")
    eventDescStringRegex := regexp.MustCompile("\\| (.)*$")
    for _,individualStringEvent := range individualStringEvents {

        submatchArray :=
            timestampStringRegex.FindStringSubmatch(individualStringEvent)
        if len(submatchArray) != 2 {
            panic ("Incorrect number of elements in submatchArray")
        }
        timestamp, _err := strconv.ParseUint(submatchArray[1],10,64)
        if _err != nil {            
            panic("Could not convert timestamp to uint")
        }

        submatchArray =
            eventDescStringRegex.FindStringSubmatch(individualStringEvent)
        if len(submatchArray) != 2 {
            panic ("Incorrect number of elements in submatchArray2")
        }
        
        timestampedEvent := TimestampedEvent {
            timestamp: timestamp,
            eventString: submatchArray[1],
        }

        toReturn = append(toReturn, &timestampedEvent)
    }
    return toReturn
}


func loggedResultsToFile(results []*LoggedReadOnlyResult, filename string) {
    fileOutputString := "["
    for index, result := range results {
        fileOutputString += result.toJSONString()
        if index != (len(results) -1) {
            fileOutputString += ","
        }
    }
    fileOutputString += "]"
    
    // automatically creates file if doesn't exist
    err := ioutil.WriteFile(filename, []byte(fileOutputString), 0777)
    if err != nil {
        panic(err)
    }
}
