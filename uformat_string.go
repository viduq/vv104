// Code generated by "stringer -type=UFormat"; DO NOT EDIT.

package vv104

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[StartDTAct-7]
	_ = x[StartDTCon-11]
	_ = x[StopDTAct-19]
	_ = x[StopDTCon-35]
	_ = x[TestFRAct-67]
	_ = x[TestFRCon-131]
}

const (
	_UFormat_name_0 = "StartDTAct"
	_UFormat_name_1 = "StartDTCon"
	_UFormat_name_2 = "StopDTAct"
	_UFormat_name_3 = "StopDTCon"
	_UFormat_name_4 = "TestFRAct"
	_UFormat_name_5 = "TestFRCon"
)

func (i UFormat) String() string {
	switch {
	case i == 7:
		return _UFormat_name_0
	case i == 11:
		return _UFormat_name_1
	case i == 19:
		return _UFormat_name_2
	case i == 35:
		return _UFormat_name_3
	case i == 67:
		return _UFormat_name_4
	case i == 131:
		return _UFormat_name_5
	default:
		return "UFormat(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
