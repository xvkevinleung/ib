package ib

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type RealTimeBarsBroker struct {
	Broker
	Contracts       map[int64]Contract
	RealTimeBarChan chan RealTimeBar
}

type RealTimeBarsRequest struct {
	Con  Contract
	Bar  int64
	Show string // what to show
	Rth  bool   // regular trading hours
	Opts string // use default "XYZ"
}

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

func (r *RealTimeBarsBroker) BarToJSON(b *RealTimeBar) ([]byte, error) {
	c := r.Contracts[b.Rid]
	return json.Marshal(struct {
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
		Time:         strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		Symbol:       c.Symbol,
		SecurityType: c.SecurityType,
		Exchange:     c.Exchange,
		Currency:     c.Currency,
		Right:        c.Right,
		Strike:       c.Strike,
		Expiry:       c.Expiry,
		BarTime:      b.Time,
		Open:         b.Open,
		High:         b.High,
		Low:          b.Low,
		Close:        b.Close,
		Volume:       b.Volume,
		WAP:          b.WAP,
		BarCount:     b.BarCount,
	})
}

func (r *RealTimeBarsBroker) BarToCSV(b *RealTimeBar) string {
	c := r.Contracts[b.Rid]
	return fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s,%.2f,%s,%s,%.2f,%.2f,%.2f,%.2f,%d,%.2f,%d",
		strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		c.Symbol,
		c.SecurityType,
		c.Exchange,
		c.Currency,
		c.Right,
		c.Strike,
		c.Expiry,
		b.Time,
		b.Open,
		b.High,
		b.Low,
		b.Close,
		b.Volume,
		b.WAP,
		b.BarCount,
	)
}

func NewRealTimeBarsBroker() RealTimeBarsBroker {
	r := RealTimeBarsBroker{
		Broker{},
		make(map[int64]Contract),
		make(chan RealTimeBar),
	}

	return r
}

func (r *RealTimeBarsBroker) SendRequest(rid int64, d RealTimeBarsRequest) {
	r.Contracts[rid] = d.Con
	r.WriteInt(REQUEST.CODE.REALTIMEBARS)
	r.WriteInt(REQUEST.VERSION.REALTIMEBARS)
	r.WriteInt(rid)
	r.WriteInt(d.Con.ContractId)
	r.WriteString(d.Con.Symbol)
	r.WriteString(d.Con.SecurityType)
	r.WriteString(d.Con.Expiry)
	r.WriteFloat(d.Con.Strike)
	r.WriteString(d.Con.Right)
	r.WriteString(d.Con.Multiplier)
	r.WriteString(d.Con.Exchange)
	r.WriteString(d.Con.PrimaryExchange)
	r.WriteString(d.Con.Currency)
	r.WriteString(d.Con.LocalSymbol)
	r.WriteString(d.Con.TradingClass)
	r.WriteInt(d.Bar)
	r.WriteString(d.Show)
	r.WriteBool(d.Rth)

	r.Broker.SendRequest()
}

func (r *RealTimeBarsBroker) Listen() {
	for {
		b, err := r.ReadString()

		if err != nil {
			continue
		}

		if b == RESPONSE.CODE.REALTIMEBARS {
			version, err := r.ReadString()

			if err != nil {
				continue
			}

			r.ReadRealTimeBar(version)
		}
	}
}

func (r *RealTimeBarsBroker) ReadRealTimeBar(version string) {
	var d RealTimeBar

	d.Rid, _ = r.ReadInt()
	d.Time, _ = r.ReadString()
	d.Open, _ = r.ReadFloat()
	d.High, _ = r.ReadFloat()
	d.Low, _ = r.ReadFloat()
	d.Close, _ = r.ReadFloat()
	d.Volume, _ = r.ReadInt()
	d.WAP, _ = r.ReadFloat()
	d.BarCount, _ = r.ReadInt()

	r.RealTimeBarChan <- d
}
