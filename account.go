package ib

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// REQUESTS
////////////////////////////////////////////////////////////////////////////////

type AccountUpdatesRequest struct {
	Subscribe   bool `json:",string"`
	AccountCode string
}

func init() {
	REQUEST_CODE["AccountUpdates"] = 6
	REQUEST_VERSION["AccountUpdates"] = 2
}

func (r *AccountUpdatesRequest) Send(id int64, b *AccountBroker) {
	_ = id

	b.WriteInt(REQUEST_CODE["AccountUpdates"])
	b.WriteInt(REQUEST_VERSION["AccountUpdates"])
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

type AccountSummaryEnd struct {
	Rid int64
}

func init() {
	RESPONSE_CODE["AccountSummaryEnd"] = "64"
}

type AccountUpdateTime struct {
	Time string
}

func init() {
	RESPONSE_CODE["AccountUpdateTime"] = "8"
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

type AccountDownloadEnd struct {
	AccountName string
}

func init() {
	RESPONSE_CODE["AccountDownloadEnd"] = "54"
}

////////////////////////////////////////////////////////////////////////////////
// BROKER
////////////////////////////////////////////////////////////////////////////////

type AccountBroker struct {
	Broker
	AccountValueChan       chan AccountValue
	PortfolioChan          chan Portfolio
	AccountUpdateTimeChan  chan AccountUpdateTime
	AccountDownloadEndChan chan AccountDownloadEnd
	AccountSummaryChan     chan AccountSummary
	AccountSummaryEndChan  chan AccountSummaryEnd
}

func NewAccountBroker() AccountBroker {
	b := AccountBroker{
		Broker{},
		make(chan AccountValue),
		make(chan Portfolio),
		make(chan AccountUpdateTime),
		make(chan AccountDownloadEnd),
		make(chan AccountSummary),
		make(chan AccountSummaryEnd),
	}

	return b
}

func (b *AccountBroker) Listen() {
	for {
		s, err := b.ReadString()

		if err != nil {
			if err.Error() == "EOF" {
				log.Println(s, err)
				break
			}
			continue
		}

		if s != RESPONSE_CODE["ErrMsg"] {
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
			case RESPONSE_CODE["AccountDownloadEnd"]:
				b.ReadAccountDownloadEnd(s, version)
			case RESPONSE_CODE["AccountSummary"]:
				b.ReadAccountSummary(s, version)
			case RESPONSE_CODE["AccountSummaryEnd"]:
				b.ReadAccountSummaryEnd(s, version)
			default:
				b.ReadString()
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
	var r AccountUpdateTime

	r.Time, _ = b.ReadString()

	b.AccountUpdateTimeChan <- r
}

func (b *AccountBroker) ReadAccountDownloadEnd(code, version string) {
	var r AccountDownloadEnd

	r.AccountName, _ = b.ReadString()

	b.AccountDownloadEndChan <- r
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

func (b *AccountBroker) ReadAccountSummaryEnd(code, version string) {
	var r AccountSummaryEnd

	r.Rid, _ = b.ReadInt()

	b.AccountSummaryEndChan <- r
}

////////////////////////////////////////////////////////////////////////////////
// SERIALIZERS
////////////////////////////////////////////////////////////////////////////////

func (b *AccountBroker) AccountValueToJSON(d *AccountValue) ([]byte, error) {
	return json.Marshal(struct {
		Time     string
		Key      string
		Value    string
		Currency string
		Account  string
	}{
		Time:     strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		Key:      d.Key,
		Value:    d.Value,
		Currency: d.Currency,
		Account:  d.Account,
	})
}

func (b *AccountBroker) AccountSummaryToJSON(d *AccountSummary) ([]byte, error) {
	return json.Marshal(struct {
		Rid      int64
		Time     string
		Account  string
		Tag      string
		Value    string
		Currency string
	}{
		Rid:      d.Rid,
		Time:     strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		Account:  d.Account,
		Tag:      d.Tag,
		Value:    d.Value,
		Currency: d.Currency,
	})
}

func (b *AccountBroker) PortfolioToJSON(d *Portfolio) ([]byte, error) {
	return json.Marshal(struct {
		Time          string
		Key           string
		Contract      Contract
		Position      int64
		MarketPrice   float64
		MarketValue   float64
		AverageCost   float64
		UnrealizedPNL float64
		RealizedPNL   float64
		AccountName   string
	}{
		Time:          strconv.FormatInt(time.Now().UTC().Add(-5*time.Hour).UnixNano(), 10),
		Key:           d.Key,
		Contract:      d.Contract,
		Position:      d.Position,
		MarketPrice:   d.MarketPrice,
		MarketValue:   d.MarketValue,
		AverageCost:   d.AverageCost,
		UnrealizedPNL: d.UnrealizedPNL,
		RealizedPNL:   d.RealizedPNL,
		AccountName:   d.AccountName,
	})
}
