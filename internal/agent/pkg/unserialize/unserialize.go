package unserialize

import (
	"encoding/json"
	"errors"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
)

// Тип общия как для получателя Server так и для отправля Agent
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"c"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type UnSerializeInterface interface {
	SetData(data *[]byte) UnSerializeInterface
	GetData(unserializeData *Metrics) UnSerializeInterface
	Errors() error
}

type UnSerializer struct {
	sourceData      *[]byte
	unserializeData *Metrics
	cfg             *config.Config
	err             error
}

func NewUnSerializer(cfg *config.Config) UnSerializeInterface {
	return &UnSerializer{cfg: cfg}
}

// Устанавливаем значение перменной
// Возвращаем интерфейс для реализации цепочки вывзовов
func (s *UnSerializer) SetData(data *[]byte) UnSerializeInterface {

	if !json.Valid(*data) {
		s.err = errors.New("error at initial client B")
	}

	s.sourceData = data

	return s
}

func (s *UnSerializer) GetData(unserializeData *Metrics) UnSerializeInterface {
	s.unserializeData = unserializeData

	switch val := s.cfg.Serializer.SerType; val {
	case "json":
		s.JSON()
	default:
		s.err = errors.New("error not config Type Serialize")
	}

	return s
}

// Преобразуем
func (s *UnSerializer) JSON() {

	err := json.Unmarshal(*s.sourceData, &s.unserializeData)

	if err != nil {
		//panic(err.Error())
		s.err = err
	}

}

func (s *UnSerializer) Errors() error {
	return s.err
}
