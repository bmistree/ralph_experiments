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
    cpusUtilized float64
    contextSwitches uint64
    cpuMigrations uint64
    pageFaults uint64
    cycles uint64
    stalledCyclesFrontend uint64
    frontendCyclesIdle float64
    stalledCyclesBackend uint64
    backendCyclesIdle float64
    instructions uint64
    branches uint64
    branchMisses uint64
}

func (perfOutput * PerfOutput) PrintAll() {
    fmt.Println(
        "%s: %f", CPUS_UTILIZED_STRING, perfOutput.cpusUtilized)
    fmt.Println(
        "%s: %d", CONTEXT_SWITCHES_STRING, perfOutput.contextSwitches)
    fmt.Println(
        "%s: %d", CPU_MIGRATIONS_STRING, perfOutput.cpuMigrations)
    fmt.Println(
        "%s: %d", PAGE_FAULTS_STRING, perfOutput.pageFaults)
    fmt.Println(
        "%s: %d", CYCLES_STRING, perfOutput.cycles)
    fmt.Println(
        "%s: %d", STALLED_CYCLES_FRONTEND_STRING, perfOutput.stalledCyclesFrontend)
    fmt.Println(
        "%s: %f", CYCLES_FRONTEND_IDLE_STRING, perfOutput.frontendCyclesIdle)
    fmt.Println(
        "%s: %d", STALLED_CYCLES_BACKEND_STRING, perfOutput.stalledCyclesBackend)
    fmt.Println(
        "%s: %f", CYCLES_BACKEND_IDLE_STRING, perfOutput.backendCyclesIdle)
    fmt.Println(
        "%s: %d", INSTRUCTIONS_STRING, perfOutput.instructions)
    fmt.Println(
        "%s: %d", BRANCHES_STRING, perfOutput.branches)
    fmt.Println(
        "%s: %d", BRANCH_MISSES_STRING, perfOutput.branchMisses)
}

func ParsePerfOutput(perfOutput string) PerfOutput{
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
        cpusUtilized: cpusUtilized,
        contextSwitches: contextSwitches,
        cpuMigrations: cpuMigrations,
        pageFaults: pageFaults,
        cycles: cycles,
        stalledCyclesFrontend: stalledCyclesFrontend,
        frontendCyclesIdle: frontendCyclesIdle,
        stalledCyclesBackend: stalledCyclesBackend,
        backendCyclesIdle: backendCyclesIdle,
        instructions: instructions,
        branches: branches,
        branchMisses: branchMisses,
    }
    
    return toReturn
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