package ib

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// REQUESTS
////////////////////////////////////////////////////////////////////////////////

type RealTimeBarsRequest struct {
	Contract Contract
	Bar      int64
	Show     string // what to show
	Rth      bool   `json:",string"` // regular trading hours
	Opts     string // use default "XYZ"
}

func init() {
	REQUEST_CODE["RealTimeBars"] = 50
	REQUEST_VERSION["RealTimeBars"] = 2
}

func (r *RealTimeBarsRequest) Send(id int64, b *RealTimeBarsBroker) {
	b.Contracts[id] = r.Contract
	b.WriteInt(REQUEST_CODE["RealTimeBars"])
	b.WriteInt(REQUEST_VERSION["RealTimeBars"])
	b.WriteInt(id)
	b.WriteInt(r.Contract.ContractId)
	b.WriteString(r.Contract.Symbol)
	b.WriteString(r.Contract.SecurityType)
	b.WriteString(r.Contract.Expiry)
	b.WriteFloat(r.Contract.Strike)
	b.WriteString(r.Contract.Right)
	b.WriteString(r.Contract.Multiplier)
	b.WriteString(r.Contract.Exchange)
	b.WriteString(r.Contract.PrimaryExchange)
	b.WriteString(r.Contract.Currency)
	b.WriteString(r.Contract.LocalSymbol)
	b.WriteString(r.Contract.TradingClass)
	b.WriteInt(r.Bar)
	b.WriteString(r.Show)
	b.WriteBool(r.Rth)

	b.Broker.SendRequest()
}

////////////////////////////////////////////////////////////////////////////////
// RESPONSES
////////////////////////////////////////////////////////////////////////////////

type RealTimeBar struct {
	Rid      int64
	Time     string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   int64
	WAP      float64
	BarCount int64
}

func init() {
	RESPONSE_CODE["RealTimeBar"] = "50"
}

////////////////////////////////////////////////////////////////////////////////
// BROKER
////////////////////////////////////////////////////////////////////////////////

type RealTimeBarsBroker struct {
	Broker
	Contracts       map[int64]Contract
	RealTimeBarChan chan RealTimeBar
}

func NewRealTimeBarsBroker() RealTimeBarsBroker {
	b := RealTimeBarsBroker{
		Broker{},
		make(map[int64]Contract),
		make(chan RealTimeBar),
	}

	return b
}

func (b *RealTimeBarsBroker) Listen() {
	for {
		s, err := b.ReadString()

		if err != nil {
			continue
		}

		if s == RESPONSE_CODE["RealTimeBar"] {
			version, err := b.ReadString()

			if err != nil {
				continue
			}

			r := b.ReadRealTimeBar(version)
			b.RealTimeBarChan <- r
		}
	}
}

func (b *RealTimeBarsBroker) ReadRealTimeBar(version string) RealTimeBar {
	var r RealTimeBar

	r.Rid, _ = b.ReadInt()
	r.Time, _ = b.ReadString()
	r.Open, _ = b.ReadFloat()
	r.High, _ = b.ReadFloat()
	r.Low, _ = b.ReadFloat()
	r.Close, _ = b.ReadFloat()
	r.Volume, _ = b.ReadInt()
	r.WAP, _ = b.ReadFloat()
	r.BarCount, _ = b.ReadInt()

	return r
}

////////////////////////////////////////////////////////////////////////////////
// SERIALIZERS
////////////////////////////////////////////////////////////////////////////////

func (b *RealTimeBarsBroker) RealTimeBarToJSON(d *RealTimeBar) ([]byte, error) {
	c := b.Contracts[d.Rid]
	return json.Marshal(struct {
		Rid          int64
		Time         string
		Symbol       string
		SecurityType string
		Exchange     string
		Currency     string
		Right        string
		Strike       float64
		Expiry       string
		BarTime      string
		Open         float64
		High         float64
		Low          float64
		Close        float64
		Volume       int64
		WAP          float64
		BarCount     int64
	}{
		Rid:          d.Rid,
		Time:         strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		Symbol:       c.Symbol,
		SecurityType: c.SecurityType,
		Exchange:     c.Exchange,
		Currency:     c.Currency,
		Right:        c.Right,
		Strike:       c.Strike,
		Expiry:       c.Expiry,
		BarTime:      d.Time,
		Open:         d.Open,
		High:         d.High,
		Low:          d.Low,
		Close:        d.Close,
		Volume:       d.Volume,
		WAP:          d.WAP,
		BarCount:     d.BarCount,
	})
}

func (b *RealTimeBarsBroker) RealTimeBarToCSV(d *RealTimeBar) string {
	c := b.Contracts[d.Rid]
	return fmt.Sprintf(
		"%d,%s,%s,%s,%s,%s,%s,%.2f,%s,%s,%.2f,%.2f,%.2f,%.2f,%d,%.2f,%d",
		d.Rid,
		strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		c.Symbol,
		c.SecurityType,
		c.Exchange,
		c.Currency,
		c.Right,
		c.Strike,
		c.Expiry,
		d.Time,
		d.Open,
		d.High,
		d.Low,
		d.Close,
		d.Volume,
		d.WAP,
		d.BarCount,
	)
}
