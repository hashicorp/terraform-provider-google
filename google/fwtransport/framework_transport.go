// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwtransport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/googleapi"
)

func SendFrameworkRequest(p *FrameworkProviderConfig, method, project, rawurl, userAgent string, body map[string]interface{}, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) (map[string]interface{}, diag.Diagnostics) {
	return SendFrameworkRequestWithTimeout(p, method, project, rawurl, userAgent, body, transport_tpg.DefaultRequestTimeout, errorRetryPredicates...)
}

func SendFrameworkRequestWithTimeout(p *FrameworkProviderConfig, method, project, rawurl, userAgent string, body map[string]interface{}, timeout time.Duration, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	reqHeaders := make(http.Header)
	reqHeaders.Set("User-Agent", userAgent)
	reqHeaders.Set("Content-Type", "application/json")

	if p.UserProjectOverride.ValueBool() && project != "" {
		// When project is "NO_BILLING_PROJECT_OVERRIDE" in the function GetCurrentUserEmail,
		// set the header X-Goog-User-Project to be empty string.
		if project == "NO_BILLING_PROJECT_OVERRIDE" {
			reqHeaders.Set("X-Goog-User-Project", "")
		} else {
			// Pass the project into this fn instead of parsing it from the URL because
			// both project names and URLs can have colons in them.
			reqHeaders.Set("X-Goog-User-Project", project)
		}
	}

	if timeout == 0 {
		timeout = time.Hour
	}

	var res *http.Response
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			var buf bytes.Buffer
			if body != nil {
				err := json.NewEncoder(&buf).Encode(body)
				if err != nil {
					return err
				}
			}

			u, err := transport_tpg.AddQueryParams(rawurl, map[string]string{"alt": "json"})
			if err != nil {
				return err
			}
			req, err := http.NewRequest(method, u, &buf)
			if err != nil {
				return err
			}

			req.Header = reqHeaders
			res, err = p.Client.Do(req)
			if err != nil {
				return err
			}

			if err := googleapi.CheckResponse(res); err != nil {
				googleapi.CloseBody(res)
				return err
			}

			return nil
		},
		Timeout:              timeout,
		ErrorRetryPredicates: errorRetryPredicates,
	})
	if err != nil {
		diags.AddError("error sending request", err.Error())
		return nil, diags
	}

	if res == nil {
		diags.AddError("Unable to parse server response.", "This is most likely a terraform problem, please file a bug at https://github.com/hashicorp/terraform-provider-google/issues.")
		return nil, diags
	}

	// The defer call must be made outside of the retryFunc otherwise it's closed too soon.
	defer googleapi.CloseBody(res)

	// 204 responses will have no body, so we're going to error with "EOF" if we
	// try to parse it. Instead, we can just return nil.
	if res.StatusCode == 204 {
		return nil, diags
	}
	result := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		diags.AddError("error decoding response body", err.Error())
		return nil, diags
	}

	return result, diags
}
