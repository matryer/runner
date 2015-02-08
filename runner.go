package runner

import (
	"sync"
)

// S is a function that will return true if the
// code should stop.
type S func() bool

// Routine is a goroutine capable of calling the
// shouldStop S function.
type Routine func(S) error

// Go executes the function in a goroutine and returns a
// Task capable of stopping the execution.
func Go(fn Routine) *Task {
	t := &Task{
		stopChan: make(chan struct{}),
		running:  true,
	}
	go func() {
		// call the target function
		err := fn(func() bool {
			// this is the shouldStop() function available to the
			// target function
			t.lock.RLock()
			shouldStop := t.shouldStop
			t.lock.RUnlock()
			return shouldStop
		})
		// stopped
		t.lock.Lock()
		t.err = err
		t.running = false
		close(t.stopChan)
		t.lock.Unlock()
	}()
	return t
}

// Task represents an interruptable goroutine.
type Task struct {
	lock       sync.RWMutex
	stopChan   chan struct{}
	shouldStop bool
	running    bool
	err        error
}

// Stop tells the goroutine to stop.
func (t *Task) Stop() {
	t.shouldStop = true
}

// StopChan gets the channel that will be closed when
// the task has finished.
func (t *Task) StopChan() <-chan struct{} {
	return t.stopChan
}

// Running gets whether the goroutine is
// running or not.
func (t *Task) Running() bool {
	t.lock.RLock()
	running := t.running
	t.lock.RUnlock()
	return running
}

// Err gets the error returned by the goroutine.
func (t *Task) Err() error {
	t.lock.RLock()
	err := t.err
	t.lock.RUnlock()
	return err
}
