package vv104

type TypeId uint8
type Num uint8
type CauseTx uint8
type OrigAddr uint8
type Casdu uint16
type Ioa uint32
type FrameFormat byte
type IntVal int32
type FloatVal float32
type UFormat byte

const (
	IFormatFrame FrameFormat = 0
	SFormatFrame FrameFormat = 1
	UFormatFrame FrameFormat = 3
)

// U-Format types
//
//go:generate stringer -type=UFormat
const (
	StartDTAct UFormat = 0x07
	StartDTCon UFormat = 0x0b
	StopDTAct  UFormat = 0x13
	StopDTCon  UFormat = 0x23
	TestFRAct  UFormat = 0x43
	TestFRCon  UFormat = 0x83
)

// for GUI
var (
	TypeIDs []string = []string{
		M_SP_NA_1.String(),
		M_DP_NA_1.String(),
		M_ST_NA_1.String(),
		M_BO_NA_1.String(),
		M_ME_NA_1.String(),
		M_ME_NB_1.String(),
		M_ME_NC_1.String(),
		M_IT_NA_1.String(),
		M_SP_TB_1.String(),
		M_DP_TB_1.String(),
		M_ST_TB_1.String(),
		M_BO_TB_1.String(),
		M_ME_TD_1.String(),
		M_ME_TE_1.String(),
		M_ME_TF_1.String(),
		M_IT_TB_1.String(),
		C_SC_NA_1.String(),
		C_DC_NA_1.String(),
		C_RC_NA_1.String(),
		C_SE_NA_1.String(),
		C_SE_NB_1.String(),
		C_SE_NC_1.String(),
		C_BO_NA_1.String(),
		C_SC_TA_1.String(),
		C_DC_TA_1.String(),
		C_RC_TA_1.String(),
		C_SE_TA_1.String(),
		C_SE_TB_1.String(),
		C_SE_TC_1.String(),
		C_BO_TA_1.String(),
		M_EI_NA_1.String(),
	}
)

// Type IDs
//
//go:generate stringer -type=TypeId
const (
	M_SP_NA_1 TypeId = 1   // single-point information
	M_DP_NA_1 TypeId = 3   // double-point information
	M_ST_NA_1 TypeId = 5   // step position information
	M_BO_NA_1 TypeId = 7   // bitstring of 32 bits
	M_ME_NA_1 TypeId = 9   // measured value, normalized value
	M_ME_NB_1 TypeId = 11  // measured value, scaled value
	M_ME_NC_1 TypeId = 13  // measured value, short floating point number
	M_IT_NA_1 TypeId = 15  // integrated totals
	M_SP_TB_1 TypeId = 30  // single-point information with time tag CP56Time2a
	M_DP_TB_1 TypeId = 31  // double-point information with time tag CP56Time2a
	M_ST_TB_1 TypeId = 32  // step position information with time tag CP56Time2a
	M_BO_TB_1 TypeId = 33  // bitstring of 32 bit with time tag CP56Time2a
	M_ME_TD_1 TypeId = 34  // measured value, normalized value with time tag CP56Time2a
	M_ME_TE_1 TypeId = 35  // measured value, scaled value with time tag CP56Time2a
	M_ME_TF_1 TypeId = 36  // measured value, short floating point number with time tag CP56Time2a
	M_IT_TB_1 TypeId = 37  // integrated totals with time tag CP56Time2a
	C_SC_NA_1 TypeId = 45  // single command
	C_DC_NA_1 TypeId = 46  // double command
	C_RC_NA_1 TypeId = 47  // regulating step command
	C_SE_NA_1 TypeId = 48  // set point command, normalized value
	C_SE_NB_1 TypeId = 49  // set point command, scaled value
	C_SE_NC_1 TypeId = 50  // set point command, short floating point number
	C_BO_NA_1 TypeId = 51  // bitstring of 32 bits
	C_SC_TA_1 TypeId = 58  // single command with time tag CP56Time2a
	C_DC_TA_1 TypeId = 59  // double command with time tag CP56Time2a
	C_RC_TA_1 TypeId = 60  // regulating step command with time tag CP56Time2a
	C_SE_TA_1 TypeId = 61  // set point command, normalized value with time tag CP56Time2a
	C_SE_TB_1 TypeId = 62  // set point command, scaled value with time tag CP56Time2a
	C_SE_TC_1 TypeId = 63  // set point command, short floating-point number with time tag CP56Time2a
	C_BO_TA_1 TypeId = 64  // bitstring of 32 bits with time tag CP56Time2a
	M_EI_NA_1 TypeId = 70  // end of initialization
	C_IC_NA_1 TypeId = 100 // interrogation command
	C_CI_NA_1 TypeId = 101 // counter interrogation command
	C_RD_NA_1 TypeId = 102 // read command
	C_CS_NA_1 TypeId = 103 // clock synchronization command
	C_RP_NA_1 TypeId = 105 // reset process command
	C_TS_TA_1 TypeId = 107 // test command with time tag CP56Time2a

	// "private" range for internal purposes, will not be sent
	INTERNAL_STATE_MACHINE_NOTIFIER TypeId = 200
)

func TypeIdFromName(id string) TypeId {
	switch id {
	case "M_SP_NA_1":
		return TypeId(M_SP_NA_1)
	case "M_DP_NA_1":
		return TypeId(M_DP_NA_1)
	case "M_ST_NA_1":
		return TypeId(M_ST_NA_1)
	case "M_BO_NA_1":
		return TypeId(M_BO_NA_1)
	case "M_ME_NA_1":
		return TypeId(M_ME_NA_1)
	case "M_ME_NB_1":
		return TypeId(M_ME_NB_1)
	case "M_ME_NC_1":
		return TypeId(M_ME_NC_1)
	case "M_IT_NA_1":
		return TypeId(M_IT_NA_1)
	case "M_SP_TB_1":
		return TypeId(M_SP_TB_1)
	case "M_DP_TB_1":
		return TypeId(M_DP_TB_1)
	case "M_ST_TB_1":
		return TypeId(M_ST_TB_1)
	case "M_BO_TB_1":
		return TypeId(M_BO_TB_1)
	case "M_ME_TD_1":
		return TypeId(M_ME_TD_1)
	case "M_ME_TE_1":
		return TypeId(M_ME_TE_1)
	case "M_ME_TF_1":
		return TypeId(M_ME_TF_1)
	case "M_IT_TB_1":
		return TypeId(M_IT_TB_1)

	case "C_SC_NA_1":
		return TypeId(C_SC_NA_1)
	case "C_DC_NA_1":
		return TypeId(C_DC_NA_1)
	case "C_RC_NA_1":
		return TypeId(C_RC_NA_1)
	case "C_SE_NA_1":
		return TypeId(C_SE_NA_1)
	case "C_SE_NB_1":
		return TypeId(C_SE_NB_1)
	case "C_SE_NC_1":
		return TypeId(C_SE_NC_1)
	case "C_BO_NA_1":
		return TypeId(C_BO_NA_1)
	case "C_SC_TA_1":
		return TypeId(C_SC_TA_1)
	case "C_DC_TA_1":
		return TypeId(C_DC_TA_1)
	case "C_RC_TA_1":
		return TypeId(C_RC_TA_1)
	case "C_SE_TA_1":
		return TypeId(C_SE_TA_1)
	case "C_SE_TB_1":
		return TypeId(C_SE_TB_1)
	case "C_SE_TC_1":
		return TypeId(C_SE_TC_1)
	case "C_BO_TA_1":
		return TypeId(C_BO_TA_1)
	case "M_EI_NA_1":
		return TypeId(M_EI_NA_1)

	case "C_IC_NA_1":
		return TypeId(C_IC_NA_1)
	case "C_CI_NA_1":
		return TypeId(C_CI_NA_1)
	case "C_RD_NA_1":
		return TypeId(C_RD_NA_1)
	case "C_CS_NA_1":
		return TypeId(C_CS_NA_1)
	case "C_RP_NA_1":
		return TypeId(C_RP_NA_1)
	case "C_TS_TA_1":
		return TypeId(C_TS_TA_1)

	}
	return TypeId(0)
}

// Cause of Transmission
//
//go:generate stringer -type=CauseTx
const (
	Per_Cyc      CauseTx = 1
	Back         CauseTx = 2
	Spont        CauseTx = 3
	Init         CauseTx = 4
	Req          CauseTx = 5
	Act          CauseTx = 6
	ActCon       CauseTx = 7
	Deact        CauseTx = 8
	DeactCon     CauseTx = 9
	ActTerm      CauseTx = 10
	Retrem       CauseTx = 11
	Retloc       CauseTx = 12
	File         CauseTx = 13
	Inrogen      CauseTx = 20
	Inro1        CauseTx = 21
	Inro2        CauseTx = 22
	Inro3        CauseTx = 23
	Inro4        CauseTx = 24
	Inro5        CauseTx = 25
	Inro6        CauseTx = 26
	Inro7        CauseTx = 27
	Inro8        CauseTx = 28
	Inro9        CauseTx = 29
	Inro10       CauseTx = 30
	Inro11       CauseTx = 31
	Inro12       CauseTx = 32
	Inro13       CauseTx = 33
	Inro14       CauseTx = 34
	Inro15       CauseTx = 35
	Inro16       CauseTx = 36
	Reqcogen     CauseTx = 37
	Reqco1       CauseTx = 38
	Reqco2       CauseTx = 39
	Reqco3       CauseTx = 40
	Reqco4       CauseTx = 41
	UkTypeId     CauseTx = 44
	UkCauseTx    CauseTx = 45
	UkComAdrASDU CauseTx = 46
	UkIOA        CauseTx = 47
)

// Qualifier of interrogation
//
//go:generate stringer -type=Qoi
const (
	notUsed             Qoi = iota
	statioInterrogation     = 20
	interrogationGroup1
	interrogationGroup2
	interrogationGroup3
	interrogationGroup4
	interrogationGroup5
	interrogationGroup6
	interrogationGroup7
	interrogationGroup8
	interrogationGroup9
	interrogationGroup10
	interrogationGroup11
	interrogationGroup12
	interrogationGroup13
	interrogationGroup14
	interrogationGroup15
)

// Qualifier of command
//
//go:generate stringer -type=Qu
const (
	noAddDef Qu = iota
	shortPulse
	longPulse
	persistent
)

// type TypeIDNames struct {
// 	Name     string
// 	abbrName string
// }

// var TypeIDMap = map[int]TypeIDNames{
// 	1:   {"M_SP_NA_1", "SP"},
// 	3:   {"M_DP_NA_1", "DP"},
// 	9:   {"M_ME_NA_1", "MV norm"},
// 	11:  {"M_ME_NB_1", "MV scal"},
// 	13:  {"M_ME_NC_1", "MV float"},
// 	15:  {"M_IT_NA_1", "IT"},
// 	30:  {"M_SP_TB_1", "SP (t)"},
// 	31:  {"M_DP_TB_1", "DP (t)"},
// 	34:  {"M_ME_TD_1", "MV norm (t)"},
// 	35:  {"M_ME_TE_1", "MV scal (t)"},
// 	36:  {"M_ME_TF_1", "MV float (t)"},
// 	37:  {"M_IT_TB_1", "IT (t)"},
// 	45:  {"C_SC_NA_1", "SC"},
// 	46:  {"C_DC_NA_1", "DC"},
// 	100: {"C_IC_NA_1", "IC"},
// 	105: {"C_RP_NA_1", "RP"},
// 	// usw
// }

// var TypeIDWithoutTimeMap = map[int]int{
// 	30: 1,
// 	31: 3,
// 	34: 9,
// 	35: 11,
// 	36: 13,
// 	37: 15,
// 	// usw
// }
