// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"

	googleoauth "golang.org/x/oauth2/google"
)

func ValidateCredentials(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil {
		return
	}
	creds := v.(string)

	// reject empty strings
	if v.(string) == "" {
		errors = append(errors,
			fmt.Errorf("expected a non-empty string"))
		return
	}

	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(creds); err == nil {
		return
	}
	if _, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(creds)); err != nil {
		errors = append(errors,
			fmt.Errorf("JSON credentials are not valid: %s", err))
	}

	return
}

func ValidateEmptyStrings(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil {
		return
	}

	if v.(string) == "" {
		errors = append(errors,
			fmt.Errorf("expected a non-empty string"))
	}

	return
}

func ValidateJWT(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil || v.(string) == "" {
		errors = append(errors, fmt.Errorf("%q cannot be empty", k))
		return
	}

	// JWT consists of 3 parts separated by dots: header.payload.signature
	parts := strings.Split(v.(string), ".")
	if len(parts) != 3 {
		errors = append(errors, fmt.Errorf("%q is not a valid JWT format", k))
		return
	}

	// Check that each part is base64 encoded
	for i, part := range parts {
		if _, err := base64.RawURLEncoding.DecodeString(part); err != nil {
			errors = append(errors, fmt.Errorf("part %d of JWT is not valid base64: %v", i+1, err))
		}
	}

	return
}

func ValidateServiceAccountEmail(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if value == "" {
		errors = append(errors, fmt.Errorf("%q cannot be empty", k))
		return
	}

	serviceAccountRegex := regexp.MustCompile(`^[a-zA-Z0-9-]+@[a-zA-Z0-9-]+\.iam\.gserviceaccount\.com$`)
	if !serviceAccountRegex.MatchString(value) {
		errors = append(errors, fmt.Errorf("%q is not a valid service account email address (expected format: name@project-id.iam.gserviceaccount.com)", k))
	}

	return
}
