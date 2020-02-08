package google

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRequestBatcher_batchSingle(t *testing.T) {
	testBasicCountBatches(t, "test-single", 1)
}

func TestRequestBatcher_batchMultiple(t *testing.T) {
	testBasicCountBatches(t, "test-multiple", 10)
}

func TestRequestBatcher_disableBatching(t *testing.T) {
	testBatcher := NewRequestBatcher(
		"testBatcher",
		context.Background(),
		&batchingConfig{
			sendAfter:      time.Duration(1) * time.Second,
			enableBatching: false,
		})

	testCombine := func(currV interface{}, toAddV interface{}) (interface{}, error) {
		return currV.(int) + toAddV.(int), nil
	}

	testSendBatch := func(name string, body interface{}) (interface{}, error) {
		return fmt.Sprintf("%s: %d", name, body), nil
	}

	wg := sync.WaitGroup{}
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go func(idx int) {
			defer wg.Done()

			req := &BatchRequest{
				DebugId:      fmt.Sprintf("Test Single Requests #%d", idx),
				ResourceName: "testNoBatching",
				Body:         1,
				CombineF:     testCombine,
				SendF:        testSendBatch,
			}

			respV, err := testBatcher.SendRequestWithTimeout(
				"testDisableBatching", req, time.Duration(1)*time.Second)
			if err != nil {
				t.Errorf("got unexpected error %s", err)
			}
			resp, ok := respV.(string)
			if !ok {
				t.Errorf("test returned an non-string response: %v", resp)
			}
			if resp != "testNoBatching: 1" {
				t.Errorf("expected single request response, got %s", resp)
			}
		}(i)
	}
}

func TestRequestBatcher_errInCombine(t *testing.T) {
	testBatcher := NewRequestBatcher(
		"testBatcher",
		context.Background(),
		&batchingConfig{
			sendAfter:      time.Duration(5) * time.Second,
			enableBatching: true,
		})

	combineErrText := "this is an expected error in combine"
	testCombine := func(_ interface{}, _ interface{}) (interface{}, error) {
		return nil, errors.New(combineErrText)
	}

	// sendBatchF is no-op
	testSendBatch := func(_ string, _ interface{}) (interface{}, error) {
		return nil, nil
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	// First call should no-op.
	go func() {
		defer wg.Done()

		req := &BatchRequest{
			DebugId:      "errInCombine first",
			ResourceName: "test-resource",
			Body:         nil,
			CombineF:     testCombine,
			SendF:        testSendBatch,
		}

		_, err := testBatcher.SendRequestWithTimeout("testCombineErr", req, time.Duration(10)*time.Second)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}
	}()

	// Second call should fail when being combined with original batch
	go func() {
		time.Sleep(time.Second)
		defer wg.Done()

		req := &BatchRequest{
			DebugId:      "errInCombine second",
			ResourceName: "test-resource",
			Body:         nil,
			CombineF:     testCombine,
			SendF:        testSendBatch,
		}

		_, err := testBatcher.SendRequestWithTimeout("testCombineErr", req, time.Duration(10)*time.Second)
		if err == nil {
			t.Errorf("expected error, got none")
		} else if !strings.Contains(err.Error(), combineErrText) {
			t.Errorf("error does not contain expected error %s. Got: %s", combineErrText, err)
		}
	}()

	wg.Wait()
}

func TestRequestBatcher_errInSend(t *testing.T) {
	testBatcher := NewRequestBatcher(
		"testBatcher",
		context.Background(),
		&batchingConfig{
			sendAfter:      time.Duration(5) * time.Second,
			enableBatching: true,
		})

	// combineF keeps track of the batched indexes
	testCombine := func(body interface{}, toAdd interface{}) (interface{}, error) {
		return append(body.([]int), toAdd.([]int)...), nil
	}

	failIdx := 0
	testResource := "RESOURCE-SEND-ERROR"
	expectedErrMsg := fmt.Sprintf("Error - batch %q contains idx %d", testResource, failIdx)

	testSendBatch := func(resourceName string, body interface{}) (interface{}, error) {
		log.Printf("[DEBUG] sendBatch body: %+v", body)
		for _, v := range body.([]int) {
			if v == failIdx {
				return nil, fmt.Errorf(expectedErrMsg)
			}
		}
		return nil, nil
	}

	numRequests := 3

	wg := sync.WaitGroup{}
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func(idx int) {
			defer wg.Done()

			req := &BatchRequest{
				DebugId:      fmt.Sprintf("sendError %d", idx),
				ResourceName: testResource,
				Body:         []int{idx},
				CombineF:     testCombine,
				SendF:        testSendBatch,
			}

			_, err := testBatcher.SendRequestWithTimeout("batchSendError", req, time.Duration(10)*time.Second)
			// Requests without index 0 should have succeeded
			if idx == failIdx {
				// We expect an error
				if err == nil {
					t.Errorf("expected error for request %d, got none", idx)
				}
				// Check error message
				expectedErrPrefix := "batch request and retry as single request failed - final error: "
				if !strings.Contains(err.Error(), expectedErrPrefix) {
					t.Errorf("expected error %q to contain %q", err, expectedErrPrefix)
				}
				if !strings.Contains(err.Error(), expectedErrMsg) {
					t.Errorf("expected error %q to contain %q", err, expectedErrMsg)
				}
			} else {

				// We shouldn't get error for non-failure index
				if err != nil {
					t.Errorf("expected request %d to succeed, got error: %v", idx, err)
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestRequestBatcher_errTimeout(t *testing.T) {
	testBatcher := NewRequestBatcher(
		"testBatcher",
		context.Background(),
		&batchingConfig{
			sendAfter:      time.Duration(5) * time.Second,
			enableBatching: true,
		})

	testResource := "resource for send error"

	// no-op
	testCombine := func(v interface{}, _ interface{}) (interface{}, error) {
		return v, nil
	}
	// no-op
	testSendBatch := func(resourceName string, cnt interface{}) (interface{}, error) {
		return nil, nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		req := &BatchRequest{
			DebugId:      fmt.Sprintf("timeout test"),
			ResourceName: testResource,
			Body:         1,
			CombineF:     testCombine,
			SendF:        testSendBatch,
		}

		_, err := testBatcher.SendRequestWithTimeout("batchTimeout", req, time.Duration(1)*time.Second)
		if err == nil {
			t.Errorf("expected error, got none")
		} else if !strings.Contains(err.Error(), "timed out") {
			t.Errorf("expected timeout error, got %v", err)
		}
	}()

	wg.Wait()
}

func testBasicCountBatches(t *testing.T, testName string, numBatches int) {
	testBatcher := NewRequestBatcher(
		"testBatcher",
		context.Background(),
		&batchingConfig{
			sendAfter:      time.Duration(1) * time.Second,
			enableBatching: true,
		})

	testCombine := func(currV interface{}, toAddV interface{}) (interface{}, error) {
		return currV.(int) + toAddV.(int), nil
	}

	testSendBatch := func(name string, body interface{}) (interface{}, error) {
		return fmt.Sprintf("%s: %d", name, body), nil
	}

	wg := sync.WaitGroup{}
	wg.Add(numBatches)

	for i := 0; i < numBatches; i++ {
		go func(idx int) {
			defer wg.Done()

			req := &BatchRequest{
				DebugId:      fmt.Sprintf("Test '%s' Request #%d", testName, idx),
				ResourceName: testName,
				Body:         1,
				CombineF:     testCombine,
				SendF:        testSendBatch,
			}

			respV, err := testBatcher.SendRequestWithTimeout("testBatching", req, time.Duration(6)*time.Second)
			if err != nil {
				t.Errorf("got unexpected error %s", err)
			}
			resp, ok := respV.(string)
			if !ok {
				t.Errorf("test returned an non-string response: %v", resp)
			}
			expected := fmt.Sprintf("%s: %d", testName, numBatches)
			if resp != expected {
				t.Errorf("expected response %s, got %s", expected, resp)
			}
		}(i)
	}
}
