package hgrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	pb "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pb"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository"
	mocks "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository/mock_repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type ServerGRPC struct {
	pb.UnimplementedMetricsServiceServer
	MetricsRepository repository.MetricsRepository
}

// Конструктор для GRPC обработчика
func NewGRPCHandlers(metricsRepository repository.MetricsRepository) *ServerGRPC {

	return &ServerGRPC{
		MetricsRepository: metricsRepository,
	}
}

func (s ServerGRPC) Metrics(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	var data map[string]interface{}

	switch req.Mtype {
	case "gauge":
		data = map[string]interface{}{
			"mtype": "gauge",
			"id":    req.Id,
			"value": req.Value,
		}

		err := s.MetricsRepository.UpdateMetricGauge(ctx, req.Id, req.Value)
		if err != nil {
			//	log.Fatalf("update metric counter %v error", err)
			return nil, fmt.Errorf("invalid mtype: %s", req.Mtype)
		}
		_, err = s.MetricsRepository.GetMetricGauge(ctx, req.Id)

		if err != nil {
			return nil, fmt.Errorf("invalid update mtype: %s, name: %s", req.Mtype, req.Id)
		}

	case "counter":
		data = map[string]interface{}{
			"mtype": "counter",
			"id":    req.Id,
			"delta": req.Delta,
		}

		err := s.MetricsRepository.UpdateMetricCounter(ctx, req.Id, req.Delta)
		if err != nil {
			//	log.Fatalf("update metric counter %v error", err)
			return nil, fmt.Errorf("invalid mtype: %s", req.Mtype)
		}
		_, err = s.MetricsRepository.GetMetricCounter(ctx, req.Id)

		if err != nil {
			return nil, fmt.Errorf("invalid update mtype: %s, name: %s", req.Mtype, req.Id)
		}

	default:
		return nil, fmt.Errorf("invalid mtype: %s", req.Mtype)
	}

	_, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	//log.Printf("Received metrics: %s", jsonData)

	return &pb.MetricsResponse{Success: true}, nil
}

func (s ServerGRPC) MetricsTest(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	var data map[string]interface{}

	switch req.Mtype {
	case "gauge":
		data = map[string]interface{}{
			"mtype": "gauge",
			"id":    req.Id,
			"value": req.Value,
		}

		t := &testing.T{}

		c := gomock.NewController(t)

		defer c.Finish()

		s := mocks.NewMockMetricsRepository(c)

		//services := s
		metricsRepository := s

		expectedReturn := float64(700.123)

		s.EXPECT().UpdateMetricGauge(ctx, "test", expectedReturn).Return(nil).AnyTimes()
		s.EXPECT().GetMetricGauge(ctx, "test").Return(expectedReturn, nil).AnyTimes()

		metricsRepository.UpdateMetricGauge(ctx, "test", expectedReturn)
		rezult, err := metricsRepository.GetMetricGauge(ctx, "test")

		if err != nil {
			log.Fatalf("update metric counter %v error", err)
		}

		fmt.Println(rezult)
		fmt.Println(assert.Equal(t, 700.123, rezult))

		if !assert.Equal(t, 700.123, rezult) {
			log.Fatalf("test gauge grpc metrics is FAIL...")
		}

	case "counter":
		data = map[string]interface{}{
			"mtype": "counter",
			"id":    req.Id,
			"delta": req.Delta,
		}

		t := &testing.T{}

		c := gomock.NewController(t)

		defer c.Finish()

		s := mocks.NewMockMetricsRepository(c)

		//services := s
		metricsRepository := s

		expectedReturn := int64(700)

		s.EXPECT().UpdateMetricCounter(ctx, "test", expectedReturn).Return(nil).AnyTimes()
		s.EXPECT().GetMetricCounter(ctx, "test").Return(expectedReturn, nil).AnyTimes()

		metricsRepository.UpdateMetricCounter(ctx, "test", expectedReturn)
		rezult, err := metricsRepository.GetMetricCounter(ctx, "test")

		if err != nil {
			log.Fatalf("update metric counter %v error", err)
		}

		fmt.Println(rezult)
		fmt.Println(assert.Equal(t, int64(700), rezult))

		if !assert.Equal(t, int64(700), rezult) {
			log.Fatalf("test counter grpc metrics is FAIL...")
		}

	default:
		return nil, fmt.Errorf("invalid mtype: %s", req.Mtype)
	}

	_, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	//log.Printf("Received metrics: %s", jsonData)

	return &pb.MetricsResponse{Success: true}, nil
}
