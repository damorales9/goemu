package goemu

import (
	"os"
)

func main() {
	os.Exit(emu_run(len(os.Args), os.Args))
}
