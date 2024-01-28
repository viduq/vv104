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

// an ASDU can contain multiple InfoObjects, append them using this function
func (asdu *Asdu) AddInfoObject(infoObj InfoObj) error {
	asdu.InfoObj = append(asdu.InfoObj, infoObj)

	asdu.Num = Num(len(asdu.InfoObj))

	return nil
}

func newInfoObj() InfoObj {
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
	bit1 byte = 1 << iota
	bit2
	bit3
	bit4
	bit5
	bit6
	bit7
	bit8
)

func setBit(b, flag byte) byte    { return b | flag }
func clearBit(b, flag byte) byte  { return b &^ flag }
func toggleBit(b, flag byte) byte { return b ^ flag }
func hasBit(b, flag byte) bool    { return b&flag != 0 }

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

func (infoObj InfoObj) writeInfo(typeId TypeId, buf *bytes.Buffer) error {
	var b byte = 0
	val := infoObj.Value.Value()

	switch typeId {
	case M_SP_NA_1:
		b |= 0x01 & byte(val)
		infoObj.writeSiqDiq(b, buf)

	case M_DP_NA_1:
		b |= 0x03 & byte(val)
		infoObj.writeSiqDiq(b, buf)

	case M_ME_TD_1: // todo add other mvs here
		// two bytes info
		// one byte quality
		// { IV | NT | SB | BL | 0 | 0 | 0 | OV }
		binary.Write(buf, binary.LittleEndian, int16(val))
		infoObj.writeQualitySeparateOctet(buf)

	case C_SC_NA_1, C_DC_NA_1:
		// command info, todo make configurable
		infoObj.CommandInfo.Quoc = Quoc{
			Select: false,
			Qu:     shortPulse,
		}
		b |= 0x01 & byte(val)
		infoObj.writeCommandInfo(b, buf)

	case C_IC_NA_1:
		var b byte = 0
		b |= byte(infoObj.CommandInfo.Qoi)
		buf.WriteByte(b)
	}

	return nil
}

// Info and Quality for SP and DP
func (infoObj InfoObj) writeSiqDiq(b byte, buf *bytes.Buffer) {

	if infoObj.Quality.Bl {
		b = setBit(b, bit5)
	} else {
		b = clearBit(b, bit5)
	}
	if infoObj.Quality.Sb {
		b = setBit(b, bit6)
	} else {
		b = clearBit(b, bit6)
	}
	if infoObj.Quality.Nt {
		b = setBit(b, bit7)
	} else {
		b = clearBit(b, bit7)
	}
	if infoObj.Quality.Iv {
		b = setBit(b, bit8)
	} else {
		b = clearBit(b, bit8)
	}

	buf.WriteByte(b)

}

// Value and CommandInfo for commands
func (infoObj InfoObj) writeCommandInfo(b byte, buf *bytes.Buffer) {
	if infoObj.CommandInfo.Quoc.Select {
		b = setBit(b, bit8)
	} else {
		b = clearBit(b, bit8)
	}
	b |= (byte(infoObj.CommandInfo.Quoc.Qu) << 2)
	buf.WriteByte(b)

}
func (infoObj InfoObj) writeQualitySeparateOctet(buf *bytes.Buffer) {
	var b byte

	if infoObj.Quality.Ov {
		b = setBit(b, bit1)
	} else {
		b = clearBit(b, bit1)
	}
	// bits 2..4 reserve
	if infoObj.Quality.Bl {
		b = setBit(b, bit5)
	} else {
		b = clearBit(b, bit5)
	}
	if infoObj.Quality.Sb {
		b = setBit(b, bit6)
	} else {
		b = clearBit(b, bit6)
	}
	if infoObj.Quality.Nt {
		b = setBit(b, bit7)
	} else {
		b = clearBit(b, bit7)
	}
	if infoObj.Quality.Iv {
		b = setBit(b, bit8)
	} else {
		b = clearBit(b, bit8)
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

func (asdu Asdu) serialize(state State, buf *bytes.Buffer) {
	// todo err
	var x byte = 0

	buf.WriteByte(byte(asdu.TypeId))

	if asdu.Sequence {
		x = setBit(x, bit8)
	} else {
		x = clearBit(x, bit8)
	}
	x |= byte(asdu.Num)
	buf.WriteByte(x)

	x = 0
	if asdu.Test {
		x = setBit(x, bit8)
	} else {
		x = clearBit(x, bit8)
	}

	if asdu.Negative {
		x = setBit(x, bit7)
	} else {
		x = clearBit(x, bit7)
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

		infoObj.writeInfo(asdu.TypeId, buf)

		if TypeIsTimeTagged(asdu.TypeId) {
			infoObj.serializeTime(state, buf)
		}

	}

}

func (infoObj InfoObj) serializeTime(state State, buf *bytes.Buffer) {

	// todo su, iv
	var iv, su bool

	var timetag time.Time
	var ivMask byte
	var suMask byte
	var milliWord uint16
	var weekDay int

	if iv {
		ivMask = setBit(ivMask, bit8)
	}

	if su {
		suMask = setBit(suMask, bit8)
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

func parseAsdu(buf *bytes.Buffer) (Asdu, error) {
	var b byte
	var b1 byte
	asdu := NewAsdu()

	b, _ = buf.ReadByte()
	asdu.TypeId = TypeId(b)

	// variable structure verifier
	b, _ = buf.ReadByte()

	asdu.Sequence = hasBit(b, bit8)
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
		newInfoObj.parseInfoObj(asdu.TypeId, buf)

		asdu.InfoObj = append(asdu.InfoObj, newInfoObj)
	}

	return *asdu, nil
}

func (infoObj *InfoObj) parseInfoObj(typeId TypeId, buf *bytes.Buffer) error {

	ioa1, _ := buf.ReadByte()
	ioa2, _ := buf.ReadByte()
	ioa3, _ := buf.ReadByte()

	infoObj.Ioa = Ioa(uint32(ioa1) + uint32(ioa2)*256 + uint32(ioa3)*65536)

	switch typeId {
	case M_SP_NA_1, M_SP_TB_1, M_DP_NA_1, M_DP_TB_1:
		// SP, DP
		infoObj.parseSiqDiq(typeId, buf)

	case M_ME_NA_1, M_ME_TD_1, M_ME_NB_1, M_ME_TE_1, M_ME_NC_1, M_ME_TF_1:
		// MV
		infoObj.parseMvValue(typeId, buf)

	case C_SC_NA_1, C_SC_TA_1, C_DC_NA_1, C_DC_TA_1:
		// SC, DC
		infoObj.parseScoDco(typeId, buf)

	case C_IC_NA_1:
		// GI
		b, _ := buf.ReadByte()

		infoObj.CommandInfo.Qoi = Qoi(uint8(b))
		infoObj.Value = IntVal(0)

	}

	return nil
}

func (infoObj *InfoObj) parseSiqDiq(typeId TypeId, buf *bytes.Buffer) {

	b := infoObj.parseQds(typeId, buf)

	switch typeId {
	case M_SP_NA_1, M_SP_TB_1:
		infoObj.Value = IntVal(b & 0x01)

	case M_DP_NA_1, M_DP_TB_1:
		infoObj.Value = IntVal(b & 0x03)
	}
}

func (infoObj *InfoObj) parseMvValue(typeId TypeId, buf *bytes.Buffer) {

	switch typeId {
	case M_ME_NA_1, M_ME_TD_1, M_ME_NB_1, M_ME_TE_1:
		// normalized value and scaled value
		b1, _ := buf.ReadByte()
		b2, _ := buf.ReadByte()

		value := uint16(b1) + 256*uint16(b2)
		infoObj.Value = IntVal(value)

		infoObj.parseQds(typeId, buf)

		if TypeIsTimeTagged(typeId) {
			infoObj.parseTimeTag(buf)
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

		infoObj.parseQds(typeId, buf)

	}

}

func (infoObj *InfoObj) parseQds(typeid TypeId, buf *bytes.Buffer) byte {
	b, _ := buf.ReadByte()

	quality := &infoObj.Quality

	if hasBit(b, bit8) {
		quality.Iv = true
	} else {
		quality.Iv = false
	}
	if hasBit(b, bit7) {
		quality.Nt = true
	} else {
		quality.Nt = false
	}
	if hasBit(b, bit6) {
		quality.Sb = true
	} else {
		quality.Sb = false
	}
	if hasBit(b, bit5) {
		quality.Bl = true
	} else {
		quality.Bl = false
	}

	switch typeid {
	case M_ME_NA_1, M_ME_TD_1, M_ME_NB_1, M_ME_TE_1, M_ME_NC_1, M_ME_TF_1:
		// all MV values have also OV flag

		if hasBit(b, bit1) {
			quality.Ov = true
		} else {
			quality.Ov = false
		}
	}

	return b
}

func (infoObj *InfoObj) parseScoDco(typeId TypeId, buf *bytes.Buffer) {
	b, _ := buf.ReadByte()

	infoObj.CommandInfo.Quoc.Select = hasBit(b, bit8)

	switch typeId {
	case C_SC_NA_1, C_SC_TA_1:
		infoObj.Value = IntVal(b & 0x01)
	case C_DC_NA_1, C_DC_TA_1:
		infoObj.Value = IntVal(b & 0x03)
	}

	infoObj.CommandInfo.Quoc.Qu = Qu(uint8(b >> 3))

}

func (infoObj *InfoObj) parseTimeTag(buf *bytes.Buffer) {

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
