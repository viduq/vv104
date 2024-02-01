package main

import (
	"fmt"

	"github.com/viduq/vv104"
)

func main() {
	fmt.Println("vv140 started.")

	state := vv104.NewState()
	state.ParseFlags()

	state.Objects.PrintObjects()

	state.Start()

	fmt.Println("vv140 finished.")
}
