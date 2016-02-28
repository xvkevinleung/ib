package ib

var (
	CLIENT_ID_INCR  int64  = 999
	DELIM_STR       string = "\000"
	DELIM_BYTE      byte   = '\000'
	REQUEST_CODE    map[string]int64
	REQUEST_VERSION map[string]int64
	RESPONSE_CODE   map[string]string
)

func init() {
	REQUEST_CODE = make(map[string]int64)
	REQUEST_VERSION = make(map[string]int64)
	RESPONSE_CODE = make(map[string]string)

	RESPONSE_CODE["ErrMsg"] = "4"
}
