package ib

type RealTimeBarsBroker struct {
	Broker
	RealTimeBarChan chan RealTimeBar
}

type RealTimeBarsRequest struct {
	Con  Contract
	Bar  int64
	Show string // what to show
	Rth  bool   // regular trading hours
	//	Opts string // use default "XYZ"
}

type RealTimeBar struct {
	Rid      int64
	Time     string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   int64
	WAP      float64
	BarCount int64
}

func NewRealTimeBarsBroker() RealTimeBarsBroker {
	r := RealTimeBarsBroker{Broker{}, make(chan RealTimeBar)}
	return r
}

func (r *RealTimeBarsBroker) SendRequest(rid int64, d RealTimeBarsRequest) {
	r.WriteInt(REQUEST.CODE.REALTIMEBARS)
	r.WriteInt(REQUEST.VERSION.REALTIMEBARS)
	r.WriteInt(rid)
	r.WriteInt(d.Con.ContractId)
	r.WriteString(d.Con.Symbol)
	r.WriteString(d.Con.SecurityType)
	r.WriteString(d.Con.Expiry)
	r.WriteFloat(d.Con.Strike)
	r.WriteString(d.Con.Right)
	r.WriteString(d.Con.Multiplier)
	r.WriteString(d.Con.Exchange)
	r.WriteString(d.Con.PrimaryExchange)
	r.WriteString(d.Con.Currency)
	r.WriteString(d.Con.LocalSymbol)
	r.WriteString(d.Con.TradingClass)
	r.WriteInt(d.Bar)
	r.WriteString(d.Show)
	r.WriteBool(d.Rth)

	r.Broker.SendRequest()
}

type RealTimeBarsAction func()

func (r *RealTimeBarsBroker) Listen(f RealTimeBarsAction) {
	go f()

	for {
		b, err := r.ReadString()

		if err != nil {
			continue
		}

		if b == RESPONSE.CODE.REALTIMEBARS {
			version, err := r.ReadString()

			if err != nil {
				continue
			}

			r.ReadRealTimeBar(version)
		}
	}
}

func (r *RealTimeBarsBroker) ReadRealTimeBar(version string) {
	var d RealTimeBar

	d.Rid, _ = r.ReadInt()
	d.Time, _ = r.ReadString()
	d.Open, _ = r.ReadFloat()
	d.High, _ = r.ReadFloat()
	d.Low, _ = r.ReadFloat()
	d.Close, _ = r.ReadFloat()
	d.Volume, _ = r.ReadInt()
	d.WAP, _ = r.ReadFloat()
	d.BarCount, _ = r.ReadInt()

	r.RealTimeBarChan <- d
}
