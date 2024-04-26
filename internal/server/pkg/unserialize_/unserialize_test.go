package unserialize

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/models/example"
)

func BenchmarkUnserializer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var metrics Metrics
		var cfg *config.Config

		cfg, err := config.ParseConfig()

		if err != nil {
			fmt.Println("Config", err)
		}

		unserializeData := NewUnSerializer(cfg)

		body := []byte(`{"ID":"PollCount","Type":"counter"}`)

		unserializeError := unserializeData.SetData(&body).GetData(&metrics)

		if unserializeError.Errors() != nil {
			panic(unserializeError.Errors().Error())
		}
	}
}

func BenchmarkReadJson(b *testing.B) {
	// Открываем файл с JSON данными
	file, err := os.Open("example.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Читаем JSON данные из файла
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	unserializeData := NewUnSerializer(cfg)

	for i := 0; i < b.N; i++ {

		response := &example.Response{}

		unSerializeInterface := unserializeData.SetData(&data).GetDataExample(response)

		if unSerializeInterface.Errors() != nil {
			fmt.Println("Error:", err)
		}
	}
}
