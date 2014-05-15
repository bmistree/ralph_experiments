package main
import "github.com/bmistree/ralph_experiments/read_only"

func main() {
    readOnlyModule := read_only.ReadOnly{}
    readOnlyModule.RunAll("jar_name","output_folder")
}