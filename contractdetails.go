package ib

type TagValue struct {
	Tag   string
	Value string
}

type ContractDetailsData struct {
	Rid                  string
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

type ContractDetailsBroker struct {
	Broker
	DataChan chan ContractDetailsData
}

func NewContractDetailsBroker() ContractDetailsBroker {
	c := ContractDetailsBroker{Broker{}, make(chan ContractDetailsData)}
	c.Broker.Initialize()
	return c
}

func (d *ContractDetailsBroker) SendRequest(rid int64, c Contract) {
	d.WriteInt(REQUEST.CODE.CONTRACT_DATA)
	d.WriteInt(REQUEST.VERSION.CONTRACT_DATA)
	d.WriteInt(rid)
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

	d.Broker.SendRequest()
}

type ContractDetailsAction func()

func (d *ContractDetailsBroker) Listen(f ContractDetailsAction) {
	go f()

	for {
		b, err := d.ReadString()

		if err != nil {
			continue
		}

		if b == RESPONSE.CODE.CONTRACT_DATA {
			version, err := d.ReadString()

			if err != nil {
				continue
			} 
			
			d.ReadContractDetailsData(version)
		}
	}
}

func (d *ContractDetailsBroker) ReadContractDetailsData(version string) {
	var c ContractDetailsData

	c.Rid, _ = d.ReadString()
	c.Symbol, _ = d.ReadString()
	c.SecurityType, _ = d.ReadString()
	c.Expiry, _ = d.ReadString()
	c.Strike, _ = d.ReadFloat()
	c.Right, _ = d.ReadString()
	c.Exchange, _ = d.ReadString()
	c.Currency, _ = d.ReadString()
	c.LocalSymbol, _ = d.ReadString()
	c.MarketName, _ = d.ReadString()
	c.TradingClass, _ = d.ReadString()
	c.ContractId, _ = d.ReadInt()
	c.MinTick, _ = d.ReadInt()
	c.Multiplier, _ = d.ReadInt()
	c.OrderTypes, _ = d.ReadString()
	c.ValidExchanges, _ = d.ReadString()
	c.PriceMagnifier, _ = d.ReadInt()
	c.UnderlyingContractId, _ = d.ReadInt()
	c.LongName, _ = d.ReadString()
	c.PrimaryExchange, _ = d.ReadString()
	c.ContractMonth, _ = d.ReadString()
	c.Industry, _ = d.ReadString()
	c.Category, _ = d.ReadString()
	c.SubCategory, _ = d.ReadString()
	c.TimeZoneId, _ = d.ReadString()
	c.TradingHours, _ = d.ReadString()
	c.LiquidHours, _ = d.ReadString()
	c.EconValueRule, _ = d.ReadString()
	c.EconValueMultiplier, _ = d.ReadFloat()
	c.SecIdListCount, _ = d.ReadInt()

	for i := 0; i < int(c.SecIdListCount); i++ {
		var t, v string

		t, _ = d.ReadString()
		v, _ = d.ReadString()
		tv := TagValue{t, v}
		c.SecIdList = append(c.SecIdList, tv)
	}

	d.DataChan <- c
}
