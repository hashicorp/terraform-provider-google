// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"time"

	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	"google.golang.org/api/iam/v1"
)

func serviceAccountKeyWaitTime(client *iam.ProjectsServiceAccountsKeysService, keyName, publicKeyType, activity string, timeout time.Duration) error {
	return resourcemanager.ServiceAccountKeyWaitTime(client, keyName, publicKeyType, activity, timeout)
}
