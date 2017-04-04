package party

import "time"

type Dispatcher struct {
	workers  chan worker
	timeout  time.Duration
	jobQueue chan job
}

func NewDispatcher(numWorkers int, queueSize int, timeout time.Duration) *Dispatcher {
	workers := make(chan worker, numWorkers)
	jobQueue := make(chan job, queueSize)
	return &Dispatcher{
		workers:  workers,
		timeout:  timeout,
		jobQueue: jobQueue,
	}
}

func (d *Dispatcher) Start() {
	for i := 0; i < cap(d.workers); i++ {
		w := newWorker(d)
		w.start()
	}

	go d.dispatch()
}

func (d *Dispatcher) Run(task func()) exceptionReport {
	er := newExceptionReport()
	go func() {
		d.jobQueue <- job{f: task, report: er, done: make(chan bool, 1)}
	}()
	return er
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case j0 := <-d.jobQueue:
			go func(j job) {
				w := <-d.workers
				w.jobChannel <- j
				select {
				case <-j.done:
					return
				case <-time.After(d.timeout):
					w.stop()
					j.notifyTimeout()
					return
				}
			}(j0)
		}
	}
}
