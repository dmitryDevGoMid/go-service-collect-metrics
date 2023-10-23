package models

// Структура хранения метрик
type MemStorage struct {
	Gauge   map[string]float64 //`json:"gauge,omitempty"`
	Counter map[string]int64   //  `json:"counter,omitempty"`
}
