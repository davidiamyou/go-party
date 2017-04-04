package party

type job struct {
	f      func()
	report exceptionReport
	done   chan bool
}

func (j job) notifyDone() {
	go func() {
		j.done <- true
	}()
}

func (j job) notifyPanic(r interface{}) {
	go func() {
		j.report.PanicChannel <- r
	}()
}

func (j job) notifyTimeout() {
	go func() {
		j.report.TimeoutChannel <- true
	}()
}

// Exception report tells the client about any abnormally that happens
type exceptionReport struct {
	PanicChannel   chan interface{}
	TimeoutChannel chan bool
}

func newExceptionReport() exceptionReport {
	return exceptionReport{
		PanicChannel:   make(chan interface{}, 1),
		TimeoutChannel: make(chan bool, 1),
	}
}

// A Worker executes the job
type worker struct {
	dispatcher *Dispatcher
	jobChannel chan job
	terminate  chan bool
}

func newWorker(dispatcher *Dispatcher) worker {
	return worker{
		dispatcher: dispatcher,
		jobChannel: make(chan job),
		terminate:  make(chan bool),
	}
}

func (w worker) start() {
	go func() {
		for {
			w.dispatcher.workers <- w

			select {
			case job := <-w.jobChannel:
				preventPanic(job)()

			case <-w.terminate:
				return
			}
		}
	}()
}

func (w worker) stop() {
	go func() {
		w.terminate <- true
	}()
}

func preventPanic(job job) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				job.notifyPanic(r)
			}
			job.notifyDone()
		}()
		job.f()
	}
}
