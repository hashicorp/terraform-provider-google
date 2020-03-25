package google

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/api/googleapi"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

const testRetryTransportCodeRetry = 500
const testRetryTransportCodeSuccess = 200
const testRetryTransportCodeFailure = 400

func setUpRetryTransportServerClient(hf http.Handler) (*httptest.Server, *http.Client) {
	ts := httptest.NewServer(hf)

	client := ts.Client()
	client.Transport = &retryTransport{
		internal:        http.DefaultTransport,
		retryPredicates: []RetryErrorPredicateFunc{testRetryTransportRetryPredicate},
	}
	return ts, client
}

// Check for no errors if the request succeeds the first time
func TestRetryTransport_SingleRequestSuccess(t *testing.T) {
	ts, client := setUpRetryTransportServerClient(
		// Request succeeds immediately
		testRetryTransportHandler_noRetries(t, testRetryTransportCodeSuccess))
	defer ts.Close()

	resp, err := client.Get(ts.URL)
	testRetryTransport_checkSuccess(t, resp, err)
}

// Check for error if the request fails the first time
func TestRetryTransport_SingleRequestError(t *testing.T) {
	ts, client := setUpRetryTransportServerClient(
		// Request fails immediately
		testRetryTransportHandler_noRetries(t, testRetryTransportCodeFailure))
	defer ts.Close()

	resp, err := client.Get(ts.URL)
	testRetryTransport_checkFailure(t, resp, err, 400)
}

func TestRetryTransport_SuccessAfterRetries(t *testing.T) {
	ts, client := setUpRetryTransportServerClient(
		// Request succeeds after a certain amount of time
		testRetryTransportHandler_returnAfter(t, time.Second*1, testRetryTransportCodeSuccess))
	defer ts.Close()

	ctx, cc := context.WithTimeout(context.Background(), time.Second*2)
	defer cc()
	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("unable to construct err: %v", err)
	}

	resp, err := client.Do(req)
	testRetryTransport_checkSuccess(t, resp, err)
}

func TestRetryTransport_FailAfterRetries(t *testing.T) {
	ts, client := setUpRetryTransportServerClient(
		// Request fails after a certain amount of time
		testRetryTransportHandler_returnAfter(t, time.Second*1, testRetryTransportCodeFailure))
	defer ts.Close()

	ctx, cc := context.WithTimeout(context.Background(), time.Second*2)
	defer cc()
	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("unable to construct err: %v", err)
	}

	resp, err := client.Do(req)
	testRetryTransport_checkFailure(t, resp, err, 400)
}

func TestRetryTransport_ContextTimeout(t *testing.T) {
	ts, client := setUpRetryTransportServerClient(
		// Request succeeds after a certain amount of time
		testRetryTransportHandler_returnAfter(t, time.Second*4, testRetryTransportCodeSuccess))
	defer ts.Close()

	ctx, cc := context.WithTimeout(context.Background(), time.Second*2)
	defer cc()
	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("unable to construct err: %v", err)
	}

	resp, err := client.Do(req)
	// Last failure should have been a retryable error since we timed out
	testRetryTransport_checkFailedWhileRetrying(t, resp, err)
}

// Check for no errors if the request succeeds after a certain amount of time
func TestRetryTransport_SuccessWithBody(t *testing.T) {
	ts, client := setUpRetryTransportServerClient(
		// Request succeeds after a certain amount of time and returns the body
		testRetryTransportHandler_returnAfter(t, time.Second*1, testRetryTransportCodeSuccess))
	defer ts.Close()

	body := "body for successful request"
	ctx, cc := context.WithTimeout(context.Background(), time.Second*2)
	defer cc()
	req, err := http.NewRequestWithContext(ctx, "GET", ts.URL, bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatalf("unable to construct err: %v", err)
	}

	resp, err := client.Do(req)
	testRetryTransport_checkSuccess(t, resp, err)
	testRetryTransport_checkBody(t, resp, body)
}

// Check for no and no retries if the request has no getBody (it should only run once)
func TestRetryTransport_DoesNotRetryEmptyGetBody(t *testing.T) {
	msg := "non empty body"
	attempted := false

	ts, client := setUpRetryTransportServerClient(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Check for request body
			dump, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(testRetryTransportCodeFailure)
				if _, werr := w.Write([]byte(fmt.Sprintf("got error: %v", err))); werr != nil {
					t.Errorf("[ERROR] unable to write to response writer: %v", err)
				}
			}
			dumpS := string(dump)
			if dumpS != msg {
				w.WriteHeader(testRetryTransportCodeFailure)
				if _, werr := w.Write([]byte(fmt.Sprintf("got unexpected body: %s", dumpS))); werr != nil {
					t.Errorf("[ERROR] unable to write to response writer: %v", err)
				}
			}
			if attempted {
				w.WriteHeader(testRetryTransportCodeFailure)
				if _, werr := w.Write([]byte("expected only one try")); werr != nil {
					t.Errorf("[ERROR] unable to write to response writer: %v", err)
				}
			}
			attempted = true
			w.WriteHeader(testRetryTransportCodeRetry)
		}))
	defer ts.Close()

	// Create request
	req, err := http.NewRequest("GET", ts.URL, strings.NewReader(msg))
	if err != nil {
		t.Errorf("[ERROR] unable to make test request: %v", err)
	}
	// remove GetBody
	req.GetBody = nil

	resp, err := client.Do(req)
	testRetryTransport_checkFailedWhileRetrying(t, resp, err)
}

// handlers
func testRetryTransportHandler_noRetries(t *testing.T, code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		if _, err := w.Write([]byte(fmt.Sprintf("Code: %d", code))); err != nil {
			t.Errorf("[ERROR] unable to write to response writer: %v", err)
		}
	})
}

func testRetryTransportHandler_returnAfter(t *testing.T, interval time.Duration, code int) http.Handler {
	var firstReqTime time.Time
	var testOnce sync.Once

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testOnce.Do(func() {
			firstReqTime = time.Now()
		})

		var slurp []byte
		if r.Body != nil && r.Body != http.NoBody {
			var err error
			slurp, err = ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(testRetryTransportCodeFailure)
				if _, err := w.Write([]byte(fmt.Sprintf("unable to read request body: %v", err))); err != nil {
					t.Errorf("[ERROR] unable to write to response writer: %v", err)
				}
				return
			}
		}

		if time.Since(firstReqTime) < interval {
			w.WriteHeader(testRetryTransportCodeRetry)
			resp := fmt.Sprintf("Code: %d\nRequest Body: %s", testRetryTransportCodeRetry, string(slurp))
			if _, err := w.Write([]byte(resp)); err != nil {
				t.Errorf("[ERROR] unable to write to response writer: %v", err)
			}
			return
		}

		w.WriteHeader(code)
		resp := fmt.Sprintf("Code: %d\nRequest Body: %s", code, string(slurp))
		if _, err := w.Write([]byte(resp)); err != nil {
			t.Errorf("[ERROR] unable to write to response writer: %v", err)
		}
	})
}

// Utils for checking
func testRetryTransport_checkSuccess(t *testing.T, resp *http.Response, respErr error) {
	if respErr != nil {
		t.Fatalf("expected no error, got: %v", respErr)
	}

	err := googleapi.CheckResponse(resp)
	if err != nil {
		t.Fatalf("expected no error, got response error: %v", err)
	}

	if resp.StatusCode != testRetryTransportCodeSuccess {
		t.Fatalf("got unexpected error code %d, expected %d", resp.StatusCode, testRetryTransportCodeSuccess)
	}
}

func testRetryTransport_checkFailure(t *testing.T, resp *http.Response, respErr error, expectedCode int) {
	if respErr != nil {
		t.Fatalf("expected response error, got actual error for doing request: %v", respErr)
	}

	err := googleapi.CheckResponse(resp)
	if err == nil {
		t.Fatalf("expected googleapi error, got no error")
	}

	gerr, ok := err.(*googleapi.Error)
	if !ok {
		t.Fatalf("expected error to be googleapi error: %v", err)
	}

	if gerr.Code != expectedCode {
		t.Errorf("expected error code %d, got error: %v", expectedCode, err)
	}

	expectedMsg := fmt.Sprintf("Code: %d", expectedCode)
	if !strings.Contains(gerr.Body, expectedMsg) {
		t.Errorf("expected error message %q, got: %v", expectedMsg, err)
	}
}

func testRetryTransport_checkFailedWhileRetrying(t *testing.T, resp *http.Response, respErr error) {
	if respErr != nil {
		t.Fatalf("expected response error, got actual error for doing request: %v", respErr)
	}

	err := googleapi.CheckResponse(resp)
	if err == nil {
		t.Fatalf("expected googleapi error, got no error")
	}

	gerr, ok := err.(*googleapi.Error)
	if !ok {
		t.Fatalf("expected error to be googleapi error: %v", err)
	}

	if gerr.Code != testRetryTransportCodeRetry {
		t.Errorf("expected error code %d, got error: %v", testRetryTransportCodeRetry, err)
	}
}

func testRetryTransport_checkBody(t *testing.T, resp *http.Response, expectedMsg string) {
	if resp == nil {
		t.Fatal("expected non-empty response")
	}

	actualBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("expected no error, unable to read response body: %v", err)
	}

	expectedBody := fmt.Sprintf("Request Body: %s", expectedMsg)
	if !strings.HasSuffix(string(actualBody), expectedBody) {
		t.Fatalf(expectedBody)
	}
}

// ERROR RETRY PREDICATE
// Retries 500.
func testRetryTransportRetryPredicate(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == testRetryTransportCodeRetry {
			return true, "retryable error"
		}
	}
	return false, ""
}
