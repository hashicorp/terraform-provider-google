// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apigee

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/googleapi"
	"io"
	"log"
	"net/http"
	"time"
)

func resourceApigeeNatAddressActivate(config *transport_tpg.Config, d *schema.ResourceData, billingProject string, userAgent string) error {
	// 1. check prepare for activation
	name := d.Get("name").(string)

	if d.Get("state").(string) != "RESERVED" {
		return fmt.Errorf("Activating NAT address requires the state to become RESERVED")
	}

	// 2. activation
	activateUrl, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}{{instance_id}}/natAddresses/{{name}}:activate")
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Activating NAT address: %s", name)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    activateUrl,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error activating NAT address: %s", err)
	}

	var opRes map[string]interface{}
	err = ApigeeOperationWaitTimeWithResponse(
		config, res, &opRes, "Activating NAT address", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error waiting to actiavte NAT address: %s", err)
	} else {
		log.Printf("[DEBUG] Finished activating NatAddress %q: %#v", d.Id(), res)
	}
	return nil
}

// sendRequestRawBodyWithTimeout is derived from sendRequestWithTimeout with direct pass through of request body
func sendRequestRawBodyWithTimeout(config *transport_tpg.Config, method, project, rawurl, userAgent string, body io.Reader, contentType string, timeout time.Duration, errorRetryPredicates ...transport_tpg.RetryErrorPredicateFunc) (map[string]interface{}, error) {
	log.Printf("[DEBUG] sendRequestRawBodyWithTimeout start")
	reqHeaders := make(http.Header)
	reqHeaders.Set("User-Agent", userAgent)
	reqHeaders.Set("Content-Type", contentType)

	if config.UserProjectOverride && project != "" {
		// Pass the project into this fn instead of parsing it from the URL because
		// both project names and URLs can have colons in them.
		reqHeaders.Set("X-Goog-User-Project", project)
	}

	if timeout == 0 {
		timeout = time.Duration(1) * time.Minute
	}

	var res *http.Response

	log.Printf("[DEBUG] sendRequestRawBodyWithTimeout sending request")

	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			req, err := http.NewRequest(method, rawurl, body)
			if err != nil {
				return err
			}

			req.Header = reqHeaders
			res, err = config.Client.Do(req)
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
		return nil, err
	}

	if res == nil {
		return nil, fmt.Errorf("Unable to parse server response. This is most likely a terraform problem, please file a bug at https://github.com/hashicorp/terraform-provider-google/issues.")
	}

	// The defer call must be made outside of the retryFunc otherwise it's closed too soon.
	defer googleapi.CloseBody(res)

	// 204 responses will have no body, so we're going to error with "EOF" if we
	// try to parse it. Instead, we can just return nil.
	if res.StatusCode == 204 {
		return nil, nil
	}
	result := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] sendRequestRawBodyWithTimeout returning")
	return result, nil
}
