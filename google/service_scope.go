package google

import (
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func canonicalizeServiceScope(scope string) string {
	return tpgresource.CanonicalizeServiceScope(scope)
}

func canonicalizeServiceScopes(scopes []string) []string {
	return tpgresource.CanonicalizeServiceScopes(scopes)
}

func stringScopeHashcode(v interface{}) int {
	return tpgresource.StringScopeHashcode(v)
}
