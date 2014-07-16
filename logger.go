package ib

import (
	"log"
	"os"
	"time"
)

type IBLog struct {
	L *log.Logger
}

var Log = IBLog{log.New(os.Stdout, "", 0)}

func (l *IBLog) Print(status string, message interface{}) {
	t := time.Now().Format("2006/01/02 15:04:05")
	l.L.Printf("%s\t%s\t%s\n", t, status, message)
}

func (l *IBLog) PrintFloat(status string, message float64) {
	t := time.Now().Format("2006/01/02 15:04:05")
	l.L.Printf("%s\t%s\t%f\n", t, status, message)
}
