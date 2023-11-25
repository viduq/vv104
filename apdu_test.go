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

	infoObj := NewInfoObj()
	infoObj.Ioa = 1235
	infoObj.Value = IntVal(32767)
	infoObj.TimeTag = time.Date(2023, 11, 1, 19, 44, 57, 23000000, time.Local)
	m_me_td_1.Asdu.AddInfoObject(infoObj)
	m_me_td_1.Apci.length = 23

	gi_act := NewApdu()
	gi_act.Apci.FrameFormat = IFormatFrame
	gi_act.Apci.Rsn = 0
	gi_act.Apci.Ssn = 0
	gi_act.Asdu.TypeId = C_IC_NA_1
	gi_act.Asdu.CauseTx = Act
	infoObj = NewInfoObj()
	infoObj.Ioa = 0
	infoObj.CommandInfo.Qoi = statioInterrogation
	infoObj.Value = IntVal(0)
	gi_act.Asdu.AddInfoObject(infoObj)
	gi_act.Apci.length = 14

	sc_act := NewApdu()
	sc_act.Apci.FrameFormat = IFormatFrame
	sc_act.Apci.Ssn = 1
	sc_act.Apci.Rsn = 12
	sc_act.Asdu.TypeId = C_SC_NA_1
	sc_act.Asdu.CauseTx = Act
	infoObj = NewInfoObj()
	infoObj.Ioa = 4500
	infoObj.Value = IntVal(1)
	sc_act.Asdu.AddInfoObject(infoObj)
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

	// 16 DPs in one frame
	m_dp_na_1_16x := NewApdu()
	m_dp_na_1_16x.Apci.Rsn = 1
	m_dp_na_1_16x.Apci.Ssn = 1
	m_dp_na_1_16x.Asdu.TypeId = M_DP_NA_1
	m_dp_na_1_16x.Asdu.CauseTx = Inrogen
	m_dp_na_1_16x.Asdu.Casdu = 1
	infoObj = InfoObj{}
	infoObj.Ioa = 35
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 70000
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 70001
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 70005
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 70002
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 70004
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 123
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 124
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 16000000
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 126
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 127
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 1000000
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 70003
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 30
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 128
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	infoObj = InfoObj{}
	infoObj.Ioa = 125
	infoObj.Value = IntVal(0)
	m_dp_na_1_16x.Asdu.AddInfoObject(infoObj)
	m_dp_na_1_16x.Apci.length = 74

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
		{"16x DP inrogen", m_dp_na_1_16x, "\x68\x4a\x02\x00\x02\x00\x03\x10\x14\x00\x01\x00\x23\x00\x00\x00\x70\x11\x01\x00\x71\x11\x01\x00\x75\x11\x01\x00\x72\x11\x01\x00\x74\x11\x01\x00\x7b\x00\x00\x00\x7c\x00\x00\x00\x00\x24\xf4\x00\x7e\x00\x00\x00\x7f\x00\x00\x00\x40\x42\x0f\x00\x73\x11\x01\x00\x1e\x00\x00\x00\x80\x00\x00\x00\x7d\x00\x00\x00"},
	}
	// execution loop: Serialize Apdus
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, _ := tt.apdu.Serialize(state)
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
