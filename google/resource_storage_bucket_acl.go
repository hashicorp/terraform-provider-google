// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/storage"
)

func getRoleEntityPair(role_entity string) (*storage.RoleEntity, error) {
	return storage.GetRoleEntityPair(role_entity)
}
