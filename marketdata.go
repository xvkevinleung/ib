package ib

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

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

type MarketDataRequest struct {
	Con             Contract
	GenericTickList string
	Snapshot        bool
}

type TickPrice struct {
	Rid            int64
	TickType       int64
	Price          float64
	Size           int64
	CanAutoExecute bool
}

func (m *MarketDataBroker) TickTypeToString(t int64) string {
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

func (m *MarketDataBroker) PriceToJSON(p *TickPrice) ([]byte, error) {
	c := m.Contracts[p.Rid]
	return json.Marshal(struct {
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
		Time:         strconv.FormatInt(time.Now().UnixNano(), 10),
		Symbol:       c.Symbol,
		SecurityType: c.SecurityType,
		Exchange:     c.Exchange,
		Currency:     c.Currency,
		Right:        c.Right,
		Strike:       c.Strike,
		Expiry:       c.Expiry,
		TickType:     m.TickTypeToString(p.TickType),
		Price:        p.Price,
		Size:         p.Size,
	})
}

func (m *MarketDataBroker) PriceToCSV(p *TickPrice) string {
	c := m.Contracts[p.Rid]
	return fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s,%.2f,%s,%s,%.2f,%d",
		strconv.FormatInt(time.Now().UnixNano(), 10),
		c.Symbol,
		c.SecurityType,
		c.Exchange,
		c.Currency,
		c.Right,
		c.Strike,
		c.Expiry,
		m.TickTypeToString(p.TickType),
		p.Price,
		p.Size,
	)
}

type TickSize struct {
	Rid      int64
	TickType int64
	Size     int64
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

type TickGeneric struct {
	Rid      int64
	TickType int64
	Value    float64
}

type TickString struct {
	Rid      int64
	TickType int64
	Value    string
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

type MarketDataType struct {
	Rid      int64
	TickType int64
}

func NewMarketDataBroker() MarketDataBroker {
	m := MarketDataBroker{
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

	return m
}

func (m *MarketDataBroker) SendRequest(rid int64, d MarketDataRequest) {
	m.Contracts[rid] = d.Con
	m.WriteInt(REQUEST.CODE.MARKET_DATA)
	m.WriteInt(REQUEST.VERSION.MARKET_DATA)
	m.WriteInt(rid)
	m.WriteInt(d.Con.ContractId)
	m.WriteString(d.Con.Symbol)
	m.WriteString(d.Con.SecurityType)
	m.WriteString(d.Con.Expiry)
	m.WriteFloat(d.Con.Strike)
	m.WriteString(d.Con.Right)
	m.WriteString(d.Con.Multiplier)
	m.WriteString(d.Con.Exchange)
	m.WriteString(d.Con.PrimaryExchange)
	m.WriteString(d.Con.Currency)
	m.WriteString(d.Con.LocalSymbol)
	m.WriteString(d.Con.TradingClass)
	m.WriteBool(false) // underlying
	m.WriteString(d.GenericTickList)
	m.WriteBool(d.Snapshot)

	m.Broker.SendRequest()
}

func (m *MarketDataBroker) Listen() {
	for {
		b, err := m.ReadString()

		if err != nil {
			continue
		}

		if b != RESPONSE.CODE.ERR_MSG {
			version, err := m.ReadString()

			if err != nil {
				continue
			}

			switch b {
			case RESPONSE.CODE.TICK_PRICE:
				m.ReadTickPrice(b, version)
			case RESPONSE.CODE.TICK_SIZE:
				m.ReadTickSize(b, version)
			case RESPONSE.CODE.TICK_OPTION_COMPUTATION:
				m.ReadTickOptComp(b, version)
			case RESPONSE.CODE.TICK_GENERIC:
				m.ReadTickGeneric(b, version)
			case RESPONSE.CODE.TICK_STRING:
				m.ReadTickString(b, version)
			case RESPONSE.CODE.TICK_EFP:
				m.ReadTickEFP(b, version)
			case RESPONSE.CODE.TICK_SNAPSHOT_END:
			case RESPONSE.CODE.MARKET_DATA_TYPE:
				m.ReadMarketDataType(b, version)
			default:
				m.ReadString()
			}
		}
	}
}

func (m *MarketDataBroker) ReadTickPrice(code, version string) {
	var p TickPrice

	p.Rid, _ = m.ReadInt()
	p.TickType, _ = m.ReadInt()
	p.Price, _ = m.ReadFloat()
	p.Size, _ = m.ReadInt()
	p.CanAutoExecute, _ = m.ReadBool()

	m.TickPriceChan <- p
}

func (m *MarketDataBroker) ReadTickSize(code, version string) {
	var s TickSize

	s.Rid, _ = m.ReadInt()
	s.TickType, _ = m.ReadInt()
	s.Size, _ = m.ReadInt()

	m.TickSizeChan <- s
}

func (m *MarketDataBroker) ReadTickOptComp(code, version string) {
	var o TickOptComp

	o.Rid, _ = m.ReadInt()
	o.TickType, _ = m.ReadInt()
	o.ImpliedVol, _ = m.ReadFloat()
	o.Delta, _ = m.ReadFloat()
	o.OptionPrice, _ = m.ReadFloat()
	o.PvDividend, _ = m.ReadFloat()
	o.Gamma, _ = m.ReadFloat()
	o.Vega, _ = m.ReadFloat()
	o.Theta, _ = m.ReadFloat()
	o.UndPrice, _ = m.ReadFloat()

	m.TickOptCompChan <- o
}

func (m *MarketDataBroker) ReadTickGeneric(code, version string) {
	var g TickGeneric

	g.Rid, _ = m.ReadInt()
	g.TickType, _ = m.ReadInt()
	g.Value, _ = m.ReadFloat()

	m.TickGenericChan <- g
}

func (m *MarketDataBroker) ReadTickString(code, version string) {
	var s TickString

	s.Rid, _ = m.ReadInt()
	s.TickType, _ = m.ReadInt()
	s.Value, _ = m.ReadString()

	m.TickStringChan <- s
}

func (m *MarketDataBroker) ReadTickEFP(code, version string) {
	var e TickEFP

	e.Rid, _ = m.ReadInt()
	e.TickType, _ = m.ReadInt()
	e.BasisPoints, _ = m.ReadFloat()
	e.FormattedBasisPoints, _ = m.ReadString()
	e.ImpliedFuturesPrice, _ = m.ReadFloat()
	e.HoldDays, _ = m.ReadInt()
	e.FuturesExpiry, _ = m.ReadString()
	e.DividendImpact, _ = m.ReadFloat()
	e.DividendsToExpiry, _ = m.ReadFloat()

	m.TickEFPChan <- e
}

func (m *MarketDataBroker) ReadMarketDataType(code, version string) {
	var d MarketDataType

	d.Rid, _ = m.ReadInt()
	d.TickType, _ = m.ReadInt()

	m.MarketDataTypeChan <- d
}
