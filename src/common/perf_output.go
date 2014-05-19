package common

import "fmt"
import "regexp"
import "strconv"
import "log"

const CPUS_UTILIZED_STRING = "CPUs utilized"
const CONTEXT_SWITCHES_STRING = "context-switches"
const CPU_MIGRATIONS_STRING = "CPU-migrations"
const PAGE_FAULTS_STRING = "page-faults"
const CYCLES_STRING = "cycles"
const STALLED_CYCLES_FRONTEND_STRING = "stalled-cycles-frontend"
const CYCLES_FRONTEND_IDLE_STRING = "% frontend cycles idle"
const STALLED_CYCLES_BACKEND_STRING = "stalled-cycles-backend"
// note: for some reason, perf displays backend cycles idle with
// additional space.
const CYCLES_BACKEND_IDLE_STRING = "% backend  cycles idle"
const INSTRUCTIONS_STRING = "instructions"
const BRANCHES_STRING = "branches"
const BRANCH_MISSES_STRING = "branch-misses"


type PerfOutput struct {
    CpusUtilized float64
    ContextSwitches uint64
    CpuMigrations uint64
    PageFaults uint64
    Cycles uint64
    StalledCyclesFrontend uint64
    FrontendCyclesIdle float64
    StalledCyclesBackend uint64
    BackendCyclesIdle float64
    Instructions uint64
    Branches uint64
    BranchMisses uint64
}

func (perfOutput* PerfOutput) String() string {
    return fmt.Sprintf(
        "%f,%d,%d,%d,%d,%d,%f,%d,%f,%d,%d,%d",
        perfOutput.CpusUtilized,perfOutput.ContextSwitches,
        perfOutput.CpuMigrations, perfOutput.PageFaults,
        perfOutput.Cycles, perfOutput.StalledCyclesFrontend,
        perfOutput.FrontendCyclesIdle, perfOutput.StalledCyclesBackend,
        perfOutput.BackendCyclesIdle, perfOutput.Instructions,
        perfOutput.Branches, perfOutput.BranchMisses)
}

func (perfOutput * PerfOutput) PrintAll() {
    fmt.Println(
        "%s: %f", CPUS_UTILIZED_STRING, perfOutput.CpusUtilized)
    fmt.Println(
        "%s: %d", CONTEXT_SWITCHES_STRING, perfOutput.ContextSwitches)
    fmt.Println(
        "%s: %d", CPU_MIGRATIONS_STRING, perfOutput.CpuMigrations)
    fmt.Println(
        "%s: %d", PAGE_FAULTS_STRING, perfOutput.PageFaults)
    fmt.Println(
        "%s: %d", CYCLES_STRING, perfOutput.Cycles)
    fmt.Println(
        "%s: %d", STALLED_CYCLES_FRONTEND_STRING, perfOutput.StalledCyclesFrontend)
    fmt.Println(
        "%s: %f", CYCLES_FRONTEND_IDLE_STRING, perfOutput.FrontendCyclesIdle)
    fmt.Println(
        "%s: %d", STALLED_CYCLES_BACKEND_STRING, perfOutput.StalledCyclesBackend)
    fmt.Println(
        "%s: %f", CYCLES_BACKEND_IDLE_STRING, perfOutput.BackendCyclesIdle)
    fmt.Println(
        "%s: %d", INSTRUCTIONS_STRING, perfOutput.Instructions)
    fmt.Println(
        "%s: %d", BRANCHES_STRING, perfOutput.Branches)
    fmt.Println(
        "%s: %d", BRANCH_MISSES_STRING, perfOutput.BranchMisses)
}

func ParsePerfOutput(perfOutput string) * PerfOutput{
    cpusUtilized, _err1 := strconv.ParseFloat(
        findNumberStringInPerfOutput(perfOutput,CPUS_UTILIZED_STRING),64)
    contextSwitches, _err2 := strconv.ParseUint(
        findNumberStringInPerfOutput(perfOutput,CONTEXT_SWITCHES_STRING),10,64)
    cpuMigrations, _err3 := strconv.ParseUint(
        findNumberStringInPerfOutput(perfOutput,CPU_MIGRATIONS_STRING),10,64)
    pageFaults, _err4 := strconv.ParseUint(
        findNumberStringInPerfOutput(perfOutput,PAGE_FAULTS_STRING),10,64)
    cycles, _err5 := strconv.ParseUint(
        findNumberStringInPerfOutput(perfOutput,CYCLES_STRING),10,64)
    stalledCyclesFrontend, _err6 := strconv.ParseUint(
        findNumberStringInPerfOutput(perfOutput,STALLED_CYCLES_FRONTEND_STRING),10,64)
    frontendCyclesIdle, _err7 := strconv.ParseFloat(
        findNumberStringInPerfOutput(perfOutput,CYCLES_FRONTEND_IDLE_STRING),64)
    stalledCyclesBackend, _err8 := strconv.ParseUint(
        findNumberStringInPerfOutput(perfOutput,STALLED_CYCLES_BACKEND_STRING),10,64)
    backendCyclesIdle, _err9 := strconv.ParseFloat(
        findNumberStringInPerfOutput(perfOutput,CYCLES_BACKEND_IDLE_STRING),64)
    instructions, _err10 := strconv.ParseUint(
        findNumberStringInPerfOutput(perfOutput,INSTRUCTIONS_STRING),10,64)
    branches, _err11 := strconv.ParseUint(
        findNumberStringInPerfOutput(perfOutput,BRANCHES_STRING),10,64)
    branchMisses, _err12 := strconv.ParseUint(
        findNumberStringInPerfOutput(perfOutput,BRANCH_MISSES_STRING),10,64)

    if ((_err1 != nil) || (_err2 != nil) || (_err3 != nil) ||
        (_err4 != nil) || (_err5 != nil) || (_err6 != nil) ||
        (_err7 != nil) || (_err8 != nil) || (_err9 != nil) ||
        (_err10 != nil) || (_err11 != nil) || (_err12 != nil)) {
        log.Fatal("Could not string convert perf stat output")
    }
     
    toReturn := PerfOutput {
        CpusUtilized: cpusUtilized,
        ContextSwitches: contextSwitches,
        CpuMigrations: cpuMigrations,
        PageFaults: pageFaults,
        Cycles: cycles,
        StalledCyclesFrontend: stalledCyclesFrontend,
        FrontendCyclesIdle: frontendCyclesIdle,
        StalledCyclesBackend: stalledCyclesBackend,
        BackendCyclesIdle: backendCyclesIdle,
        Instructions: instructions,
        Branches: branches,
        BranchMisses: branchMisses,
    }
    
    return &toReturn
}



/**
 @param perfOutput --- statistics reported from calling perf stat on a
 binary.  Formatted something like this:

 2.253278 task-clock         #    0.806 CPUs utilized          
 2 context-switches          #    0.001 M/sec                  
 0 CPU-migrations            #    0.000 M/sec                  
 268 page-faults             #    0.119 M/sec                  
 2,685,321 cycles            #    1.192 GHz          

 @param category --- Get any number to the left of this string.

 @returns --- The number associated with this category Eg., for
 perfOutput above and category "cycles," would return "2,685,321"
*/
func findNumberStringInPerfOutput(perfOutput, category string) string {
    r := regexp.MustCompile("([.0-9,]+)\\s*" + category)
    submatchArray :=  r.FindStringSubmatch(perfOutput)
    if len(submatchArray) != 2 {
        panic ("Incorrect number of elements in submatchArray")
    }

    return submatchArray[1]
}