// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwtransport

import (
	"fmt"
	"strings"

	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Methods to create new services from config
// Some base paths below need the version and possibly more of the path
// set on them. The client libraries are inconsistent about which values they need;
// while most only want the host URL, some older ones also want the version and some
// of those "projects" as well. You can find out if this is required by looking at
// the basePath value in the client library file.

func (p *FrameworkProviderConfig) NewDnsClient(userAgent string, diags *diag.Diagnostics) *dns.Service {
	dnsClientBasePath := transport_tpg.RemoveBasePathVersion(p.DNSBasePath)
	dnsClientBasePath = strings.ReplaceAll(dnsClientBasePath, "/dns/", "")
	tflog.Info(p.Context, fmt.Sprintf("Instantiating Google Cloud DNS client for path %s", dnsClientBasePath))
	clientDns, err := dns.NewService(p.Context, option.WithHTTPClient(p.Client))
	if err != nil {
		diags.AddWarning("error creating client dns", err.Error())
		return nil
	}
	clientDns.UserAgent = userAgent
	clientDns.BasePath = dnsClientBasePath

	return clientDns
}
