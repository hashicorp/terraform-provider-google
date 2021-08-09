package google

import (
	"strings"
	"testing"
)

func TestUnitMtls_urlSwitching(t *testing.T) {
	t.Parallel()
	for key, bp := range DefaultBasePaths {
		url := getMtlsEndpoint(bp)
		if !strings.Contains(url, ".mtls.") {
			t.Errorf("%s: mtls conversion unsuccessful preconv - %s postconv - %s", key, bp, url)
		}
	}
}
