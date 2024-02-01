package vv104

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
)

// holds three strings: obj name, typeId, IOA
// all are strings, so the toml marshaller will produce toml arrays (otherwise tables)
type NameTypeIdIoa [3]string

type Objects struct {
	sync.RWMutex
	// maps keep obj Names and refer to Asdu objects
	MoniObjects ObjectsMap
	CtrlObjects ObjectsMap

	// lists are mainly for GUI use (could be merged with configuredObjects?)
	MoniList ObjectList
	CtrlList ObjectList

	// configured Objects are for import/export in toml file
	configuredObjects ConfiguredObjects
}

type ObjectList []string
type ObjectsMap map[string]Asdu

// for toml conf file
type ConfiguredObjects struct {
	MonitoringList []NameTypeIdIoa `toml:"monitoringList,multiline,omitempty"`
	ControlList    []NameTypeIdIoa `toml:"controlList,multiline,omitempty"`
}

func NewConfiguredObjects() *ConfiguredObjects {
	return &ConfiguredObjects{
		MonitoringList: []NameTypeIdIoa{},
		ControlList:    []NameTypeIdIoa{},
	}
}

func NewObjects() *Objects {
	return &Objects{
		RWMutex:  sync.RWMutex{},
		MoniList: []string{},
		CtrlList: []string{},

		MoniObjects: map[string]Asdu{},
		CtrlObjects: map[string]Asdu{},
	}
}

func (objects *Objects) AddObjectByName(objectName string, asdu Asdu) error {
	if objectName == "" {
		return errors.New("obj name cant be empty string")
	}

	objects.RWMutex.Lock()
	defer objects.RWMutex.Unlock()

	if isMonitoringObject(int(asdu.TypeId)) {
		// monitoring direction
		if objects.MoniObjects.ObjectExists(objectName, int(asdu.InfoObj[0].Ioa), int(asdu.TypeId)) {
			return errors.New("moni object already exists")
		}
		// put in moni map, moni list and configured objects
		objects.MoniObjects[objectName] = asdu
		objects.MoniList = append(objects.MoniList, describeObject(objectName, asdu))

		typeId := strconv.Itoa(int(asdu.TypeId))
		ioa := strconv.Itoa(int(asdu.InfoObj[0].Ioa))
		objects.configuredObjects.MonitoringList = append(objects.configuredObjects.MonitoringList, [3]string{objectName, typeId, ioa})
	} else if isControlObject(int(asdu.TypeId)) {
		// control direction
		if objects.CtrlObjects.ObjectExists(objectName, int(asdu.InfoObj[0].Ioa), int(asdu.TypeId)) {
			return errors.New("ctrl object already exists")
		}
		// put in ctrl map, ctrl list and configured objects
		objects.CtrlObjects[objectName] = asdu
		objects.CtrlList = append(objects.CtrlList, describeObject(objectName, asdu))

		typeId := strconv.Itoa(int(asdu.TypeId))
		ioa := strconv.Itoa(int(asdu.InfoObj[0].Ioa))
		objects.configuredObjects.ControlList = append(objects.configuredObjects.ControlList, [3]string{objectName, typeId, ioa})
	}

	return nil
}

func (objects *Objects) addObjectsFromList(list ConfiguredObjects) error {
	objects.rangeOverListAndAdd(list.MonitoringList)
	objects.rangeOverListAndAdd(list.ControlList)
	return nil
}

func (objects *Objects) rangeOverListAndAdd(list []NameTypeIdIoa) {
	for _, name := range list {
		objName := name[0]
		typeIdStr := name[1]
		ioaStr := name[2]

		typeIdInt, err := strconv.Atoi(typeIdStr)

		if err != nil {
			fmt.Println(err)
		}
		ioaInt, err := strconv.Atoi(ioaStr)
		if err != nil {
			fmt.Println(err)
		}

		asdu := Asdu{}
		asdu.TypeId = TypeId(typeIdInt)
		infoObj := NewInfoObj()
		infoObj.Ioa = Ioa(ioaInt)
		asdu.AddInfoObject(infoObj)

		err = objects.AddObjectByName(objName, asdu)
		if err != nil {
			fmt.Println(err)
		}
	}
	// return nil

}

func isMonitoringObject(typeId int) bool {
	return typeId < 45
}

func isControlObject(typeId int) bool {
	return typeId >= 45 && typeId < 70
}

func describeObject(objectName string, asdu Asdu) string {
	return fmt.Sprintln(objectName + " | " + asdu.TypeId.String() + " | " + asdu.InfoObj[0].Ioa.String())
}

func (objects *Objects) objNameOrIoa(asdu Asdu) string {
	if len(asdu.InfoObj) == 0 {
		// not an i-format
		return ""
	}

	for name, asduFromMap := range objects.MoniObjects {
		if asdu.TypeId == asduFromMap.TypeId {
			if asdu.InfoObj[0].Ioa == asduFromMap.InfoObj[0].Ioa {
				return name
			}
		}
	}

	for name, asduFromMap := range objects.CtrlObjects {
		if asdu.TypeId == asduFromMap.TypeId {
			if asdu.InfoObj[0].Ioa == asduFromMap.InfoObj[0].Ioa {
				return name
			}
		}
	}

	return asdu.InfoObj[0].Ioa.String()
}

func (objects *Objects) RemoveObject(objectName string) error {

	objects.RWMutex.Lock()
	defer objects.RWMutex.Unlock()

	if objects.MoniObjects.ObjectExists(objectName, 0, 0) {
		delete(objects.MoniObjects, objectName)
		index, err := objects.MoniList.FindObjectInList(objectName)
		if err != nil {
			// fmt.Println(objects.MoniObjects)
			return errors.New("obj does not exist in object List, cant be removed there")
		}
		objects.MoniList = slices.Delete(objects.MoniList, index, index+1)
	} else if objects.CtrlObjects.ObjectExists(objectName, 0, 0) {
		delete(objects.CtrlObjects, objectName)
		index, err := objects.CtrlList.FindObjectInList(objectName)
		if err != nil {
			// fmt.Println(objects.MoniObjects)
			return errors.New("obj does not exist in object List, cant be removed there")
		}
		objects.CtrlList = slices.Delete(objects.CtrlList, index, index+1)

	} else {
		return errors.New("not found in map, can't remove")
	}

	return nil
}

func (objectList ObjectList) FindObjectInList(objName string) (int, error) {

	for i, name := range objectList {
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

// func (objects Objects) ObjectExists(objName string, ioa int, typeId int) bool {
// 	return objects.MoniObjects.ObjectExists(objName, ioa, typeId) || objects.CtrlObjects.ObjectExists(objName, ioa, typeId)

// }

func (objects Objects) PrintObjects() {
	fmt.Println("=============  Control Objects  =============")

	for objName, asdu := range objects.CtrlObjects {
		fmt.Println(objName, asdu.String())
	}

	fmt.Println("============= Monitoring Objects =============")

	for objName, asdu := range objects.MoniObjects {
		fmt.Println(objName, asdu.String())
	}
}
