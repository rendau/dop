package dopTools

import (
	"context"
	"sync"
	"time"
)

type workerPoolTaskT func(context.Context) error

type WorkerPool struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	taskChan   chan workerPoolTaskT
	finishChan chan error
	startTime  time.Time
	duration   time.Duration
}

func NewWorkerPool(
	ctx context.Context,
	workerCount int,
	bufferSize int,
) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)

	wp := &WorkerPool{
		ctx:        ctx,
		cancel:     cancel,
		wg:         &sync.WaitGroup{},
		taskChan:   make(chan workerPoolTaskT, bufferSize),
		finishChan: make(chan error, 1),
		startTime:  time.Now(),
	}

	// start workers
	for i := 0; i < workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}

	return wp
}

func (w *WorkerPool) worker() {
	defer w.wg.Done()

	var err error
	var task workerPoolTaskT
	var ok bool

	for {
		select {
		case <-w.ctx.Done():
			return
		default:
		}

		select {
		case <-w.ctx.Done():
			return
		case task, ok = <-w.taskChan:
			if !ok {
				return
			}

			err = task(w.ctx)
			if err != nil {
				w.Cancel(err)
				return
			}
		}
	}
}

func (w *WorkerPool) Submit(task workerPoolTaskT) bool {
	select {
	case <-w.ctx.Done():
		return false
	default:
	}

	select {
	case <-w.ctx.Done():
		return false
	case w.taskChan <- task:
		return true
	}
}

func (w *WorkerPool) Finish() {
	close(w.taskChan)
}

func (w *WorkerPool) Cancel(err error) {
	w.cancel()
	w.finishChan <- err
}

func (w *WorkerPool) Wait() error {
	w.wg.Wait()
	w.duration = time.Since(w.startTime)
	close(w.finishChan)
	return <-w.finishChan
}

func (w *WorkerPool) FinishAndWait() error {
	w.Finish()
	return w.Wait()
}

func (w *WorkerPool) GetDuration() time.Duration {
	return w.duration
}
