package wpool

// Обьявляем структуру
type WorkerPool struct {
	maxWorkers int
	tasks      chan func()
}

// Создаем структуру и возвращаем указатель на нее
func NewWorkerPool(maxWorkers int) *WorkerPool {
	w := &WorkerPool{
		//Кол-во воркеров
		maxWorkers: maxWorkers,
		//Канал с тасками поступающие на исполнение
		tasks: make(chan func()),
	}
	//Запускаем горотины с работниками
	w.run()

	return w
}

// Запускаем воркеры
func (w *WorkerPool) run() {
	for i := 0; i < w.maxWorkers; i++ {
		go func() {
			for task := range w.tasks {
				task()
			}
		}()
	}
}

// Передаем таски на выполнение в канал
func (w *WorkerPool) AddTask(task func(), i int) {
	w.tasks <- task
}
