package vv104

import "sync"

type Objects struct {
	sync.RWMutex
	MoniObjects ObjectsMap
	CtrlObjects ObjectsMap
}

type ObjectsMap map[string]Asdu

func NewObjects() *Objects {
	return &Objects{
		RWMutex:     sync.RWMutex{},
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

	return nil
}
