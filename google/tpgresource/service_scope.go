// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpgresource

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func CanonicalizeServiceScope(scope string) string {
	// This is a convenience map of short names used by the gcloud tool
	// to the GCE auth endpoints they alias to.
	scopeMap := map[string]string{
		"bigquery":              "https://www.googleapis.com/auth/bigquery",
		"cloud-platform":        "https://www.googleapis.com/auth/cloud-platform",
		"cloud-source-repos":    "https://www.googleapis.com/auth/source.full_control",
		"cloud-source-repos-ro": "https://www.googleapis.com/auth/source.read_only",
		"compute-ro":            "https://www.googleapis.com/auth/compute.readonly",
		"compute-rw":            "https://www.googleapis.com/auth/compute",
		"datastore":             "https://www.googleapis.com/auth/datastore",
		"logging-write":         "https://www.googleapis.com/auth/logging.write",
		"monitoring":            "https://www.googleapis.com/auth/monitoring",
		"monitoring-read":       "https://www.googleapis.com/auth/monitoring.read",
		"monitoring-write":      "https://www.googleapis.com/auth/monitoring.write",
		"pubsub":                "https://www.googleapis.com/auth/pubsub",
		"service-control":       "https://www.googleapis.com/auth/servicecontrol",
		"service-management":    "https://www.googleapis.com/auth/service.management.readonly",
		"sql":                   "https://www.googleapis.com/auth/sqlservice",
		"sql-admin":             "https://www.googleapis.com/auth/sqlservice.admin",
		"storage-full":          "https://www.googleapis.com/auth/devstorage.full_control",
		"storage-ro":            "https://www.googleapis.com/auth/devstorage.read_only",
		"storage-rw":            "https://www.googleapis.com/auth/devstorage.read_write",
		"taskqueue":             "https://www.googleapis.com/auth/taskqueue",
		"trace":                 "https://www.googleapis.com/auth/trace.append",
		"useraccounts-ro":       "https://www.googleapis.com/auth/cloud.useraccounts.readonly",
		"useraccounts-rw":       "https://www.googleapis.com/auth/cloud.useraccounts",
		"userinfo-email":        "https://www.googleapis.com/auth/userinfo.email",
	}

	if matchedURL, ok := scopeMap[scope]; ok {
		return matchedURL
	}

	return scope
}

func CanonicalizeServiceScopes(scopes []string) []string {
	cs := make([]string, len(scopes))
	for i, scope := range scopes {
		cs[i] = CanonicalizeServiceScope(scope)
	}
	return cs
}

func StringScopeHashcode(v interface{}) int {
	v = CanonicalizeServiceScope(v.(string))
	return schema.HashString(v)
}
