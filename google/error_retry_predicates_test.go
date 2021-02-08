package google

import (
	"testing"

	"google.golang.org/api/googleapi"
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
