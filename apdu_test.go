package vv104

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestVariousFrames(t *testing.T) {
	state := State{
		Config: Config{},
	}

	m_me_td_1 := NewApdu()
	m_me_td_1.Apci.FrameFormat = IFormatFrame
	m_me_td_1.Apci.Rsn = 1
	m_me_td_1.Apci.Ssn = 10
	m_me_td_1.Asdu.TypeId = M_ME_TD_1
	m_me_td_1.Asdu.CauseTx = Spont
	m_me_td_1.Asdu.Casdu = 1
	m_me_td_1.Asdu.InfoObj.Ioa = 1235
	m_me_td_1.Asdu.InfoObj.Value = IntVal(32767)
	m_me_td_1.Asdu.InfoObj.TimeTag = time.Date(2023, 11, 1, 19, 44, 57, 23000000, time.Local)
	m_me_td_1.Apci.length = 23

	gi_act := NewApdu()
	gi_act.Apci.FrameFormat = IFormatFrame
	gi_act.Apci.Rsn = 0
	gi_act.Apci.Ssn = 0
	gi_act.Asdu.TypeId = C_IC_NA_1
	gi_act.Asdu.CauseTx = Act
	gi_act.Asdu.InfoObj.Ioa = 0
	gi_act.Asdu.InfoObj.CommandInfo.Qoi = statioInterrogation
	gi_act.Apci.length = 14
	gi_act.Asdu.InfoObj.Value = IntVal(0)

	sc_act := NewApdu()
	sc_act.Apci.FrameFormat = IFormatFrame
	sc_act.Apci.Ssn = 1
	sc_act.Apci.Rsn = 12
	sc_act.Asdu.TypeId = C_SC_NA_1
	sc_act.Asdu.CauseTx = Act
	sc_act.Asdu.InfoObj.Ioa = 4500
	sc_act.Asdu.InfoObj.Value = IntVal(1)
	sc_act.Apci.length = 14

	startDtAct := NewApdu()
	startDtAct.Apci.FrameFormat = UFormatFrame
	startDtAct.Apci.UFormat = StartDTAct
	startDtAct.Apci.length = 4

	startDtCon := startDtAct
	startDtCon.Apci.UFormat = StartDTCon

	stopDtAct := startDtAct
	stopDtAct.Apci.UFormat = StopDTAct

	stopDtCon := startDtAct
	stopDtCon.Apci.UFormat = StopDTCon

	testFrAct := startDtAct
	testFrAct.Apci.UFormat = TestFRAct

	testFrCon := startDtAct
	testFrCon.Apci.UFormat = TestFRCon

	sframe_0 := NewApdu()
	sframe_0.Apci.FrameFormat = SFormatFrame
	sframe_0.Apci.Rsn = 0
	sframe_0.Apci.length = 4

	sframe_3 := sframe_0
	sframe_3.Apci.Rsn = 3

	sframe_32767 := sframe_0
	sframe_32767.Apci.Rsn = 32767

	// test table adjusted from here: https://blog.jetbrains.com/go/2022/11/22/comprehensive-guide-to-testing-in-go/
	var tests = []struct {
		name      string
		apdu      Apdu
		apduBytes string
	}{
		{"M_ME_TD_1 spont", m_me_td_1, "\x68\x17\x14\x00\x02\x00\x22\x01\x03\x00\x01\x00\xd3\x04\x00\xff\x7f\x00\xbf\xde\x2c\x13\x61\x0b\x17"},
		{"C_IC_NA_1 act", gi_act, "\x68\x0e\x00\x00\x00\x00\x64\x01\x06\x00\x01\x00\x00\x00\x00\x14"},
		{"C_SC_NA_1 act", sc_act, "\x68\x0e\x02\x00\x18\x00\x2d\x01\x06\x00\x01\x00\x94\x11\x00\x01"},
		{"start dt act", startDtAct, "\x68\x04\x07\x00\x00\x00"},
		{"start dt con", startDtCon, "\x68\x04\x0b\x00\x00\x00"},
		{"stop dt act", stopDtAct, "\x68\x04\x13\x00\x00\x00"},
		{"stop dt con", stopDtCon, "\x68\x04\x23\x00\x00\x00"},
		{"testfr act", testFrAct, "\x68\x04\x43\x00\x00\x00"},
		{"testfr con", testFrCon, "\x68\x04\x83\x00\x00\x00"},
		{"sframe 0", sframe_0, "\x68\x04\x01\x00\x00\x00"},
		{"sframe 3", sframe_3, "\x68\x04\x01\x00\x06\x00"},
		{"sframe 32767", sframe_32767, "\x68\x04\x01\x00\xfe\xff"},
	}
	// execution loop: Serialize Apdus
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := tt.apdu.Serialize(state)
			if !bytes.Equal(ans, []byte(tt.apduBytes)) {
				t.Errorf("\ngot : %x \nwant: %x\n", ans, tt.apduBytes)
			}
		})
	}

	// execution loop: Parse Frames
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// andersrum als oben, nehme string, bekomme apdu
			var buf bytes.Buffer
			buf.Write([]byte(tt.apduBytes))
			ans, err := ParseApdu(&buf)
			if err != nil {
				fmt.Println(err)
			}

			if !reflect.DeepEqual(ans, tt.apdu) {
				t.Errorf("\ngot : %#v \nwant: %#v\n", ans, tt.apdu)
			}

		})
	}
}
