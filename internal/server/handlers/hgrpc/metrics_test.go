package hgrpc_test

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/handlers/hgrpc"
	pb "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pb"
	mocks "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository/mock_repository"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestMetric(t *testing.T) {

	ctx := context.Background()

	//Иницаиалзация тестирования
	c := gomock.NewController(t)

	defer c.Finish()

	s := mocks.NewMockMetricsRepository(c)

	handlerGRPC := hgrpc.NewGRPCHandlers(s)

	lis = bufconn.Listen(bufSize)

	serverGRPC := grpc.NewServer()

	defer serverGRPC.GracefulStop()

	pb.RegisterMetricsServiceServer(serverGRPC, handlerGRPC)

	go func() {
		if err := serverGRPC.Serve(lis); err != nil {
			log.Fatalf("Server exited with error3: %v", err)
		}
	}()

	resolver.SetDefaultScheme("passthrough")

	opts := []grpc.DialOption{
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient("bufnet", opts...)

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	if err != nil {
		t.Errorf("Test failed")
	}

	m := pb.MetricsRequest{Id: "test", Mtype: "counter", Delta: 123}

	client := pb.NewMetricsServiceClient(conn)

	resp, _ := client.MetricsTest(ctx, &m)

	assert.Equal(t, true, resp.Success)

	m = pb.MetricsRequest{Id: "test", Mtype: "gauge", Value: 123.123}

	client = pb.NewMetricsServiceClient(conn)

	resp, _ = client.MetricsTest(ctx, &m)

	assert.Equal(t, true, resp.Success)
}
