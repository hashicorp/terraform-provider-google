// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package transport

import (
	"net/http"
)

// adapted from https://stackoverflow.com/questions/51325704/adding-a-default-http-header-in-go
type headerTransportLayer struct {
	http.Header
	baseTransit http.RoundTripper
}

func NewTransportWithHeaders(baseTransit http.RoundTripper) headerTransportLayer {
	if baseTransit == nil {
		baseTransit = http.DefaultTransport
	}

	headers := make(http.Header)

	return headerTransportLayer{Header: headers, baseTransit: baseTransit}
}

func (h headerTransportLayer) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range h.Header {
		// only set headers that are not previously defined
		if _, ok := req.Header[key]; !ok {
			req.Header[key] = value
		}
	}
	return h.baseTransit.RoundTrip(req)
}
