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

	infoObj := vv104.InfoObj{}
	infoObj.Ioa = 1235
	infoObj.Value = vv104.IntVal(32767)
	infoObj.TimeTag = time.Date(2023, 11, 1, 19, 44, 57, 23000000, time.Local)
	m_me_td_1.Asdu.AddInfoObject(infoObj)

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

	// m_dp_na_1_16x := vv104.NewApdu()
	// m_dp_na_1_16x.Apci.Rsn = 1
	// m_dp_na_1_16x.Apci.Ssn = 10
	// m_dp_na_1_16x.Asdu.TypeId = vv104.M_DP_NA_1
	// m_dp_na_1_16x.Asdu.CauseTx = vv104.Inrogen
	// m_dp_na_1_16x.Asdu.Casdu = 1
	// m_dp_na_1_16x.Asdu.InfoObj.Ioa = 35
	// m_dp_na_1_16x.Asdu.InfoObj.Value = vv104.IntVal(0)

	// sixteen_dps := `\x68\x4a\x02\x00\x02\x00\x03\x10\x14\x00\x01\x00\x23\x00
	// \x00\x00\x70\x11\x01\x00\x71\x11\x01\x00\x75\x11\x01\x00\x72\x11
	// \x01\x00\x74\x11\x01\x00\x7b\x00\x00\x00\x7c\x00\x00\x00\x00\x24
	// \xf4\x00\x7e\x00\x00\x00\x7f\x00\x00\x00\x40\x42\x0f\x00\x73\x11
	// \x01\x00\x1e\x00\x00\x00\x80\x00\x00\x00\x7d\x00\x00\x00
	// `

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
