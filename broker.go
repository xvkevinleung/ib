package ib 

import (
	"bytes"
	"net"
	"strconv"
	"bufio"
)

type Broker struct {
	ClientId int64
	Conn net.Conn
	ReqId int64
	OutStream *bytes.Buffer
	InStream *bytes.Buffer
}

// GLOBAL
var ClientIdIncr int64 = 999 

func NextClientId() int64 {
	ClientIdIncr += 1
	return ClientIdIncr
}

func (b *Broker) NextReqId() int64 {
	b.ReqId += 1
	return b.ReqId
}

func (b *Broker) Initialize() {
	b.ClientId = NextClientId()
	b.OutStream = bytes.NewBuffer(make([]byte, 0, 4096))
	b.InStream = bytes.NewBuffer(make([]byte, 0, 4096))
}

func (b *Broker) Connect() error {
	conn, err := net.Dial("tcp", Conf.Host + ":" + Conf.Port)
	
	if err != nil {
		Log.Print("error", "unable to connect to IB via tcp")
		return err
	}

	b.Conn = conn

	Log.Print("okay", "connected via tcp")

	err = b.ServerShake()

	if err != nil {
		Log.Print("error", "unable to perform server shake")
		return err
	}

	return err 
}

func (b *Broker) ServerShake() error {
	b.WriteInt(Conf.ClientVersion)
	b.WriteInt(b.ClientId)

	_, err := b.SendRequest()

	return err
}

func (b *Broker) Disconnect() error {
	return b.Conn.Close()
}

func (b *Broker) SendRequest() (int, error) {
	Log.Print("request", b.OutStream.String())

	b.WriteString("\000")

	i, err := b.Conn.Write(b.OutStream.Bytes())

	b.OutStream.Reset()

	return i, err
}

func (b *Broker) Listen() {
	buf := bufio.NewReader(b.Conn)

	for {
		b, err := buf.ReadString('\000')

		if err != nil {
			continue
		}

		Log.Print("response", string(b))
	}
}

func (b *Broker) SetServerLogLevel(i int64) {
	b.WriteInt(14)
	b.WriteInt(1)
	b.WriteInt(i)

	b.SendRequest()
}

// GLOBAL
const Delim = "\000"

func (b *Broker) WriteString(s string) (int, error) {
	return b.OutStream.WriteString(s + Delim)
}

func (b *Broker) WriteInt(i int64) (int, error) {
	return b.OutStream.WriteString(strconv.FormatInt(i, 10) + Delim)
}

func (b *Broker) WriteFloat(f float64) (int, error) {
	return b.OutStream.WriteString(strconv.FormatFloat(f, 'g', 10, 64) + Delim)
}

func (b *Broker) WriteBool(boo bool) (int, error) {
	if boo {
		return b.OutStream.WriteString("1")
	}

	return b.OutStream.WriteString("0")
}


