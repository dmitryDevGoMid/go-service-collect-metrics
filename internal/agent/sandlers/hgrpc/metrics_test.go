package hgrpc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	pb "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pb"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type ServerGRPC struct {
	pb.UnimplementedMetricsServiceServer
	metricsRepository repository.MetricsRepository
}

// Конструктор для GRPC обработчика
func NewGRPCHandlers(metricsRepository repository.MetricsRepository) *ServerGRPC {
	return &ServerGRPC{
		metricsRepository: metricsRepository,
	}
}

func (s *ServerGRPC) Metrics(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	var data map[string]interface{}

	switch req.Mtype {
	case "gauge":
		data = map[string]interface{}{
			"mtype": "gauge",
			"id":    req.Id,
			"value": req.Value,
		}

		s.metricsRepository.UpdateMetricGauge(ctx, req.Id, req.Value)
		_, err := s.metricsRepository.GetMetricGauge(ctx, req.Id)

		if err != nil {
			return nil, fmt.Errorf("invalid update mtype: %s, name: %s", req.Mtype, req.Id)
		}

	case "counter":
		data = map[string]interface{}{
			"mtype": "counter",
			"id":    req.Id,
			"delta": req.Delta,
		}

		s.metricsRepository.UpdateMetricCounter(ctx, req.Id, req.Delta)
		_, err := s.metricsRepository.GetMetricGauge(ctx, req.Id)

		if err != nil {
			return nil, fmt.Errorf("invalid update mtype: %s, name: %s", req.Mtype, req.Id)
		}
	default:
		return nil, fmt.Errorf("invalid mtype: %s", req.Mtype)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	log.Printf("Received metrics: %s", jsonData)

	return &pb.MetricsResponse{Success: true}, nil
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	handlerGRPC := NewGRPCHandlers(nil)
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterMetricsServiceServer(s, handlerGRPC)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

/*import "google.golang.org/grpc/test/bufconn"

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
    lis = bufconn.Listen(bufSize)
    s := grpc.NewServer()
    pb.RegisterGreeterServer(s, &server{})
    go func() {
        if err := s.Serve(lis); err != nil {
            log.Fatalf("Server exited with error: %v", err)
        }
    }()
}

func bufDialer(context.Context, string) (net.Conn, error) {
    return lis.Dial()
}

func TestSayHello(t *testing.T) {
    ctx := context.Background()
    conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
    if err != nil {
        t.Fatalf("Failed to dial bufnet: %v", err)
    }
    defer conn.Close()
    client := pb.NewGreeterClient(conn)
    resp, err := client.SayHello(ctx, &pb.HelloRequest{"Dr. Seuss"})
    if err != nil {
        t.Fatalf("SayHello failed: %v", err)
    }
    log.Printf("Response: %+v", resp)
    // Test for output here.
}*/

func TestMetrics(t *testing.T) {
	// набор тестовых данных
	metrics := []*pb.MetricsRequest{
		{Id: "my_gauge", Mtype: "gauge", Value: 123.234},
		{Id: "my_counter", Mtype: "counter", Delta: 456},
		// при добавлении этой записи должна вернуться ошибка:
		// неверный тип метрики
		{Id: "my_wrong", Mtype: "wrong", Value: 123},
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient("127.0.0.1:50051", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMetricsServiceClient(conn)

	// добавляем метрики
	for _, metric := range metrics {
		resp, err := client.Metrics(context.Background(), metric)
		if err != nil {
			splits := strings.Split(err.Error(), "=")
			last := len(splits) - 1
			fmt.Println(splits[len(splits)-1])

			assert.Equal(t, "invalid mtype: wrong", strings.TrimSpace(splits[last]))
			//break
		} else {

			assert.Equal(t, true, resp.Success)
		}
	}
}
