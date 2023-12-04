package models

var MetricsList string = `
Alloc   ,
BuckHashSys ,
Frees       ,
GCCPFraction,
GCSys       ,
HeapAlloc   ,
HeapIdle    ,
HeapInuse   ,
HeapObjects ,
HeapReleased,
HeapSys     ,
LastGC      ,
Lookups     ,
MCacheInuse ,
MCacheSys   ,
MSpanInuse  ,
MSpanSys    ,
Mallocs     ,
NextGC      ,
NumForceGC  ,
NumGC       ,
OtherSys    ,
PauseTotalNs,
StackInuse  ,
StackSys    ,
Sys         ,
TotalAlloc  ,
PollCount   ,
RandomValue`

type AllMetrics struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

// Структура хранения метрик
type MemStorage struct {
	Gauge   Gauge
	Counter Counter
}

type Gauge struct {
	//gopsutil
	TotalMemory    float64
	FreeMemory     float64
	CPUtilization1 float64

	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64
	RandomValue   float64
}

type Counter struct {
	PollCount int64
}
