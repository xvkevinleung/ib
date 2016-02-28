package ib

var (
	CLIENT_ID_INCR int64  = 999
	DELIM_STR      string = "\000"
	DELIM_BYTE     byte   = '\000'
)

var REQUEST_CODE = make(map[string]int64)
var REQUEST_VERSION = make(map[string]int64)
var RESPONSE_CODE = make(map[string]string)

func init() {
	RESPONSE_CODE["ErrMsg"] = "4"
}
