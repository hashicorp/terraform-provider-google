package google

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Deprecated: For backward compatibility DefaultRequestTimeout is still working,
// but all new code should use DefaultRequestTimeout in the transport_tpg package instead.
var DefaultRequestTimeout = transport_tpg.DefaultRequestTimeout

// Deprecated: For backward compatibility SendRequest is still working,
// but all new code should use SendRequest in the transport_tpg package instead.
func SendRequest(config *transport_tpg.Config, method, project, rawurl, userAgent string, body map[string]interface{}, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) (map[string]interface{}, error) {
	return transport_tpg.SendRequest(config, method, project, rawurl, userAgent, body, errorRetryPredicates...)
}

// Deprecated: For backward compatibility SendRequestWithTimeout is still working,
// but all new code should use SendRequestWithTimeout in the transport_tpg package instead.
func SendRequestWithTimeout(config *transport_tpg.Config, method, project, rawurl, userAgent string, body map[string]interface{}, timeout time.Duration, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) (map[string]interface{}, error) {
	return transport_tpg.SendRequestWithTimeout(config, method, project, rawurl, userAgent, body, DefaultRequestTimeout, errorRetryPredicates...)
}

// Deprecated: For backward compatibility AddQueryParams is still working,
// but all new code should use AddQueryParams in the transport_tpg package instead.
func AddQueryParams(rawurl string, params map[string]string) (string, error) {
	return transport_tpg.AddQueryParams(rawurl, params)
}

// Deprecated: For backward compatibility handleNotFoundError is still working,
// but all new code should use HandleNotFoundError in the transport_tpg package instead.
func handleNotFoundError(err error, d *schema.ResourceData, resource string) error {
	return transport_tpg.HandleNotFoundError(err, d, resource)
}

// Deprecated: For backward compatibility IsGoogleApiErrorWithCode is still working,
// but all new code should use IsGoogleApiErrorWithCode in the transport_tpg package instead.
func IsGoogleApiErrorWithCode(err error, errCode int) bool {
	return transport_tpg.IsGoogleApiErrorWithCode(err, errCode)
}

// Deprecated: For backward compatibility isApiNotEnabledError is still working,
// but all new code should use IsApiNotEnabledError in the transport_tpg package instead.
func isApiNotEnabledError(err error) bool {
	return transport_tpg.IsApiNotEnabledError(err)
}
