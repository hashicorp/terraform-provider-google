package google

import (
	"net/http"
	"os"
)

// adapted from https://stackoverflow.com/questions/51325704/adding-a-default-http-header-in-go
type headerTransportLayer struct {
	http.Header
	baseTransit http.RoundTripper
}

func newTransportWithHeaders(baseTransit http.RoundTripper, c *Config) headerTransportLayer {
	if baseTransit == nil {
		baseTransit = http.DefaultTransport
	}

	headers := make(http.Header)

	if requestReason := os.Getenv("CLOUDSDK_CORE_REQUEST_REASON"); requestReason != "" {
		headers.Set("X-Goog-Request-Reason", requestReason)
	}

	// Ensure $userProject is set for all HTTP requests using the client if specified by the provider config
	// See https://cloud.google.com/apis/docs/system-parameters
	if c.UserProjectOverride && c.BillingProject != "" {
		headers.Set("X-Goog-User-Project", c.BillingProject)
	}

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
