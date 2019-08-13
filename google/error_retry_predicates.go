package google

import (
	"strings"

	"google.golang.org/api/googleapi"
)

// If a permission necessary to provision a resource is created in the same config
// as the resource itself, the permission may not have propagated by the time terraform
// attempts to create the resource. This allows those errors to be retried until the timeout expires
func iamMemberMissing(err error) (bool, string) {
	if gerr, ok := err.(*googleapi.Error); ok {
		if gerr.Code == 400 && strings.Contains(gerr.Body, "permission") {
			return true, "Waiting for IAM member permissions to propagate."
		}
	}
	return false, ""
}
