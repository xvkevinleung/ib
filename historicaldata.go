package ib

type HistoricalData struct {
	Broker
}

func HistoricalDataBroker() HistoricalData {
	h := HistoricalData{Broker{}}
	h.Broker.Initialize()
	return h
} 
