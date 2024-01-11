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
	objects.AddObject("dp1", *asdu)

	infoObject.Ioa = 101

	objects.AddObject("dp2", *asdu)
	infoObject.Ioa = 102
	objects.AddObject("dp3", *asdu)
	infoObject.Ioa = 103
	objects.AddObject("dp4", *asdu)

	asdu = vv104.NewAsdu()

	infoObjects = nil
	infoObject.Ioa = 10
	infoObject.Value = vv104.IntVal(0)
	infoObjects = append(infoObjects, infoObject)
	asdu.Casdu = 1
	asdu.TypeId = vv104.C_SC_NA_1
	asdu.InfoObj = infoObjects
	objects.AddObject("sc1", *asdu)

	objects.PrintObjects()

	// go sendAFrame(&state)
	state.Start()

	fmt.Println("vv140 finished.")
}
