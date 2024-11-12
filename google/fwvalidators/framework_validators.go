// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwvalidators

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

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
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
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

// Non Empty String Validator
type nonEmptyStringValidator struct {
}

// Description describes the validation in plain text formatting.
func (v nonEmptyStringValidator) Description(_ context.Context) string {
	return "value expected to be a string that isn't an empty string"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v nonEmptyStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v nonEmptyStringValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if value == "" {
		response.Diagnostics.AddError("expected a non-empty string", fmt.Sprintf("%s was set to `%s`", request.Path, value))
	}
}

func NonEmptyStringValidator() validator.String {
	return nonEmptyStringValidator{}
}

// Define the possible service account name patterns
var ServiceAccountEmailPatterns = []string{
	`^.+@.+\.iam\.gserviceaccount\.com$`,                     // Standard IAM service account
	`^.+@developer\.gserviceaccount\.com$`,                   // Legacy developer service account
	`^.+@appspot\.gserviceaccount\.com$`,                     // App Engine service account
	`^.+@cloudservices\.gserviceaccount\.com$`,               // Google Cloud services service account
	`^.+@cloudbuild\.gserviceaccount\.com$`,                  // Cloud Build service account
	`^service-[0-9]+@.+-compute\.iam\.gserviceaccount\.com$`, // Compute Engine service account
}

// Create a custom validator for service account names
type ServiceAccountEmailValidator struct{}

func (v ServiceAccountEmailValidator) Description(ctx context.Context) string {
	return "value must be a valid service account email address"
}

func (v ServiceAccountEmailValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ServiceAccountEmailValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	// Check for empty string
	if value == "" {
		resp.Diagnostics.AddError("Invalid Service Account Name", "Service account name must not be empty")
		return
	}

	valid := false
	for _, pattern := range ServiceAccountEmailPatterns {
		if matched, _ := regexp.MatchString(pattern, value); matched {
			valid = true
			break
		}
	}

	if !valid {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Service Account Name",
			"Service account name must match one of the expected patterns for Google service accounts",
		)
	}
}

// Create a custom validator for duration
type BoundedDuration struct {
	MinDuration time.Duration
	MaxDuration time.Duration
}

func (v BoundedDuration) Description(ctx context.Context) string {
	return fmt.Sprintf("value must be a valid duration string between %v and %v", v.MinDuration, v.MaxDuration)
}

func (v BoundedDuration) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v BoundedDuration) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	duration, err := time.ParseDuration(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Duration Format",
			"Duration must be a valid duration string (e.g., '3600s', '1h')",
		)
		return
	}

	if duration < v.MinDuration || duration > v.MaxDuration {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Duration",
			fmt.Sprintf("Duration must be between %v and %v", v.MinDuration, v.MaxDuration),
		)
	}
}
