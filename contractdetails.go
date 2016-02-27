package ib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// REQUESTS
////////////////////////////////////////////////////////////////////////////////

type ContractDetailsRequest struct {
	Contract Contract
}

func init() {
	REQUEST_CODE["ContractDetails"] = 9
	REQUEST_VERSION["ContractDeetails"] = 7
}

func (r *ContractDetailsRequest) Send(id int64, b *ContractDetailsBroker) {
	b.Contracts[id] = r.Contract
	b.WriteInt(REQUEST_CODE["ContractDetails"])
	b.WriteInt(REQUEST_VERSION["ContractDetails"])
	b.WriteInt(id)
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
	b.WriteBool(r.Contract.IncludeExpired)
	b.WriteString(r.Contract.SecIdType)
	b.WriteString(r.Contract.SecId)

	b.Broker.SendRequest()
}

////////////////////////////////////////////////////////////////////////////////
// RESPONSES
////////////////////////////////////////////////////////////////////////////////

type ContractDetails struct {
	Rid                  int64
	Symbol               string
	SecurityType         string
	Expiry               string
	Strike               float64
	Right                string
	Exchange             string
	Currency             string
	LocalSymbol          string
	MarketName           string
	TradingClass         string
	ContractId           int64
	MinTick              int64
	Multiplier           int64
	OrderTypes           string
	ValidExchanges       string
	PriceMagnifier       int64
	UnderlyingContractId int64
	LongName             string
	PrimaryExchange      string
	ContractMonth        string
	Industry             string
	Category             string
	SubCategory          string
	TimeZoneId           string
	TradingHours         string
	LiquidHours          string
	EconValueRule        string
	EconValueMultiplier  float64
	SecIdListCount       int64
	SecIdList            []TagValue
}

type TagValue struct {
	Tag   string
	Value string
}

func init() {
	RESPONSE_CODE["ContractDetails"] = "10"
}

////////////////////////////////////////////////////////////////////////////////
// BROKER
////////////////////////////////////////////////////////////////////////////////

type ContractDetailsBroker struct {
	Broker
	Contracts           map[int64]Contract
	ContractDetailsChan chan ContractDetails
}

func NewContractDetailsBroker() ContractDetailsBroker {
	c := ContractDetailsBroker{
		Broker{},
		make(map[int64]Contract),
		make(chan ContractDetails),
	}
	c.Broker.Initialize()
	return c
}

func (b *ContractDetailsBroker) Listen() {
	for {
		s, err := b.ReadString()

		if err != nil {
			continue
		}

		if s == RESPONSE_CODE["ContractDetails"] {
			version, err := b.ReadString()

			if err != nil {
				continue
			}

			b.ReadContractDetails(version)
		}
	}
}

func (b *ContractDetailsBroker) ReadContractDetails(version string) {
	var c ContractDetails

	c.Rid, _ = b.ReadInt()
	c.Symbol, _ = b.ReadString()
	c.SecurityType, _ = b.ReadString()
	c.Expiry, _ = b.ReadString()
	c.Strike, _ = b.ReadFloat()
	c.Right, _ = b.ReadString()
	c.Exchange, _ = b.ReadString()
	c.Currency, _ = b.ReadString()
	c.LocalSymbol, _ = b.ReadString()
	c.MarketName, _ = b.ReadString()
	c.TradingClass, _ = b.ReadString()
	c.ContractId, _ = b.ReadInt()
	c.MinTick, _ = b.ReadInt()
	c.Multiplier, _ = b.ReadInt()
	c.OrderTypes, _ = b.ReadString()
	c.ValidExchanges, _ = b.ReadString()
	c.PriceMagnifier, _ = b.ReadInt()
	c.UnderlyingContractId, _ = b.ReadInt()
	c.LongName, _ = b.ReadString()
	c.PrimaryExchange, _ = b.ReadString()
	c.ContractMonth, _ = b.ReadString()
	c.Industry, _ = b.ReadString()
	c.Category, _ = b.ReadString()
	c.SubCategory, _ = b.ReadString()
	c.TimeZoneId, _ = b.ReadString()
	c.TradingHours, _ = b.ReadString()
	c.LiquidHours, _ = b.ReadString()
	c.EconValueRule, _ = b.ReadString()
	c.EconValueMultiplier, _ = b.ReadFloat()
	c.SecIdListCount, _ = b.ReadInt()

	for i := 0; i < int(c.SecIdListCount); i++ {
		var t, v string

		t, _ = b.ReadString()
		v, _ = b.ReadString()
		tv := TagValue{t, v}
		c.SecIdList = append(c.SecIdList, tv)
	}

	b.ContractDetailsChan <- c
}

////////////////////////////////////////////////////////////////////////////////
// SERIALIZERS
////////////////////////////////////////////////////////////////////////////////

func (b *ContractDetailsBroker) ContractDetailsToJSON(d *ContractDetails) ([]byte, error) {
	r, err := json.Marshal(struct {
		Time                 string
		Symbol               string
		SecurityType         string
		Expiry               string
		Strike               float64
		Right                string
		Exchange             string
		Currency             string
		LocalSymbol          string
		MarketName           string
		TradingClass         string
		ContractId           string
		MinTick              int64
		Multiplier           int64
		OrderTypes           string
		ValidExchanges       string
		PriceMagnifier       int64
		UnderlyingContractId int64
		LongName             string
		PrimaryExchange      string
		ContractMonth        string
		Industry             string
		Category             string
		SubCategory          string
		TimeZoneId           string
		TradingHours         string
		LiquidHours          string
		EconValueRule        string
		EconValueMultiplier  float64
		//  	SecIdListCount       int64
		//  	SecIdList            []TagValue
	}{
		Time:                 strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		Symbol:               d.Symbol,
		SecurityType:         d.SecurityType,
		Expiry:               d.Expiry,
		Strike:               d.Strike,
		Right:                d.Right,
		Exchange:             d.Exchange,
		Currency:             d.Currency,
		LocalSymbol:          d.LocalSymbol,
		MarketName:           d.MarketName,
		TradingClass:         d.TradingClass,
		ContractId:           strconv.FormatInt(d.ContractId, 10),
		MinTick:              d.MinTick,
		Multiplier:           d.Multiplier,
		OrderTypes:           d.OrderTypes,
		ValidExchanges:       d.ValidExchanges,
		PriceMagnifier:       d.PriceMagnifier,
		UnderlyingContractId: d.UnderlyingContractId,
		LongName:             d.LongName,
		PrimaryExchange:      d.PrimaryExchange,
		ContractMonth:        d.ContractMonth,
		Industry:             d.Industry,
		Category:             d.Category,
		SubCategory:          d.SubCategory,
		TimeZoneId:           d.TimeZoneId,
		TradingHours:         d.TradingHours,
		LiquidHours:          d.LiquidHours,
		EconValueRule:        d.EconValueRule,
		EconValueMultiplier:  d.EconValueMultiplier,
		//  	SecIdListCount:       d.SecIdListCount,
		//  	SecIdList:            d.SecIdList,
	})

	return bytes.Replace(r, []byte("\\u0026"), []byte("&"), -1), err
}

func (b *ContractDetailsBroker) ContractDetailsToCSV(d *ContractDetails) string {
	return fmt.Sprintf(
		"%s,%s,%s,%s,%.2f,%s,%s,%s,%s,%s,%s,%d,%d,%d,%s,%s,%d,%d,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%.2f",
		strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		d.Symbol,
		d.SecurityType,
		d.Expiry,
		d.Strike,
		d.Right,
		d.Exchange,
		d.Currency,
		d.LocalSymbol,
		d.MarketName,
		d.TradingClass,
		d.ContractId,
		d.MinTick,
		d.Multiplier,
		d.OrderTypes,
		d.ValidExchanges,
		d.PriceMagnifier,
		d.UnderlyingContractId,
		d.LongName,
		d.PrimaryExchange,
		d.ContractMonth,
		d.Industry,
		d.Category,
		d.SubCategory,
		d.TimeZoneId,
		d.TradingHours,
		d.LiquidHours,
		d.EconValueRule,
		d.EconValueMultiplier,
		//    d.SecIdListCount,
		//    d.SecIdList,
	)
}
