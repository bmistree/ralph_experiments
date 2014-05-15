package main
import "read_only"

func main() {
    readOnlyModule := read_only.ReadOnly{}
    readOnlyModule.RunAll("jar_name","output_folder")
}