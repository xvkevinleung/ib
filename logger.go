package ib

import (
	"log"
	"os"
	"time"
)

type IBLog struct{
	L *log.Logger
}

var Log = IBLog{log.New(os.Stdout, "", 0)}

func (l *IBLog) Print(status, message string) {
	t := time.Now().Format("2006/01/02 15:04:05")
	l.L.Printf("%s\t%s\t%s\n", t, status, message)
}
