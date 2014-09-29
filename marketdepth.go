package ib

type MarketDepthBroker struct {
	Broker
	MarketDepthChan         chan MarketDepth
	MarketDepthLevelTwoChan chan MarketDepthLevelTwo
}

type MarketDepthRequest struct {
	Con     Contract
	NumRows int64
}

type MarketDepth struct {
	Rid       string
	Position  int64
	Operation int64
	Side      int64
	Price     float64
	Size      int64
}

type MarketDepthLevelTwo struct {
	Rid         string
	Position    int64
	MarketMaker string
	Operation   int64
	Side        int64
	Price       float64
	Size        int64
}

func NewMarketDepthBroker() MarketDepthBroker {
	m := MarketDepthBroker{
		Broker{},
		make(chan MarketDepth),
		make(chan MarketDepthLevelTwo),
	}

	return m
}

func (m *MarketDepthBroker) SendRequest(rid int64, d MarketDepthRequest) {
	m.WriteInt(REQUEST.CODE.MARKET_DEPTH)
	m.WriteInt(REQUEST.VERSION.MARKET_DEPTH)
	m.WriteInt(rid)
	m.WriteInt(d.Con.ContractId)
	m.WriteString(d.Con.Symbol)
	m.WriteString(d.Con.SecurityType)
	m.WriteString(d.Con.Expiry)
	m.WriteFloat(d.Con.Strike)
	m.WriteString(d.Con.Right)
	m.WriteString(d.Con.Multiplier)
	m.WriteString(d.Con.Exchange)
	m.WriteString(d.Con.Currency)
	m.WriteString(d.Con.LocalSymbol)
	m.WriteString(d.Con.TradingClass)
	m.WriteInt(d.NumRows)

	m.Broker.SendRequest()
}

type MarketDepthAction func()

func (m *MarketDepthBroker) Listen(f MarketDepthAction) {
	go f()

	for {
		b, err := m.ReadString()

		if err != nil {
			continue
		}

		if b != RESPONSE.CODE.ERR_MSG {
			version, err := m.ReadString()

			if err != nil {
				continue
			}

			switch b {
			case RESPONSE.CODE.MARKET_DEPTH:
				m.ReadMarketDepth(b, version)
			case RESPONSE.CODE.MARKET_DEPTH_LEVEL_TWO:
				m.ReadMarketDepthLevelTwo(b, version)
			default:
				m.ReadString()
			}
		}
	}
}

func (m *MarketDepthBroker) ReadMarketDepth(code, version string) {
	var d MarketDepth

	d.Rid, _ = m.ReadString()
	d.Position, _ = m.ReadInt()
	d.Operation, _ = m.ReadInt()
	d.Side, _ = m.ReadInt()
	d.Price, _ = m.ReadFloat()
	d.Size, _ = m.ReadInt()

	m.MarketDepthChan <- d
}

func (m *MarketDepthBroker) ReadMarketDepthLevelTwo(code, version string) {
	var d MarketDepthLevelTwo

	d.Rid, _ = m.ReadString()
	d.Position, _ = m.ReadInt()
	d.MarketMaker, _ = m.ReadString()
	d.Operation, _ = m.ReadInt()
	d.Side, _ = m.ReadInt()
	d.Price, _ = m.ReadFloat()
	d.Size, _ = m.ReadInt()

	m.MarketDepthLevelTwoChan <- d
}
