// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"fmt"
	"log"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func oldValue(old, new interface{}) interface{} {
	return old
}

func handleNotFoundDCLError(err error, d *schema.ResourceData, resourceName string) error {
	if dcl.IsNotFound(err) {
		log.Printf("[WARN] Removing %s because it's gone", resourceName)
		// The resource doesn't exist anymore
		d.SetId("")
		return nil
	}

	return errwrap.Wrapf(
		fmt.Sprintf("Error when reading or editing %s: {{err}}", resourceName), err)
}
