package ib

import (
	"encoding/json"
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

func (m *MarketDepthBroker) MarshalDepth(d *MarketDepth) ([]byte, error) {
	var s string
	switch d.Side {
	case 0:
		s = "ASK"
	case 1:
		s = "BID"
	default:
		s = strconv.FormatInt(d.Side, 10)
	}

	var o string
	switch d.Side {
	case 0:
		o = "INSERT"
	case 1:
		o = "UPDATE"
	case 2:
		o = "DELETE"
	default:
		o = strconv.FormatInt(d.Side, 10)
	}

	c := m.Contracts[d.Rid]
	return json.Marshal(struct {
		Time         string
		Symbol       string
		SecurityType string
		Exchange     string
		Currency     string
		Position     int64
		Operation    string
		Side         string
		Price        float64
		Size         int64
	}{
		Time:         strconv.FormatInt(time.Now().UnixNano(), 10),
		Symbol:       c.Symbol,
		SecurityType: c.SecurityType,
		Exchange:     c.Exchange,
		Currency:     c.Currency,
		Position:     d.Position,
		Operation:    o,
		Side:         s,
		Price:        d.Price,
		Size:         d.Size,
	})
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
