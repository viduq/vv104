// Code generated by "stringer -type=TypeId"; DO NOT EDIT.

package vv104

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[M_SP_NA_1-1]
	_ = x[M_DP_NA_1-3]
	_ = x[M_ST_NA_1-5]
	_ = x[M_BO_NA_1-7]
	_ = x[M_ME_NA_1-9]
	_ = x[M_ME_NB_1-11]
	_ = x[M_ME_NC_1-13]
	_ = x[M_IT_NA_1-15]
	_ = x[M_SP_TB_1-30]
	_ = x[M_DP_TB_1-31]
	_ = x[M_ST_TB_1-32]
	_ = x[M_BO_TB_1-33]
	_ = x[M_ME_TD_1-34]
	_ = x[M_ME_TE_1-35]
	_ = x[M_ME_TF_1-36]
	_ = x[M_IT_TB_1-37]
	_ = x[C_SC_NA_1-45]
	_ = x[C_DC_NA_1-46]
	_ = x[C_RC_NA_1-47]
	_ = x[C_SE_NA_1-48]
	_ = x[C_SE_NB_1-49]
	_ = x[C_SE_NC_1-50]
	_ = x[C_BO_NA_1-51]
	_ = x[C_SC_TA_1-58]
	_ = x[C_DC_TA_1-59]
	_ = x[C_RC_TA_1-60]
	_ = x[C_SE_TA_1-61]
	_ = x[C_SE_TB_1-62]
	_ = x[C_SE_TC_1-63]
	_ = x[C_BO_TA_1-64]
	_ = x[M_EI_NA_1-70]
	_ = x[C_IC_NA_1-100]
	_ = x[C_CI_NA_1-101]
	_ = x[C_RD_NA_1-102]
	_ = x[C_CS_NA_1-103]
	_ = x[C_RP_NA_1-105]
	_ = x[C_TS_TA_1-107]
}

const _TypeId_name = "M_SP_NA_1M_DP_NA_1M_ST_NA_1M_BO_NA_1M_ME_NA_1M_ME_NB_1M_ME_NC_1M_IT_NA_1M_SP_TB_1M_DP_TB_1M_ST_TB_1M_BO_TB_1M_ME_TD_1M_ME_TE_1M_ME_TF_1M_IT_TB_1C_SC_NA_1C_DC_NA_1C_RC_NA_1C_SE_NA_1C_SE_NB_1C_SE_NC_1C_BO_NA_1C_SC_TA_1C_DC_TA_1C_RC_TA_1C_SE_TA_1C_SE_TB_1C_SE_TC_1C_BO_TA_1M_EI_NA_1C_IC_NA_1C_CI_NA_1C_RD_NA_1C_CS_NA_1C_RP_NA_1C_TS_TA_1"

var _TypeId_map = map[TypeId]string{
	1:   _TypeId_name[0:9],
	3:   _TypeId_name[9:18],
	5:   _TypeId_name[18:27],
	7:   _TypeId_name[27:36],
	9:   _TypeId_name[36:45],
	11:  _TypeId_name[45:54],
	13:  _TypeId_name[54:63],
	15:  _TypeId_name[63:72],
	30:  _TypeId_name[72:81],
	31:  _TypeId_name[81:90],
	32:  _TypeId_name[90:99],
	33:  _TypeId_name[99:108],
	34:  _TypeId_name[108:117],
	35:  _TypeId_name[117:126],
	36:  _TypeId_name[126:135],
	37:  _TypeId_name[135:144],
	45:  _TypeId_name[144:153],
	46:  _TypeId_name[153:162],
	47:  _TypeId_name[162:171],
	48:  _TypeId_name[171:180],
	49:  _TypeId_name[180:189],
	50:  _TypeId_name[189:198],
	51:  _TypeId_name[198:207],
	58:  _TypeId_name[207:216],
	59:  _TypeId_name[216:225],
	60:  _TypeId_name[225:234],
	61:  _TypeId_name[234:243],
	62:  _TypeId_name[243:252],
	63:  _TypeId_name[252:261],
	64:  _TypeId_name[261:270],
	70:  _TypeId_name[270:279],
	100: _TypeId_name[279:288],
	101: _TypeId_name[288:297],
	102: _TypeId_name[297:306],
	103: _TypeId_name[306:315],
	105: _TypeId_name[315:324],
	107: _TypeId_name[324:333],
}

func (i TypeId) String() string {
	if str, ok := _TypeId_map[i]; ok {
		return str
	}
	return "TypeId(" + strconv.FormatInt(int64(i), 10) + ")"
}
