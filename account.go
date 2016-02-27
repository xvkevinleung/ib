package ib

////////////////////////////////////////////////////////////////////////////////
// REQUESTS
////////////////////////////////////////////////////////////////////////////////

type AccountUpdatesRequest struct {
	Subscribe   bool
	AccountCode string
}

func init() {
	REQUEST_CODE["AccountUpdates"] = 63
	REQUEST_VERSION["AccountUpdates"] = 1
}

func (r *AccountUpdatesRequest) Send(id int64, b *AccountBroker) {
	b.WriteInt(REQUEST_CODE["AccountUpdates"])
	b.WriteInt(REQUEST_CODE["AccountUpdates"])
	b.WriteInt(id)
	b.WriteBool(r.Subscribe)
	b.WriteString(r.AccountCode)

	b.Broker.SendRequest()
}

type AccountSummaryRequest struct {
	Rid       int64
	GroupName string
	Tags      string
}

func init() {
	REQUEST_CODE["AccountSummary"] = 62
	REQUEST_VERSION["AccountSummary"] = 1
}

func (r *AccountSummaryRequest) Send(id int64, b *AccountBroker) {
	b.WriteInt(REQUEST_CODE["AccountSummary"])
	b.WriteInt(REQUEST_VERSION["AccountSummary"])
	b.WriteInt(id)
	b.WriteString(r.GroupName)
	b.WriteString(r.Tags)

	b.Broker.SendRequest()
}

////////////////////////////////////////////////////////////////////////////////
// RESPONSES
////////////////////////////////////////////////////////////////////////////////

type AccountValue struct {
	Key      string
	Value    string
	Currency string
	Account  string
}

func init() {
	RESPONSE_CODE["AccountValue"] = "6"
}

type AccountSummary struct {
	Rid      int64
	Account  string
	Tag      string
	Value    string
	Currency string
}

func init() {
	RESPONSE_CODE["AccountSummary"] = "63"
}

type AccountTime struct {
	Time string
}

func init() {
	RESPONSE_CODE["AccountTime"] = "8"
}

type Portfolio struct {
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

func init() {
	RESPONSE_CODE["Portfolio"] = "7"
}

////////////////////////////////////////////////////////////////////////////////
// BROKER
////////////////////////////////////////////////////////////////////////////////

type AccountBroker struct {
	Broker
	AccountValueChan   chan AccountValue
	PortfolioChan      chan Portfolio
	AccountTimeChan    chan AccountTime
	AccountSummaryChan chan AccountSummary
}

func NewAccountBroker() AccountBroker {
	a := AccountBroker{
		Broker{},
		make(chan AccountValue),
		make(chan Portfolio),
		make(chan AccountTime),
		make(chan AccountSummary),
	}

	return a
}

func (a *AccountBroker) Listen() {
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
			case RESPONSE_CODE["AccountValue"]:
				a.ReadAccountValue(b, version)
			case RESPONSE_CODE["Portfolio"]:
				a.ReadPortfolio(b, version)
			case RESPONSE_CODE["AccountUpdateTime"]:
				a.ReadAccountUpdateTime(b, version)
			case RESPONSE_CODE["AccountSummary"]:
				a.ReadAccountSummary(b, version)
			}
		}
	}
}

func (a *AccountBroker) ReadAccountValue(code, version string) {
	var d AccountValue

	d.Key, _ = a.ReadString()
	d.Value, _ = a.ReadString()
	d.Currency, _ = a.ReadString()
	d.Account, _ = a.ReadString()

	a.AccountValueChan <- d
}

func (a *AccountBroker) ReadPortfolio(code, version string) {
	var d Portfolio

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

	a.PortfolioChan <- d
}

func (a *AccountBroker) ReadAccountUpdateTime(code, version string) {
	var d AccountTime

	d.Time, _ = a.ReadString()

	a.AccountTimeChan <- d
}

func (a *AccountBroker) ReadAccountSummary(code, version string) {
	var d AccountSummary

	d.Rid, _ = a.ReadInt()
	d.Account, _ = a.ReadString()
	d.Tag, _ = a.ReadString()
	d.Value, _ = a.ReadString()
	d.Currency, _ = a.ReadString()

	a.AccountSummaryChan <- d
}

////////////////////////////////////////////////////////////////////////////////
// SERIALIZERS
////////////////////////////////////////////////////////////////////////////////
