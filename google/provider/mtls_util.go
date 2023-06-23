// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package provider

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"google.golang.org/api/option/internaloption"
	"google.golang.org/api/transport"
)

// The transport libaray does not natively expose logic to determine whether
// the user is within mtls mode or not. They do return the mtls endpoint if
// it is enabled during client creation so we will use this logic to determine
// the mode the user is in and throw away the client they give us back.
func isMtls() bool {
	regularEndpoint := "https://mockservice.googleapis.com/v1/"
	mtlsEndpoint := getMtlsEndpoint(regularEndpoint)
	_, endpoint, err := transport.NewHTTPClient(context.Background(),
		internaloption.WithDefaultEndpoint(regularEndpoint),
		internaloption.WithDefaultMTLSEndpoint(mtlsEndpoint),
	)
	if err != nil {
		return false
	}
	isMtls := endpoint == mtlsEndpoint
	return isMtls
}

func getMtlsEndpoint(baseEndpoint string) string {
	u, err := url.Parse(baseEndpoint)
	if err != nil {
		if strings.Contains(baseEndpoint, ".googleapis") {
			return strings.Replace(baseEndpoint, ".googleapis", ".mtls.googleapis", 1)
		}
		return baseEndpoint
	}
	domainParts := strings.Split(u.Host, ".")
	if len(domainParts) > 1 {
		u.Host = fmt.Sprintf("%s.mtls.%s", domainParts[0], strings.Join(domainParts[1:], "."))
	} else {
		u.Host = fmt.Sprintf("%s.mtls", domainParts[0])
	}
	return u.String()
}
