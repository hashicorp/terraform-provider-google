// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudidentity

import (
	"log"
	"strings"

	"github.com/hashicorp/errwrap"
	"google.golang.org/api/googleapi"
)

func transformCloudIdentityGroupMembershipReadError(err error) error {
	if gErr, ok := errwrap.GetType(err, &googleapi.Error{}).(*googleapi.Error); ok {
		if gErr.Code == 403 && strings.Contains(gErr.Message, "(or it may not exist)") {
			// This error occurs when either the group membership does not exist, or permission is denied. It is
			// deliberately ambiguous so that existence information is not revealed to the caller. However, for
			// the Read function, we can only assume that the membership does not exist, and proceed with attempting
			// other operations. Since HandleNotFoundError(...) expects an error code of 404 when a resource does not
			// exist, to get the desired behavior, we modify the error code to be 404.
			gErr.Code = 404
		}

		log.Printf("[DEBUG] Transformed CloudIdentityGroupMembership error")
		return gErr
	}

	return err
}
