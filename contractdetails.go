package ib

type ContractDetails struct {
	Broker
}

func ContractDetailsBroker() ContractDetails {
	c := ContractDetails{Broker{}}
	c.Broker.Initialize()
	return c
}

func (d *ContractDetails) CreateRequest(c Contract) {
	d.WriteInt(REQ_CONTRACT_DETAILS)
	d.WriteInt(7)
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
