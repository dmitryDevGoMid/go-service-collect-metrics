package wpool

import (
	"context"
	"fmt"
)

type WorkerStack struct {
	workersCount int
	jobs         chan Job
	results      chan Result
	ListObjJobs  chan []Job
	RunJobs      chan bool
	Done         chan struct{}
}

func New(workCount int) WorkerStack {
	return WorkerStack{
		workersCount: workCount,
		// канал с задачами
		jobs: make(chan Job, workCount),
		// канал с результатами
		results:     make(chan Result),
		ListObjJobs: make(chan []Job),
		RunJobs:     make(chan bool),
		Done:        make(chan struct{}),
	}
}

// Запускаем стек работников для выполнения отправки данных
func (ws *WorkerStack) WorkerRun(ctx context.Context) {

	for {
		select {
		case _, ok := <-ws.RunJobs:
			if ok {
				for i := 0; i < ws.workersCount; i++ {
					go ws.workers(ctx, i, ws.jobs, ws.results)
				}
			}
		case <-ctx.Done():
			fmt.Println("===============Worker STOP All===============")
			//Закрываем каналы там где отправлял
			close(ws.Done)
			close(ws.results)
			close(ws.RunJobs)
			return
		}
	}

}

// func (ws WorkerStack) GenerateJob([]Job) {
func (ws *WorkerStack) GenerateJob(ctx context.Context) {

	for {
		select {
		case job, ok := <-ws.ListObjJobs:
			//fmt.Println("GenerateJob 1 Отправили в канал ws.jobs", job)
			if ok {
				for _, v := range job {
					//fmt.Println("GenerateJob 2 Отправили в канал ws.jobs", v)
					ws.jobs <- v
				}
			}
		case <-ctx.Done():
			fmt.Println("GenerateJob -> Stop")
			//Закрываем каналы там где отправляли
			close(ws.ListObjJobs)
			close(ws.jobs)
			return
			//default:
			//time.Sleep(time.Duration(1) * time.Millisecond)
		}
	}
}

// Возвращаем результаты работы (ответ сервера)
func (ws *WorkerStack) ListResults() chan Result {
	return ws.results
}

// func workersAnd(wg *sync.WaitGroup, id int, ctx context.Context, jobs chan Job, results chan Result) {
func (ws *WorkerStack) workers(ctx context.Context, id int, jobs chan Job, results chan Result) {
	// Выполняем работу слушаем каналы и отправляем результат (ответ сервера)
	for {
		select {
		case job, ok := <-jobs:
			if ok {
				results <- job.execution(ctx)
			}
		case <-ctx.Done():
			//выводим сообщение о том что завершили работу с выходом из воркера
			fmt.Println("Close Worker => ", id)
			return
			//default:
			//time.Sleep(time.Duration(1) * time.Millisecond)
		}
	}
}
