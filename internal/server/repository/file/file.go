package file

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/serialize"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/pkg/unserialize"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository"
)

type WorkerFile interface {
	SaveMetricsByTime(ctx context.Context)
	SaveAllMetrics(ctx context.Context)
	RunWorker(ctx context.Context)
}

var mutex sync.Mutex

type workerFile struct {
	metricsRepository repository.MetricsRepository
	cfg               *config.Config
	ctx               context.Context
	mutex             *sync.Mutex
	f                 *os.File
}

func NewWorkFile(metricsRepository repository.MetricsRepository, cfg *config.Config, ctx context.Context) WorkerFile {
	return &workerFile{metricsRepository: metricsRepository, cfg: cfg, ctx: ctx}
}

func (wf *workerFile) setFolder() {
	folderPath := wf.cfg.File.FileStoragePath

	d, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	dir := filepath.Dir(folderPath)

	err = os.MkdirAll(d+dir, 0775)

	if err != nil {
		fmt.Println(err)
	}

	filename := d + folderPath

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

	if err != nil {
		fmt.Println(err)
	}

	wf.f = f
}

func (wf *workerFile) RunWorker(ctx context.Context) {
	wf.mutex = &mutex

	if wf.cfg.File.Restore {
		wf.GetAllMetricsByFile(ctx)
	}

	if wf.cfg.File.StoreInterval > 0 {
		go wf.SaveMetricsByTime(ctx)
	}
}

func (wf *workerFile) SaveMetricsByTime(ctx context.Context) {
	secondChange := time.Duration(wf.cfg.File.StoreInterval)
	//i := 0
	for {
		select {
		case <-wf.ctx.Done():
			fmt.Println("SaveMetricsByTime Stop")
			return
		// Run change metrics before sleep 2 seconds
		default:
			{
				time.Sleep(secondChange * time.Second)
				wf.SaveAllMetrics(ctx)
			}
		}
	}
}

// Пишем метрики в файл
func (wf *workerFile) SaveAllMetrics(ctx context.Context) {

	if wf.cfg.File.FileStoragePath == "" {
		return
	}

	defer func() {
		wf.mutex.Unlock()
	}()

	wf.mutex.Lock()

	wf.setFolder()

	allMetrics, _ := wf.metricsRepository.GetAllMetrics(context.TODO())

	typeMetric := "gauge"

	serializer := serialize.NewSerializer(wf.cfg)

	for key, val := range allMetrics.Gauge {
		metricsSData := serialize.Metrics{ID: key, MType: typeMetric, Delta: nil, Value: &val}
		sendData := serializer.SerializerResponse(&metricsSData)

		if _, err := wf.f.WriteString(fmt.Sprintf("%s\n", *sendData)); err != nil {
			log.Println(err)
		}
	}

	typeMetric = "counter"

	for key, val := range allMetrics.Counter {
		metricsSData := serialize.Metrics{ID: key, MType: typeMetric, Delta: &val, Value: nil}
		sendData := serializer.SerializerResponse(&metricsSData)

		if _, err := wf.f.WriteString(fmt.Sprintf("%s\n", *sendData)); err != nil {
			log.Println(err)
		}
	}

}

// Поднимаем метрики из файл и пишем в хранилище
func (wf *workerFile) GetAllMetricsByFile(ctx context.Context) {

	d, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		return
	}

	folderPath := d + wf.cfg.File.FileStoragePath
	_, err = os.Stat(folderPath)

	if errors.Is(err, os.ErrNotExist) {
		return
	}

	file, err := os.Open(folderPath)
	if err != nil {
		log.Println(err)
		return
	}

	defer file.Close()

	unserializeData := unserialize.NewUnSerializer(wf.cfg)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		metrics := unserializeData.UnSerialize([]byte(scanner.Text()))

		switch val := metrics.MType; val {
		case "gauge":
			wf.metricsRepository.UpdateMetricGauge(context.TODO(), metrics.ID, *metrics.Value)
		case "counter":
			wf.metricsRepository.UpdateMetricCounter(context.TODO(), metrics.ID, *metrics.Delta)
		default:
			fmt.Println("Нет такого типа метрики!")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
		return
	}

}
