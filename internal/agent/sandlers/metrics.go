package sandlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pkg/compress"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/pkg/serialize"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/agent/repository"

	"github.com/go-resty/resty/v2"
)

type SandlerMetrics interface {
	ChangeMetricsByTime()
	SendMetricsByTime()
	SendMetrics()
}

type sandlerMetrics struct {
	repository repository.RepositoryMetrics
	client     *resty.Client
	ctx        context.Context
	cfg        *config.Config
}

func NewMetricsSendler(repository repository.RepositoryMetrics, client *resty.Client,
	ctx context.Context, cfg *config.Config) SandlerMetrics {
	return &sandlerMetrics{repository: repository, client: client, ctx: ctx, cfg: cfg}
}

// Change metrics
func (rm *sandlerMetrics) ChangeMetricsByTime() {
	secondChange := time.Duration(rm.cfg.Metrics.PollInterval)
	for {
		select {
		case <-rm.ctx.Done():
			fmt.Println("ChangeMetricsByTime -> Эй! Энштейн! Спасибо, что остановили мою горутину :)")
			return
		// Run change metrics before sleep 2 seconds
		default:
			{
				time.Sleep(secondChange * time.Second)
				rm.repository.ChangeMetrics()
			}
		}
	}
}

// Send metrics
func (rm *sandlerMetrics) SendMetricsByTime() {
	secondSend := time.Duration(rm.cfg.Metrics.ReportInterval)

	for {
		select {
		case <-rm.ctx.Done():
			fmt.Println("SendMetricsByTime -> Эй! Энштейн! Спасибо, что остановили мою горутину :)")
			return
		// Run change metrics before sleep 2 seconds
		default:
			{
				// Run change metrics before sleep 2 seconds
				time.Sleep(secondSend * time.Second)
				rm.SendMetrics()
			}
		}
	}
}

// Send metrics to server http://localhost:8080
func (rm *sandlerMetrics) SendMetrics() {
	cfg := rm.cfg

	serializer := serialize.NewSerializer(rm.cfg)

	fmt.Printf("Запустили отправку метрик: http://%s\n", cfg.Server.Address)

	defer fmt.Printf("Завершили отправку метрик: http://%s\n", cfg.Server.Address)

	client := rm.client

	metrics, _ := rm.repository.GetGaugeMetricsAll()

	var sendStringMetrics string

	//Send metrics GAUGE
	for key, val := range metrics.Gauge {

		metricsSData := serialize.Metrics{ID: key, MType: "gauge", Delta: nil, Value: &val}
		serializeErr := serializer.SetData(&metricsSData).GetData(&sendStringMetrics)

		if serializeErr.Errors() != nil {
			panic(serializeErr.Errors().Error())
		}

		url := fmt.Sprintf("http://%s/update", cfg.Server.Address)
		//url = fmt.Sprintf("http://%s/update/gauge/%s/%v", cfg.Server.Address, key, val)

		var err error
		if cfg.Gzip.Enable {
			sendDataCompress, _ := compress.CompressGzip([]byte(sendStringMetrics))
			_, err = client.R().
				SetBody(sendDataCompress).
				Post(url)
		} else {
			_, err = client.R().
				SetBody(sendStringMetrics).
				Post(url)
		}

		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				log.Print("Вот такая вот ошибка: This is broken pipe error")
				if !errors.Is(err, syscall.ECONNREFUSED) {
					log.Print("Ошибка не: ECONNREFUSED", err.Error())
					//log.Fatal(err)
				} else {
					log.Println("ECONNREFUSED")
					break
				}
			}
			//fmt.Println(response)
		}

		metrics, _ = rm.repository.GetCounterMetricsAll()

		//Send metrics COUNTER
		for key, val := range metrics.Counter {
			metricsSData := serialize.Metrics{ID: key, MType: "counter", Delta: &val, Value: nil}
			serializeErr := serializer.SetData(&metricsSData).GetData(&sendStringMetrics)

			if serializeErr.Errors() != nil {
				panic(serializeErr.Errors().Error())
			}

			url := fmt.Sprintf("http://%s/update", cfg.Server.Address)
			//url = fmt.Sprintf("http://%s/update/counter/%s/%v", cfg.Server.Address, key, val)

			var err error
			if cfg.Gzip.Enable {
				sendDataCompress, _ := compress.CompressGzip([]byte(sendStringMetrics))
				_, err = client.R().
					SetBody(sendDataCompress).
					Post(url)
			} else {
				_, err = client.R().
					SetBody(sendStringMetrics).
					Post(url)
			}
			if err != nil {
				if errors.Is(err, syscall.EPIPE) {
					log.Println("Вот такая вот ошибка: This is broken pipe error")

					if !errors.Is(err, syscall.ECONNREFUSED) {
						//log.Fatal(err)
						log.Println("Ошибка не: ECONNREFUSED", err.Error())
					} else {
						log.Println("ECONNREFUSED")
						break
					}
				}
			}
			//fmt.Println(response)
		}
	}
}
