package ib

type MarketData struct {
	Broker
	TickPriceChan chan TickPrice
	TickSizeChan chan TickSize
	TickOptCompChan chan TickOptComp
	TickGenericChan chan TickGeneric
	TickStringChan chan TickString
	TickEFPChan chan TickEFP
	MarketDataTypeChan chan MarketDataType
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

type TickOptComp struct {
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
	m := MarketData{
		Broker{},
		make(chan TickPrice),
		make(chan TickSize),
		make(chan TickOptComp),
		make(chan TickGeneric),
		make(chan TickString),
		make(chan TickEFP),
		make(chan MarketDataType),
	}
	m.Broker.Initialize()
	return m
}

func (m *MarketData) SendRequest(c Contract) {
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

	m.Broker.SendRequest()
}

type MarketDataAction func()

func (m *MarketData) Listen(f MarketDataAction) {
	go f()

	for {
		b, err := m.ReadString()

		if err != nil {
			continue
		}

		if b != RESPONSE.CODE.ERR_MSG {
			version, err := m.ReadString()
			
			if err != nil {
				Log.Print("error", err.Error())
			} else {
				switch b {
					case RESPONSE.CODE.TICK_PRICE:
						m.GetTickPrice(b, version)
					case RESPONSE.CODE.TICK_SIZE:
						m.GetTickSize(b,version)		
					case RESPONSE.CODE.TICK_OPTION_COMPUTATION:
						m.GetTickOptComp(b, version)
					case RESPONSE.CODE.TICK_GENERIC:
						m.GetTickGeneric(b, version)
					case RESPONSE.CODE.TICK_STRING:
						m.GetTickString(b, version)
					case RESPONSE.CODE.TICK_EFP:
						m.GetTickEFP(b, version)
					case RESPONSE.CODE.TICK_SNAPSHOT_END:
					case RESPONSE.CODE.MARKET_DATA_TYPE:
						m.GetMarketDataType(b, version)
					default:
						m.ReadString()
				}
			}
		}
	}
}

func (m *MarketData) GetTickPrice(code, version string) {
	var p TickPrice
	var err error

	p.Rid, err = m.ReadString()
	p.TickType, err = m.ReadInt()
	p.Price, err = m.ReadFloat()
	p.Size, err = m.ReadInt()
	p.CanAutoExecute, err = m.ReadBool()

	if err != nil {
		Log.Print("error", err.Error())
	} else {
		m.TickPriceChan <- p
	}
}

func (m *MarketData) GetTickSize(code, version string) {
	var s TickSize
	var err error

	s.Rid, err = m.ReadString()
	s.TickType, err = m.ReadInt()
	s.Size, err = m.ReadInt()
	                              
	if err != nil {
		Log.Print("error", err.Error())
	} else {
		m.TickSizeChan <- s
	}}

func (m *MarketData) GetTickOptComp(code, version string) {
	var o TickOptComp
	var err error

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

	if err != nil {
		Log.Print("error", err.Error())
	} else {
		m.TickOptCompChan <- o
	}
}

func (m *MarketData) GetTickGeneric(code, version string) {
	var g TickGeneric
	var err error

	g.Rid, err = m.ReadString()
	g.TickType, err = m.ReadInt()
	g.Value, err = m.ReadFloat()

	if err != nil {
		Log.Print("error", err.Error())
	} else {
		m.TickGenericChan <- g
	}
}

func (m *MarketData) GetTickString(code, version string) {
	var s TickString
	var err error

	s.Rid, err = m.ReadString()
	s.TickType, err = m.ReadInt()
	s.Value, err = m.ReadString()

	if err != nil {
		Log.Print("error", err.Error())
	} else {
		m.TickStringChan <- s
	}
}

func (m *MarketData) GetTickEFP(code, version string) {
	var e TickEFP
	var err error

	e.Rid, err = m.ReadString()
	e.TickType, err = m.ReadInt()
	e.BasisPoints, err = m.ReadFloat()
	e.FormattedBasisPoints, err = m.ReadString()
	e.ImpliedFuturesPrice, err = m.ReadFloat()
	e.HoldDays, err = m.ReadInt()
	e.FuturesExpiry, err = m.ReadString()
	e.DividendImpact, err = m.ReadFloat()
	e.DividendsToExpiry, err = m.ReadFloat()

	if err != nil {
		Log.Print("error", err.Error())
	} else {
		m.TickEFPChan <- e
	}
}

func (m *MarketData) GetMarketDataType(code, version string) {
	var d MarketDataType
	var err error

	d.Rid, err  = m.ReadString()
	d.TickType, err = m.ReadInt()

	if err != nil {
		Log.Print("error", err.Error())
	} else {
		m.MarketDataTypeChan <- d
	}
}

