// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package tpgdclresource

import (
	"context"
	"fmt"
	"log"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func OldValue(old, new interface{}) interface{} {
	return old
}

func HandleNotFoundDCLError(err error, d *schema.ResourceData, resourceName string) error {
	if dcl.IsNotFound(err) {
		log.Printf("[WARN] Removing %s because it's gone", resourceName)
		// The resource doesn't exist anymore
		d.SetId("")
		return nil
	}

	return errwrap.Wrapf(
		fmt.Sprintf("Error when reading or editing %s: {{err}}", resourceName), err)
}

func ResourceContainerAwsNodePoolCustomizeDiffFunc(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	count := diff.Get("update_settings.#").(int)
	if count < 1 {
		return nil
	}

	oMaxSurge, nMaxSurge := diff.GetChange("update_settings.0.surge_settings.0.max_surge")
	oMaxUnavailable, nMaxUnavailable := diff.GetChange("update_settings.0.surge_settings.0.max_unavailable")

	// Server default of maxSurge = 1 and maxUnavailable = 0 is not returned
	// Clear the diff if trying to resolve these specific values
	if oMaxSurge == 0 && nMaxSurge == 1 && oMaxUnavailable == 0 && nMaxUnavailable == 0 {
		err := diff.Clear("update_settings")
		if err != nil {
			return err
		}
	}

	return nil
}
