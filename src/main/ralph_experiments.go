package main
import "read_only"
import "os"
import "log"

func main() {
    if len(os.Args) != 3 {
        log.Fatal("Require argument for read performance test jar dir")
    }
    readOnlyJarDir := os.Args[1]
    outputDir := os.Args[2]
    // create output directory if doesn't already exist
    os.MkdirAll(outputDir,0777)
    readOnlyModule := read_only.ReadOnly{}
    readOnlyModule.RunAll(readOnlyJarDir,outputDir)
}