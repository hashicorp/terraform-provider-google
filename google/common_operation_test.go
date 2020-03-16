package google

import (
	"net/url"
	"testing"
	"time"
)

type TestWaiter struct {
	runCount int
}

func (w *TestWaiter) State() string {
	if w.runCount == 2 {
		return "DONE"
	}
	return "RUNNING"
}

func (TestWaiter) IsRetryable(err error) bool {
	return false
}

func (TestWaiter) Error() error {
	return nil
}

func (TestWaiter) SetOp(interface{}) error {
	return nil
}

func (w *TestWaiter) QueryOp() (interface{}, error) {
	w.runCount++
	if w.runCount == 1 {
		return nil, &url.Error{
			Err: &TimeoutError{timeout: true},
		}
	}
	return "my return value", nil
}

func (TestWaiter) OpName() string {
	return "my-operation-name"
}

func (TestWaiter) PendingStates() []string {
	return []string{}
}

func (TestWaiter) TargetStates() []string {
	return []string{"DONE"}
}

func TestOperationWait_TimeoutsShouldRetry(t *testing.T) {
	testWaiter := TestWaiter{
		runCount: 0,
	}
	err := OperationWait(&testWaiter, "my-activity", 1, 0*time.Second)
	if err != nil {
		t.Fatalf("unexpected error waiting for operation: got '%v', want 'nil'", err)
	}
	expectedRunCount := 2
	if testWaiter.runCount != expectedRunCount {
		t.Errorf("expected the retryFunc to be called %v time(s), instead was called %v time(s)",
			expectedRunCount, testWaiter.runCount)
	}
}
