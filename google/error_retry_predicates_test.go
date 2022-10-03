package google

import (
	"strconv"
	"testing"

	"google.golang.org/api/googleapi"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestIsAppEngineRetryableError_operationInProgress(t *testing.T) {
	err := googleapi.Error{
		Code: 409,
		Body: "Operation is already in progress",
	}
	isRetryable, _ := isAppEngineRetryableError(&err)
	if !isRetryable {
		t.Errorf("Error not detected as retryable")
	}
}

func TestIsAppEngineRetryableError_p4saPropagation(t *testing.T) {
	err := googleapi.Error{
		Code: 404,
		Body: "Unable to retrieve P4SA: [service-111111111111@gcp-gae-service.iam.gserviceaccount.com] from GAIA. Could be GAIA propagation delay or request from deleted apps.",
	}
	isRetryable, _ := isAppEngineRetryableError(&err)
	if !isRetryable {
		t.Errorf("Error not detected as retryable")
	}
}

func TestIsAppEngineRetryableError_missingPage(t *testing.T) {
	err := googleapi.Error{
		Code: 404,
		Body: "Missing page",
	}
	isRetryable, _ := isAppEngineRetryableError(&err)
	if isRetryable {
		t.Errorf("Error incorrectly detected as retryable")
	}
}

func TestIsAppEngineRetryableError_serverError(t *testing.T) {
	err := googleapi.Error{
		Code: 500,
		Body: "Unable to retrieve P4SA because of a bad thing happening",
	}
	isRetryable, _ := isAppEngineRetryableError(&err)
	if isRetryable {
		t.Errorf("Error incorrectly detected as retryable")
	}
}

func TestIsCommonRetryableErrorCode_retryableErrorCode(t *testing.T) {
	codes := []int{429, 500, 502, 503}
	for _, code := range codes {
		code := code
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			err := googleapi.Error{
				Code: code,
				Body: "some text describing error",
			}
			isRetryable, _ := isCommonRetryableErrorCode(&err)
			if !isRetryable {
				t.Errorf("Error not detected as retryable")
			}
		})
	}
}

func TestIsCommonRetryableErrorCode_otherError(t *testing.T) {
	err := googleapi.Error{
		Code: 404,
		Body: "Some unretryable issue",
	}
	isRetryable, _ := isCommonRetryableErrorCode(&err)
	if isRetryable {
		t.Errorf("Error incorrectly detected as retryable")
	}
}

func TestIsOperationReadQuotaError_quotaExceeded(t *testing.T) {
	err := googleapi.Error{
		Code: 403,
		Body: "Quota exceeded for quota metric 'OperationReadGroup' and limit 'Operation read requests per minute' of service 'compute.googleapis.com' for consumer 'project_number:11111111'.",
	}
	isRetryable, _ := is403QuotaExceededPerMinuteError(&err)
	if !isRetryable {
		t.Errorf("Error not detected as retryable")
	}
}

func TestIs403QuotaExceededPerMinuteError_perMinuteQuotaExceeded(t *testing.T) {
	err := googleapi.Error{
		Code: 403,
		Body: "Quota exceeded for quota metric 'Queries' and limit 'Queries per minute' of service 'compute.googleapis.com' for consumer 'project_number:11111111'.",
	}
	isRetryable, _ := is403QuotaExceededPerMinuteError(&err)
	if !isRetryable {
		t.Errorf("Error not detected as retryable")
	}
}

func TestIs403QuotaExceededPerMinuteError_perDayQuotaExceededNotRetryable(t *testing.T) {
	err := googleapi.Error{
		Code: 403,
		Body: "Quota exceeded for quota metric 'Queries' and limit 'Queries per day' of service 'compute.googleapis.com' for consumer 'project_number:11111111'.",
	}
	isRetryable, _ := is403QuotaExceededPerMinuteError(&err)
	if isRetryable {
		t.Errorf("Error incorrectly detected as retryable")
	}
}

// An error with retry info is retryable.
func TestBigtableError_retryable(t *testing.T) {
	retryInfo := &errdetails.RetryInfo{
		RetryDelay: &durationpb.Duration{Seconds: 10, Nanos: 10},
	}
	status, _ := status.New(codes.FailedPrecondition, "is retryable").WithDetails(retryInfo)
	isRetryable, _ := isBigTableRetryableError(status.Err())
	if !isRetryable {
		t.Errorf("Error not detected as retryable")
	}
}

// An error without retry info is not retryable.
func TestBigtableError_withoutRetryInfoNotRetryable(t *testing.T) {
	status := status.New(codes.FailedPrecondition, "is not retryable")
	isRetryable, _ := isBigTableRetryableError(status.Err())
	if isRetryable {
		t.Errorf("Error incorrectly detected as retryable")
	}
}

// An OK status with retry info is not retryable.
func TestBigtableError_okIsNotRetryable(t *testing.T) {
	retryInfo := &errdetails.RetryInfo{
		RetryDelay: &durationpb.Duration{Seconds: 10, Nanos: 10},
	}
	status, _ := status.New(codes.OK, "is not retryable").WithDetails(retryInfo)
	isRetryable, _ := isBigTableRetryableError(status.Err())
	if isRetryable {
		t.Errorf("Error incorrectly detected as retryable")
	}
}
