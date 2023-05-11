package google

import (
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Deprecated: For backward compatibility MultiEnvSearch is still working,
// but all new code should use MultiEnvSearch in the transport_tpg package instead.
func MultiEnvSearch(ks []string) string {
	return transport_tpg.MultiEnvSearch(ks)
}
