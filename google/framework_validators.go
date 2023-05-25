// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package google

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	googleoauth "golang.org/x/oauth2/google"
)

// Credentials Validator
var _ validator.String = credentialsValidator{}

// credentialsValidator validates that a string Attribute's is valid JSON credentials.
type credentialsValidator struct {
}

// Description describes the validation in plain text formatting.
func (v credentialsValidator) Description(_ context.Context) string {
	return "value must be a path to valid JSON credentials or valid, raw, JSON credentials"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v credentialsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v credentialsValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() || request.ConfigValue.Equal(types.StringValue("")) {
		return
	}

	value := request.ConfigValue.ValueString()

	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(value); err == nil {
		return
	}
	if _, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(value)); err != nil {
		response.Diagnostics.AddError("JSON credentials are not valid", err.Error())
	}
}

func CredentialsValidator() validator.String {
	return credentialsValidator{}
}

// Non Negative Duration Validator
type nonnegativedurationValidator struct {
}

// Description describes the validation in plain text formatting.
func (v nonnegativedurationValidator) Description(_ context.Context) string {
	return "value expected to be a string representing a non-negative duration"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v nonnegativedurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v nonnegativedurationValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	dur, err := time.ParseDuration(value)
	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("expected %s to be a duration", value), err.Error())
		return
	}

	if dur < 0 {
		response.Diagnostics.AddError("duration must be non-negative", fmt.Sprintf("duration provided: %d", dur))
	}
}

func NonNegativeDurationValidator() validator.String {
	return nonnegativedurationValidator{}
}
