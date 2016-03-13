package ib

import (
	"encoding/json"
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// REQUESTS
////////////////////////////////////////////////////////////////////////////////

type HistoricalDataRequest struct {
	Contract Contract
	End      string
	Bar      string
	Dur      string
	Rth      bool `json:",string"`
	Show     string
	Datef    int64
}

func init() {
	REQUEST_CODE["HistoricalData"] = 20
	REQUEST_VERSION["HistoricalData"] = 5
}

func (r *HistoricalDataRequest) Send(id int64, b *HistoricalDataBroker) {
	b.Contract = r.Contract
	b.WriteInt(REQUEST_CODE["HistoricalData"])
	b.WriteInt(REQUEST_VERSION["HistoricalData"])
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
	b.WriteInt(0) // include expired
	b.WriteString(r.End)
	b.WriteString(r.Bar)
	b.WriteString(r.Dur)
	b.WriteBool(r.Rth)
	b.WriteString(r.Show)
	b.WriteInt(r.Datef)

	b.Broker.SendRequest()
}

////////////////////////////////////////////////////////////////////////////////
// RESPONSES
////////////////////////////////////////////////////////////////////////////////

type HistoricalData struct {
	Rid   string
	Start string
	End   string
	Count int64
	Data  []HistoricalDataItem
}

type HistoricalDataItem struct {
	Date     string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   int64
	WAP      float64
	HasGaps  bool
	BarCount int64
}

func init() {
	RESPONSE_CODE["HistoricalData"] = "17"
}

////////////////////////////////////////////////////////////////////////////////
// BROKER
////////////////////////////////////////////////////////////////////////////////

type HistoricalDataBroker struct {
	Broker
	Contract           Contract
	HistoricalDataChan chan HistoricalData
}

func NewHistoricalDataBroker() HistoricalDataBroker {
	b := HistoricalDataBroker{Broker{}, Contract{}, make(chan HistoricalData)}
	return b
}

func (b *HistoricalDataBroker) Listen() {
	for {
		s, err := b.ReadString()

		if err != nil {
			continue
		}

		if s == RESPONSE_CODE["HistoricalData"] {
			version, err := b.ReadString()

			if err != nil {
				continue
			}

			r := b.ReadHistoricalData(version)
			b.HistoricalDataChan <- r
		}
	}
}

func (b *HistoricalDataBroker) ReadHistoricalData(version string) HistoricalData {
	var r HistoricalData

	r.Rid, _ = b.ReadString()
	r.Start, _ = b.ReadString()
	r.End, _ = b.ReadString()
	r.Count, _ = b.ReadInt()

	r.Data = make([]HistoricalDataItem, r.Count)

	for i := range r.Data {
		r.Data[i].Date, _ = b.ReadString()
		r.Data[i].Open, _ = b.ReadFloat()
		r.Data[i].High, _ = b.ReadFloat()
		r.Data[i].Low, _ = b.ReadFloat()
		r.Data[i].Close, _ = b.ReadFloat()
		r.Data[i].Volume, _ = b.ReadInt()
		r.Data[i].WAP, _ = b.ReadFloat()
		r.Data[i].HasGaps, _ = b.ReadBool()
		r.Data[i].BarCount, _ = b.ReadInt()
	}

	return r
}

////////////////////////////////////////////////////////////////////////////////
// SERIALIZERS
////////////////////////////////////////////////////////////////////////////////

func (b *HistoricalDataBroker) HistoricalDataItemToJSON(d *HistoricalDataItem) ([]byte, error) {
	return json.Marshal(struct {
		Date         string
		Symbol       string
		Exchange     string
		SecurityType string
		Currency     string
		Right        string
		Strike       float64
		Expiry       string
		Open         float64
		High         float64
		Low          float64
		Close        float64
		Volume       int64
		WAP          float64
		HasGaps      bool
		BarCount     int64
	}{
		Date:         d.Date,
		Symbol:       b.Contract.Symbol,
		Exchange:     b.Contract.Exchange,
		SecurityType: b.Contract.SecurityType,
		Currency:     b.Contract.Currency,
		Right:        b.Contract.Right,
		Strike:       b.Contract.Strike,
		Expiry:       b.Contract.Expiry,
		Open:         d.Open,
		High:         d.High,
		Low:          d.Low,
		Close:        d.Close,
		Volume:       d.Volume,
		WAP:          d.WAP,
		HasGaps:      d.HasGaps,
		BarCount:     d.BarCount,
	})
}

func (b *HistoricalDataBroker) HistoricalDataItemToCSV(d *HistoricalDataItem) string {
	return fmt.Sprintf(
		"%v,%s,%s,%s,%s,%s,%.2f,%s,%.2f,%.2f,%.2f,%.2f,%d,%.2f,%t,%d",
		d.Date,
		b.Contract.Symbol,
		b.Contract.Exchange,
		b.Contract.SecurityType,
		b.Contract.Currency,
		b.Contract.Right,
		b.Contract.Strike,
		b.Contract.Expiry,
		d.Open,
		d.High,
		d.Low,
		d.Close,
		d.Volume,
		d.WAP,
		d.HasGaps,
		d.BarCount,
	)
}
