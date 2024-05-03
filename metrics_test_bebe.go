package handlers_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/handlers"
	mocks "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository/mock_repository"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

// ExampleGetMetricsCounter демонстрирует использование функции GetMetricsCounter.
func ExampleMetricsHandlersType_GetMetricsCounter() {
	// Вычисляем квадрат числа 4.
	t := &testing.T{}

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tt := struct {
		name         string
		mType        string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		counterValue string
	}{
		name:       "get metrics by Post request with params into url 200",
		mType:      "counter",
		nameMetric: "TestCounter",
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
			expectedReturn := int64(700)

			mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
		},
		statusCode:   200,
		counterValue: "700",
	}
	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	//Обработчики
	handler := handlers.NewMetricsHandlers(nil, services, cfg)

	//handler := metricsHandlers{metricsRepository: services, cfg: cfg}

	//Init Point Handlers
	r := gin.Default()
	r.GET("/value/:type/:metric", func(context *gin.Context) {
		tt.mockBehavior(context, s, tt.nameMetric)
		context.Set("type", tt.mType)
		context.Set("metric", tt.nameMetric)
	}, handler.Value)

	//Create request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/value/%s/%s", tt.mType, tt.nameMetric), nil)

	r.ServeHTTP(w, req)

	// Print the result.
	fmt.Printf("Status code: %v\n", w.Code)
	fmt.Printf("Body: %v\n", w.Body.String())

	// Output:
	// Status code: 200
	// Body: 700
}
