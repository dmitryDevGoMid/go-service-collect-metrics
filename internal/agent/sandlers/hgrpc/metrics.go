package hgrpc

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
	pb "github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetClient(ctx context.Context, cancel context.CancelFunc, cfg *config.Config) (pb.MetricsServiceClient, error) {

	//rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient("127.0.0.1:50051", opts...)
	if err != nil {
		//log.Fatalf("did not connect: %v", err)
		return nil, err
	}

	client := pb.NewMetricsServiceClient(conn)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancel()
		fmt.Println("Close Client Connect!")
		conn.Close()
		os.Exit(0)
	}()

	return client, nil
}
