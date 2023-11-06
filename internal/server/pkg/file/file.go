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
	SaveMetricsByTime()
	SaveAllMetrics()
	SetFolder()
	RunWorker()
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

func (wf *workerFile) SetFolder() {
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

func (wf *workerFile) RunWorker() {
	wf.mutex = &mutex

	if wf.cfg.File.Restore {
		wf.GetAllMetricsByFile()
	}

	if wf.cfg.File.StoreInterval > 0 {
		go wf.SaveMetricsByTime()
	}
}

func (wf *workerFile) SaveMetricsByTime() {
	secondChange := time.Duration(wf.cfg.File.StoreInterval)
	//i := 0
	for {
		select {
		case <-wf.ctx.Done():
			fmt.Println("SaveMetricsByTime -> Эй! Энштейн! Спасибо, что остановили мою горутину :)")
			return
		// Run change metrics before sleep 2 seconds
		default:
			{
				time.Sleep(secondChange * time.Second)
				wf.SaveAllMetrics()
				//wf.GetAllMetricsByFile()
				//i++
				//wf.WriteToFile(i)
			}
		}
	}
}

// Пишем метрики в файл
func (wf *workerFile) SaveAllMetrics() {

	if wf.cfg.File.FileStoragePath == "" {
		return
	}

	defer func() {
		wf.mutex.Unlock()
	}()

	wf.mutex.Lock()

	wf.SetFolder()

	allMetrics, _ := wf.metricsRepository.GetAllMetrics()

	typeMetric := "gauge"

	serializer := serialize.NewSerializer(wf.cfg)

	for key, val := range allMetrics.Gauge {
		metricsSData := serialize.Metrics{ID: key, MType: typeMetric, Delta: nil, Value: &val}
		sendData := serializer.SerializerResponse(&metricsSData)

		//fmt.Fprintf(wf.f, "%s\n", *sendData)
		//curcurrentTime := time.Now()
		//timeAndDate := fmt.Sprintf("YYYY-MM-DD hh:mm:ss : ", curcurrentTime.Format("2006-01-02 15:04:05"))
		//if _, err := wf.f.WriteString(fmt.Sprintf("%s<=>%s\n", timeAndDate, *sendData)); err != nil {
		if _, err := wf.f.WriteString(fmt.Sprintf("%s\n", *sendData)); err != nil {
			log.Println(err)
		}
	}

	typeMetric = "counter"

	for key, val := range allMetrics.Counter {
		metricsSData := serialize.Metrics{ID: key, MType: typeMetric, Delta: &val, Value: nil}
		sendData := serializer.SerializerResponse(&metricsSData)
		/*if _, err := f.WriteString(fmt.Sprintf("%s\n", *sendData)); err != nil {
			log.Println(err)
		}*/
		//curcurrentTime := time.Now()
		//timeAndDate := fmt.Sprintf("YYYY-MM-DD hh:mm:ss : ", curcurrentTime.Format("2006-01-02 15:04:05"))
		//if _, err := wf.f.WriteString(fmt.Sprintf("%s<=>%s\n", timeAndDate, *sendData)); err != nil {
		if _, err := wf.f.WriteString(fmt.Sprintf("%s\n", *sendData)); err != nil {
			log.Println(err)
		}
	}

}

// Поднимаем метрики из файл и пишем в хранилище
func (wf *workerFile) GetAllMetricsByFile() {

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
			wf.metricsRepository.UpdateMetricGauge(metrics.ID, *metrics.Value)
		case "counter":
			wf.metricsRepository.UpdateMetricCounter(metrics.ID, *metrics.Delta)
		default:
			fmt.Println("Нет такого типа метрики!")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
		return
	}

}
