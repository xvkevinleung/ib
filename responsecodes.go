package ib

var RESPONSE ResStruct = ResStruct{
	CODE: ResCodeStruct{
		ERR_MSG: "4" + Conf.Delim,
		CONTRACT_DATA: "10" + Conf.Delim,
		TICK_PRICE: "1" + Conf.Delim,
		TICK_SIZE: "2" + Conf.Delim,
		TICK_OPTION_COMPUTATION: "21" + Conf.Delim,
		TICK_GENERIC: "45" + Conf.Delim,
		TICK_STRING: "46" + Conf.Delim,
		TICK_EFP: "47" + Conf.Delim,
		TICK_SNAPSHOT_END: "57" + Conf.Delim,
		MARKET_DATA_TYPE: "58" + Conf.Delim,
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
}
