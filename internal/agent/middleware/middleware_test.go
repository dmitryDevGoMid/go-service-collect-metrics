package middleware_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/middleware"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pkg/asimencrypt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

// Тестовая структура
type RequestBody struct {
	Name string `json:"name"`
}

// Тестируем отправку запроса на сервер предварительно зашифровав и возвращаем расшифрованный ответ клиенту
func TestAssimEncryptBody(t *testing.T) {
	t.Run("encrypt body successfully", func(t *testing.T) {
		// Указываем конкретный порт
		listener, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
		defer listener.Close()

		// Создаем тестовый HTTP-сервер
		server := &httptest.Server{
			Listener: listener,
			Config: &http.Server{
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

					// Указываем тип запроса
					if r.Method != http.MethodPost {
						http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
						return
					}

					// Читаем тело запроса
					bodyBytes, err := ioutil.ReadAll(r.Body)
					if err != nil {
						http.Error(w, "error reading request body", http.StatusInternalServerError)
						return
					}

					defer r.Body.Close()
					// Разбираем тело запроса в структуру
					var reqBody RequestBody

					//Создаем экземпляр конфиг файла
					cfg := &config.Config{
						PathEncrypt: config.PathEncrypt{
							KeyEncryptEnbled: true,
							PathEncryptKey:   "keys_test/private.pem",
						},
					}

					//Создаем экземпляр обьект для шифрования, который принимает конфиг файл
					asme := asimencrypt.NewAsimEncrypt(cfg)
					errSetPrivateKey := asme.SetPrivateKey()
					if errSetPrivateKey != nil {
						log.Println("error set private key:", errSetPrivateKey)
					}
					// Расшифровываем тело запроса
					bodyString, err := asme.Decrypt(bodyBytes)
					if err != nil {
						fmt.Println(err)
					}

					err = json.Unmarshal([]byte(bodyString), &reqBody)
					if err != nil {
						http.Error(w, "error unmarshalling request body", http.StatusBadRequest)
						return
					}

					// Выводим тело запроса в консоль
					//fmt.Fprintf(w, "Request body: %+v\n", reqBody)

					//Отправляем ответ с сервера
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(bodyString))
				}),
			},
		}

		//Запускаем сервер
		server.Start()

		//По заврешению отсанавливаем сервер
		defer server.Close()

		//Создаем экземпляр конфиг файла
		cfg := &config.Config{
			PathEncrypt: config.PathEncrypt{
				KeyEncryptEnbled: true,
				PathEncryptKey:   "keys_test/public.pem",
			},
		}

		//Создаем экземпляр обьект для шифрования, который принимает конфиг файл
		asme := asimencrypt.NewAsimEncrypt(cfg)
		errSetPrivateKey := asme.SetPublicKey()
		if errSetPrivateKey != nil {
			log.Println("error set private key:", errSetPrivateKey)
		}

		// Создаем клиент resty
		client := resty.New()

		//Создаем middleware
		clientMiddleware := middleware.NewClientMiddleware(client, cfg, nil, asme)

		clientMiddleware.AssimEncryptBody()

		// Создаем тело запроса
		reqBody := &RequestBody{Name: "John"}
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		// Выполняем HTTP-запрос к тестовому серверу
		resp, err := client.R().
			SetBody(bodyBytes).
			Post(server.URL)

		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		fmt.Println(resp)

		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, string(bodyBytes), string(resp.Body()))

	})
}
