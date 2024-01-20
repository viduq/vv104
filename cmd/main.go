package main

import (
	"fmt"

	"github.com/viduq/vv104"
)

func main() {
	fmt.Println("vv140 started.")

	objects := vv104.NewObjects()
	state := vv104.NewState()
	state.Config.ParseFlags(objects)

	objects.PrintObjects()

	state.Start()

	fmt.Println("vv140 finished.")
}
