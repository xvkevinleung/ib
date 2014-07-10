package ib

var RESPONSE ResStruct = ResStruct{
	CODE: ResCodeStruct{
		CONTRACT_DATA: "10" + Conf.Delim,
	},
}

// types for bucketing response codes
type ResStruct struct {
	CODE ResCodeStruct 
}

type ResCodeStruct struct {
	CONTRACT_DATA string 
}
