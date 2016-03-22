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

type MarketDataRequest struct {
	Rid             int64
	Contract        Contract
	GenericTickList string
	Snapshot        bool
}

func init() {
	REQUEST_CODE["MarketData"] = 1
	REQUEST_VERSION["MarketData"] = 10
}

func (r *MarketDataRequest) Send(b *MarketDataBroker) {
	b.Contracts[r.Rid] = r.Contract
	b.WriteInt(REQUEST_CODE["MarketData"])
	b.WriteInt(REQUEST_VERSION["MarketData"])
	b.WriteInt(r.Rid)
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
	b.WriteBool(false) // underlying
	b.WriteString(r.GenericTickList)
	b.WriteBool(r.Snapshot)

	b.Broker.SendRequest()
}

type CancelMarketDataRequest struct {
	Rid int64
}

func init() {
	REQUEST_CODE["CancelMarketData"] = 2
	REQUEST_VERSION["CancelMarketData"] = 1
}

func (r *CancelMarketDataRequest) Send(b *MarketDataBroker) {
	b.WriteInt(REQUEST_CODE["CancelMarketData"])
	b.WriteInt(REQUEST_VERSION["CancelMarketData"])
	b.WriteInt(r.Rid)

	b.Broker.SendRequest()

	delete(b.Contracts, r.Rid)
}

////////////////////////////////////////////////////////////////////////////////
// RESPONSES
////////////////////////////////////////////////////////////////////////////////

type TickPrice struct {
	Rid            int64
	Symbol         string
	SecurityType   string
	Exchange       string
	Currency       string
	Right          string
	Strike         float64
	Expiry         string
	TickType       int64
	Price          float64
	Size           int64
	CanAutoExecute bool
}

func init() {
	RESPONSE_CODE["TickPrice"] = "1"
}

type TickSize struct {
	Rid      int64
	TickType int64
	Size     int64
}

func init() {
	RESPONSE_CODE["TickSize"] = "2"
}

type TickOptComp struct {
	Rid         int64
	TickType    int64
	ImpliedVol  float64
	Delta       float64
	OptionPrice float64
	PvDividend  float64
	Gamma       float64
	Vega        float64
	Theta       float64
	UndPrice    float64
}

func init() {
	RESPONSE_CODE["TicOptComp"] = "21"
}

type TickGeneric struct {
	Rid      int64
	TickType int64
	Value    float64
}

func init() {
	RESPONSE_CODE["TickGeneric"] = "45"
}

type TickString struct {
	Rid      int64
	TickType int64
	Value    string
}

func init() {
	RESPONSE_CODE["TickString"] = "46"
}

type TickEFP struct {
	Rid                  int64
	TickType             int64
	BasisPoints          float64
	FormattedBasisPoints string
	ImpliedFuturesPrice  float64
	HoldDays             int64
	FuturesExpiry        string
	DividendImpact       float64
	DividendsToExpiry    float64
}

func init() {
	RESPONSE_CODE["TickEFP"] = "47"
}

type MarketDataType struct {
	Rid      int64
	TickType int64
}

func init() {
	RESPONSE_CODE["MarketDataType"] = "58"
}

////////////////////////////////////////////////////////////////////////////////
// BROKER
////////////////////////////////////////////////////////////////////////////////

type MarketDataBroker struct {
	Broker
	Contracts          map[int64]Contract
	TickPriceChan      chan TickPrice
	TickSizeChan       chan TickSize
	TickOptCompChan    chan TickOptComp
	TickGenericChan    chan TickGeneric
	TickStringChan     chan TickString
	TickEFPChan        chan TickEFP
	MarketDataTypeChan chan MarketDataType
}

func NewMarketDataBroker() MarketDataBroker {
	b := MarketDataBroker{
		Broker{},
		make(map[int64]Contract),
		make(chan TickPrice),
		make(chan TickSize),
		make(chan TickOptComp),
		make(chan TickGeneric),
		make(chan TickString),
		make(chan TickEFP),
		make(chan MarketDataType),
	}

	return b
}

func (b *MarketDataBroker) Listen() {
	for {
		s, err := b.ReadString()

		if err != nil {
			continue
		}

		if s != RESPONSE_CODE["ErrMsg"] {
			version, err := b.ReadString()

			if err != nil {
				continue
			}

			switch s {
			case RESPONSE_CODE["TickPrice"]:
				r := b.ReadTickPrice(s, version)
				b.TickPriceChan <- r
			case RESPONSE_CODE["TickSize"]:
				r := b.ReadTickSize(s, version)
				b.TickSizeChan <- r
			case RESPONSE_CODE["TickOptComp"]:
				r := b.ReadTickOptComp(s, version)
				b.TickOptCompChan <- r
			case RESPONSE_CODE["TickGeneric"]:
				r := b.ReadTickGeneric(s, version)
				b.TickGenericChan <- r
			case RESPONSE_CODE["TickString"]:
				r := b.ReadTickString(s, version)
				b.TickStringChan <- r
			case RESPONSE_CODE["TickEFP"]:
				r := b.ReadTickEFP(s, version)
				b.TickEFPChan <- r
				//			case RESPONSE.CODE.TICK_SNAPSHOT_END:
			case RESPONSE_CODE["MarketDataType"]:
				r := b.ReadMarketDataType(s, version)
				b.MarketDataTypeChan <- r
			default:
				b.ReadString()
			}
		}
	}
}

func (b *MarketDataBroker) ReadTickPrice(code, version string) TickPrice {
	var r TickPrice

	r.Rid, _ = b.ReadInt()

	c := b.Contracts[r.Rid]

	r.Symbol = c.Symbol
	r.SecurityType = c.SecurityType
	r.Exchange = c.Exchange
	r.Currency = c.Currency
	r.Right = c.Right
	r.Strike = c.Strike
	r.Expiry = c.Expiry
	r.TickType, _ = b.ReadInt()
	r.Price, _ = b.ReadFloat()
	r.Size, _ = b.ReadInt()
	r.CanAutoExecute, _ = b.ReadBool()

	return r
}

func (b *MarketDataBroker) ReadTickSize(code, version string) TickSize {
	var r TickSize

	r.Rid, _ = b.ReadInt()
	r.TickType, _ = b.ReadInt()
	r.Size, _ = b.ReadInt()

	return r
}

func (b *MarketDataBroker) ReadTickOptComp(code, version string) TickOptComp {
	var r TickOptComp

	r.Rid, _ = b.ReadInt()
	r.TickType, _ = b.ReadInt()
	r.ImpliedVol, _ = b.ReadFloat()
	r.Delta, _ = b.ReadFloat()
	r.OptionPrice, _ = b.ReadFloat()
	r.PvDividend, _ = b.ReadFloat()
	r.Gamma, _ = b.ReadFloat()
	r.Vega, _ = b.ReadFloat()
	r.Theta, _ = b.ReadFloat()
	r.UndPrice, _ = b.ReadFloat()

	return r
}

func (b *MarketDataBroker) ReadTickGeneric(code, version string) TickGeneric {
	var r TickGeneric

	r.Rid, _ = b.ReadInt()
	r.TickType, _ = b.ReadInt()
	r.Value, _ = b.ReadFloat()

	return r
}

func (b *MarketDataBroker) ReadTickString(code, version string) TickString {
	var r TickString

	r.Rid, _ = b.ReadInt()
	r.TickType, _ = b.ReadInt()
	r.Value, _ = b.ReadString()

	return r
}

func (b *MarketDataBroker) ReadTickEFP(code, version string) TickEFP {
	var r TickEFP

	r.Rid, _ = b.ReadInt()
	r.TickType, _ = b.ReadInt()
	r.BasisPoints, _ = b.ReadFloat()
	r.FormattedBasisPoints, _ = b.ReadString()
	r.ImpliedFuturesPrice, _ = b.ReadFloat()
	r.HoldDays, _ = b.ReadInt()
	r.FuturesExpiry, _ = b.ReadString()
	r.DividendImpact, _ = b.ReadFloat()
	r.DividendsToExpiry, _ = b.ReadFloat()

	return r
}

func (b *MarketDataBroker) ReadMarketDataType(code, version string) MarketDataType {
	var r MarketDataType

	r.Rid, _ = b.ReadInt()
	r.TickType, _ = b.ReadInt()

	return r
}

////////////////////////////////////////////////////////////////////////////////
// SERIALIZERS
////////////////////////////////////////////////////////////////////////////////

func (b *MarketDataBroker) TickTypeToString(t int64) string {
	switch t {
	case 0:
		return "BID SIZE"
	case 1:
		return "BID"
	case 2:
		return "ASK"
	case 3:
		return "ASK_SIZE"
	case 4:
		return "LAST"
	case 5:
		return "LAST SIZE"
	case 6:
		return "HIGH"
	case 7:
		return "LOW"
	case 8:
		return "VOLUME"
	case 9:
		return "CLOSE"
	default:
		return strconv.FormatInt(t, 10)
	}
}

func (b *MarketDataBroker) PriceToJSON(d *TickPrice) ([]byte, error) {
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
		TickType     string
		Price        float64
		Size         int64
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
		TickType:     b.TickTypeToString(d.TickType),
		Price:        d.Price,
		Size:         d.Size,
	})
}

func (b *MarketDataBroker) PriceToCSV(d *TickPrice) string {
	c := b.Contracts[d.Rid]
	return fmt.Sprintf(
		"%d,%s,%s,%s,%s,%s,%s,%.2f,%s,%s,%.2f,%d",
		d.Rid,
		strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		c.Symbol,
		c.SecurityType,
		c.Exchange,
		c.Currency,
		c.Right,
		c.Strike,
		c.Expiry,
		b.TickTypeToString(d.TickType),
		d.Price,
		d.Size,
	)
}
