package ib

type TagValue struct {
	Tag string
	Value string
}

type ContractDetailsData struct {
	Rid string
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
		b, err := d.ReadString()

		if err != nil {
			continue
		}

		if b == RESPONSE.CODE.CONTRACT_DATA {
			version, err := d.ReadString()
			c, err := d.ReadMsg(version)
			
			if err != nil {
				Log.Print("error", err.Error())
			} else {
				Log.Print("response", c)
			}
		}
	}
}

func (d *ContractDetails) ReadMsg(version string) (ContractDetailsData, error) {
	var c ContractDetailsData
	var err error

	c.Rid, err = d.ReadString()
	c.Symbol, err = d.ReadString()
	c.SecurityType, err = d.ReadString()
	c.Expiry, err = d.ReadString()
	c.Strike, err = d.ReadFloat() 
	c.Right, err = d.ReadString()
	c.Exchange, err = d.ReadString()
	c.Currency, err = d.ReadString()
	c.LocalSymbol, err = d.ReadString()
	c.MarketName, err = d.ReadString()
	c.TradingClass, err = d.ReadString()
	c.ContractId, err = d.ReadInt() 
	c.MinTick, err = d.ReadInt() 
	c.Multiplier, err = d.ReadInt() 
	c.OrderTypes, err = d.ReadString()
	c.ValidExchanges, err = d.ReadString()
	c.PriceMagnifier, err = d.ReadInt() 
	c.UnderlyingContractId, err = d.ReadInt() 
	c.LongName, err = d.ReadString()
	c.PrimaryExchange, err = d.ReadString()
	c.ContractMonth, err = d.ReadString()
	c.Industry, err = d.ReadString()
	c.Category, err = d.ReadString()
	c.SubCategory, err = d.ReadString()
	c.TimeZoneId, err = d.ReadString()
	c.TradingHours, err = d.ReadString()
	c.LiquidHours, err = d.ReadString()
	c.EconValueRule, err = d.ReadString()
	c.EconValueMultiplier, err = d.ReadFloat() 
	c.SecIdListCount, err = d.ReadInt()

	for i := 0; i < int(c.SecIdListCount); i++ {
		var t, v string
		
		t, err = d.ReadString()
		v, err = d.ReadString()	
		tv := TagValue{t, v}
		c.SecIdList = append(c.SecIdList, tv)
	}

	return c, err
}

