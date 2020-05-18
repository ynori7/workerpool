package workerpool

import (
	"context"
	"fmt"
	"reflect"
)

type (
	successFunc func(result interface{})
	errorFunc   func(err error)
	doWorkFunc  func(job interface{}) (result interface{}, err error)
)

// WorkerPool abstracts the setup around creating worker pools.
type WorkerPool struct {
	onSuccess successFunc
	onError   errorFunc
	doWork    doWorkFunc
}

// NewWorkerPool creates a new WorkerPool instance with the given onSuccess, onError, and doWork callbacks.
func NewWorkerPool(
	onSuccess successFunc,
	onError errorFunc,
	doWork doWorkFunc,
) *WorkerPool {
	return &WorkerPool{
		onSuccess: onSuccess,
		onError:   onError,
		doWork:    doWork,
	}
}

// Work spawns the workers and creates the concurrency control channels, and then distributes the given jobs to each worker.
// When the given context is canceled, the work will be halted. An error is returned if the given jobSlice is not a slice.
func (w *WorkerPool) Work(ctx context.Context, workerCount int, jobsSlice interface{}) error {
	//validate input
	jobs, err := interfaceToSlice(jobsSlice)
	if err != nil {
		return err
	}
	if workerCount < 1 {
		return fmt.Errorf("there must be at least one worker")
	}
	jobCount := len(jobs)
	if jobCount < 1 {
		return nil // no work to do, so just return
	}

	resultsChan := make(chan interface{}, jobCount)
	errorChan := make(chan error, jobCount)

	//Spawn workers to process in parallel
	workers := make([]chan interface{}, workerCount)
	for i := 0; i < workerCount; i++ {
		workers[i] = make(chan interface{}, len(jobs)/workerCount)
		go w.worker(resultsChan, errorChan, workers[i])
	}

	//Assign an equal number of releases to be checked by each worker
	var i = 0
	for _, s := range jobs {
		workers[i] <- s
		i = (i + 1) % workerCount
	}

	//Process results
WORK:
	for i := 0; i < len(jobs); i++ {
		select {
		case <-ctx.Done():
			break WORK //Stop processing. The workers will all be closed
		case r := <-resultsChan:
			w.onSuccess(r)
		case err := <-errorChan:
			w.onError(err)
		}
	}

	//Signal workers to stop working
	for _, worker := range workers {
		close(worker)
	}

	return nil
}

func (w *WorkerPool) worker(successes chan interface{}, errors chan error, jobs chan interface{}) {
	for j := range jobs {
		res, err := w.doWork(j)
		if err != nil {
			errors <- err
		} else {
			successes <- res
		}
	}
}

func interfaceToSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input is not a slice")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}
