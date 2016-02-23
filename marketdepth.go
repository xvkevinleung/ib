package ib

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type MarketDepthBroker struct {
	Broker
	Contracts               map[int64]Contract
	MarketDepthChan         chan MarketDepth
	MarketDepthLevelTwoChan chan MarketDepthLevelTwo
}

type MarketDepthRequest struct {
	Con     Contract
	NumRows int64
}

type MarketDepth struct {
	Rid       int64
	Position  int64
	Operation int64
	Side      int64
	Price     float64
	Size      int64
}

func (m *MarketDepthBroker) SideToString(s int64) string {
	switch s {
	case 0:
		return "ASK"
	case 1:
		return "BID"
	default:
		return strconv.FormatInt(s, 10)
	}
}

func (m *MarketDepthBroker) OperationToString(o int64) string {
	switch o {
	case 0:
		return "INSERT"
	case 1:
		return "UPDATE"
	case 2:
		return "DELETE"
	default:
		return strconv.FormatInt(o, 10)
	}
}

func (m *MarketDepthBroker) DepthToJSON(d *MarketDepth) ([]byte, error) {
	c := m.Contracts[d.Rid]
	return json.Marshal(struct {
		Time         string
		Symbol       string
		SecurityType string
		Exchange     string
		Currency     string
		Right        string
		Strike       float64
		Expiry       string
		Position     int64
		Operation    string
		Side         string
		Price        float64
		Size         int64
	}{
		Time:         strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		Symbol:       c.Symbol,
		SecurityType: c.SecurityType,
		Exchange:     c.Exchange,
		Currency:     c.Currency,
		Right:        c.Right,
		Strike:       c.Strike,
		Expiry:       c.Expiry,
		Position:     d.Position,
		Operation:    m.OperationToString(d.Operation),
		Side:         m.SideToString(d.Side),
		Price:        d.Price,
		Size:         d.Size,
	})
}

func (m *MarketDepthBroker) DepthToCSV(d *MarketDepth) string {
	c := m.Contracts[d.Rid]
	return fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s,%.2f,%s,,%d,%s,%s,%.2f,%d",
		strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		c.Symbol,
		c.SecurityType,
		c.Exchange,
		c.Currency,
		c.Right,
		c.Strike,
		c.Expiry,
		d.Position,
		m.OperationToString(d.Operation),
		m.SideToString(d.Side),
		d.Price,
		d.Size,
	)
}

type MarketDepthLevelTwo struct {
	Rid         int64
	Position    int64
	MarketMaker string
	Operation   int64
	Side        int64
	Price       float64
	Size        int64
}

func NewMarketDepthBroker() MarketDepthBroker {
	m := MarketDepthBroker{
		Broker{},
		make(map[int64]Contract),
		make(chan MarketDepth),
		make(chan MarketDepthLevelTwo),
	}

	return m
}

func (m *MarketDepthBroker) SendRequest(rid int64, d MarketDepthRequest) {
	m.Contracts[rid] = d.Con
	m.WriteInt(REQUEST.CODE.MARKET_DEPTH)
	m.WriteInt(REQUEST.VERSION.MARKET_DEPTH)
	m.WriteInt(rid)
	m.WriteInt(d.Con.ContractId)
	m.WriteString(d.Con.Symbol)
	m.WriteString(d.Con.SecurityType)
	m.WriteString(d.Con.Expiry)
	m.WriteFloat(d.Con.Strike)
	m.WriteString(d.Con.Right)
	m.WriteString(d.Con.Multiplier)
	m.WriteString(d.Con.Exchange)
	m.WriteString(d.Con.Currency)
	m.WriteString(d.Con.LocalSymbol)
	m.WriteString(d.Con.TradingClass)
	m.WriteInt(d.NumRows)

	m.Broker.SendRequest()
}

func (m *MarketDepthBroker) Listen() {
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
			case RESPONSE.CODE.MARKET_DEPTH:
				m.ReadMarketDepth(b, version)
			case RESPONSE.CODE.MARKET_DEPTH_LEVEL_TWO:
				m.ReadMarketDepthLevelTwo(b, version)
			default:
				m.ReadString()
			}
		}
	}
}

func (m *MarketDepthBroker) ReadMarketDepth(code, version string) {
	var d MarketDepth

	d.Rid, _ = m.ReadInt()
	d.Position, _ = m.ReadInt()
	d.Operation, _ = m.ReadInt()
	d.Side, _ = m.ReadInt()
	d.Price, _ = m.ReadFloat()
	d.Size, _ = m.ReadInt()

	m.MarketDepthChan <- d
}

func (m *MarketDepthBroker) ReadMarketDepthLevelTwo(code, version string) {
	var d MarketDepthLevelTwo

	d.Rid, _ = m.ReadInt()
	d.Position, _ = m.ReadInt()
	d.MarketMaker, _ = m.ReadString()
	d.Operation, _ = m.ReadInt()
	d.Side, _ = m.ReadInt()
	d.Price, _ = m.ReadFloat()
	d.Size, _ = m.ReadInt()

	m.MarketDepthLevelTwoChan <- d
}
