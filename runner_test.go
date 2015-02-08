package runner_test

import (
	"errors"
	"testing"
	"time"

	"github.com/cheekybits/is"
	"github.com/matryer/runner"
)

func TestRun(t *testing.T) {
	is := is.New(t)
	var ticker []time.Time
	task := runner.Go(func(shouldStop runner.S) error {
		for {
			ticker = append(ticker, time.Now())
			time.Sleep(100 * time.Millisecond)
			if shouldStop() {
				break
			}
		}
		return nil
	})
	is.Equal(true, task.Running())
	time.Sleep(1 * time.Second)
	task.Stop()
	select {
	case <-task.StopChan():
	case <-time.After(2 * time.Second):
		is.Fail("timed out")
	}
	is.Equal(false, task.Running())
	is.Equal(10, len(ticker))
}

func TestRunErr(t *testing.T) {
	is := is.New(t)

	err := errors.New("something went wrong")
	task := runner.Go(func(shouldStop runner.S) error {
		return err
	})

	time.Sleep(100 * time.Millisecond)
	is.Equal(false, task.Running())
	is.Equal(err, task.Err())

	task.Stop()
	select {
	case <-task.StopChan():
	case <-time.After(2 * time.Second):
		is.Fail("timed out")
	}

}
