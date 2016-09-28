package runner

import (
	"sync"
)

// S is a function that will return true if the
// goroutine should stop executing.
type S func() bool

// Go executes the function in a goroutine and returns a
// Task capable of stopping the execution.
func Go(fn func(S) error) *Task {
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
	// When task is stopped from a different go-routine other than the one
	// that actually started it.
	t.lock.Lock()
	t.shouldStop = true
	t.lock.Unlock()
}

// StopChan gets the stop channel for this task.
// Reading from this channel will block while the task is running, and will
// unblock once the task has stopped (because the channel gets closed).
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
