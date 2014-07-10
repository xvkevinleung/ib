package ib

import (
	"bufio"
)

type TagValue struct {
	Tag string
	Value string
}

type ContractDetailsData struct {
	Symbol string
	SecurityType string
	Expiry string
	Strike float64
	Right string
	Exchange string
	Currency string
	LocalSymbol string
	MarketName string
	TradingClass string
	ContractId int64 
	MinTick int64
	Multiplier int64
	OrderTypes string
	ValidExchanges string
	PriceMagnifier int64
	UnderlyingContractId int64
	LongName string
	PrimaryExchange string
	ContractMonth string
	Industry string
	Category string
	SubCategory string
	TimeZoneId string
	TradingHours string
	LiquidHours string
	EconValueRule string
	EconValueMultiplier float64
	SecIdListCount int64
	SecIdList []TagValue
}

type ContractDetails struct {
	Broker
}

func ContractDetailsBroker() ContractDetails {
	c := ContractDetails{Broker{}}
	c.Broker.Initialize()
	return c
}

func (d *ContractDetails) CreateRequest(c Contract) {
	d.WriteInt(REQUEST.CODE.CONTRACT_DATA)
	d.WriteInt(REQUEST.VERSION.CONTRACT_DATA)
	d.WriteInt(d.NextReqId())
	d.WriteInt(c.ContractId)
	d.WriteString(c.Symbol)
	d.WriteString(c.SecurityType)
	d.WriteString(c.Expiry)
	d.WriteFloat(c.Strike)
	d.WriteString(c.Right)
	d.WriteString(c.Multiplier)
	d.WriteString(c.Exchange)
	d.WriteString(c.Currency)
	d.WriteString(c.LocalSymbol)
	d.WriteString(c.TradingClass)
	d.WriteBool(c.IncludeExpired)
	d.WriteString(c.SecIdType)
	d.WriteString(c.SecId)
}

func (d *ContractDetails) Listen() {
	for {
		b, err := d.InStream.ReadString('\000')

		if err != nil {
			continue
		}

		if b == RESPONSE.CODE.CONTRACT_DATA {
			version, err := d.InStream.ReadString('\000')
			rid, err := d.InStream.ReadString('\000')
//			d.ReadMsg(version)
		}
	}
}
/*
func (d *ContractDetails) ReadMsg(version string) {
	var c ContractDetailsData

	c.Symbol string
	c.SecurityType string
	c.Expiry string
	c.Strike float64
	c.Right string
	c.Exchange string
	c.Currency string
	c.LocalSymbol string
	c.MarketName string
	c.TradingClass string
	c.ContractId int64 
	c.MinTick int64
	c.Multiplier int64
	c.OrderTypes string
	c.ValidExchanges string
	c.PriceMagnifier int64
	c.UnderlyingContractId int64
	c.LongName string
	c.PrimaryExchange string
	c.ContractMonth string
	c.Industry string
	c.Category string
	c.SubCategory string
	c.TimeZoneId string
	c.TradingHours string
	c.LiquidHours string
	c.EconValueRule string
	c.EconValueMultiplier float64
	c.SecIdListCount int64
	c.SecIdList []TagValue
	
	Log.Print("contract", string(b))
}
*/
