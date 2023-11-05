package main

import (
	"bytes"
	"fmt"
	"rth/vv104"
	"time"
)

func main() {
	fmt.Println("vv140 started.")

	state := vv104.State{
		Config: vv104.Config{},
	}

	m_me_td_1 := vv104.NewApdu()

	m_me_td_1.Apci.FrameFormat = vv104.IFormatFrame
	m_me_td_1.Apci.Rsn = 1
	m_me_td_1.Apci.Ssn = 10
	m_me_td_1.Asdu.TypeId = vv104.M_ME_TD_1
	m_me_td_1.Asdu.CauseTx = vv104.Spont
	m_me_td_1.Asdu.Casdu = 1
	m_me_td_1.Asdu.InfoObj.Ioa = 1235
	m_me_td_1.Asdu.InfoObj.Value = vv104.IntVal(32767)
	m_me_td_1.Asdu.InfoObj.TimeTag = time.Date(2023, 11, 1, 19, 44, 57, 23000000, time.Local)

	apduBytes := m_me_td_1.Serialize(state)
	fmt.Printf("have: %x\n", apduBytes)

	m_me_td_1_bytes := "\x68\x17\x14\x00\x02\x00\x22\x01\x03\x00\x01\x00\xd3\x04\x00\xff\x7f\x00\xbf\xde\x2c\x13\x61\x0b\x17"
	fmt.Printf("equal?: %v\n", bytes.Equal(apduBytes, []byte(m_me_td_1_bytes)))

	sframebytes := "\x68\x04\x01\x00\xfe\xff"
	var buf bytes.Buffer
	buf.Write([]byte(sframebytes))
	sframe, err := vv104.ParseApdu(&buf)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(sframe)

}

// Generic aproach for InfoValue
// var b vv104.FrameFormat = vv104.UFormat

// fmt.Println(b)

// a := vv104.IFormat
// fmt.Println(a)

// info := vv104.Info[int32]{
// 	Value:   0,
// 	Quality: vv104.Quality{},
// }
// as := vv104.Asdu[int32]{
// 	TypeId:   0,
// 	Num:      0,
// 	Sequence: false,
// 	CauseTx:  0,
// 	Negative: false,
// 	Test:     false,
// 	OrigAddr: 0,
// 	Casdu:    0,
// 	Ioa:      0,
// 	Info:     info,
// }

// fmt.Println(as)

// // interface aproach
// var b_ vv104.FrameFormat = vv104.UFormat

// fmt.Println(b_)

// a_ := vv104.IFormat
// fmt.Println(a_)

// var val vv104.Floatvalue = 23

// info_ := vv104.Info_{
// 	Value:   val,
// 	Quality: vv104.Quality{},
// }

// as_ := vv104.Asdu_{
// 	TypeId:   0,
// 	Num:      0,
// 	Sequence: false,
// 	CauseTx:  0,
// 	Negative: false,
// 	Test:     false,
// 	OrigAddr: 0,
// 	Casdu:    0,
// 	Ioa:      0,
// 	Info:     info_,
// }

// fmt.Println(as_)

/*

	info.Value = intval
	info.Quality = q

	fmt.Println("Print Info: ", info)
	info.Value = floatval

	fmt.Println("Print Info: ", info)

	info.Value = vv104.IntVal(1)
	info.Ioa = 33001
	var asdu = vv104.Asdu{
		TypeId:   vv104.M_ME_TF_1,
		Num:      1,
		Sequence: false,
		CauseTx:  vv104.Spont,
		Negative: false,
		Test:     false,
		OrigAddr: 0,
		Casdu:    1185,
		InfoObj:  info,
	}

	fmt.Println(asdu)

	var q1 vv104.Quality
	fmt.Println(q1)

	asdu.Serialize(state, buf)

	fmt.Println(buf.Bytes())

	info.Value = vv104.IntVal(2)
	info.Ioa = 33001
	var asdu1 = vv104.Asdu{
		TypeId:   vv104.M_DP_NA_1,
		Num:      1,
		Sequence: false,
		CauseTx:  vv104.Spont,
		Negative: false,
		Test:     false,
		OrigAddr: 0,
		Casdu:    1185,
		InfoObj:  info,
	}
	fmt.Println(asdu1)



*/
