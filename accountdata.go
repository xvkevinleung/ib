package ib

type AccountDataRequest struct {
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

type AccountTimeData struct {
	Time string
}

type AccountDataBroker struct {
	Broker
	AccountValueDataChan chan AccountValueData
	PortfolioDataChan chan PortfolioData
	AccountTimeDataChan chan AccountTimeData
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

func (a *AccountDataBroker) SendRequest(d AccountDataRequest) {
	a.WriteInt(REQUEST.CODE.ACCOUNT_VALUE)
	a.WriteInt(REQUEST.VERSION.ACCOUNT_VALUE)
	a.WriteInt(a.NextReqId())
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
				Log.Print("error", err.Error())
			} else {
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
}

func (a *AccountDataBroker) ReadAccountValueData(code, version string) {
	var d AccountValueData
	var err error

	d.Key, err = a.ReadString()
	d.Value, err = a.ReadString()
	d.Currency, err = a.ReadString()
	d.Account, err = a.ReadString()

	if err != nil {
		Log.Print("error", err.Error)
	} else {
		a.AccountValueDataChan <- d
	}
}

func (a *AccountDataBroker) ReadPortfolioData(code, version string) {
	var d PortfolioData
	var err error

	d.Contract.ContractId, err = a.ReadInt()
	d.Contract.Symbol, err = a.ReadString()
	d.Contract.SecurityType, err = a.ReadString()
	d.Contract.Expiry, err = a.ReadString()
	d.Contract.Strike, err = a.ReadFloat()
	d.Contract.Right, err = a.ReadString()
	d.Contract.Multiplier, err = a.ReadString()
	d.Contract.PrimaryExchange, err = a.ReadString()
	d.Contract.Currency, err = a.ReadString()
	d.Contract.LocalSymbol, err = a.ReadString()
	d.Contract.TradingClass, err = a.ReadString()
	d.Position, err = a.ReadInt()
	d.MarketPrice, err = a.ReadFloat()
	d.MarketValue, err = a.ReadFloat()
	d.AverageCost, err = a.ReadFloat()
	d.UnrealizedPNL, err = a.ReadFloat()
	d.AccountName, err = a.ReadString()
	
	if err != nil {
		Log.Print("error", err.Error)
	} else {
		a.PortfolioDataChan <- d
	}
}

func (a *AccountDataBroker) ReadAccountUpdateTime(code, version string) {
	var d AccountTimeData
	var err error

	d.Time, err = a.ReadString()

	if err != nil {
		Log.Print("error", err.Error)
	} else {
		a.AccountTimeDataChan <- d
	}
}
