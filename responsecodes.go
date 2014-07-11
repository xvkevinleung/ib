package ib

var RESPONSE ResStruct = ResStruct{
	CODE: ResCodeStruct{
		ERR_MSG: "4",
		CONTRACT_DATA: "10",
		TICK_PRICE: "1",
		TICK_SIZE: "2",
		TICK_OPTION_COMPUTATION: "21",
		TICK_GENERIC: "45",
		TICK_STRING: "46",
		TICK_EFP: "47",
		TICK_SNAPSHOT_END: "57",
		MARKET_DATA_TYPE: "58",
		HISTORICAL_DATA: "17",
	},
}

// types for bucketing response codes
type ResStruct struct {
	CODE ResCodeStruct 
}

type ResCodeStruct struct {
	ERR_MSG string
	CONTRACT_DATA string 
	TICK_PRICE string
	TICK_SIZE string
	TICK_OPTION_COMPUTATION string
	TICK_GENERIC string
	TICK_STRING string
	TICK_EFP string
	TICK_SNAPSHOT_END string
	MARKET_DATA_TYPE string
	HISTORICAL_DATA string
}
