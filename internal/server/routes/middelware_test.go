package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParseCIDRAndCheckIPOk(t *testing.T) {
	// Создаем конфигурацию с доверенным диапазоном IP-адресов
	cfg := &config.Config{
		TrustedSubnet: config.TrustedSubnet{
			TrustedSubnetCIDR: "192.168.0.0/24",
		},
	}

	// Создаем маршрутизатор Gin
	r := gin.Default()

	// Добавляем обработчик ParseCIDRAndCheckIP в цепочку обработчиков
	r.Use(routes.ParseCIDRAndCheckIP(cfg))

	// Создаем запрос с IP-адресом из доверенного диапазона
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "192.168.0.100")

	// Создаем ответ для запроса
	w := httptest.NewRecorder()

	// Выполняем запрос с помощью маршрутизатора Gin
	r.ServeHTTP(w, req)

	// Проверяем, что ответ имеет код состояния 404 Not found
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", w.Code)
	}

	// Создаем запрос с IP-адресом из не доверенного диапазона
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "192.168.0.100")

	// Создаем ответ для запроса
	w = httptest.NewRecorder()

	// Выполняем запрос с помощью маршрутизатора Gin
	r.ServeHTTP(w, req)

	//Ожидаем 404 так как адрес входит в доверительный диапазон, а так как такой роута нет то ответ сервера 404
	assert.Equal(t, 404, w.Code)
}

func TestParseCIDRAndCheckIPForbidden(t *testing.T) {
	// Создаем конфигурацию с доверенным диапазоном IP-адресов
	cfg := &config.Config{
		TrustedSubnet: config.TrustedSubnet{
			TrustedSubnetCIDR: "192.168.0.0/24",
		},
	}

	// Создаем маршрутизатор Gin
	r := gin.Default()

	// Добавляем обработчик ParseCIDRAndCheckIP в цепочку обработчиков
	r.Use(routes.ParseCIDRAndCheckIP(cfg))

	// Создаем запрос с IP-адресом из доверенного диапазона
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "192.168.0.100")

	// Создаем ответ для запроса
	w := httptest.NewRecorder()

	// Выполняем запрос с помощью маршрутизатора Gin
	r.ServeHTTP(w, req)

	// Проверяем, что ответ имеет код состояния 404 Not found
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", w.Code)
	}

	// Создаем запрос с IP-адресом из не доверенного диапазона
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "192.169.0.100")

	// Создаем ответ для запроса
	w = httptest.NewRecorder()

	// Выполняем запрос с помощью маршрутизатора Gin
	r.ServeHTTP(w, req)

	//Ожидаем 403 так как адрес не входит в доверительный диапазон
	assert.Equal(t, 403, w.Code)
}
