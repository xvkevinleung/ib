package ib

import "log"

type Order struct {
	OrderID                       int64
	ClientID                      int64
	PermID                        int64
	Action                        string
	TotalQty                      int64
	OrderType                     string
	LimitPrice                    float64
	AuxPrice                      float64
	TIF                           string
	OCAGroup                      string
	Account                       string
	OpenClose                     string
	Origin                        int64
	OrderRef                      string
	Transmit                      bool
	ParentID                      int64
	BlockOrder                    bool
	SweepToFill                   bool
	DisplaySize                   int64
	TriggerMethod                 int64
	OutsideRTH                    bool
	Hidden                        bool
	DiscretionaryAmount           float64
	GoodAfterTime                 string
	GoodTillDate                  string
	FAGroup                       string
	FAMethod                      string
	FAPercentage                  string
	FAProfile                     string
	ShortSaleSlot                 int64
	DesignatedLocation            string
	ExemptCode                    int64
	OCAType                       int64
	Rule80A                       string
	SettlingFirm                  string
	AllOrNone                     bool
	MinQty                        int64
	PercentOffset                 float64
	ETradeOnly                    bool
	FirmQuoteOnly                 bool
	NBBOPriceCap                  float64
	AuctionStrategy               int64
	StartingPrice                 float64
	StockRefPrice                 float64
	Delta                         float64
	StockRangeLower               float64
	StockRangeUpper               float64
	OverridePercentageConstraints bool
	Volatility                    float64
	VolatilityType                int64
	DeltaNeutralOrderType         string
	DeltaNeutralAuxPrice          float64
	ContinuousUpdate              int64
	ReferencePriceType            int64
	TrailStopPrice                float64
	TrailingPercent               float64
	ScaleInitLevelSize            int64
	ScaleSubsLevelSize            int64
	ScalePriceIncrement           float64
	ScaleTable                    string
	ActiveStartTime               string
	ActiveStopTime                string
	HedgeType                     string
	OptOutSmartRouting            bool
	ClearingAccount               string
	ClearingIntent                string
	NotHeld                       bool
	AlgoStrategy                  string
	WhatIf                        bool
	OrderMiscOptions              string // []TagValue
}

////////////////////////////////////////////////////////////////////////////////
// REQUESTS
////////////////////////////////////////////////////////////////////////////////

type PlaceOrderRequest struct {
	Contract Contract
	Order    Order
}

func init() {
	REQUEST_CODE["PlaceOrder"] = 3
	REQUEST_VERSION["PlaceOrder"] = 42
}

func (r *PlaceOrderRequest) Send(id int64, b *OrderBroker) {
	b.WriteInt(REQUEST_CODE["PlaceOrder"])
	b.WriteInt(REQUEST_VERSION["PlaceOrder"])
	b.WriteInt(id)

	////////////////////
	// contract fields

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
	b.WriteBool(r.Contract.IncludeExpired)
	b.WriteString(r.Contract.SecIdType)
	b.WriteString(r.Contract.SecId)

	////////////////////
	// order fields

	b.WriteString(r.Order.Action)
	b.WriteInt(r.Order.TotalQty)
	b.WriteString(r.Order.OrderType)
	b.WriteFloat(r.Order.LimitPrice)
	b.WriteFloat(r.Order.AuxPrice)
	b.WriteString(r.Order.TIF)
	b.WriteString(r.Order.OCAGroup)
	b.WriteString(r.Order.Account)
	b.WriteString(r.Order.OpenClose)
	b.WriteInt(r.Order.Origin)
	b.WriteString(r.Order.OrderRef)
	b.WriteBool(r.Order.Transmit)
	b.WriteInt(r.Order.ParentID)
	b.WriteBool(r.Order.BlockOrder)
	b.WriteBool(r.Order.SweepToFill)
	b.WriteInt(r.Order.DisplaySize)
	b.WriteInt(r.Order.TriggerMethod)
	b.WriteBool(r.Order.OutsideRTH)
	b.WriteBool(r.Order.Hidden)

	// ignore combo legs by default
	// TODO implement combo legs

	// ignore smart combo routing options by default
	// TODO smart combo routing options

	// send deprecated shares allocation field
	b.WriteString("")

	b.WriteFloat(r.Order.DiscretionaryAmount)
	b.WriteString(r.Order.GoodAfterTime)
	b.WriteString(r.Order.GoodTillDate)
	b.WriteString(r.Order.FAGroup)
	b.WriteString(r.Order.FAMethod)
	b.WriteString(r.Order.FAPercentage)
	b.WriteString(r.Order.FAProfile)
	b.WriteInt(r.Order.ShortSaleSlot)
	b.WriteString(r.Order.DesignatedLocation)
	b.WriteInt(r.Order.ExemptCode)
	b.WriteInt(r.Order.OCAType)
	b.WriteString(r.Order.Rule80A)
	b.WriteString(r.Order.SettlingFirm)
	b.WriteBool(r.Order.AllOrNone)
	b.WriteInt(r.Order.MinQty)
	b.WriteFloat(r.Order.PercentOffset)
	b.WriteBool(r.Order.ETradeOnly)
	b.WriteBool(r.Order.FirmQuoteOnly)
	b.WriteFloat(r.Order.NBBOPriceCap)
	b.WriteInt(r.Order.AuctionStrategy)
	b.WriteFloat(r.Order.StartingPrice)
	b.WriteFloat(r.Order.StockRefPrice)
	b.WriteFloat(r.Order.Delta)
	b.WriteFloat(r.Order.StockRangeLower)
	b.WriteFloat(r.Order.StockRangeUpper)
	b.WriteBool(r.Order.OverridePercentageConstraints)
	b.WriteFloat(r.Order.Volatility)
	b.WriteInt(r.Order.VolatilityType)
	b.WriteString(r.Order.DeltaNeutralOrderType)
	b.WriteFloat(r.Order.DeltaNeutralAuxPrice)
	b.WriteInt(r.Order.ContinuousUpdate)
	b.WriteInt(r.Order.ReferencePriceType)
	b.WriteFloat(r.Order.TrailStopPrice)
	b.WriteFloat(r.Order.TrailingPercent)
	b.WriteInt(r.Order.ScaleInitLevelSize)
	b.WriteInt(r.Order.ScaleSubsLevelSize)
	b.WriteFloat(r.Order.ScalePriceIncrement)

	// ignore scale price fields by default
	// TODO implement scale price fields

	b.WriteString(r.Order.ScaleTable)
	b.WriteString(r.Order.ActiveStartTime)
	b.WriteString(r.Order.ActiveStopTime)
	b.WriteString(r.Order.HedgeType)

	// ignore hedge param by default
	// TODO implement hedge param

	b.WriteBool(r.Order.OptOutSmartRouting)
	b.WriteString(r.Order.ClearingAccount)
	b.WriteString(r.Order.ClearingIntent)
	b.WriteBool(r.Order.NotHeld)

	// ignore contract undercomp by default
	// TODO implement contract undercomp
	b.WriteBool(false)
	b.WriteString(r.Order.AlgoStrategy)

	// ignore algo params by default
	// TODO implement algo params

	b.WriteBool(r.Order.WhatIf)

	// ignore misc options by default
	// TODO implement misc options
	b.WriteString(r.Order.OrderMiscOptions)

	b.Broker.SendRequest()
}

type CancelOrderRequest struct {
	Rid int64
}

func init() {
	REQUEST_CODE["CancelOrder"] = 4
	REQUEST_VERSION["CancelOrder"] = 1
}

func (r *CancelOrderRequest) Send(id int64, b *OrderBroker) {
	_ = id

	b.WriteInt(REQUEST_CODE["CancelOrder"])
	b.WriteInt(REQUEST_VERSION["CancelOrder"])
	b.WriteInt(r.Rid)

	b.Broker.SendRequest()
}

type NextValidIdRequest struct {
	Num int64
}

func init() {
	REQUEST_CODE["NextRid"] = 8
	REQUEST_VERSION["NextRid"] = 1
}

func (r *NextValidIdRequest) Send(id int64, b *OrderBroker) {
	_ = id

	b.WriteInt(REQUEST_CODE["NextRid"])
	b.WriteInt(REQUEST_VERSION["NextRid"])
	b.WriteInt(r.Num)

	b.Broker.SendRequest()
}

////////////////////////////////////////////////////////////////////////////////
// RESPONSES
////////////////////////////////////////////////////////////////////////////////

type OrderStatus struct {
	Rid             int64
	Status          string
	Filled          int64
	Remaining       int64
	AvgFillPrice    float64
	PermId          int64
	ParentId        int64
	LastFilledPrice float64
	ClientId        int64
	WhyHeld         string
}

func init() {
	RESPONSE_CODE["OrderStatus"] = "3"
}

type OrderState struct {
	Commission         float64
	CommissionCurrency string
	EquityWithLoan     string
	InitMargin         string
	MaintMargin        string
	MaxCommission      float64
	MinCommission      float64
	Status             string
	WarningText        string
}

type OpenOrder struct {
	OrderId    int64
	Contract   Contract
	Order      Order
	OrderState OrderState
}

func init() {
	RESPONSE_CODE["OpenOrder"] = "5"
}

type NextValidId struct {
	OrderId int64
}

func init() {
	RESPONSE_CODE["NextValidId"] = "9"
}

////////////////////////////////////////////////////////////////////////////////
// BROKER
////////////////////////////////////////////////////////////////////////////////

type OrderBroker struct {
	Broker
	OrderStatusChan chan OrderStatus
	OpenOrderChan   chan OpenOrder
	NextValidIdChan chan NextValidId
}

func NewOrderBroker() OrderBroker {
	b := OrderBroker{
		Broker{},
		make(chan OrderStatus),
		make(chan OpenOrder),
		make(chan NextValidId),
	}

	return b
}

func (b *OrderBroker) NewOrder() Order {
	return Order{
		LimitPrice:           MAX_FLOAT,
		AuxPrice:             MAX_FLOAT,
		Transmit:             true,
		MinQty:               MAX_INT,
		PercentOffset:        MAX_FLOAT,
		TrailStopPrice:       MAX_FLOAT,
		TrailingPercent:      MAX_FLOAT,
		OpenClose:            "0",
		ExemptCode:           -1,
		ETradeOnly:           true,
		FirmQuoteOnly:        true,
		NBBOPriceCap:         MAX_FLOAT,
		StartingPrice:        MAX_FLOAT,
		StockRefPrice:        MAX_FLOAT,
		Delta:                MAX_FLOAT,
		StockRangeLower:      MAX_FLOAT,
		StockRangeUpper:      MAX_FLOAT,
		Volatility:           MAX_FLOAT,
		VolatilityType:       MAX_INT,
		DeltaNeutralAuxPrice: MAX_FLOAT,
		ReferencePriceType:   MAX_INT,
		ScaleInitLevelSize:   MAX_INT,
		ScaleSubsLevelSize:   MAX_INT,
		ScalePriceIncrement:  MAX_FLOAT,
	}
}

func (b *OrderBroker) Listen() {
	for {
		s, err := b.ReadString()

		if err != nil {
			if err.Error() == "EOF" {
				log.Println(s, err)
				break
			}
			continue
		}
		log.Println(s)
		if s != RESPONSE_CODE["ErrMsg"] {
			version, err := b.ReadString()

			if err != nil {
				continue
			}

			switch s {
			case RESPONSE_CODE["OrderStatus"]:
				r := b.ReadOrderStatus(s, version)
				b.OrderStatusChan <- r
				//      case RESPONSE_CODE["OpenOrder"]:
				//        r := b.ReadOpenOrder(s, version)
				//        b.OpenOrderChan <- r
			case RESPONSE_CODE["NextValidId"]:
				r := b.ReadNextValidId(s, version)
				b.NextValidIdChan <- r
      default:
        b.ReadString()
			}
		}
	}
}

func (b *OrderBroker) ReadOrderStatus(code, version string) OrderStatus {
	var r OrderStatus

	r.Rid, _ = b.ReadInt()
	r.Status, _ = b.ReadString()
	r.Filled, _ = b.ReadInt()
	r.Remaining, _ = b.ReadInt()
	r.AvgFillPrice, _ = b.ReadFloat()
	r.PermId, _ = b.ReadInt()
	r.ParentId, _ = b.ReadInt()
	r.LastFilledPrice, _ = b.ReadFloat()
	r.ClientId, _ = b.ReadInt()
	r.WhyHeld, _ = b.ReadString()

	return r
}

// func (b *OrderBroker) ReadOpenOrder(code, version string) OpenOrder {
//   var r OpenOrder
//
//   return r
// }

func (b *OrderBroker) ReadNextValidId(code, version string) NextValidId {
	var r NextValidId

	r.OrderId, _ = b.ReadInt()

	return r
}
