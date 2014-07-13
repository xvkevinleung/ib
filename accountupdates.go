package ib

type AccountValueDataReq struct {
	Subscribe bool
	AccountCode string
}

type AccountValueData struct {
	Key string
	Value string
	Currency string
	Account string
}

type PortfolioData struct {
	Key string
	Contract Contract
	Position int64
	MarketPrice float64
	MarketValue float64
	AverageCost float64
	UnrealizedPNL float64
	RealizedPNL float64
	AccountName string
}

type AccountValueBroker struc {
	Broker
	AccountValueDataChan chan AccountValueData
	PortfolioDataChan chan PortfolioDAta
}

func NewAccountBroker() AccountValueBroker {
	a := AccountValueBroker{Broker{}, make(chan AccountValueData), make(chan PortfolioData)}
	a.Broker.Initialize()
	return a
}

func (a *AccountValueBroker) SendRequest(d AccountValueDataReq) {
	a.WriteInt(REQUEST.CODE.ACCOUNT_VALUE)
	a.WriteInt(REQUEST.VERSION.ACCOUNT_VALUE)
	a.WriteInt(a.NextReqId())
	a.WriteBool(d.Subscribe)
	a.WriteString(d.AccountCode)

	a.Broker.SendRequest()
}

type AccountValueAction func()

func (a *AccountValueBroker) Listen(f AccountValueDataReq) {
	go f()

	for {
		b, err := a.ReadString()

		if err != nil {
			continue
		}

		if b != RESPONSE.CODE.ERR_MSG {
			version, err := a.ReadString()

			if err != nil {
				Log.Print("error", err.Error())
			} else {
				switch b {
					case RESPONSE.CODE.ACCOUNT_VALUE:
						a.GetAccountValueData(b, version)
					case RESPONSE.CODE.PORTFOLIO_VALUE:
						a.GetPortfolioData(b, version)
					case RESPONSE.CODE.ACCOUNT_UPDATE_TIME:
						a.GetAccountUpdateTime(b, version)
				}
			}
		}
	}
}

func (a *AccountValueBroker) GetAccountValueData(code, version string) {
	var d AccountValueData
	var err error

	d.Rid, err = a.ReadString()
	d.Key, err = a.ReadString()
	d.Value, err = a.ReadString()
	d.Currency, err = a.ReadString()
	d.Account, err = a.ReadString()

	if err != nil {
		Log.Print("error", err.Error))
	} else {
		a.AccountValueDataChan <- d
	}
}

func (a *AccountValueBroker) GetPortfolioData(code, version string) {
	var d PortfolioData
	var err error

	d.Rid, err = a.ReadString()
	
	d.Contract.ContractId = a.ReadInt()
	d.Contract.Symbol = a.ReadString()
	d.Contract.SecurityType = a.ReadString()
	d.Contract.Expiry = a.ReadString()
	d.Contract.Strike = a.ReadFloat()
	d.Contract.Right = a.ReadString()
	d.Contract.Multiplier = a.ReadString()
	d.Contract.PrimaryExchange = a.ReadString()
	d.Contract.Currency = a.ReadString()
	d.Contract.LocalSymbol = a.ReadString()
	d.Contract.TradingClass = a.ReadString()

	
	if err != nil {
		Log.Print("error", err.Error))
	} else {
		a.PortfolioDataChan <- d
	}
}

func (a *AccountValueBroker) GetAccountUpdateTime(code, version string) {

}
