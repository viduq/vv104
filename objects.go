package vv104

import (
	"errors"
	"fmt"
	"slices"
	"strings"
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
	if objectName == "" {
		return errors.New("obj name cant be empty string")
	}
	if objects.ObjectExists(objectName, int(asdu.InfoObj[0].Ioa), int(asdu.TypeId)) {
		return errors.New("obj already exists")
	}

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

	objects.ObjectsList = append(objects.ObjectsList, describeObject(objectName, asdu))

	return nil
}

func describeObject(objectName string, asdu Asdu) string {
	return fmt.Sprintln(objectName + " | " + asdu.TypeId.String() + " | " + asdu.InfoObj[0].Ioa.String())
}

func (objects *Objects) RemoveObject(objectName string) error {

	objects.RWMutex.Lock()
	defer objects.RWMutex.Unlock()

	if objects.MoniObjects.ObjectExists(objectName, 0, 0) {
		delete(objects.MoniObjects, objectName)
	} else if objects.CtrlObjects.ObjectExists(objectName, 0, 0) {
		delete(objects.CtrlObjects, objectName)
	} else {
		return errors.New("not found in map, can't remove")
	}

	index, err := objects.FindObjectInList(objectName)
	if err != nil {
		// fmt.Println(objects.MoniObjects)
		return errors.New("obj does not exist in object List, cant be removed there")
	}
	objects.ObjectsList = slices.Delete(objects.ObjectsList, index, index+1)

	return nil
}

func (objects Objects) FindObjectInList(objName string) (int, error) {

	for i, name := range objects.ObjectsList {
		// ObjectsList contains full description with TypeID | IOA .. therefore we need to split
		if objName == strings.Split(name, " ")[0] {
			return i, nil
		}
	}
	return -1, errors.New("not found")

}

// type Exister interface {
// 	ObjectExists(string, int, int) bool
// }

// search for object by providing name, or ioa and typeid
func (objectsMap ObjectsMap) ObjectExists(objName string, ioa int, typeId int) bool {
	_, ok := objectsMap[objName]
	if ok {
		return true
	}

	for _, asdu := range objectsMap {
		// if asdu.TypeId == TypeId(typeId) {
		if asdu.InfoObj[0].Ioa == Ioa(ioa) {
			return true
			// }
		}
	}

	return false
}

func (objects Objects) ObjectExists(objName string, ioa int, typeId int) bool {
	return objects.MoniObjects.ObjectExists(objName, ioa, typeId) || objects.CtrlObjects.ObjectExists(objName, ioa, typeId)

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
