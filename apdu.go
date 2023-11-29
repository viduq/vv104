package vv104

import (
	"bytes"
	"errors"
	"fmt"
)

type Apdu struct {
	Apci Apci
	Asdu Asdu
}

// CheckApdu checks all Apdu Fields if they comply to IEC 104 standard
// todo: extend
func (apdu Apdu) CheckApdu() (bool, error) {
	if apdu.Apci.FrameFormat == IFormatFrame {
		if apdu.Asdu.Num < 1 {
			return false, errors.New("Num < 1 for IFormat")
		}
		if apdu.Asdu.Casdu == 0 {
			return false, errors.New("Casdu = 0 is not permitted")
		}

		// ...
	}
	return true, nil
}

func (apdu Apdu) String() string {
	switch apdu.Apci.FrameFormat {

	case IFormatFrame:
		return fmt.Sprintf("(%d/%d) ", apdu.Apci.Rsn, apdu.Apci.Ssn) + apdu.Asdu.String()
	case SFormatFrame:
		return fmt.Sprintf("S-Format (%d/%d)\n", apdu.Apci.Rsn, apdu.Apci.Ssn)
	case UFormatFrame:
		return apdu.Apci.UFormat.String()

	}
	return ""
}

func (apdu *Apdu) Serialize(state State) ([]byte, error) { // TODO error
	asduBuf := new(bytes.Buffer)
	apciBuf := new(bytes.Buffer)

	switch apdu.Apci.FrameFormat {
	case IFormatFrame:
		apdu.Asdu.serialize(state, asduBuf)
		apdu.Apci.serialize(state, apciBuf, uint8(asduBuf.Len()))

		s := [][]byte{apciBuf.Bytes(), asduBuf.Bytes()}
		var emptySep []byte
		return bytes.Join(s, emptySep), nil

	case UFormatFrame, SFormatFrame:
		apdu.Apci.serialize(state, apciBuf, 0)
		return apciBuf.Bytes(), nil

	}

	return []byte{}, nil

}

func NewApdu() Apdu {
	apdu := Apdu{
		Apci: Apci{
			length:      6,
			FrameFormat: 0,
			Rsn:         0,
			Ssn:         0,
			UFormat:     0,
		},
		Asdu: Asdu{
			TypeId:   0,
			Num:      1,
			Sequence: false,
			CauseTx:  0,
			Negative: false,
			Test:     false,
			OrigAddr: 0,
			Casdu:    1,
			InfoObj:  []InfoObj{},
		},
	}
	return apdu
}

func ParseApdu(buf *bytes.Buffer) (Apdu, error) {
	var b byte
	var b1 byte
	var b2 byte
	// var err error
	apdu := NewApdu()

	if buf.Len() < 6 {
		return apdu, errors.New("buffer < 6 bytes, can't parse")
	}
	if buf.Len() > 253 { // todo check if 255?
		// hier müssen wahrscheinlich mehrere frames ausgewertet werden
		// todo
	}

	b, _ = buf.ReadByte()
	if b != STARTBYTE {
		return apdu, errors.New("startbyte is not first byte, todo")
	}

	b, _ = buf.ReadByte()
	if (uint8(b) > 253) || (uint8(b) < 4) {
		return apdu, errors.New("apdu len is not within range")
	}
	apdu.Apci.length = uint8(b)

	// ctrl field 1
	b, _ = buf.ReadByte()
	if (b & 0b0000_0001) == 0b0000_0000 {
		// i frame
		apdu.Apci.FrameFormat = IFormatFrame
		b1 = b
		b2, _ = buf.ReadByte()
		apdu.Apci.Ssn = SeqNumber((uint16(b1) >> 1) + (uint16(b2) << 7))

		b1, _ = buf.ReadByte()
		b2, _ = buf.ReadByte()
		apdu.Apci.Rsn = SeqNumber((uint16(b1) >> 1) + (uint16(b2) << 7))

		apdu.Asdu, _ = parseAsdu(buf)

	} else if (b & 0b0000_0011) == 0b0000_0001 {
		// s frame
		apdu.Apci.FrameFormat = SFormatFrame
		//ctrl field 2
		_, _ = buf.ReadByte() // empty byte
		//ctrl field 3
		b1, _ = buf.ReadByte()
		//ctrl field 4
		b2, _ = buf.ReadByte()

		apdu.Apci.Rsn = SeqNumber((uint16(b1) >> 1) + (uint16(b2) << 7))

	} else if (b & 0b0000_0011) == 0b0000_0011 {
		// u frame
		apdu.Apci.FrameFormat = UFormatFrame

		if b&byte(StartDTAct) == byte(StartDTAct) {
			apdu.Apci.UFormat = StartDTAct
		} else if b&byte(StartDTCon) == byte(StartDTCon) {
			apdu.Apci.UFormat = StartDTCon
		} else if b&byte(StopDTAct) == byte(StopDTAct) {
			apdu.Apci.UFormat = StopDTAct
		} else if b&byte(StopDTCon) == byte(StopDTCon) {
			apdu.Apci.UFormat = StopDTCon
		} else if b&byte(TestFRAct) == byte(TestFRAct) {
			apdu.Apci.UFormat = TestFRAct
		} else if b&byte(TestFRCon) == byte(TestFRCon) {
			apdu.Apci.UFormat = TestFRCon
		}

	}

	return apdu, nil
}
