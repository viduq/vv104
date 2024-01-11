package vv104

import (
	"fmt"
	"sync"
)

type Objects struct {
	sync.RWMutex
	ObjectsList []string
	MoniObjects ObjectsMap
	CtrlObjects ObjectsMap
}

type ObjectsMap map[string]Asdu

func NewObjects() *Objects {
	return &Objects{
		RWMutex:     sync.RWMutex{},
		ObjectsList: []string{},
		MoniObjects: map[string]Asdu{},
		CtrlObjects: map[string]Asdu{},
	}
}

func (objects *Objects) AddObject(objectName string, asdu Asdu) error {
	objects.RWMutex.Lock()
	defer objects.RWMutex.Unlock()

	switch typeId := asdu.TypeId; {
	case typeId < 45:
		// monitoring direction
		objects.MoniObjects[objectName] = asdu

	case typeId >= 45 && typeId < 70:
		// control direction
		objects.CtrlObjects[objectName] = asdu

	}

	objects.ObjectsList = append(objects.ObjectsList, objectName)

	return nil
}

func (objects Objects) PrintObjects() {
	fmt.Println("============= Control Objects =============")

	for objName, asdu := range objects.CtrlObjects {
		fmt.Println(objName, asdu.String())
	}

	fmt.Println("============= Monitoring Objects =============")

	for objName, asdu := range objects.MoniObjects {
		fmt.Println(objName, asdu.String())
	}
}
