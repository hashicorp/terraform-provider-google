package google

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var DefaultRequestTimeout = transport_tpg.DefaultRequestTimeout

func SendRequest(config *transport_tpg.Config, method, project, rawurl, userAgent string, body map[string]interface{}, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) (map[string]interface{}, error) {
	return transport_tpg.SendRequest(config, method, project, rawurl, userAgent, body, errorRetryPredicates...)
}

func SendRequestWithTimeout(config *transport_tpg.Config, method, project, rawurl, userAgent string, body map[string]interface{}, timeout time.Duration, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) (map[string]interface{}, error) {
	return transport_tpg.SendRequestWithTimeout(config, method, project, rawurl, userAgent, body, DefaultRequestTimeout, errorRetryPredicates...)
}

func AddQueryParams(rawurl string, params map[string]string) (string, error) {
	return transport_tpg.AddQueryParams(rawurl, params)
}

func handleNotFoundError(err error, d *schema.ResourceData, resource string) error {
	return transport_tpg.HandleNotFoundError(err, d, resource)
}

func IsGoogleApiErrorWithCode(err error, errCode int) bool {
	return transport_tpg.IsGoogleApiErrorWithCode(err, errCode)
}

func isApiNotEnabledError(err error) bool {
	return transport_tpg.IsApiNotEnabledError(err)
}
