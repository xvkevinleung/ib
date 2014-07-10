package ib

var REQUEST ReqStruct = ReqStruct{
	CODE: ReqCodeStruct{
		CONTRACT_DATA: 9,
	},
	VERSION: ReqVersionStruct{
		CONTRACT_DATA: 7,
	},
}

// types for bucketing request codes
type ReqStruct struct {
	CODE ReqCodeStruct 
	VERSION ReqVersionStruct
}

type ReqCodeStruct struct {
	CONTRACT_DATA int64
}

type ReqVersionStruct struct {
	CONTRACT_DATA int64
}
