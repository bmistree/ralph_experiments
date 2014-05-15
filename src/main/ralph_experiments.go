package main
import "read_only"
import "os"
import "log"

func main() {
    if len(os.Args) != 2 {
        log.Fatal("Require argument for read performance test jar dir")
    }
    readOnlyJarDir := os.Args[1]
    readOnlyModule := read_only.ReadOnly{}
    readOnlyModule.RunAll(readOnlyJarDir,"output_folder")
}