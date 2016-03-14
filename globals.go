package ib

import "math"

var (
	CLIENT_ID_INCR  int64   = 999
	DELIM_STR       string  = "\000"
	DELIM_BYTE      byte    = '\000'
	REQUEST_CODE            = make(map[string]int64)
	REQUEST_VERSION         = make(map[string]int64)
	RESPONSE_CODE           = make(map[string]string)
	MAX_INT         int64   = math.MaxInt64
	MAX_FLOAT       float64 = math.MaxFloat64
)

func init() {
	RESPONSE_CODE["ErrMsg"] = "4"
}
