package ib

type MarketData struct {
	Broker
}

func MarketDataBroker() MarketData {
	m := MarketData{Broker{}}
	m.Broker.Initialize()
	return m
}


