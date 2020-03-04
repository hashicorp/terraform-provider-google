package google

import (
	"context"
	"fmt"
	"github.com/hashicorp/errwrap"
	"log"
	"sync"
	"time"
)

const defaultBatchSendIntervalSec = 3

// RequestBatcher keeps track of batched requests globally.
// It should be created at a provider level. In general, one
// should be created per service that requires batching to:
//   - prevent blocking batching for one service due to another,
//   - minimize the possibility of overlap in batchKey formats (see SendRequestWithTimeout)
type RequestBatcher struct {
	sync.Mutex

	*batchingConfig
	parentCtx context.Context
	batches   map[string]*startedBatch
	debugId   string
}

// These types are meant to be the public interface to batchers. They define
// batch data format and logic to send/combine batches, i.e. they require
// specific implementations per type of request.
type (
	// BatchRequest represents a single request to a global batcher.
	BatchRequest struct {
		// ResourceName represents the underlying resource for which
		// a request is made. Its format is determined by what SendF expects, but
		// typically should be the name of the parent GCP resource being changed.
		ResourceName string

		// Body is this request's data to be passed to SendF, and may be combined
		// with other bodies using CombineF.
		Body interface{}

		// CombineF function determines how to combine bodies from two batches.
		CombineF BatcherCombineFunc

		// SendF function determines how to actually send a batched request to a
		// third party service. The arguments given to this function are
		// (ResourceName, Body) where Body may have been combined with other request
		// Bodies.
		SendF BatcherSendFunc

		// ID for debugging request. This should be specific to a single request
		// (i.e. per Terraform resource)
		DebugId string
	}

	// BatcherCombineFunc is a function type for combine existing batches and additional batch data
	BatcherCombineFunc func(body interface{}, toAdd interface{}) (interface{}, error)

	// BatcherSendFunc is a function type for sending a batch request
	BatcherSendFunc func(resourceName string, body interface{}) (interface{}, error)
)

// batchResponse bundles an API response (data, error) tuple.
type batchResponse struct {
	body interface{}
	err  error
}

func (br *batchResponse) IsError() bool {
	return br.err != nil
}

// startedBatch refers to a registered batch to group batch requests coming in.
// The timer manages the time after which a given batch is sent.
type startedBatch struct {
	batchKey string

	// Combined Batch Request
	*BatchRequest

	// subscribers is a registry of the requests (batchSubscriber) combined into this batcher.

	subscribers []batchSubscriber

	timer *time.Timer
}

// batchSubscriber contains information required for a single request for a startedBatch.
type batchSubscriber struct {
	// singleRequest is the original request this subscriber represents
	singleRequest *BatchRequest

	// respCh is the channel created to communicate the result to a waiting goroutine.s
	respCh chan batchResponse
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

	// Start goroutine to managing stopping the batcher if the provider-level parent context is closed.
	go func(b *RequestBatcher) {
		// Block until parent context is closed
		<-b.parentCtx.Done()

		log.Printf("[DEBUG] parent context canceled, cleaning up batcher batches")
		b.stop()
	}(batcher)

	return batcher
}

func (b *RequestBatcher) stop() {
	b.Lock()
	defer b.Unlock()

	log.Printf("[DEBUG] Stopping batcher %q", b.debugId)
	for batchKey, batch := range b.batches {
		log.Printf("[DEBUG] Cancelling started batch for batchKey %q", batchKey)
		batch.timer.Stop()
		for _, l := range batch.subscribers {
			close(l.respCh)
		}
	}
}

// SendRequestWithTimeout is a blocking call for making a single request, run alone or as part of a batch.
// It manages registering the single request with the batcher and waiting on the result.
//
// Params:
// batchKey: A string to group batchable requests. It should be unique to the API request being sent, similar to
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
			return nil, errwrap.Wrapf(
				fmt.Sprintf("Request %q returned error: {{err}}", request.DebugId),
				resp.err)
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

	// Batch doesn't exist for given batch key - create a new batch.

	log.Printf("[DEBUG] Creating new batch %q from request %q", newRequest.DebugId, batchKey)

	// The calling goroutine will need a channel to wait on for a response.
	respCh := make(chan batchResponse, 1)
	sub := batchSubscriber{
		singleRequest: newRequest,
		respCh:        respCh,
	}

	// Create a new batch with copy of the given batch request.
	b.batches[batchKey] = &startedBatch{
		BatchRequest: &BatchRequest{
			ResourceName: newRequest.ResourceName,
			Body:         newRequest.Body,
			CombineF:     newRequest.CombineF,
			SendF:        newRequest.SendF,
			DebugId:      fmt.Sprintf("Combined batch for started batch %q", batchKey),
		},
		batchKey:    batchKey,
		subscribers: []batchSubscriber{sub},
	}

	// Start a timer to send the request
	b.batches[batchKey].timer = time.AfterFunc(b.sendAfter, func() {
		batch := b.popBatch(batchKey)
		if batch == nil {
			log.Printf("[ERROR] batch should have been added to saved batches - just run as single request %q", newRequest.DebugId)
			respCh <- newRequest.send()
			close(respCh)
		} else {
			b.sendBatchWithSingleRetry(batchKey, batch)
		}
	})

	return respCh, nil
}

func (b *RequestBatcher) sendBatchWithSingleRetry(batchKey string, batch *startedBatch) {
	log.Printf("[DEBUG] Sending batch %q combining %d requests)", batchKey, len(batch.subscribers))
	resp := batch.send()

	// If the batch failed and combines more than one request, retry each single request.
	if resp.IsError() && len(batch.subscribers) > 1 {
		log.Printf("[DEBUG] Batch failed with error: %v", resp.err)
		log.Printf("[DEBUG] Sending each request in batch separately")
		for _, sub := range batch.subscribers {
			log.Printf("[DEBUG] Retrying single request %q", sub.singleRequest.DebugId)
			singleResp := sub.singleRequest.send()
			log.Printf("[DEBUG] Retried single request %q returned response: %v", sub.singleRequest.DebugId, singleResp)

			if singleResp.IsError() {
				singleResp.err = errwrap.Wrapf(
					fmt.Sprintf("Batch request and retried single request %q both failed. Final error: {{err}}", sub.singleRequest.DebugId),
					singleResp.err)
			}
			sub.respCh <- singleResp
			close(sub.respCh)
		}
	} else {
		// Send result to all subscribers
		for _, sub := range batch.subscribers {
			sub.respCh <- resp
			close(sub.respCh)
		}
	}
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
	sub := batchSubscriber{
		singleRequest: newRequest,
		respCh:        respCh,
	}
	batch.subscribers = append(batch.subscribers, sub)
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
