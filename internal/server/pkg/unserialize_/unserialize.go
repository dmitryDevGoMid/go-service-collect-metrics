package unserialize

import (
	"encoding/json"
	"errors"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/models/example"
)

// Тип общия как для получателя Server так и для отправля Agent
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type UnSerializeInterface interface {
	SetData(data *[]byte) UnSerializeInterface
	GetData(unserializeData *Metrics) UnSerializeInterface
	GetDataExample(unserializeData *example.Response) UnSerializeInterface
	GetDataBatch(unserializeData *[]Metrics) UnSerializeInterface
	UnSerialize(body []byte) Metrics
	Errors() error
}

type UnSerializer struct {
	sourceData             *[]byte
	unserializeData        *Metrics
	unserializeDataBatch   *[]Metrics
	unserializeDataExample *example.Response
	cfg                    *config.Config
	err                    error
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

func (s *UnSerializer) GetDataExample(unserializeData *example.Response) UnSerializeInterface {
	s.unserializeDataExample = unserializeData

	switch val := s.cfg.Serializer.SerType; val {
	case "json":
		s.JSONExample()
	default:
		s.err = errors.New("error not config Type Serialize")
	}

	return s
}

func (s *UnSerializer) GetDataBatch(unserializeData *[]Metrics) UnSerializeInterface {
	s.unserializeDataBatch = unserializeData

	switch val := s.cfg.Serializer.SerType; val {
	case "json":
		s.JSONBatch()
	default:
		s.err = errors.New("error not config Type Serialize")
	}

	return s
}

// Преобразуем
func (s *UnSerializer) JSONBatch() {

	err := json.Unmarshal(*s.sourceData, &s.unserializeDataBatch)

	if err != nil {
		//panic(err.Error())
		s.err = err
	}

}

// Преобразуем
func (s *UnSerializer) JSON() {

	err := json.Unmarshal(*s.sourceData, &s.unserializeData)

	if err != nil {
		//panic(err.Error())
		s.err = err
	}

}

// Преобразуем используем обычный json
func (s *UnSerializer) JSONExample() {

	err := json.Unmarshal(*s.sourceData, &s.unserializeDataExample)

	if err != nil {
		//panic(err.Error())
		s.err = err
	}

}

func (s *UnSerializer) Errors() error {
	return s.err
}

func (s *UnSerializer) UnSerialize(body []byte) Metrics {
	var metrics Metrics

	unserializeError := s.SetData(&body).GetData(&metrics)

	if unserializeError.Errors() != nil {
		panic(unserializeError.Errors().Error())
	}

	return metrics
}
