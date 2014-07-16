package ib

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

type HistoricalDataBroker struct {
	Broker
	HistoricalDataChan chan HistoricalData
}

func NewHistoricalDataBroker() HistoricalDataBroker {
	h := HistoricalDataBroker{Broker{}, make(chan HistoricalData)}
	return h
}

func (h *HistoricalDataBroker) SendRequest(d HistoricalDataRequest) {
	h.WriteInt(REQUEST.CODE.HISTORICAL_DATA)
	h.WriteInt(REQUEST.VERSION.HISTORICAL_DATA)
	h.WriteInt(h.NextReqId())
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

type HistoricalDataAction func()

func (h *HistoricalDataBroker) Listen(f HistoricalDataAction) {
	go f()

	for {
		b, err := h.ReadString()

		if err != nil {
			continue
		}

		if b == RESPONSE.CODE.HISTORICAL_DATA {
			version, err := h.ReadString()

			if err != nil {
				Log.Print("error", err.Error())
			} else {
				h.ReadHistoricalData(version)
			}
		}
	}
}

func (h *HistoricalDataBroker) ReadHistoricalData(version string) {
	var d HistoricalData
	var err error

	d.Rid, err = h.ReadString()
	d.Start, err = h.ReadString()
	d.End, err = h.ReadString()
	d.Count, err = h.ReadInt()

	d.Data = make([]HistoricalDataItem, d.Count)

	for i := range d.Data {
		d.Data[i].Date, err = h.ReadString()
		d.Data[i].Open, err = h.ReadFloat()
		d.Data[i].High, err = h.ReadFloat()
		d.Data[i].Low, err = h.ReadFloat()
		d.Data[i].Close, err = h.ReadFloat()
		d.Data[i].Volume, err = h.ReadInt()
		d.Data[i].WAP, err = h.ReadFloat()
		d.Data[i].HasGaps, err = h.ReadBool()
		d.Data[i].BarCount, err = h.ReadInt()

		if err != nil {
			h.HistoricalDataChan <- d
			break
		}
	}

	if err == nil {
		h.HistoricalDataChan <- d
	}
}
