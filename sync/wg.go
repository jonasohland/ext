package sync

import (
	"fmt"
	"sync"
	"time"
)

type TimeoutError struct {
	Time time.Duration
}

func (t *TimeoutError) Error() string {
	return fmt.Sprintf("timed out after %s", t.Time)
}

func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) error {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-time.After(timeout):
		return &TimeoutError{timeout}
	case <-done:
		return nil
	}
}
