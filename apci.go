package vv104

import (
	"bytes"
)

// Apci header of the iec 104 frame

const (
	startbyte byte = 0x68
)

type SeqNumber int

type Apci struct {
	length      uint8
	FrameFormat FrameFormat
	Rsn         SeqNumber
	Ssn         SeqNumber
	UFormat     UFormat
}

func (apci *Apci) serialize(state State, buf *bytes.Buffer, asduLength uint8) {
	buf.WriteByte(startbyte)
	apci.length = asduLength + 4
	buf.WriteByte(byte(apci.length)) // todo check for overfow
	apci.Ssn = state.sendAck.seqNumber
	apci.Rsn = state.recvAck.seqNumber
	apci.writeCtrlFields(state, buf)
}

func (apci Apci) writeCtrlFields(state State, buf *bytes.Buffer) {
	// var b byte
	switch apci.FrameFormat {

	case IFormatFrame:
		buf.WriteByte(byte(apci.Ssn) << 1) // & 0xFF)
		buf.WriteByte(byte(apci.Ssn >> 7))
		buf.WriteByte(byte(apci.Rsn) << 1) //& 0xFF
		buf.WriteByte(byte(apci.Rsn >> 7))
	case SFormatFrame:
		buf.WriteByte(0b00000001)
		buf.WriteByte(0)
		buf.WriteByte(byte(apci.Rsn) << 1) //& 0xFF
		buf.WriteByte(byte(apci.Rsn >> 7))

	case UFormatFrame:
		buf.WriteByte(byte(apci.UFormat) | 0b00000011)
		buf.WriteByte(0)
		buf.WriteByte(0)
		buf.WriteByte(0)
	}
}
