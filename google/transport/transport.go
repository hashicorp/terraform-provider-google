package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/googleapi"
)

var DefaultRequestTimeout = 5 * time.Minute

type SendRequestOptions struct {
	Config               *Config
	Method               string
	Project              string
	RawURL               string
	UserAgent            string
	Body                 map[string]any
	Timeout              time.Duration
	ErrorRetryPredicates []RetryErrorPredicateFunc
}

func SendRequest(opt SendRequestOptions) (map[string]interface{}, error) {
	reqHeaders := make(http.Header)
	reqHeaders.Set("User-Agent", opt.UserAgent)
	reqHeaders.Set("Content-Type", "application/json")

	if opt.Config.UserProjectOverride && opt.Project != "" {
		// When opt.Project is "NO_BILLING_PROJECT_OVERRIDE" in the function GetCurrentUserEmail,
		// set the header X-Goog-User-Project to be empty string.
		if opt.Project == "NO_BILLING_PROJECT_OVERRIDE" {
			reqHeaders.Set("X-Goog-User-Project", "")
		} else {
			// Pass the project into this fn instead of parsing it from the URL because
			// both project names and URLs can have colons in them.
			reqHeaders.Set("X-Goog-User-Project", opt.Project)
		}
	}

	if opt.Timeout == 0 {
		opt.Timeout = DefaultRequestTimeout
	}

	var res *http.Response
	err := RetryTimeDuration(
		func() error {
			var buf bytes.Buffer
			if opt.Body != nil {
				err := json.NewEncoder(&buf).Encode(opt.Body)
				if err != nil {
					return err
				}
			}

			u, err := AddQueryParams(opt.RawURL, map[string]string{"alt": "json"})
			if err != nil {
				return err
			}
			req, err := http.NewRequest(opt.Method, u, &buf)
			if err != nil {
				return err
			}

			req.Header = reqHeaders
			res, err = opt.Config.Client.Do(req)
			if err != nil {
				return err
			}

			if err := googleapi.CheckResponse(res); err != nil {
				googleapi.CloseBody(res)
				return err
			}

			return nil
		},
		opt.Timeout,
		opt.ErrorRetryPredicates...,
	)
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

	return result, nil
}

func AddQueryParams(rawurl string, params map[string]string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func HandleNotFoundError(err error, d *schema.ResourceData, resource string) error {
	if IsGoogleApiErrorWithCode(err, 404) {
		log.Printf("[WARN] Removing %s because it's gone", resource)
		// The resource doesn't exist anymore
		d.SetId("")

		return nil
	}

	return errwrap.Wrapf(
		fmt.Sprintf("Error when reading or editing %s: {{err}}", resource), err)
}

func IsGoogleApiErrorWithCode(err error, errCode int) bool {
	gerr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
	return ok && gerr != nil && gerr.Code == errCode
}

func IsApiNotEnabledError(err error) bool {
	gerr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error)
	if !ok {
		return false
	}
	if gerr == nil {
		return false
	}
	if gerr.Code != 403 {
		return false
	}
	for _, e := range gerr.Errors {
		if e.Reason == "accessNotConfigured" {
			return true
		}
	}
	return false
}
