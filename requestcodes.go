package ib

var REQUEST ReqStruct = ReqStruct{
	CODE: ReqCodeStruct{
		MARKET_DATA: 1,
		CONTRACT_DATA: 9,
	},
	VERSION: ReqVersionStruct{
		MARKET_DATA: 10,
		CONTRACT_DATA: 7,
	},
}

// types for bucketing request codes
type ReqStruct struct {
	CODE ReqCodeStruct 
	VERSION ReqVersionStruct
}

type ReqCodeStruct struct {
	MARKET_DATA int64
	CONTRACT_DATA int64
}

type ReqVersionStruct struct {
	MARKET_DATA int64
	CONTRACT_DATA int64
}
