package models

// Структура хранения метрик
type MemStorage struct {
	Gauge   map[string]float64 //`json:"gauge,omitempty"`
	Counter map[string]int64   //  `json:"counter,omitempty"`
}

// Удобная структура данных по метрикам для передачи и для gRPC
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"mtype"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
