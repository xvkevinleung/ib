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
	ActiveStartTime               string
	ActiveStopTime                string
	OCAGroup                      string
	OCAType                       int64
	OrderRef                      string
	Transmit                      bool
	ParentID                      int64
	BlockOrder                    bool
	SweepToFill                   bool
	DisplaySize                   int64
	TriggerMethod                 int64
	OutsideRTH                    bool
	Hidden                        bool
	GoodAfterTime                 string
	GoodTillDate                  string
	OverridePercentageConstraints bool
	Rule80A                       string
	AllOrNone                     bool
	MinQty                        int64
	PercentOffset                 float64
	TrailStopPrice                float64
	TrailingPercent               float64
	FAGroup                       string
	FAProfile                     string
	FAMethod                      string
	FAPercentage                  string
	OpenClose                     string
	Origin                        int64
	ShortSaleSlot                 int64
	DesignatedLocation            string
	ExemptCode                    int64
	DiscretionaryAmount           float64
	ETradeOnly                    int64
	FirmQuoteOnly                 bool
	NBBOPriceCap                  float64
	OptOutSmartRouting            bool
	AuctionStrategy               int64
	StartingPrice                 float64
	StockRefPrice                 float64
	Delta                         float64
	StockRangeLower               float64
	StockRangeUpper               float64
	Volatility                    float64
	VolatilityType                int64
	ContinuousUpdate              int64
	ReferencePriceType            int64
	DeltaNeutralOrderType         string
	DeltaNeutralAuxPrice          float64
	BasisPoints                   float64
	BasisPointsType               int64
	ScaleInitLevelSize            int64   // max
	ScaleSubsLevelSize            int64   // max
	ScalePriceIncrement           float64 // max
	ScalePriceAdjustValue         float64
	ScalePriceAdjustInterval      int64
	ScaleProfitOffset             float64
	ScaleAutoReset                bool
	ScaleInitPosition             int64
	ScaleInitFillQty              int64
	ScaleRandomPercent            bool
	ScaleTable                    string
	HedgeType                     string
	HedgeParam                    string
	Account                       string
	SettlingFirm                  string
	ClearingAccount               string
	ClearingIntent                string
	AlgoStrategy                  string
	AlgoParams                    string // []TagValue
	WhatIf                        bool
	NotHeld                       bool
	SmartComboRoutingParams       string // []TagValue
	OrderComboLegs                string // []TagValue
	OrderMiscOptions              string // []TagValue
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

	b.WriteString(r.Order.Action)
	b.WriteInt(r.Order.TotalQty)
	b.WriteString(r.Order.OrderType)
	b.WriteFloat(r.Order.LimitPrice)
	b.WriteFloat(r.Order.AuxPrice)
	b.WriteString(r.Order.TIF)
	b.WriteString(r.Order.ActiveStartTime)
	b.WriteString(r.Order.ActiveStopTime)
	b.WriteString(r.Order.OCAGroup)
	b.WriteInt(r.Order.OCAType)
	b.WriteString(r.Order.OrderRef)
	b.WriteBool(r.Order.Transmit)
	b.WriteInt(r.Order.ParentID)
	b.WriteBool(r.Order.BlockOrder)
	b.WriteBool(r.Order.SweepToFill)
	b.WriteInt(r.Order.DisplaySize)
	b.WriteInt(r.Order.TriggerMethod)
	b.WriteBool(r.Order.OutsideRTH)
	b.WriteBool(r.Order.Hidden)
	b.WriteString(r.Order.GoodAfterTime)
	b.WriteString(r.Order.GoodTillDate)
	b.WriteBool(r.Order.OverridePercentageConstraints)
	b.WriteString(r.Order.Rule80A)
	b.WriteBool(r.Order.AllOrNone)
	b.WriteInt(r.Order.MinQty)
	b.WriteFloat(r.Order.PercentOffset)
	b.WriteFloat(r.Order.TrailStopPrice)
	b.WriteFloat(r.Order.TrailingPercent)
	b.WriteString(r.Order.FAGroup)
	b.WriteString(r.Order.FAProfile)
	b.WriteString(r.Order.FAMethod)
	b.WriteString(r.Order.FAPercentage)
	b.WriteString(r.Order.OpenClose)
	b.WriteInt(r.Order.Origin)
	b.WriteInt(r.Order.ShortSaleSlot)
	b.WriteString(r.Order.DesignatedLocation)
	b.WriteInt(r.Order.ExemptCode)
	b.WriteFloat(r.Order.DiscretionaryAmount)
	b.WriteInt(r.Order.ETradeOnly)
	b.WriteBool(r.Order.FirmQuoteOnly)
	b.WriteFloat(r.Order.NBBOPriceCap)
	b.WriteBool(r.Order.OptOutSmartRouting)
	b.WriteInt(r.Order.AuctionStrategy)
	b.WriteFloat(r.Order.StartingPrice)
	b.WriteFloat(r.Order.StockRefPrice)
	b.WriteFloat(r.Order.Delta)
	b.WriteFloat(r.Order.StockRangeLower)
	b.WriteFloat(r.Order.StockRangeUpper)
	b.WriteFloat(r.Order.Volatility)
	b.WriteInt(r.Order.VolatilityType)
	b.WriteInt(r.Order.ContinuousUpdate)
	b.WriteInt(r.Order.ReferencePriceType)
	b.WriteString(r.Order.DeltaNeutralOrderType)
	b.WriteFloat(r.Order.DeltaNeutralAuxPrice)
	b.WriteFloat(r.Order.BasisPoints)
	b.WriteInt(r.Order.BasisPointsType)
	b.WriteInt(r.Order.ScaleInitLevelSize)
	b.WriteInt(r.Order.ScaleSubsLevelSize)
	b.WriteFloat(r.Order.ScalePriceIncrement)
	b.WriteFloat(r.Order.ScalePriceAdjustValue)
	b.WriteInt(r.Order.ScalePriceAdjustInterval)
	b.WriteFloat(r.Order.ScaleProfitOffset)
	b.WriteBool(r.Order.ScaleAutoReset)
	b.WriteInt(r.Order.ScaleInitPosition)
	b.WriteInt(r.Order.ScaleInitFillQty)
	b.WriteBool(r.Order.ScaleRandomPercent)
	b.WriteString(r.Order.ScaleTable)
	b.WriteString(r.Order.HedgeType)
	b.WriteString(r.Order.HedgeParam)
	b.WriteString(r.Order.Account)
	b.WriteString(r.Order.SettlingFirm)
	b.WriteString(r.Order.ClearingAccount)
	b.WriteString(r.Order.ClearingIntent)
	b.WriteString(r.Order.AlgoStrategy)
	b.WriteString(r.Order.AlgoParams)
	b.WriteBool(r.Order.WhatIf)
	b.WriteBool(r.Order.NotHeld)
	b.WriteString(r.Order.SmartComboRoutingParams)
	b.WriteString(r.Order.OrderComboLegs)
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
				b.ReadOrderStatus(s, version)
				//      case RESPONSE_CODE["OpenOrder"]:
				//        b.ReadOpenOrder(s, version)
			case RESPONSE_CODE["NextValidId"]:
				b.ReadNextValidId(s, version)
			}
		}
	}
}

func (b *OrderBroker) ReadOrderStatus(code, version string) {
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

	b.OrderStatusChan <- r
}

// func (b *OrderBroker) ReadOpenOrder(code, version string) {
//   var r OpenOrder
//
//   b.OpenOrderChan <- r
// }

func (b *OrderBroker) ReadNextValidId(code, version string) {
	var r NextValidId

	r.OrderId, _ = b.ReadInt()

	b.NextValidIdChan <- r
}
