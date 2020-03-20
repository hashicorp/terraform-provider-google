// A http.RoundTripper that retries common errors, with convenience constructors.
//
// NOTE: This meant for TEMPORARY, TRANSIENT ERRORS.
// Do not use for waiting on operations or polling of resource state,
// especially if the expected state (operation done, resource ready, etc)
// takes longer to reach than the default client Timeout.
// In those cases, retryTimeDuration(...)/resource.Retry with appropriate timeout
// and error predicates/handling should be used as a wrapper around the request
// instead.
//
// Example Usage:
// For handwritten/Go clients, the retry transport should be provided via
// the main client or a shallow copy of the HTTP resources, depending on the
// API-specific retry predicates.
// Example Usage in Terraform Config:
//	client := oauth2.NewClient(ctx, tokenSource)
//	// Create with default retry predicates
//	client.Transport := NewTransportWithDefaultRetries(client.Transport, defaultTimeout)
//
//	// If API uses just default retry predicates:
//	c.clientCompute, err = compute.NewService(ctx, option.WithHTTPClient(client))
//	...
//	// If API needs custom additional retry predicates:
//	sqlAdminHttpClient := ClientWithAdditionalRetries(client, retryTransport,
//			isTemporarySqlError1,
//			isTemporarySqlError2)
//	c.clientSqlAdmin, err = compute.NewService(ctx, option.WithHTTPClient(sqlAdminHttpClient))
// ...

package google

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"google.golang.org/api/googleapi"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

const defaultRetryTransportTimeoutSec = 30

// NewTransportWithDefaultRetries constructs a default retryTransport that will retry common temporary errors
func NewTransportWithDefaultRetries(t http.RoundTripper) *retryTransport {
	return &retryTransport{
		retryPredicates: defaultErrorRetryPredicates,
		internal:        t,
	}
}

// Helper method to create a shallow copy of an HTTP client with a shallow-copied retryTransport
// s.t. the base HTTP transport is the same (i.e. client connection pools are shared, retryPredicates are different)
func ClientWithAdditionalRetries(baseClient *http.Client, baseRetryTransport *retryTransport, predicates ...RetryErrorPredicateFunc) *http.Client {
	copied := *baseClient
	if baseRetryTransport == nil {
		baseRetryTransport = NewTransportWithDefaultRetries(baseClient.Transport)
	}
	copied.Transport = baseRetryTransport.WithAddedPredicates(predicates...)
	return &copied
}

// Returns a shallow copy of the retry transport with additional retry
// predicates but same wrapped http.RoundTripper
func (t *retryTransport) WithAddedPredicates(predicates ...RetryErrorPredicateFunc) *retryTransport {
	copyT := *t
	copyT.retryPredicates = append(t.retryPredicates, predicates...)
	return &copyT
}

type retryTransport struct {
	retryPredicates []RetryErrorPredicateFunc
	internal        http.RoundTripper
}

// RoundTrip implements the RoundTripper interface method.
// It retries the given HTTP request based on the retry predicates
// registered under the retryTransport.
func (t *retryTransport) RoundTrip(req *http.Request) (resp *http.Response, respErr error) {
	// Set timeout to default value.
	ctx := req.Context()
	var ccancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		ctx, ccancel = context.WithTimeout(ctx, defaultRetryTransportTimeoutSec*time.Second)
		defer func() {
			if ctx.Err() == nil {
				// Cleanup child context created for retry loop if ctx not done.
				ccancel()
			}
		}()
	}

	attempts := 0
	backoff := time.Millisecond * 500
	nextBackoff := time.Millisecond * 500

	log.Printf("[DEBUG] Retry Transport: starting RoundTrip retry loop")
Retry:
	for {
		log.Printf("[DEBUG] Retry Transport: request attempt %d", attempts)

		// Copy the request - we dont want to use the original request as
		// RoundTrip contract says request body can/will be consumed
		newRequest, copyErr := copyHttpRequest(req)
		if copyErr != nil {
			respErr = errwrap.Wrapf("unable to copy invalid http.Request for retry: {{err}}", copyErr)
			break Retry
		}

		// Do the wrapped Roundtrip. This is one request in the retry loop.
		resp, respErr = t.internal.RoundTrip(newRequest)
		attempts++

		retryErr := t.checkForRetryableError(resp, respErr)
		if retryErr == nil {
			log.Printf("[DEBUG] Retry Transport: Stopping retries, last request was successful")
			break Retry
		}
		if !retryErr.Retryable {
			log.Printf("[DEBUG] Retry Transport: Stopping retries, last request failed with non-retryable error: %s", retryErr.Err)
			break Retry
		}

		log.Printf("[DEBUG] Retry Transport: Waiting %s before trying request again", backoff)
		select {
		case <-ctx.Done():
			log.Printf("[DEBUG] Retry Transport: Stopping retries, context done: %v", ctx.Err())
			break Retry
		case <-time.After(backoff):
			log.Printf("[DEBUG] Retry Transport: Finished waiting %s before next retry", backoff)

			// Fibonnaci backoff - 0.5, 1, 1.5, 2.5, 4, 6.5, 10.5, ...
			lastBackoff := backoff
			backoff = backoff + nextBackoff
			nextBackoff = lastBackoff
			continue
		}
	}

	// VCR depends on the original request body being consumed, so consume it here
	_, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Printf("[DEBUG] Retry Transport: Reading request failed: %v", err)
	}

	log.Printf("[DEBUG] Retry Transport: Returning after %d attempts", attempts)
	return resp, respErr
}

// copyHttpRequest provides an copy of the given HTTP request for one RoundTrip.
// If the request has a non-empty body (io.ReadCloser), the body is deep copied
// so it can be consumed.
func copyHttpRequest(req *http.Request) (*http.Request, error) {
	newRequest := *req

	if req.Body == nil || req.Body == http.NoBody {
		return &newRequest, nil
	}

	// Helpers like http.NewRequest add a GetBody for copying.
	// If not given, we should reject the request.
	if req.GetBody == nil {
		return nil, errors.New("invalid HTTP request for transport, expected request.GetBody for non-empty Body")
	}

	bd, err := req.GetBody()
	if err != nil {
		return nil, err
	}

	newRequest.Body = bd
	return &newRequest, nil
}

// checkForRetryableError uses the googleapi.CheckResponse util to check for
// errors in the response, and determines whether there is a retryable error.
// in response/response error.
func (t *retryTransport) checkForRetryableError(resp *http.Response, respErr error) *resource.RetryError {
	var errToCheck error

	if respErr != nil {
		errToCheck = respErr
	} else {
		respToCheck := *resp
		// The RoundTrip contract states that the HTTP response/response error
		// returned cannot be edited. We need to consume the Body to check for
		// errors, so we need to create a copy if the Response has a body.
		if resp.Body != nil && resp.Body != http.NoBody {
			// Use httputil.DumpResponse since the only important info is
			// error code and messages in the response body.
			dumpBytes, err := httputil.DumpResponse(resp, true)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("unable to check response for error: %v", err))
			}
			respToCheck.Body = ioutil.NopCloser(bytes.NewReader(dumpBytes))
		}
		errToCheck = googleapi.CheckResponse(&respToCheck)
	}

	if isRetryableError(errToCheck, t.retryPredicates...) {
		return resource.RetryableError(errToCheck)
	}
	return resource.NonRetryableError(errToCheck)
}
