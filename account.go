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
	b := AccountBroker{
		Broker{},
		make(chan AccountValue),
		make(chan Portfolio),
		make(chan AccountTime),
		make(chan AccountSummary),
	}

	return b
}

func (b *AccountBroker) Listen() {
	for {
		s, err := b.ReadString()

		if err != nil {
			continue
		}

		if s != RESPONSE.CODE.ERR_MSG {
			version, err := b.ReadString()

			if err != nil {
				continue
			}

			switch s {
			case RESPONSE_CODE["AccountValue"]:
				b.ReadAccountValue(s, version)
			case RESPONSE_CODE["Portfolio"]:
				b.ReadPortfolio(s, version)
			case RESPONSE_CODE["AccountUpdateTime"]:
				b.ReadAccountUpdateTime(s, version)
			case RESPONSE_CODE["AccountSummary"]:
				b.ReadAccountSummary(s, version)
			}
		}
	}
}

func (b *AccountBroker) ReadAccountValue(code, version string) {
	var r AccountValue

	r.Key, _ = b.ReadString()
	r.Value, _ = b.ReadString()
	r.Currency, _ = b.ReadString()
	r.Account, _ = b.ReadString()

	b.AccountValueChan <- r
}

func (b *AccountBroker) ReadPortfolio(code, version string) {
	var r Portfolio

	r.Contract.ContractId, _ = b.ReadInt()
	r.Contract.Symbol, _ = b.ReadString()
	r.Contract.SecurityType, _ = b.ReadString()
	r.Contract.Expiry, _ = b.ReadString()
	r.Contract.Strike, _ = b.ReadFloat()
	r.Contract.Right, _ = b.ReadString()
	r.Contract.Multiplier, _ = b.ReadString()
	r.Contract.PrimaryExchange, _ = b.ReadString()
	r.Contract.Currency, _ = b.ReadString()
	r.Contract.LocalSymbol, _ = b.ReadString()
	r.Contract.TradingClass, _ = b.ReadString()
	r.Position, _ = b.ReadInt()
	r.MarketPrice, _ = b.ReadFloat()
	r.MarketValue, _ = b.ReadFloat()
	r.AverageCost, _ = b.ReadFloat()
	r.UnrealizedPNL, _ = b.ReadFloat()
	r.AccountName, _ = b.ReadString()

	b.PortfolioChan <- r
}

func (b *AccountBroker) ReadAccountUpdateTime(code, version string) {
	var r AccountTime

	r.Time, _ = b.ReadString()

	b.AccountTimeChan <- r
}

func (b *AccountBroker) ReadAccountSummary(code, version string) {
	var r AccountSummary

	r.Rid, _ = b.ReadInt()
	r.Account, _ = b.ReadString()
	r.Tag, _ = b.ReadString()
	r.Value, _ = b.ReadString()
	r.Currency, _ = b.ReadString()

	b.AccountSummaryChan <- r
}

////////////////////////////////////////////////////////////////////////////////
// SERIALIZERS
////////////////////////////////////////////////////////////////////////////////
