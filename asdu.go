package vv104

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

type Asdu struct {
	TypeId   TypeId // Type Identification
	Num      Num
	Sequence bool
	CauseTx  CauseTx // Cause of Transmission
	Negative bool
	Test     bool
	OrigAddr OrigAddr // Originator Address
	Casdu    Casdu    // Common ASDU
	InfoObj  []InfoObj
}

func NewAsdu() *Asdu {
	// infoObj := NewInfoObj()

	asdu := Asdu{
		TypeId:   0,
		Num:      0,
		Sequence: false,
		CauseTx:  0,
		Negative: false,
		Test:     false,
		OrigAddr: 0,
		Casdu:    0,
		InfoObj:  []InfoObj{},
	}
	return &asdu
}

func (asdu *Asdu) AddInfoObject(infoObj InfoObj) error {
	asdu.InfoObj = append(asdu.InfoObj, infoObj)

	return nil
}

func NewInfoObj() InfoObj {
	ttag := time.Time{}
	val := IntVal(0)

	infoObj := InfoObj{
		Ioa:   0,
		Value: val,
		Quality: Quality{
			Bl: false,
			Sb: false,
			Nt: false,
			Iv: false,
			Ov: false,
		},
		CommandInfo: CommandInfo{
			Quoc: Quoc{
				Select: false,
				Qu:     0,
			},
			Qoi: 0,
		},
		TimeTag: ttag,
	}
	return infoObj
}

type InfoValue interface {
	Value() float32
	String() string
}

func (f FloatVal) Value() float32 {
	return float32(f)
}

func (f FloatVal) String() string {
	return fmt.Sprintf("%f", f)
}

func (i IntVal) Value() float32 {
	return float32(i)
}

func (i IntVal) String() string {
	return fmt.Sprintf("%d", i)
}

const (
	Bit1 byte = 1 << iota
	Bit2
	Bit3
	Bit4
	Bit5
	Bit6
	Bit7
	Bit8
)

func SetBit(b, flag byte) byte    { return b | flag }
func ClearBit(b, flag byte) byte  { return b &^ flag }
func ToggleBit(b, flag byte) byte { return b ^ flag }
func HasBit(b, flag byte) bool    { return b&flag != 0 }

type InfoObj struct {
	Ioa         Ioa // Info Object Address
	Value       InfoValue
	Quality     Quality
	CommandInfo CommandInfo
	TimeTag     time.Time // todo: add su, iv bits
}

// Quality flags
type Quality struct {
	Bl bool
	Sb bool
	Nt bool
	Iv bool
	Ov bool
}

type CommandInfo struct {
	Quoc Quoc
	Qoi  Qoi
}

type Quoc struct {
	Select bool
	Qu     Qu
}

type Qu uint8
type Qoi uint8

func (ioa Ioa) String() string {
	return fmt.Sprintf("IOA: %d", ioa)
}

func (q Quality) String() string {
	var s string = ""
	if q.Bl {
		s += "BL "
	}
	if q.Sb {
		s += "SB "
	}
	if q.Nt {
		s += "NT "
	}
	if q.Iv {
		s += "IV "
	}
	if q.Ov {
		s += "OV "
	}

	// remove last space
	if len(s) > 1 {
		s = s[:len(s)-1]
	}
	return s
}

func (casdu Casdu) String() string {

	return fmt.Sprintf("CASDU: %d", casdu)
}

func (asdu Asdu) String() string {
	var s string = ""
	s += asdu.TypeId.String()
	s += " "
	s += asdu.CauseTx.String()
	s += " "
	if asdu.Negative {
		s += "NEGATIVE "
	}

	if asdu.Num > 1 {
		// todo
		s += fmt.Sprintf(" (%d)", asdu.Num)
	}

	if asdu.OrigAddr != 0 {
		s += "OA: %s" + string(asdu.OrigAddr)
	}

	if asdu.Sequence {
		// todo
		s += "sequence todo"
	}

	if asdu.Test {
		s += "TEST "
	}

	for _, infoObj := range asdu.InfoObj {
		s += infoObj.Ioa.String()
		s += " "
		s += infoObj.Value.String()
		s += " "
		s += infoObj.Quality.String()
		s += " "

	}

	return s
}

func (infoObj InfoObj) WriteInfo(typeId TypeId, buf *bytes.Buffer) error {
	var b byte = 0
	val := infoObj.Value.Value()

	switch typeId {
	case M_SP_NA_1:
		b |= 0x01 & byte(val)
		infoObj.WriteSiqDiq(b, buf)

	case M_DP_NA_1:
		b |= 0x03 & byte(val)
		infoObj.WriteSiqDiq(b, buf)

	case M_ME_TD_1: // todo add other mvs here
		// two bytes info
		// one byte quality
		// { IV | NT | SB | BL | 0 | 0 | 0 | OV }
		binary.Write(buf, binary.LittleEndian, int16(val))
		infoObj.WriteQualitySeparateOctet(buf)

	case C_SC_NA_1:
		b |= 0x01 & byte(val)
		infoObj.WriteCommandInfo(b, buf)

	case C_IC_NA_1:
		var b byte = 0
		b |= byte(infoObj.CommandInfo.Qoi)
		buf.WriteByte(b)
	}

	return nil
}

// Info and Quality for SP and DP
func (infoObj InfoObj) WriteSiqDiq(b byte, buf *bytes.Buffer) {

	if infoObj.Quality.Bl {
		b = SetBit(b, Bit5)
	} else {
		b = ClearBit(b, Bit5)
	}
	if infoObj.Quality.Sb {
		b = SetBit(b, Bit6)
	} else {
		b = ClearBit(b, Bit6)
	}
	if infoObj.Quality.Nt {
		b = SetBit(b, Bit7)
	} else {
		b = ClearBit(b, Bit7)
	}
	if infoObj.Quality.Iv {
		b = SetBit(b, Bit8)
	} else {
		b = ClearBit(b, Bit8)
	}

	buf.WriteByte(b)

}

// Value and CommandInfo for commands
func (infoObj InfoObj) WriteCommandInfo(b byte, buf *bytes.Buffer) {
	if infoObj.CommandInfo.Quoc.Select {
		b = SetBit(b, Bit8)
	} else {
		b = ClearBit(b, Bit8)
	}
	b |= (byte(infoObj.CommandInfo.Quoc.Qu) << 2)
	buf.WriteByte(b)

}
func (infoObj InfoObj) WriteQualitySeparateOctet(buf *bytes.Buffer) {
	var b byte

	if infoObj.Quality.Ov {
		b = SetBit(b, Bit1)
	} else {
		b = ClearBit(b, Bit1)
	}
	// bits 2..4 reserve
	if infoObj.Quality.Bl {
		b = SetBit(b, Bit5)
	} else {
		b = ClearBit(b, Bit5)
	}
	if infoObj.Quality.Sb {
		b = SetBit(b, Bit6)
	} else {
		b = ClearBit(b, Bit6)
	}
	if infoObj.Quality.Nt {
		b = SetBit(b, Bit7)
	} else {
		b = ClearBit(b, Bit7)
	}
	if infoObj.Quality.Iv {
		b = SetBit(b, Bit8)
	} else {
		b = ClearBit(b, Bit8)
	}

	buf.WriteByte(b)

}

func TypeIsTimeTagged(typeId TypeId) bool {
	switch typeId {
	case

		M_SP_TB_1, // single-point information with time tag CP56Time2a
		M_DP_TB_1, // double-point information with time tag CP56Time2a
		M_ST_TB_1, // step position information with time tag CP56Time2a
		M_BO_TB_1, // bitstring of 32 bit with time tag CP56Time2a
		M_ME_TD_1, // measured value, normalized value with time tag CP56Time2a
		M_ME_TE_1, // measured value, scaled value with time tag CP56Time2a
		M_ME_TF_1, // measured value, short floating point number with time tag CP56Time2a
		M_IT_TB_1, // integrated totals with time tag CP56Time2a
		C_SC_TA_1, // single command with time tag CP56Time2a
		C_DC_TA_1, // double command with time tag CP56Time2a
		C_RC_TA_1, // regulating step command with time tag CP56Time2a
		C_SE_TA_1, // set point command, normalized value with time tag CP56Time2a
		C_SE_TB_1, // set point command, scaled value with time tag CP56Time2a
		C_SE_TC_1, // set point command, short floating-point number with time tag CP56Time2a
		C_BO_TA_1, // bitstring of 32 bits with time tag CP56Time2a
		C_TS_TA_1: // test command with time tag CP56Time2a
		return true

	default:
		return false
	}
}

func (asdu Asdu) Serialize(state State, buf *bytes.Buffer) {
	// todo err
	var x byte = 0

	buf.WriteByte(byte(asdu.TypeId))

	if asdu.Sequence {
		x = SetBit(x, Bit8)
	} else {
		x = ClearBit(x, Bit8)
	}
	x |= byte(asdu.Num)
	buf.WriteByte(x)

	x = 0
	if asdu.Test {
		x = SetBit(x, Bit8)
	} else {
		x = ClearBit(x, Bit8)
	}

	if asdu.Negative {
		x = SetBit(x, Bit7)
	} else {
		x = ClearBit(x, Bit7)
	}

	x |= byte(asdu.CauseTx)
	buf.WriteByte(x)
	buf.WriteByte(byte(asdu.OrigAddr))

	binary.Write(buf, binary.LittleEndian, int16(asdu.Casdu))

	for _, infoObj := range asdu.InfoObj {

		// ioa is three bytes long. convert..
		var b [3]byte
		ioa := infoObj.Ioa
		b[0] = byte(ioa & 0xFF)
		b[1] = byte((ioa & 0xFF00) >> 8)
		b[2] = byte((ioa & 0xFF0000) >> 16)
		binary.Write(buf, binary.LittleEndian, b)

		infoObj.WriteInfo(asdu.TypeId, buf)

		if TypeIsTimeTagged(asdu.TypeId) {
			infoObj.SerializeTime(state, buf)
		}

	}

}

func (infoObj InfoObj) SerializeTime(state State, buf *bytes.Buffer) {

	// todo su, iv
	var iv, su bool

	var timetag time.Time
	var ivMask byte
	var suMask byte
	var milliWord uint16
	var weekDay int

	if iv {
		ivMask = SetBit(ivMask, Bit8)
	}

	if su {
		suMask = SetBit(suMask, Bit8)
	}

	if !infoObj.TimeTag.IsZero() {
		// has a time tag already, use it
		timetag = infoObj.TimeTag
	} else {
		if state.Config.UseLocalTime {
			timetag = time.Now()
		} else {
			timetag = time.Now().UTC()
		}
	}

	milliWord = uint16(timetag.Second()*1000 + timetag.Nanosecond()/1000000)
	weekDay = int(timetag.Weekday())
	if weekDay == 0 {
		// Sunday equals 7
		weekDay = 7
	}
	day := timetag.Day()
	buf.WriteByte(byte(milliWord & 0x00FF))
	buf.WriteByte(byte((milliWord & 0xFF00) >> 8))
	buf.WriteByte((byte(timetag.Minute()) & 0b111111) | ivMask)
	buf.WriteByte((byte(timetag.Hour()) & 0b11111) | suMask)
	buf.WriteByte(((byte(weekDay) & 0b111) << 5) | (byte(day) & 0b11111))
	buf.WriteByte(byte(timetag.Month()) & 0xF)
	buf.WriteByte(byte(timetag.Year()-2000) & 0x7F)

}

func ParseAsdu(buf *bytes.Buffer) (Asdu, error) {
	var b byte
	var b1 byte
	asdu := NewAsdu()

	b, _ = buf.ReadByte()
	asdu.TypeId = TypeId(b)

	// variable structure verifier
	b, _ = buf.ReadByte()

	asdu.Sequence = HasBit(b, Bit8)
	asdu.Num = Num(b & 0b0111_1111)

	// cot
	b, _ = buf.ReadByte()
	asdu.CauseTx = CauseTx(b)

	// oa
	b, _ = buf.ReadByte()
	asdu.OrigAddr = OrigAddr(b)

	// casdu
	b, _ = buf.ReadByte()
	b1, _ = buf.ReadByte()
	asdu.Casdu = Casdu(uint16(b) + uint16(b1)*256)

	for i := 0; i < int(asdu.Num); i++ {

		var newInfoObj InfoObj
		newInfoObj.ParseInfoObj(asdu.TypeId, buf)

		asdu.InfoObj = append(asdu.InfoObj, newInfoObj)
	}

	return *asdu, nil
}

func (infoObj *InfoObj) ParseInfoObj(typeId TypeId, buf *bytes.Buffer) error {

	ioa1, _ := buf.ReadByte()
	ioa2, _ := buf.ReadByte()
	ioa3, _ := buf.ReadByte()

	infoObj.Ioa = Ioa(uint32(ioa1) + uint32(ioa2)*256 + uint32(ioa3)*65536)

	switch typeId {
	case M_SP_NA_1, M_SP_TB_1, M_DP_NA_1, M_DP_TB_1:
		// SP, DP
		infoObj.ParseSiqDiq(typeId, buf)

	case M_ME_NA_1, M_ME_TD_1, M_ME_NB_1, M_ME_TE_1, M_ME_NC_1, M_ME_TF_1:
		// MV
		infoObj.ParseMvValue(typeId, buf)

	case C_SC_NA_1, C_SC_TA_1, C_DC_NA_1, C_DC_TA_1:
		// SC, DC
		infoObj.ParseScoDco(typeId, buf)

	case C_IC_NA_1:
		// GI
		b, _ := buf.ReadByte()

		infoObj.CommandInfo.Qoi = Qoi(uint8(b))
		infoObj.Value = IntVal(0)

	}

	return nil
}

func (infoObj *InfoObj) ParseSiqDiq(typeId TypeId, buf *bytes.Buffer) {

	b := infoObj.ParseQds(typeId, buf)

	switch typeId {
	case M_SP_NA_1, M_SP_TB_1:
		infoObj.Value = IntVal(b & 0x01)

	case M_DP_NA_1, M_DP_TB_1:
		infoObj.Value = IntVal(b & 0x03)
	}
	fmt.Println("parsesiqdiq:", infoObj.Value)

}

func (infoObj *InfoObj) ParseMvValue(typeId TypeId, buf *bytes.Buffer) {

	switch typeId {
	case M_ME_NA_1, M_ME_TD_1, M_ME_NB_1, M_ME_TE_1:
		// normalized value and scaled value
		b1, _ := buf.ReadByte()
		b2, _ := buf.ReadByte()

		value := uint16(b1) + 256*uint16(b2)
		infoObj.Value = IntVal(value)

		infoObj.ParseQds(typeId, buf)

		if TypeIsTimeTagged(typeId) {
			infoObj.ParseTimeTag(buf)
		}

	case M_ME_NC_1, M_ME_TF_1:
		// float value
		bb := make([]byte, 4)

		bb[0], _ = buf.ReadByte()
		bb[1], _ = buf.ReadByte()
		bb[2], _ = buf.ReadByte()
		bb[3], _ = buf.ReadByte()

		bits := binary.LittleEndian.Uint32(bb)
		infoObj.Value = FloatVal(math.Float32frombits(bits))

		infoObj.ParseQds(typeId, buf)

	}

}

func (infoObj *InfoObj) ParseQds(typeid TypeId, buf *bytes.Buffer) byte {
	b, _ := buf.ReadByte()

	quality := &infoObj.Quality

	if HasBit(b, Bit8) {
		quality.Iv = true
	} else {
		quality.Iv = false
	}
	if HasBit(b, Bit7) {
		quality.Nt = true
	} else {
		quality.Nt = false
	}
	if HasBit(b, Bit6) {
		quality.Sb = true
	} else {
		quality.Sb = false
	}
	if HasBit(b, Bit5) {
		quality.Bl = true
	} else {
		quality.Bl = false
	}

	switch typeid {
	case M_ME_NA_1, M_ME_TD_1, M_ME_NB_1, M_ME_TE_1, M_ME_NC_1, M_ME_TF_1:
		// all MV values have also OV flag

		if HasBit(b, Bit1) {
			quality.Ov = true
		} else {
			quality.Ov = false
		}
	}

	return b
}

func (infoObj *InfoObj) ParseScoDco(typeId TypeId, buf *bytes.Buffer) {
	b, _ := buf.ReadByte()

	infoObj.CommandInfo.Quoc.Select = HasBit(b, Bit8)

	switch typeId {
	case C_SC_NA_1, C_SC_TA_1:
		infoObj.Value = IntVal(b & 0x01)
	case C_DC_NA_1, C_DC_TA_1:
		infoObj.Value = IntVal(b & 0x03)
	}

	infoObj.CommandInfo.Quoc.Qu = Qu(uint8(b >> 3))

}

func (infoObj *InfoObj) ParseTimeTag(buf *bytes.Buffer) {

	b1, _ := buf.ReadByte()
	b2, _ := buf.ReadByte()
	b3, _ := buf.ReadByte()
	b4, _ := buf.ReadByte()
	b5, _ := buf.ReadByte()
	b6, _ := buf.ReadByte()
	b7, _ := buf.ReadByte()

	millisec := int(b1) + (int(b2) << 8)
	min := int(b3 & 0b111111)
	hour := int(b4 & 0b11111)
	day := int(b5 & 0b11111)
	month := int(b6 & 0xF)
	year := int(b7&0x7F) + 2000

	sec := millisec / 1000
	ns := (millisec - (sec * 1000)) * 1000000
	infoObj.TimeTag = time.Date(year, time.Month(month), day, hour, min, sec, ns, time.Local)
}

// type Asdu[T NumberValue] struct {
// 	TypeId   TypeId // Type Identification
// 	Num      Num
// 	Sequence bool
// 	CauseTx  CauseTx // Cause of Transmission
// 	Negative bool
// 	Test     bool
// 	OrigAddr OrigAddr // Originator Address
// 	Casdu    Casdu    // Common ASDU
// 	Ioa      Ioa      // Info Object Address
// 	Info     Info[T]
// }
// type Info[T NumberValue] struct {
// 	Value   T
// 	Quality Quality
// }

// type Floatvalue float32
// type Intvalue int32

// type InfoValue interface {
// 	value()
// }

// func (f Floatvalue) value() float32 {
// 	return float32(f)
// }

// func (i Intvalue) value() float32 {
// 	return float32(i)
// }

// type Asdu_ struct {
// 	TypeId   TypeId // Type Identification
// 	Num      Num
// 	Sequence bool
// 	CauseTx  CauseTx // Cause of Transmission
// 	Negative bool
// 	Test     bool
// 	OrigAddr OrigAddr // Originator Address
// 	Casdu    Casdu    // Common ASDU
// 	Ioa      Ioa      // Info Object Address
// 	Info     Info_
// }
// type Info_ struct {
// 	Value   InfoValue
// 	Quality Quality
// }
