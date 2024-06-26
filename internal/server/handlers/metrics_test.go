// internal/user/delivery_test.go
package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/handlers"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/unserialize"
	mocks "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository/mock_repository"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/storage"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPingDataBase(t *testing.T) {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		statusCode   int
		counterValue string
	}{
		{
			name: "get ping database connect to posgress 200",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository) {
				mocks.EXPECT().PingDatabase(ctx).Return(nil).AnyTimes()
			},
			statusCode:   200,
			counterValue: "Success to ping database",
		},
		{
			name: "get ping database connect to posgress 500",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository) {
				mocks.EXPECT().PingDatabase(ctx).Return(validator.ErrBadRequest).AnyTimes()
			},
			statusCode:   500,
			counterValue: "Failed to ping database",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Иницаиалзация тестирования
			c := gomock.NewController(t)
			defer c.Finish()

			s := mocks.NewMockMetricsRepository(c)

			services := s

			cfg, err := config.ParseConfig()

			if err != nil {
				fmt.Println("Config", err)
			}

			handler := handlers.NewMetricsHandlers(nil, services, cfg)

			//Init Point Handlers
			r := gin.Default()
			r.GET("/ping", func(context *gin.Context) {
				tt.mockBehavior(context, s)
			}, handler.Ping)

			//Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/ping", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.counterValue, w.Body.String())
		})
	}
}

func TestRequestUpdatesBatch(t *testing.T) {

	Delta := int64(500)
	Value := float64(500.456)

	metrics := []unserialize.Metrics{
		{
			ID:    "test",
			MType: "counter",
			Delta: &Delta,
		},
		{
			ID:    "test",
			MType: "gauge",
			Value: &Value,
		},
		{
			ID:    "test",
			MType: "counter",
			Delta: &Delta,
		},
		{
			ID:    "test",
			MType: "gauge",
			Value: &Value,
		},
	}
	metricsMarshal, err := json.Marshal(metrics)
	if err != nil {
		fmt.Println("Couldn't marshal metrics by test:", err)
	}
	stringMetricsMarshal := string(metricsMarshal)

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, metrics []unserialize.Metrics)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		statusCode   int
		counterValue string
	}{
		{
			name: "get all metrics 200",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, metrics []unserialize.Metrics) {

				expectedMetrics := metrics

				mocks.EXPECT().SaveMetricsBatch(ctx, expectedMetrics).Return(nil).AnyTimes()
			},
			statusCode:   200,
			counterValue: "Success get to Updates",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Иницаиалзация тестирования
			c := gomock.NewController(t)
			defer c.Finish()

			s := mocks.NewMockMetricsRepository(c)

			services := s

			cfg, err := config.ParseConfig()

			if err != nil {
				fmt.Println("Config", err)
			}

			handler := handlers.NewMetricsHandlers(nil, services, cfg)

			//Init Point Handlers
			r := gin.Default()
			r.POST("/updates", func(context *gin.Context) {
				tt.mockBehavior(context, s, metrics)
			}, handler.Updates)

			//Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/updates", bytes.NewBufferString(stringMetricsMarshal))

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.counterValue, w.Body.String())
		})
	}
}

func TestRequestListMetrics(t *testing.T) {

	metricsModel := storage.NewMemStorage()

	metricsModel.Gauge["metrics1"] = 543.4657
	metricsModel.Gauge["metrics2"] = 456.4657
	metricsModel.Gauge["metrics3"] = 432.4657

	metricsModel.Counter["counter"] = 500

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		statusCode   int
		counterValue string
	}{
		{
			name: "get all metrics 200",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository) {
				expectedReturn := metricsModel

				mocks.EXPECT().GetAllMetrics(ctx).Return(expectedReturn, nil).AnyTimes()
			},
			statusCode: 200,
			//counterValue: "<div>counter => 500 </div><div>metrics1 => 543.4657 </div><div>metrics2 => 456.4657 </div><div>metrics3 => 432.4657 </div>",
		},
		{
			name: "get all metrics bad request 400",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository) {
				mocks.EXPECT().GetAllMetrics(ctx).Return(nil, validator.ErrBadRequest).AnyTimes()
			},
			statusCode:   400,
			counterValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Иницаиалзация тестирования
			c := gomock.NewController(t)
			defer c.Finish()

			s := mocks.NewMockMetricsRepository(c)

			services := s

			cfg, err := config.ParseConfig()

			if err != nil {
				fmt.Println("Config", err)
			}

			handler := handlers.NewMetricsHandlers(nil, services, cfg)

			//Init Point Handlers
			r := gin.Default()
			r.GET("/", func(context *gin.Context) {
				tt.mockBehavior(context, s)
			}, handler.GetAllMetricsHTML)

			//Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			//assert.Equal(t, tt.counterValue, w.Body.String())
		})
	}
}

func TestRequestGetMetricsCounter(t *testing.T) {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tests := []struct {
		name         string
		mType        string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		counterValue string
	}{
		{
			name:       "get metrics by Post request with params into url 200",
			mType:      "counter",
			nameMetric: "TestCounter",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := int64(700)

				mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
			},
			statusCode:   200,
			counterValue: "700",
		},
		{
			name:       "get metrics by Post request with params into url 404",
			mType:      "counter",
			nameMetric: "TestCounter",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := int64(0)

				mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
			},
			statusCode:   404,
			counterValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Иницаиалзация тестирования
			c := gomock.NewController(t)
			defer c.Finish()

			s := mocks.NewMockMetricsRepository(c)

			services := s

			cfg, err := config.ParseConfig()

			if err != nil {
				fmt.Println("Config", err)
			}

			handler := handlers.NewMetricsHandlers(nil, services, cfg)

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

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.counterValue, w.Body.String())
		})
	}
}

func TestRequestGetMetricsGauge(t *testing.T) {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tests := []struct {
		name         string
		mType        string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		counterValue string
	}{
		{
			name:       "get metrics by Post request with params into url 200",
			mType:      "gauge",
			nameMetric: "TestCounter",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := float64(700.123)

				mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
			},
			statusCode:   200,
			counterValue: "700.123",
		},
		{
			name:       "get metrics by Post request with params into url 404",
			mType:      "gauge",
			nameMetric: "TestCounter",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := float64(0)

				mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
			},
			statusCode:   404,
			counterValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Иницаиалзация тестирования
			c := gomock.NewController(t)
			defer c.Finish()

			s := mocks.NewMockMetricsRepository(c)

			services := s

			cfg, err := config.ParseConfig()

			if err != nil {
				fmt.Println("Config", err)
			}

			handler := handlers.NewMetricsHandlers(nil, services, cfg)

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

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.counterValue, w.Body.String())
		})
	}
}

func TestRequestPostMetricsGauge(t *testing.T) {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tests := []struct {
		name         string
		mType        string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		counterValue string
	}{
		{
			name:       "get metrics by Post request with params into url 200",
			mType:      "gauge",
			nameMetric: "TestCounter",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := float64(700.123)

				mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
			},
			statusCode:   200,
			counterValue: "700.123",
		},
		{
			name:       "get metrics by Post request with params into url 404",
			mType:      "gauge",
			nameMetric: "TestCounter",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := float64(0)

				mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
			},
			statusCode:   404,
			counterValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Иницаиалзация тестирования
			c := gomock.NewController(t)
			defer c.Finish()

			s := mocks.NewMockMetricsRepository(c)

			services := s

			cfg, err := config.ParseConfig()

			if err != nil {
				fmt.Println("Config", err)
			}

			handler := handlers.NewMetricsHandlers(nil, services, cfg)

			//Init Point Handlers
			r := gin.Default()
			r.POST("/value/:type/:metric", func(context *gin.Context) {
				tt.mockBehavior(context, s, tt.nameMetric)
				context.Set("type", tt.mType)
				context.Set("metric", tt.nameMetric)
			}, handler.Value)

			//Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/value/%s/%s", tt.mType, tt.nameMetric), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.counterValue, w.Body.String())
		})
	}
}

func TestRequestPostMetricsCounter(t *testing.T) {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tests := []struct {
		name         string
		mType        string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		counterValue string
	}{
		{
			name:       "get metrics by Post request with params into url 200",
			mType:      "counter",
			nameMetric: "TestCounter",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := int64(700)

				mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
			},
			statusCode:   200,
			counterValue: "700",
		},
		{
			name:       "get metrics by Post request with params into url 404",
			mType:      "counter",
			nameMetric: "TestCounter",
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := int64(0)

				mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
			},
			statusCode:   404,
			counterValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Иницаиалзация тестирования
			c := gomock.NewController(t)
			defer c.Finish()

			s := mocks.NewMockMetricsRepository(c)

			services := s

			cfg, err := config.ParseConfig()

			if err != nil {
				fmt.Println("Config", err)
			}

			handler := handlers.NewMetricsHandlers(nil, services, cfg)

			//Init Point Handlers
			r := gin.Default()
			r.POST("/value/:type/:metric", func(context *gin.Context) {
				tt.mockBehavior(context, s, tt.nameMetric)
				context.Set("type", tt.mType)
				context.Set("metric", tt.nameMetric)
			}, handler.Value)

			//Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/value/%s/%s", tt.mType, tt.nameMetric), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.counterValue, w.Body.String())
		})
	}
}

func TestGetMetricsCounter(t *testing.T) {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tests := []struct {
		name         string
		body         string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name: "get Metrics by name",
			body: fmt.Sprintf(`{"id":"%s","type":"%s"}`, "TestCounter", "counter"),
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := int64(700)

				mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
			},
			statusCode:   200,
			responseBody: fmt.Sprintf(`{"id":"%s","type":"%s","delta":700}`, "TestCounter", "counter"),
		},
		{
			name: "get Metrics by name",
			body: fmt.Sprintf(`{"id":"%s","type":"%s"}`, "TestCounter", "counter"),
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := int64(0)
				mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
			},
			statusCode:   404,
			responseBody: "",
		},
		{
			name: "get Metrics by name",
			body: fmt.Sprintf(`{"id":"%s","type":"%s"}`, "TestCounter", "unter"),
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := int64(0)
				mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
			},
			statusCode:   400,
			responseBody: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Иницаиалзация тестирования
			c := gomock.NewController(t)
			defer c.Finish()

			s := mocks.NewMockMetricsRepository(c)

			services := s

			cfg, err := config.ParseConfig()

			if err != nil {
				fmt.Println("Config", err)
			}

			handler := handlers.NewMetricsHandlers(nil, services, cfg)

			//Init Point Handlers
			r := gin.Default()
			r.POST("/value", func(context *gin.Context) {
				tt.mockBehavior(context, s, "TestCounter")
				//Прописыаем заголовки
				context.Request.Header.Set("Content-Type", "application/json")
				context.Request.Header.Set("Accept", "application/json")
			}, handler.ValuePostJSON)

			//Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/value", bytes.NewBufferString(tt.body))

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestGetMetricsGauge(t *testing.T) {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tests := []struct {
		name         string
		body         string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name: "get Metrics by name",
			body: fmt.Sprintf(`{"id":"%s","type":"%s"}`, "TestCounter", "gauge"),
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := float64(700.123)

				mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
			},
			statusCode:   200,
			responseBody: fmt.Sprintf(`{"id":"%s","type":"%s","value":700.123}`, "TestCounter", "gauge"),
		},
		{
			name: "get Metrics by name",
			body: fmt.Sprintf(`{"id":"%s","type":"%s"}`, "TestCounter", "gauge"),
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := float64(0)
				mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
			},
			statusCode:   404,
			responseBody: "",
		},
		{
			name: "get Metrics by name",
			body: fmt.Sprintf(`{"id":"%s","type":"%s"}`, "TestCounter", "gaug"),
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
				expectedReturn := float64(0)
				mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
			},
			statusCode:   400,
			responseBody: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Иницаиалзация тестирования
			c := gomock.NewController(t)
			defer c.Finish()

			s := mocks.NewMockMetricsRepository(c)

			services := s

			cfg, err := config.ParseConfig()

			if err != nil {
				fmt.Println("Config", err)
			}

			handler := handlers.NewMetricsHandlers(nil, services, cfg)

			//Init Point Handlers
			r := gin.Default()
			r.POST("/value", func(context *gin.Context) {
				tt.mockBehavior(context, s, "TestCounter")
				//Прописыаем заголовки
				context.Request.Header.Set("Content-Type", "application/json")
				context.Request.Header.Set("Accept", "application/json")
			}, handler.ValuePostJSON)

			//Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/value", bytes.NewBufferString(tt.body))

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestUpdateJsonPostCounter(t *testing.T) {

	type mockBehaviorUpdateCounter func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string, value int64)

	tests := []struct {
		name         string
		body         string
		nameMetric   string
		mockBehavior mockBehaviorUpdateCounter
		statusCode   int
		responseBody string
	}{
		{
			name: "Update metrics counter 200",
			body: fmt.Sprintf(`{"id":"%s","type":"%s","delta": "%v"}`, "TestCounter", "counter", "500"),
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string, value int64) {
				expectedReturn := int64(500)

				mocks.EXPECT().UpdateMetricCounter(ctx, nameMetrics, value).Return(nil).AnyTimes()
				mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
			},

			statusCode:   200,
			responseBody: fmt.Sprintf(`{"id":"%s","type":"%s","delta":500}`, "TestCounter", "counter"),
		},
		{
			name: "get Metrics by name",
			body: fmt.Sprintf(`{"id":"%s","type":"%s","delta": "%v"}`, "TestCounter", "count", "500"),
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string, value int64) {

				mocks.EXPECT().UpdateMetricCounter(ctx, nameMetrics, value).Return(validator.ErrNotFoundType).AnyTimes()
			},
			statusCode:   400,
			responseBody: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Иницаиалзация тестирования
			c := gomock.NewController(t)
			defer c.Finish()

			s := mocks.NewMockMetricsRepository(c)

			services := s

			cfg, err := config.ParseConfig()

			if err != nil {
				fmt.Println("Config", err)
			}

			handler := handlers.NewMetricsHandlers(nil, services, cfg)

			//Init Point Handlers
			r := gin.Default()
			r.POST("/update", func(context *gin.Context) {
				tt.mockBehavior(context, s, "TestCounter", 500)
				//Прописыаем заголовки
				context.Request.Header.Set("Content-Type", "application/json")
				context.Request.Header.Set("Accept", "application/json")
			}, handler.ValuePostJSON)

			//Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/update", bytes.NewBufferString(tt.body))

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

func TestUpdateJsonPostGauge(t *testing.T) {

	type mockBehaviorUpdateGauge func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string, value float64)

	tests := []struct {
		name         string
		body         string
		nameMetric   string
		mockBehavior mockBehaviorUpdateGauge
		statusCode   int
		responseBody string
	}{
		{
			name: "Update metrics gauge 200",
			body: fmt.Sprintf(`{"id":"%s","type":"%s","value": "%v"}`, "TestCounter", "gauge", "500.123"),
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string, value float64) {
				expectedReturn := float64(500.123)

				mocks.EXPECT().UpdateMetricGauge(ctx, nameMetrics, value).Return(nil).AnyTimes()
				mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
			},

			statusCode:   200,
			responseBody: fmt.Sprintf(`{"id":"%s","type":"%s","value":500.123}`, "TestCounter", "gauge"),
		},
		{
			name: "get Metrics by name",
			body: fmt.Sprintf(`{"id":"%s","type":"%s","value": "%v"}`, "TestCounter", "ga", "500.123"),
			mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string, value float64) {

				mocks.EXPECT().UpdateMetricGauge(ctx, nameMetrics, value).Return(validator.ErrNotFoundType).AnyTimes()
			},
			statusCode:   400,
			responseBody: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Иницаиалзация тестирования
			c := gomock.NewController(t)
			defer c.Finish()

			s := mocks.NewMockMetricsRepository(c)

			services := s

			cfg, err := config.ParseConfig()

			if err != nil {
				fmt.Println("Config", err)
			}

			handler := handlers.NewMetricsHandlers(nil, services, cfg)

			//Init Point Handlers
			r := gin.Default()
			r.POST("/update", func(context *gin.Context) {
				tt.mockBehavior(context, s, "TestCounter", 500.123)
				//Прописыаем заголовки
				context.Request.Header.Set("Content-Type", "application/json")
				context.Request.Header.Set("Accept", "application/json")
			}, handler.ValuePostJSON)

			//Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/update", bytes.NewBufferString(tt.body))

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.responseBody, w.Body.String())
		})
	}
}

// ExampleMetricsHandlersType_GetMetricsCounter демонстрирует использование функции GetMetricsCounter.
func ExampleMetricsHandlersType_GetMetricsCounter_ok() {
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
		nameMetric: "TestGauge",
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
			expectedReturn := int64(700)

			mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
		},
		statusCode:   200,
		counterValue: "700",
	}

	t := &testing.T{}

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

// ExampleMetricsHandlersType_GetMetricsGauge демонстрирует использование функции GetMetricsGauge с ответов 404.
func ExampleMetricsHandlersType_GetMetricsGauge_notFound() {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tt := struct {
		name         string
		mType        string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		name:       "get Metrics by name",
		mType:      "gauge",
		nameMetric: "TestGauge",
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
			expectedReturn := float64(0)
			mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
		},
		statusCode:   404,
		responseBody: "",
	}

	t := &testing.T{}

	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	handler := handlers.NewMetricsHandlers(nil, services, cfg)

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
	// Status code: 404
	// Body:

}

// ExampleMetricsHandlersType_GetMetricsGauge демонстрирует использование функции GetMetricsGauge с ответов 200.
func ExampleMetricsHandlersType_GetMetricsGauge_ok() {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tt := struct {
		name         string
		mType        string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		name:       "get Metrics by name",
		mType:      "gauge",
		nameMetric: "TestCounter",
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
			expectedReturn := float64(700.123)

			mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
		},
		statusCode:   200,
		responseBody: fmt.Sprintf(`{"id":"%s","type":"%s","value":700.123}`, "TestCounter", "gauge"),
	}

	t := &testing.T{}

	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	handler := handlers.NewMetricsHandlers(nil, services, cfg)

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
	// Body: 700.123

}

// ExampleMetricsHandlersType_ValuePostJSON_Gauge_400 демонстрирует использование функции ValuePostJSON с ответов 400.
func ExampleMetricsHandlersType_ValuePostJSON_gauge_bad() {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tt := struct {
		name         string
		body         string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		name: "get Metrics by name",
		body: fmt.Sprintf(`{"id":"%s","type":"%s"}`, "TestCounter", "gaug"),
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
			expectedReturn := float64(0)
			mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
		},
		statusCode:   400,
		responseBody: "",
	}

	t := &testing.T{}
	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	//handler := metricsHandlers{metricsRepository: services, cfg: cfg}
	handler := handlers.NewMetricsHandlers(nil, services, cfg)

	//Init Point Handlers
	r := gin.Default()
	r.POST("/value", func(context *gin.Context) {
		tt.mockBehavior(context, s, "TestCounter")
		//Прописыаем заголовки
		context.Request.Header.Set("Content-Type", "application/json")
		context.Request.Header.Set("Accept", "application/json")
	}, handler.ValuePostJSON)

	//Create request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/value", bytes.NewBufferString(tt.body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	r.ServeHTTP(w, req)

	// Print the result.
	fmt.Printf("Status code: %v\n", w.Code)
	fmt.Printf("Body: %v\n", w.Body.String())

	// Output:
	// Status code: 400
	// Body:

}

// ExampleMetricsHandlersType_ValuePostJSON_Gauge_404 демонстрирует использование функции ValuePostJSON с ответов 404.
func ExampleMetricsHandlersType_ValuePostJSON_gauge_notFound() {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tt := struct {
		name         string
		body         string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		name: "get Metrics by name",
		body: fmt.Sprintf(`{"id":"%s","type":"%s"}`, "TestCounter", "gauge"),
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
			expectedReturn := float64(0)
			mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
		},
		statusCode:   404,
		responseBody: "",
	}

	t := &testing.T{}
	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	//handler := metricsHandlers{metricsRepository: services, cfg: cfg}
	handler := handlers.NewMetricsHandlers(nil, services, cfg)

	//Init Point Handlers
	r := gin.Default()
	r.POST("/value", func(context *gin.Context) {
		tt.mockBehavior(context, s, "TestCounter")
		//Прописыаем заголовки
		context.Request.Header.Set("Content-Type", "application/json")
		context.Request.Header.Set("Accept", "application/json")
	}, handler.ValuePostJSON)

	//Create request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/value", bytes.NewBufferString(tt.body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	r.ServeHTTP(w, req)

	// Print the result.
	fmt.Printf("Status code: %v\n", w.Code)
	fmt.Printf("Body: %v\n", w.Body.String())

	// Output:
	// Status code: 404
	// Body:

}

// ExampleMetricsHandlersTypeValuePostJSON_Gauge_200 демонстрирует использование функции ValuePostJSON с ответов 200.
func ExampleMetricsHandlersType_ValuePostJSON_gauge_ok() {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tt := struct {
		name         string
		body         string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{

		name: "get Metrics by name",
		body: fmt.Sprintf(`{"id":"%s","type":"%s"}`, "TestCounter", "gauge"),
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
			expectedReturn := float64(700.123)

			mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
		},
		statusCode:   200,
		responseBody: fmt.Sprintf(`{"id":"%s","type":"%s","value":700.123}`, "TestCounter", "gauge"),
	}

	t := &testing.T{}
	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	//handler := metricsHandlers{metricsRepository: services, cfg: cfg}
	handler := handlers.NewMetricsHandlers(nil, services, cfg)

	//Init Point Handlers
	r := gin.Default()
	r.POST("/value", func(context *gin.Context) {
		tt.mockBehavior(context, s, "TestCounter")
		//Прописыаем заголовки
		context.Request.Header.Set("Content-Type", "application/json")
		context.Request.Header.Set("Accept", "application/json")
	}, handler.ValuePostJSON)

	//Create request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/value", bytes.NewBufferString(tt.body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	r.ServeHTTP(w, req)

	// Print the result.
	fmt.Printf("Status code: %v\n", w.Code)
	fmt.Printf("Body: %v\n", w.Body.String())

	// Output:
	// Status code: 200
	// Body: {"id":"TestCounter","type":"gauge","value":700.123}

}

// ExampleMetricsHandlersType_ValuePostJSON_Counter_400 демонстрирует использование функции ValuePostJSON с ответов 400.
func ExampleMetricsHandlersType_ValuePostJSON_counter_bad() {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tt := struct {
		name         string
		body         string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		name: "get Metrics by name",
		body: fmt.Sprintf(`{"id":"%s","type":"%s"}`, "TestCounter", "unter"),
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
			expectedReturn := int64(0)
			mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
		},
		statusCode:   400,
		responseBody: "",
	}

	t := &testing.T{}
	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	//handler := metricsHandlers{metricsRepository: services, cfg: cfg}
	handler := handlers.NewMetricsHandlers(nil, services, cfg)

	//Init Point Handlers
	r := gin.Default()
	r.POST("/value", func(context *gin.Context) {
		tt.mockBehavior(context, s, "TestCounter")
		//Прописыаем заголовки
		context.Request.Header.Set("Content-Type", "application/json")
		context.Request.Header.Set("Accept", "application/json")
	}, handler.ValuePostJSON)

	//Create request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/value", bytes.NewBufferString(tt.body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	r.ServeHTTP(w, req)

	// Print the result.
	fmt.Printf("Status code: %v\n", w.Code)
	fmt.Printf("Body: %v\n", w.Body.String())

	// Output:
	// Status code: 400
	// Body:

}

// ExampleMetricsHandlersType_ValuePostJSON_Counter_200 демонстрирует использование функции ValuePostJSON с ответов 200.
func ExampleMetricsHandlersType_ValuePostJSON_counter_ok() {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tt := struct {
		name         string
		body         string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{

		name: "get Metrics by name",
		body: fmt.Sprintf(`{"id":"%s","type":"%s"}`, "TestCounter", "counter"),
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
			expectedReturn := int64(700)

			mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
		},
		statusCode:   200,
		responseBody: fmt.Sprintf(`{"id":"%s","type":"%s","delta":700}`, "TestCounter", "counter"),
	}

	t := &testing.T{}
	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	//handler := metricsHandlers{metricsRepository: services, cfg: cfg}
	handler := handlers.NewMetricsHandlers(nil, services, cfg)

	//Init Point Handlers
	r := gin.Default()
	r.POST("/value", func(context *gin.Context) {
		tt.mockBehavior(context, s, "TestCounter")
		//Прописыаем заголовки
		context.Request.Header.Set("Content-Type", "application/json")
		context.Request.Header.Set("Accept", "application/json")
	}, handler.ValuePostJSON)

	//Create request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/value", bytes.NewBufferString(tt.body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	r.ServeHTTP(w, req)

	// Print the result.
	fmt.Printf("Status code: %v\n", w.Code)
	fmt.Printf("Body: %v\n", w.Body.String())

	// Output:
	// Status code: 200
	// Body: {"id":"TestCounter","type":"counter","delta":700}

}

// ExampleMetricsHandlersType_ValuePostJSON_Counter_404 демонстрирует использование функции ValuePostJSON с ответом 404.
func ExampleMetricsHandlersType_ValuePostJSON_counter_notFound() {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string)

	tt := struct {
		name         string
		body         string
		nameMetric   string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		name: "get Metrics by name",
		body: fmt.Sprintf(`{"id":"%s","type":"%s"}`, "TestCounter", "counter"),
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string) {
			expectedReturn := int64(0)
			mocks.EXPECT().GetMetricCounter(ctx, nameMetrics).Return(expectedReturn, validator.ErrMetricsKeyNotFound).AnyTimes()
		},
		statusCode:   404,
		responseBody: "",
	}

	t := &testing.T{}
	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	//handler := metricsHandlers{metricsRepository: services, cfg: cfg}
	handler := handlers.NewMetricsHandlers(nil, services, cfg)

	//Init Point Handlers
	r := gin.Default()
	r.POST("/value", func(context *gin.Context) {
		tt.mockBehavior(context, s, "TestCounter")
		//Прописыаем заголовки
		context.Request.Header.Set("Content-Type", "application/json")
		context.Request.Header.Set("Accept", "application/json")
	}, handler.ValuePostJSON)

	//Create request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/value", bytes.NewBufferString(tt.body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	r.ServeHTTP(w, req)

	// Print the result.
	fmt.Printf("Status code: %v\n", w.Code)
	fmt.Printf("Body: %v\n", w.Body.String())

	// Output:
	// Status code: 404
	// Body:

}

// ExampleMetricsHandlersType_UpdateMetrics демонстрирует использование функции UpdateMetrics.
func ExampleMetricsHandlersType_UpdateMetrics() {

	type mockBehaviorUpdateGauge func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string, value float64)

	tt := struct {
		name         string
		body         string
		nameMetric   string
		mockBehavior mockBehaviorUpdateGauge
		statusCode   int
		responseBody string
	}{

		name: "Update metrics gauge 200",
		body: fmt.Sprintf(`{"id":"%s","type":"%s","value": %v}`, "TestGauge", "gauge", 500.123),
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository, nameMetrics string, value float64) {
			expectedReturn := float64(500.123)

			mocks.EXPECT().UpdateMetricGauge(ctx, nameMetrics, value).Return(nil).AnyTimes()
			mocks.EXPECT().GetMetricGauge(ctx, nameMetrics).Return(expectedReturn, nil).AnyTimes()
		},

		statusCode:   200,
		responseBody: fmt.Sprintf(`{"id":"%s","type":"%s","value":500.123}`, "TestGauge", "gauge"),
	}

	t := &testing.T{}
	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	handler := handlers.NewMetricsHandlers(nil, services, cfg)

	//Init Point Handlers
	r := gin.Default()
	r.POST("/update", func(context *gin.Context) {
		tt.mockBehavior(context, s, "TestGauge", 500.123)
		//Прописыаем заголовки
		context.Request.Header.Set("Content-Type", "application/json")
		context.Request.Header.Set("Accept", "application/json")
	}, handler.UpdateMetrics)

	//Create request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/update", bytes.NewBufferString(tt.body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	r.ServeHTTP(w, req)

	// Print the result.
	fmt.Printf("Status code: %v\n", w.Code)
	fmt.Printf("Body: %v\n", w.Body.String())

	// Output:
	// Status code: 200
	// Body: {"id":"TestGauge","type":"gauge","value":500.123}
}

// ExampleMetricsHandlersType_GetAllMetricsHTML демонстрирует использование функции GetAllMetricsHTML.
func ExampleMetricsHandlersType_GetAllMetricsHTML() {

	metricsModel := storage.NewMemStorage()

	metricsModel.Gauge["metrics1"] = 543.4657

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository)

	tt := struct {
		name         string
		mockBehavior mockBehavior
		statusCode   int
		counterValue string
	}{
		name: "get all metrics 200",
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository) {
			expectedReturn := metricsModel

			mocks.EXPECT().GetAllMetrics(ctx).Return(expectedReturn, nil).AnyTimes()
		},
		statusCode: 200,
		//counterValue: "<div>counter => 500 </div><div>metrics1 => 543.4657 </div><div>metrics2 => 456.4657 </div><div>metrics3 => 432.4657 </div>",
	}

	t := &testing.T{}

	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	handler := handlers.NewMetricsHandlers(nil, services, cfg)

	//Init Point Handlers
	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		tt.mockBehavior(context, s)
	}, handler.GetAllMetricsHTML)

	//Create request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	r.ServeHTTP(w, req)

	// Print the result.
	fmt.Printf("Status code: %v\n", w.Code)
	fmt.Printf("Body: %v\n", w.Body.String())

	// Output:
	// Status code: 200
	// Body: <div>metrics1 => 543.4657 </div>
}

// ExampleMetricsHandlersType_Ping демонстрирует использование функции Ping с ответом 200.
func ExampleMetricsHandlersType_Ping_success() {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository)

	tt := struct {
		name         string
		mockBehavior mockBehavior
		statusCode   int
		counterValue string
	}{

		name: "get ping database connect to posgress 200",
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository) {
			mocks.EXPECT().PingDatabase(ctx).Return(nil).AnyTimes()
		},
		statusCode:   200,
		counterValue: "Success to ping database",
	}

	t := &testing.T{}

	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	handler := handlers.NewMetricsHandlers(nil, services, cfg)

	//Init Point Handlers
	r := gin.Default()
	r.GET("/ping", func(context *gin.Context) {
		tt.mockBehavior(context, s)
	}, handler.Ping)

	//Create request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ping", nil)

	r.ServeHTTP(w, req)

	// Print the result.
	fmt.Printf("Status code: %v\n", w.Code)
	fmt.Printf("Body: %v\n", w.Body.String())

	// Output:
	// Status code: 200
	// Body: Success to ping database
}

// ExampleMetricsHandlersType_Ping демонстрирует использование функции Ping с ответом 500.
func ExampleMetricsHandlersType_Ping_failed() {

	type mockBehavior func(ctx context.Context, mocks *mocks.MockMetricsRepository)

	tt := struct {
		name         string
		mockBehavior mockBehavior
		statusCode   int
		counterValue string
	}{
		name: "get ping database connect to posgress 500",
		mockBehavior: func(ctx context.Context, mocks *mocks.MockMetricsRepository) {
			mocks.EXPECT().PingDatabase(ctx).Return(validator.ErrBadRequest).AnyTimes()
		},
		statusCode:   500,
		counterValue: "Failed to ping database",
	}

	t := &testing.T{}

	//Иницаиалзация тестирования
	c := gomock.NewController(t)
	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	services := s

	cfg, err := config.ParseConfig()

	if err != nil {
		fmt.Println("Config", err)
	}

	handler := handlers.NewMetricsHandlers(nil, services, cfg)

	//Init Point Handlers
	r := gin.Default()
	r.GET("/ping", func(context *gin.Context) {
		tt.mockBehavior(context, s)
	}, handler.Ping)

	//Create request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ping", nil)

	r.ServeHTTP(w, req)

	// Print the result.
	fmt.Printf("Status code: %v\n", w.Code)
	fmt.Printf("Body: %v\n", w.Body.String())

	// Output:
	// Status code: 500
	// Body: Failed to ping database
}
