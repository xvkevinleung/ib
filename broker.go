package ib 

import (
	"bytes"
	"net"
	"strconv"
	"bufio"
	"strings"
)

type Broker struct {
	ClientId int64
	Conn net.Conn
	ReqId int64
	OutStream *bytes.Buffer
	InStream *bufio.Reader
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
}

func (b *Broker) Connect() error {
	conn, err := net.Dial("tcp", Conf.Host + ":" + Conf.Port)
	
	if err != nil {
		Log.Print("error", "unable to connect to IB via tcp")
		return err
	}

	b.Conn = conn

	b.InStream = bufio.NewReader(b.Conn)

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
		b, err := buf.ReadString(DelimByte)

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
const DelimStr = "\000"
const DelimByte = '\000'

func (b *Broker) WriteString(s string) (int, error) {
	return b.OutStream.WriteString(s + DelimStr)
}

func (b *Broker) WriteInt(i int64) (int, error) {
	return b.OutStream.WriteString(strconv.FormatInt(i, 10) + DelimStr)
}

func (b *Broker) WriteFloat(f float64) (int, error) {
	return b.OutStream.WriteString(strconv.FormatFloat(f, 'g', 10, 64) + DelimStr)
}

func (b *Broker) WriteBool(boo bool) (int, error) {
	if boo {
		return b.OutStream.WriteString("1")
	}

	return b.OutStream.WriteString("0")
}

func (b *Broker) ReadString() (string, error) {
	if str, err := b.InStream.ReadString(DelimByte); err != nil {
		return "", err
	} else {
		return strings.TrimRight(str, DelimStr), err
	}
}

func (b *Broker) ReadInt() (int64, error) {
	if str, err := b.ReadString(); err != nil {
		return 0, err
	} else {
		return strconv.ParseInt(strings.TrimRight(str, DelimStr), 10, 64)
	}
}

func (b *Broker) ReadFloat() (float64, error) {
	if str, err := b.ReadString(); err != nil {
		return 0, err
	} else {
		return strconv.ParseFloat(strings.TrimRight(str, DelimStr), 64)
	}
}

func (b *Broker) ReadBool() (bool, error) {
	if int, err := b.ReadInt(); err != nil {
		return false, err
	} else {
		if int != 0 {
			return true, err
		}

		return false, err
	}
}
