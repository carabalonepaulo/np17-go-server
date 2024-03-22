package main

const WorkerTaskBufferSize = 1024

type Task struct {
	fn   func()
	sync func()
}

type TaskOutput func()

type WorkerPool struct {
	input  chan Task
	output chan TaskOutput
}

func NewWorkerPool(cap int) *WorkerPool {
	w := &WorkerPool{
		input:  make(chan Task, WorkerTaskBufferSize),
		output: make(chan TaskOutput, WorkerTaskBufferSize),
	}
	for i := 0; i < cap; i++ {
		go worker(w.input, w.output)
	}
	return w
}

/*
This function serves as a mechanism to execute potentially
blocking operations asynchronously in a separate worker pool while maintaining
responsiveness in the main thread. It takes two functions as arguments: the first
function is executed in a worker pool, allowing for blocking interactions such as
database or file system operations, while the second function is executed in the
main thread upon completion of the first function. This approach enables tasks
such as user authentication to be handled asynchronously without blocking the
main thread, enhancing the overall responsiveness and performance of the application.
*/
func DoAndSync[T any](w *WorkerPool, fn func() T, sync func(T)) {
	temp := new(T)
	w.input <- Task{
		fn:   func() { *temp = fn() },
		sync: func() { sync(*temp) },
	}
}

func (w *WorkerPool) Poll() {
	select {
	case r := <-w.output:
		r()
	default:
	}
}

func worker(input chan Task, output chan TaskOutput) {
	for task := range input {
		task.fn()
		output <- task.sync
	}
}
