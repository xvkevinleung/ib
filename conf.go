package ib

var Conf = ConfStruct{
	Host: "127.0.0.1",
	Port: "4001",
	ClientVersion: 63,
	Delim: "\000",
}

type ConfStruct struct {
	Host string
	Port string
	ClientVersion int64
	Delim string
}
