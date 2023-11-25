package serialize

import (
	"encoding/json"
	"errors"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
)

// Тип общия как для получателя Server так и для отправля Agent
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"mtype"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type SerializeInterface interface {
	SetData(data *Metrics) SerializeInterface
	GetData(serializeData *string) SerializeInterface
	Errors() error
}

type Serializer struct {
	sourceData    *Metrics
	serializeData *string
	cfg           *config.Config
	err           error
}

func NewSerializer(cfg *config.Config) SerializeInterface {
	return &Serializer{cfg: cfg}
}

// Устанавливаем значение перменной
// Возвращаем интерфейся для реализации цепочки вывзовов
func (s *Serializer) SetData(data *Metrics) SerializeInterface {

	s.sourceData = data

	return s
}

func (s *Serializer) GetData(serializeData *string) SerializeInterface {
	s.serializeData = serializeData

	switch val := s.cfg.Serializer.SerType; val {
	case "json":
		s.JSON()
	default:
		s.err = errors.New("error not config Type Serialize")
	}

	return s
}

// Преобразуем
func (s *Serializer) JSON() {
	json, _ := json.Marshal(s.sourceData)

	jsonString := string(json)

	*s.serializeData = jsonString
}

func (s *Serializer) Errors() error {
	return s.err
}
