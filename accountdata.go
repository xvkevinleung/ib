package ib

type AccountDataRequest struct {
	Subscribe   bool
	AccountCode string
}

type AccountValueData struct {
	Key      string
	Value    string
	Currency string
	Account  string
}

type PortfolioData struct {
	Key           string
	Contract      Contract
	Position      int64
	MarketPrice   float64
	MarketValue   float64
	AverageCost   float64
	UnrealizedPNL float64
	RealizedPNL   float64
	AccountName   string
}

type AccountTimeData struct {
	Time string
}

type AccountDataBroker struct {
	Broker
	AccountValueDataChan chan AccountValueData
	PortfolioDataChan    chan PortfolioData
	AccountTimeDataChan  chan AccountTimeData
}

func NewAccountDataBroker() AccountDataBroker {
	a := AccountDataBroker{
		Broker{},
		make(chan AccountValueData),
		make(chan PortfolioData),
		make(chan AccountTimeData),
	}

	return a
}

func (a *AccountDataBroker) SendRequest(rid int64, d AccountDataRequest) {
	a.WriteInt(REQUEST.CODE.ACCOUNT_VALUE)
	a.WriteInt(REQUEST.VERSION.ACCOUNT_VALUE)
	a.WriteInt(rid)
	a.WriteBool(d.Subscribe)
	a.WriteString(d.AccountCode)

	a.Broker.SendRequest()
}

type AccountDataAction func()

func (a *AccountDataBroker) Listen(f AccountDataAction) {
	go f()

	for {
		b, err := a.ReadString()

		if err != nil {
			continue
		}

		if b != RESPONSE.CODE.ERR_MSG {
			version, err := a.ReadString()

			if err != nil {
				continue
			}

			switch b {
			case RESPONSE.CODE.ACCOUNT_VALUE:
				a.ReadAccountValueData(b, version)
			case RESPONSE.CODE.PORTFOLIO_VALUE:
				a.ReadPortfolioData(b, version)
			case RESPONSE.CODE.ACCOUNT_UPDATE_TIME:
				a.ReadAccountUpdateTime(b, version)
			}
		}
	}
}

func (a *AccountDataBroker) ReadAccountValueData(code, version string) {
	var d AccountValueData

	d.Key, _ = a.ReadString()
	d.Value, _ = a.ReadString()
	d.Currency, _ = a.ReadString()
	d.Account, _ = a.ReadString()

	a.AccountValueDataChan <- d
}

func (a *AccountDataBroker) ReadPortfolioData(code, version string) {
	var d PortfolioData

	d.Contract.ContractId, _ = a.ReadInt()
	d.Contract.Symbol, _ = a.ReadString()
	d.Contract.SecurityType, _ = a.ReadString()
	d.Contract.Expiry, _ = a.ReadString()
	d.Contract.Strike, _ = a.ReadFloat()
	d.Contract.Right, _ = a.ReadString()
	d.Contract.Multiplier, _ = a.ReadString()
	d.Contract.PrimaryExchange, _ = a.ReadString()
	d.Contract.Currency, _ = a.ReadString()
	d.Contract.LocalSymbol, _ = a.ReadString()
	d.Contract.TradingClass, _ = a.ReadString()
	d.Position, _ = a.ReadInt()
	d.MarketPrice, _ = a.ReadFloat()
	d.MarketValue, _ = a.ReadFloat()
	d.AverageCost, _ = a.ReadFloat()
	d.UnrealizedPNL, _ = a.ReadFloat()
	d.AccountName, _ = a.ReadString()

	a.PortfolioDataChan <- d
}

func (a *AccountDataBroker) ReadAccountUpdateTime(code, version string) {
	var d AccountTimeData

	d.Time, _ = a.ReadString()

	a.AccountTimeDataChan <- d
}
