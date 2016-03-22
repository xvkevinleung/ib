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

type MarketDepthRequest struct {
	Rid      int64
	Contract Contract
	NumRows  int64
}

func init() {
	REQUEST_CODE["MarketDepth"] = 10
	REQUEST_VERSION["MarketDepth"] = 4
}

func (r *MarketDepthRequest) Send(b *MarketDepthBroker) {
	b.Contracts[r.Rid] = r.Contract
	b.WriteInt(REQUEST_CODE["MarketDepth"])
	b.WriteInt(REQUEST_VERSION["MarketDepth"])
	b.WriteInt(r.Rid)
	b.WriteInt(r.Contract.ContractId)
	b.WriteString(r.Contract.Symbol)
	b.WriteString(r.Contract.SecurityType)
	b.WriteString(r.Contract.Expiry)
	b.WriteFloat(r.Contract.Strike)
	b.WriteString(r.Contract.Right)
	b.WriteString(r.Contract.Multiplier)
	b.WriteString(r.Contract.Exchange)
	b.WriteString(r.Contract.Currency)
	b.WriteString(r.Contract.LocalSymbol)
	b.WriteString(r.Contract.TradingClass)
	b.WriteInt(r.NumRows)

	b.Broker.SendRequest()
}

type CancelMarketDepthRequest struct {
	Rid int64
}

func init() {
	REQUEST_CODE["CancelMarketDepth"] = 11
	REQUEST_VERSION["CancelMarketDepth"] = 1
}

func (r *CancelMarketDepthRequest) Send(b *MarketDepthBroker) {
	b.WriteInt(REQUEST_CODE["CancelMarketDepth"])
	b.WriteInt(REQUEST_VERSION["CancelMarketDepth"])
	b.WriteInt(r.Rid)

	b.Broker.SendRequest()

	delete(b.Contracts, r.Rid)
}

////////////////////////////////////////////////////////////////////////////////
// RESPONSES
////////////////////////////////////////////////////////////////////////////////

type MarketDepth struct {
	Rid          int64
	Symbol       string
	SecurityType string
	Exchange     string
	Currency     string
	Right        string
	Strike       float64
	Expiry       string
	Position     int64
	Operation    int64
	Side         int64
	Price        float64
	Size         int64
}

func init() {
	RESPONSE_CODE["MarketDepth"] = "12"
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

func init() {
	RESPONSE_CODE["MarketDepthLevelTwo"] = "13"
}

////////////////////////////////////////////////////////////////////////////////
// BROKER
////////////////////////////////////////////////////////////////////////////////

type MarketDepthBroker struct {
	Broker
	Contracts               map[int64]Contract
	MarketDepthChan         chan MarketDepth
	MarketDepthLevelTwoChan chan MarketDepthLevelTwo
}

func NewMarketDepthBroker() MarketDepthBroker {
	b := MarketDepthBroker{
		Broker{},
		make(map[int64]Contract),
		make(chan MarketDepth),
		make(chan MarketDepthLevelTwo),
	}

	return b
}

func (b *MarketDepthBroker) Listen() {
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
			case RESPONSE_CODE["MarketDepth"]:
				r := b.ReadMarketDepth(s, version)
				b.MarketDepthChan <- r
			case RESPONSE_CODE["MarketDepthLevelTwo"]:
				r := b.ReadMarketDepthLevelTwo(s, version)
				b.MarketDepthLevelTwoChan <- r
			default:
				b.ReadString()
			}
		}
	}
}

func (b *MarketDepthBroker) ReadMarketDepth(code, version string) MarketDepth {
	var r MarketDepth

	r.Rid, _ = b.ReadInt()

	c := b.Contracts[r.Rid]

	r.Symbol = c.Symbol
	r.SecurityType = c.SecurityType
	r.Exchange = c.Exchange
	r.Currency = c.Currency
	r.Right = c.Right
	r.Strike = c.Strike
	r.Expiry = c.Expiry
	r.Position, _ = b.ReadInt()
	r.Operation, _ = b.ReadInt()
	r.Side, _ = b.ReadInt()
	r.Price, _ = b.ReadFloat()
	r.Size, _ = b.ReadInt()

	return r
}

func (b *MarketDepthBroker) ReadMarketDepthLevelTwo(code, version string) MarketDepthLevelTwo {
	var r MarketDepthLevelTwo

	r.Rid, _ = b.ReadInt()
	r.Position, _ = b.ReadInt()
	r.MarketMaker, _ = b.ReadString()
	r.Operation, _ = b.ReadInt()
	r.Side, _ = b.ReadInt()
	r.Price, _ = b.ReadFloat()
	r.Size, _ = b.ReadInt()

	return r
}

////////////////////////////////////////////////////////////////////////////////
// SERIALIZERS
////////////////////////////////////////////////////////////////////////////////

func (b *MarketDepthBroker) SideToString(s int64) string {
	switch s {
	case 0:
		return "ASK"
	case 1:
		return "BID"
	default:
		return strconv.FormatInt(s, 10)
	}
}

func (b *MarketDepthBroker) OperationToString(o int64) string {
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

func (b *MarketDepthBroker) DepthToJSON(d *MarketDepth) ([]byte, error) {
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
		Position     int64
		Operation    string
		Side         string
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
		Position:     d.Position,
		Operation:    b.OperationToString(d.Operation),
		Side:         b.SideToString(d.Side),
		Price:        d.Price,
		Size:         d.Size,
	})
}

func (b *MarketDepthBroker) DepthToCSV(d *MarketDepth) string {
	c := b.Contracts[d.Rid]
	return fmt.Sprintf(
		"%d,%s,%s,%s,%s,%s,%s,%.2f,%s,,%d,%s,%s,%.2f,%d",
		d.Rid,
		strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		c.Symbol,
		c.SecurityType,
		c.Exchange,
		c.Currency,
		c.Right,
		c.Strike,
		c.Expiry,
		d.Position,
		b.OperationToString(d.Operation),
		b.SideToString(d.Side),
		d.Price,
		d.Size,
	)
}
