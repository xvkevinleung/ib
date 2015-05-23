package ib

// NOTE: interactive brokers historical data is on central time

type HistoricalDataBroker struct {
	Broker
	HistoricalDataChan chan HistoricalData
}

type HistoricalDataRequest struct {
	Con   Contract
	End   string
	Bar   string
	Dur   string
	Rth   bool
	Show  string
	Datef int64
}

type HistoricalDataItem struct {
	Date     string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   int64
	WAP      float64
	HasGaps  bool
	BarCount int64
}

type HistoricalData struct {
	Rid   string
	Start string
	End   string
	Count int64
	Data  []HistoricalDataItem
}

func NewHistoricalDataBroker() HistoricalDataBroker {
	h := HistoricalDataBroker{Broker{}, make(chan HistoricalData)}
	return h
}

func (h *HistoricalDataBroker) SendRequest(rid int64, d HistoricalDataRequest) {
	h.WriteInt(REQUEST.CODE.HISTORICAL_DATA)
	h.WriteInt(REQUEST.VERSION.HISTORICAL_DATA)
	h.WriteInt(rid)
	h.WriteInt(d.Con.ContractId)
	h.WriteString(d.Con.Symbol)
	h.WriteString(d.Con.SecurityType)
	h.WriteString(d.Con.Expiry)
	h.WriteFloat(d.Con.Strike)
	h.WriteString(d.Con.Right)
	h.WriteString(d.Con.Multiplier)
	h.WriteString(d.Con.Exchange)
	h.WriteString(d.Con.PrimaryExchange)
	h.WriteString(d.Con.Currency)
	h.WriteString(d.Con.LocalSymbol)
	h.WriteString(d.Con.TradingClass)
	h.WriteInt(0) // include expired
	h.WriteString(d.End)
	h.WriteString(d.Bar)
	h.WriteString(d.Dur)
	h.WriteBool(d.Rth)
	h.WriteString(d.Show)
	h.WriteInt(d.Datef)

	h.Broker.SendRequest()
}

func (h *HistoricalDataBroker) Listen() {
	for {
		b, err := h.ReadString()

		if err != nil {
			continue
		}

		if b == RESPONSE.CODE.HISTORICAL_DATA {
			version, err := h.ReadString()

			if err != nil {
				continue
			}

			h.ReadHistoricalData(version)
		}
	}
}

func (h *HistoricalDataBroker) ReadHistoricalData(version string) {
	var d HistoricalData

	d.Rid, _ = h.ReadString()
	d.Start, _ = h.ReadString()
	d.End, _ = h.ReadString()
	d.Count, _ = h.ReadInt()

	d.Data = make([]HistoricalDataItem, d.Count)

	for i := range d.Data {
		d.Data[i].Date, _ = h.ReadString()
		d.Data[i].Open, _ = h.ReadFloat()
		d.Data[i].High, _ = h.ReadFloat()
		d.Data[i].Low, _ = h.ReadFloat()
		d.Data[i].Close, _ = h.ReadFloat()
		d.Data[i].Volume, _ = h.ReadInt()
		d.Data[i].WAP, _ = h.ReadFloat()
		d.Data[i].HasGaps, _ = h.ReadBool()
		d.Data[i].BarCount, _ = h.ReadInt()
	}

	h.HistoricalDataChan <- d
}
