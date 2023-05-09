package google

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

// Deprecated: For backward compatibility CommonRefreshFunc is still working,
// but all new code should use CommonRefreshFunc in the tpgresource package instead.
func CommonRefreshFunc(w tpgresource.Waiter) resource.StateRefreshFunc {
	return tpgresource.CommonRefreshFunc(w)
}

// Deprecated: For backward compatibility OperationWait is still working,
// but all new code should use OperationWait in the tpgresource package instead.
func OperationWait(w tpgresource.Waiter, activity string, timeout time.Duration, pollInterval time.Duration) error {
	return tpgresource.OperationWait(w, activity, timeout, pollInterval)
}
