package ib

type MarketData struct {
	Broker
}

type TickPrice struct {
	Rid string
	TickType int64
	Price float64
	Size int64
	CanAutoExecute bool 
}

type TickSize struct {
	Rid string
	TickType int64
	Size int64
}

type TickOptionComputation struct {
	Rid string
	TickType int64
	ImpliedVol float64
	Delta float64
	OptionPrice float64
	PvDividend float64
	Gamma float64
	Vega float64
	Theta float64
	SpotPrice float64
}

type TickGeneric struct {
	Rid string
	TickType int64
	Value float64
}

type TickString struct {
	Rid string
	TickType int64
	Value string
}

type TickEFP struct {
	Rid string
	TickType int64
	BasisPoints float64
	FormattedBasisPoints string
	ImpliedFuturesPrice float64
	HoldDays int64
	FuturesExpiry string
	DividendImpact float64
	DividendsToExpiry float64
}

type MarketDataType struct {
	Rid string
	TickType int64
}

func MarketDataBroker() MarketData {
	m := MarketData{Broker{}}
	m.Broker.Initialize()
	return m
}

func (m *MarketData) CreateRequest(c Contract) {
	m.WriteInt(REQUEST.CODE.MARKET_DATA)
	m.WriteInt(REQUEST.VERSION.MARKET_DATA)
	m.WriteInt(m.NextReqId())
	m.WriteInt(c.ContractId)
	m.WriteString(c.Symbol)
	m.WriteString(c.SecurityType)
	m.WriteString(c.Expiry)
	m.WriteFloat(c.Strike)
	m.WriteString(c.Right)
	m.WriteString(c.Multiplier)
	m.WriteString(c.Exchange)
	m.WriteString(c.PrimaryExchange)
	m.WriteString(c.Currency)
	m.WriteString(c.LocalSymbol)
	m.WriteString(c.TradingClass)
	m.WriteBool(false) // underlying
	m.WriteString(c.GenericTickList)
	m.WriteBool(c.Snapshot)
}

func (m *MarketData) Listen() {
	for {
		b, err := m.ReadString()

		if err != nil {
			continue
		}
		
//		Log.Print("mkt", string(b))

		if b != RESPONSE.CODE.ERR_MSG {
			version, err := m.ReadString()
			c, err := m.ReadMsg(b, version)
			
			if err != nil {
				Log.Print("error", err.Error())
			} else {
				Log.Print("response", c)
			}
		}
	}
}

func (m *MarketData) ReadMsg(code, version string) (interface{}, error) {
	var err error

	switch code {
		case RESPONSE.CODE.TICK_PRICE:
			var p TickPrice

			p.Rid, err = m.ReadString()
			p.TickType, err = m.ReadInt()
			p.Price, err = m.ReadFloat()
			p.Size, err = m.ReadInt()
			p.CanAutoExecute, err = m.ReadBool()

			return p, err

		case RESPONSE.CODE.TICK_SIZE:
			var s TickSize

			s.Rid, err = m.ReadString()
			s.TickType, err = m.ReadInt()
			s.Size, err = m.ReadInt()

			return s, err

		case RESPONSE.CODE.TICK_OPTION_COMPUTATION:
			var o TickOptionComputation

			o.Rid, err = m.ReadString()
			o.TickType, err = m.ReadInt() 
			o.ImpliedVol, err = m.ReadFloat()
			o.Delta, err = m.ReadFloat()
			o.OptionPrice, err = m.ReadFloat()
			o.PvDividend, err = m.ReadFloat()
			o.Gamma, err = m.ReadFloat()
			o.Vega, err = m.ReadFloat()
			o.Theta, err = m.ReadFloat()
			o.SpotPrice, err = m.ReadFloat()

			return o, err

		case RESPONSE.CODE.TICK_GENERIC:
			var g TickGeneric

			g.Rid, err = m.ReadString()
			g.TickType, err = m.ReadInt()
			g.Value, err = m.ReadFloat()

			return g, err

		case RESPONSE.CODE.TICK_STRING:
			var s TickString

			s.Rid, err = m.ReadString()
			s.TickType, err = m.ReadInt()
			s.Value, err = m.ReadString()

			return s, err

		case RESPONSE.CODE.TICK_EFP:
			var e TickEFP

			e.Rid, err = m.ReadString()
			e.TickType, err = m.ReadInt()
			e.BasisPoints, err = m.ReadFloat()
			e.FormattedBasisPoints, err = m.ReadString()
			e.ImpliedFuturesPrice, err = m.ReadFloat()
			e.HoldDays, err = m.ReadInt()
			e.FuturesExpiry, err = m.ReadString()
			e.DividendImpact, err = m.ReadFloat()
			e.DividendsToExpiry, err = m.ReadFloat()

			return e, err

		case RESPONSE.CODE.TICK_SNAPSHOT_END:
		case RESPONSE.CODE.MARKET_DATA_TYPE:
			var d MarketDataType

			d.Rid, err  = m.ReadString()
			d.TickType, err = m.ReadInt()

			return d, err
	}

	return 0, err
}
