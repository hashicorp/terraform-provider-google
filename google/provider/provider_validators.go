// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider

import (
	"context"
	"fmt"
	"os"

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
