package sandlers

import (
	"context"
	"encoding/json"
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
	SendMetrics([]string)
}

type sandlerMetrics struct {
	repository       repository.RepositoryMetrics
	client           *resty.Client
	ctx              context.Context
	cfg              *config.Config
	listMetrics      []string
	listMetricsBatch []string
	urlMetrics       string
	sendBatch        bool
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
			fmt.Println("ChangeMetricsByTime stop")
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
			fmt.Println("SendMetricsByTime stop")
			return
		// Run change metrics before sleep 2 seconds
		default:
			{
				// Run change metrics before sleep 2 seconds
				time.Sleep(secondSend * time.Second)
				rm.setMetrics()
				rm.SendMetrics(rm.listMetrics)
				rm.SendMetrics(rm.listMetricsBatch)
				rm.repository.SetZeroPollCount()
			}
		}
	}
}

// SendMetrics sends Batched metrics to local server
func (rm *sandlerMetrics) GetBatchStringMetrics() []string {
	listMetrics := rm.getListMetrics()

	json, err := json.Marshal(listMetrics)
	if err != nil {
		fmt.Println("Marshal failed: ", err)
	}
	jsonString := string(json)

	fmt.Println(jsonString)

	return []string{jsonString}
}

// SendMetrics sends Slices to local server
func (rm *sandlerMetrics) GetSliceStringMetrics() []string {
	listMetrics := rm.getListMetrics()

	//fmt.Println("listMetrics")
	//fmt.Println(listMetrics)
	//fmt.Println("listMetrics")

	var listMetricsString []string

	for _, m := range listMetrics {
		//fmt.Println("listMetrics=>>>>", *m.Value)

		sendStringMetrics := rm.serializeMetrics(m)
		listMetricsString = append(listMetricsString, sendStringMetrics)
	}

	fmt.Println("listMetrics=>>>>", listMetricsString)

	return listMetricsString
}

/*func (rm *sandlerMetrics) serverPing() {

	rm.sendBatch = false

	url := fmt.Sprintf("http://%s/ping", rm.cfg.Server.Address)
	resp, err := rm.client.R().Get(url)
	if err != nil {
		log.Println("Ошибка не: ECONNREFUSED", err.Error())
	}

	if resp.IsError() {
		fmt.Println("Status Error:", resp.StatusCode()) // prints 404
	}

	if resp.StatusCode() == 200 {
		rm.sendBatch = true
	}

}*/

func (rm *sandlerMetrics) setMetrics() {
	cfg := rm.cfg

	//var listMetrics []string

	//rm.serverPing()

	//if rm.sendBatch {
	fmt.Println("Send One Batch metrics")
	rm.urlMetrics = fmt.Sprintf("http://%s/updates", cfg.Server.Address)
	rm.listMetricsBatch = rm.GetBatchStringMetrics()
	//} else {
	fmt.Println("Send Single request metrics")
	rm.urlMetrics = fmt.Sprintf("http://%s/update", cfg.Server.Address)
	rm.listMetrics = rm.GetSliceStringMetrics()
	//}

	//rm.listMetrics = listMetrics
}

func (rm *sandlerMetrics) SendMetrics(listMetrics []string) {
	cfg := rm.cfg

	//rm.setMetrics()

	fmt.Printf("Запустили отправку метрик: %s\n", rm.urlMetrics)

	defer fmt.Printf("Завершили отправку метрик: %s\n", rm.urlMetrics)

	client := rm.client

	//url := fmt.Sprintf("http://%s/update", cfg.Server.Address)

	var err error

	for _, sendStringMetrics := range listMetrics {
		if cfg.Gzip.Enable {
			sendDataCompress, _ := compress.CompressGzip([]byte(sendStringMetrics))
			_, err = client.R().
				SetBody(sendDataCompress).
				Post(rm.urlMetrics)
		} else {
			_, err = client.R().
				SetBody(sendStringMetrics).
				Post(rm.urlMetrics)
		}

		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				log.Print("Вот такая вот ошибка: This is broken pipe error")
				if !errors.Is(err, syscall.ECONNREFUSED) {
					log.Print("Ошибка не: ECONNREFUSED", err.Error())
					//log.Fatal(err)
				} else {
					log.Println("ECONNREFUSED")
				}
			}
			//fmt.Println(response)
		}
	}

}

// Send metrics to server http://localhost:8080
func (rm *sandlerMetrics) getListMetrics() []serialize.Metrics {

	metrics, _ := rm.repository.GetGaugeMetricsAll()

	var collectMetrics []serialize.Metrics

	//Send metrics GAUGE
	for key, val := range metrics.Gauge {

		//fmt.Println("value:", val)

		valNew := val

		metricsSData := serialize.Metrics{ID: key, MType: "gauge", Delta: nil, Value: &valNew}

		//fmt.Println(metricsSData)
		//fmt.Println("value =>>>", *metricsSData.Value)

		collectMetrics = append(collectMetrics, metricsSData)
	}

	//os.Exit(1)

	metrics, _ = rm.repository.GetCounterMetricsAll()

	//Send metrics COUNTER
	for key, val := range metrics.Counter {

		deltaNew := val

		metricsSData := serialize.Metrics{ID: key, MType: "counter", Delta: &deltaNew, Value: nil}

		collectMetrics = append(collectMetrics, metricsSData)

	}

	fmt.Println(collectMetrics)
	//var collectMetrics []serialize.Metrics

	return collectMetrics

}

func (rm *sandlerMetrics) serializeMetrics(metricsSData serialize.Metrics) string {

	//data := metricsSData

	serializer := serialize.NewSerializer(rm.cfg)

	var sendStringMetrics string

	serializeErr := serializer.SetData(&metricsSData).GetData(&sendStringMetrics)

	if serializeErr.Errors() != nil {
		panic(serializeErr.Errors().Error())
	}

	//fmt.Println("serializeMetrics====>", sendStringMetrics)

	return sendStringMetrics
}

/*

 */
