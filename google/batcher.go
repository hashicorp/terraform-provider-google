package google

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
)

const defaultBatchSendIntervalSec = 3

// RequestBatcher is a global batcher object that keeps track of
// existing batches.
// In general, a batcher should be created per service that requires batching
// in order to prevent blocking batching for one service due to another,
// and to minimize the possibility of overlap in batchKey formats
// (see SendRequestWithTimeout)
type RequestBatcher struct {
	sync.Mutex

	*batchingConfig
	parentCtx context.Context
	batches   map[string]*startedBatch
	debugId   string
}

// BatchRequest represents a single request to a global batcher.
type BatchRequest struct {
	// ResourceName represents the underlying resource for which
	// a request is made. Its format is determined by what SendF expects, but
	// typically should be the name of the parent GCP resource being changed.
	ResourceName string

	// Body is this request's data to be passed to SendF, and may be combined
	// with other bodies using CombineF.
	Body interface{}

	// CombineF function determines how to combine bodies from two batches.
	CombineF batcherCombineFunc

	// SendF function determines how to actually send a batched request to a
	// third party service. The arguments given to this function are
	// (ResourceName, Body) where Body may have been combined with other request
	// Bodies.
	SendF batcherSendFunc

	// ID for debugging request. This should be specific to a single request
	// (i.e. per Terraform resource)
	DebugId string
}

// These types are meant to be the public interface to batchers. They define
// logic to manage batch data type and behavior, and require service-specific
// implementations per type of request per service.
// Function type for combine existing batches and additional batch data
type batcherCombineFunc func(body interface{}, toAdd interface{}) (interface{}, error)

// Function type for sending a batch request
type batcherSendFunc func(resourceName string, body interface{}) (interface{}, error)

// batchResponse bundles an API response (data, error) tuple.
type batchResponse struct {
	body interface{}
	err  error
}

// startedBatch refers to a processed batch whose timer to send the request has
// already been started. The responses for the request is sent to each listener
// channel, representing parallel callers that are waiting on requests
// combined into this batch.
type startedBatch struct {
	batchKey string
	*BatchRequest

	listeners []chan batchResponse
	timer     *time.Timer
}

// batchingConfig contains user configuration for controlling batch requests.
type batchingConfig struct {
	sendAfter      time.Duration
	enableBatching bool
}

// Initializes a new batcher.
func NewRequestBatcher(debugId string, ctx context.Context, config *batchingConfig) *RequestBatcher {
	batcher := &RequestBatcher{
		debugId:        debugId,
		parentCtx:      ctx,
		batchingConfig: config,
		batches:        make(map[string]*startedBatch),
	}

	go func(b *RequestBatcher) {
		<-ctx.Done()
		b.stop()
	}(batcher)

	return batcher
}

func (b *RequestBatcher) stop() {
	b.Lock()
	defer b.Unlock()

	log.Printf("[DEBUG] Stopping batcher %q", b.debugId)
	for batchKey, batch := range b.batches {
		log.Printf("[DEBUG] Cleaning up batch request %q", batchKey)
		batch.timer.Stop()
		for _, l := range batch.listeners {
			close(l)
		}
	}
}

// SendRequestWithTimeout is expected to be called per parallel call.
// It manages waiting on the result of a batch request.
//
// Batch requests are grouped by the given batchKey. batchKey
// should be unique to the API request being sent, most likely similar to
// the HTTP request URL with GCP resource ID included in the URL (the caller
// may choose to use a key with method if needed to diff GET/read and
// POST/create)
//
// As an example, for google_project_service, the
// batcher is called to batch services.batchEnable() calls for a project
// $PROJECT. The calling code uses the template
// "serviceusage:projects/$PROJECT/services:batchEnable", which mirrors the HTTP request:
// POST https://serviceusage.googleapis.com/v1/projects/$PROJECT/services:batchEnable
func (b *RequestBatcher) SendRequestWithTimeout(batchKey string, request *BatchRequest, timeout time.Duration) (interface{}, error) {
	if request == nil {
		return nil, fmt.Errorf("error, cannot request batching for nil BatchRequest")
	}
	if request.CombineF == nil {
		return nil, fmt.Errorf("error, cannot request batching for BatchRequest with nil CombineF")
	}
	if request.SendF == nil {
		return nil, fmt.Errorf("error, cannot request batching for BatchRequest with nil SendF")
	}
	if !b.enableBatching {
		log.Printf("[DEBUG] Batching is disabled, sending single request for %q", request.DebugId)
		return request.SendF(request.ResourceName, request.Body)
	}

	respCh, err := b.registerBatchRequest(batchKey, request)
	if err != nil {
		return nil, fmt.Errorf("error adding request to batch: %s", err)
	}

	ctx, cancel := context.WithTimeout(b.parentCtx, timeout)
	defer cancel()

	select {
	case resp := <-respCh:
		if resp.err != nil {
			// use wrapf so we can potentially extract the original error type
			return nil, errwrap.Wrapf(fmt.Sprintf("Batch %q for request %q returned error: {{err}}", batchKey, request.DebugId), resp.err)
		}
		return resp.body, nil
	case <-ctx.Done():
		break
	}
	return nil, fmt.Errorf("Request %s timed out after %v", batchKey, timeout)
}

// registerBatchRequest safely sees if an existing batch has been started
// with the given batchKey. If a batch exists, this will combine the new
// request into this existing batch. Else, this method manages starting a new
// batch and adding it to the RequestBatcher's started batches.
func (b *RequestBatcher) registerBatchRequest(batchKey string, newRequest *BatchRequest) (<-chan batchResponse, error) {
	b.Lock()
	defer b.Unlock()

	// If batch already exists, combine this request into existing request.
	if batch, ok := b.batches[batchKey]; ok {
		return batch.addRequest(newRequest)
	}

	log.Printf("[DEBUG] Creating new batch %q from request %q", newRequest.DebugId, batchKey)
	// The calling goroutine will need a channel to wait on for a response.
	respCh := make(chan batchResponse, 1)

	// Create a new batch.
	b.batches[batchKey] = &startedBatch{
		BatchRequest: newRequest,
		batchKey:     batchKey,
		listeners:    []chan batchResponse{respCh},
	}

	// Start a timer to send the request
	b.batches[batchKey].timer = time.AfterFunc(b.sendAfter, func() {
		batch := b.popBatch(batchKey)

		var resp batchResponse
		if batch == nil {
			log.Printf("[DEBUG] Batch not found in saved batches, running single request batch %q", batchKey)
			resp = newRequest.send()
		} else {
			log.Printf("[DEBUG] Sending batch %q combining %d requests)", batchKey, len(batch.listeners))
			resp = batch.send()
		}

		// Send message to all goroutines waiting on result.
		for _, ch := range batch.listeners {
			ch <- resp
			close(ch)
		}
	})

	return respCh, nil
}

// popBatch safely gets and removes a batch with given batchkey from the
// RequestBatcher's started batches.
func (b *RequestBatcher) popBatch(batchKey string) *startedBatch {
	b.Lock()
	defer b.Unlock()

	batch, ok := b.batches[batchKey]
	if !ok {
		log.Printf("[DEBUG] Batch with ID %q not found in batcher", batchKey)
		return nil
	}

	delete(b.batches, batchKey)
	return batch
}

func (batch *startedBatch) addRequest(newRequest *BatchRequest) (<-chan batchResponse, error) {
	log.Printf("[DEBUG] Adding batch request %q to existing batch %q", newRequest.DebugId, batch.batchKey)
	if batch.CombineF == nil {
		return nil, fmt.Errorf("Provider Error: unable to add request %q to batch %q with no CombineF", newRequest.DebugId, batch.batchKey)
	}
	newBody, err := batch.CombineF(batch.Body, newRequest.Body)
	if err != nil {
		return nil, fmt.Errorf("Provider Error: Unable to combine request %q data into existing batch %q: %v", newRequest.DebugId, batch.batchKey, err)
	}
	batch.Body = newBody

	log.Printf("[DEBUG] Added batch request %q to batch. New batch body: %v", newRequest.DebugId, batch.Body)

	respCh := make(chan batchResponse, 1)
	batch.listeners = append(batch.listeners, respCh)
	return respCh, nil
}

func (req *BatchRequest) send() batchResponse {
	if req.SendF == nil {
		return batchResponse{
			err: fmt.Errorf("provider error: Batch request has no SendBatch function"),
		}
	}
	v, err := req.SendF(req.ResourceName, req.Body)
	return batchResponse{v, err}
}
