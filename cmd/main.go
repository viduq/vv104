package main

import (
	"fmt"

	"github.com/viduq/vv104"
)

func main() {
	fmt.Println("vv140 started.")

	state := vv104.NewState()
	config := vv104.NewState().Config

	config.ParseFlags()
	state.Config = config

	objects := vv104.NewObjects()
	asdu := vv104.NewAsdu()
	var infoObject vv104.InfoObj
	infoObject.Ioa = 100
	infoObject.Value = vv104.IntVal(2)
	var infoObjects []vv104.InfoObj
	infoObjects = append(infoObjects, infoObject)
	asdu.Casdu = 1
	asdu.TypeId = vv104.M_DP_NA_1
	asdu.InfoObj = infoObjects
	objects.AddObject("dp", *asdu)

	// go sendAFrame(&state)
	state.Start()

	fmt.Println("vv140 finished.")
}
