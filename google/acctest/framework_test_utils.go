// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"context"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func GetFwTestProvider(t *testing.T) *frameworkTestProvider {
	configsLock.RLock()
	fwProvider, ok := fwProviders[t.Name()]
	configsLock.RUnlock()
	if ok {
		return fwProvider
	}

	var diags diag.Diagnostics
	p := NewFrameworkTestProvider(t.Name())
	configureApiClient(context.Background(), &p.FrameworkProvider, &diags)
	if diags.HasError() {
		log.Fatalf("%d errors when configuring test provider client: first is %s", diags.ErrorsCount(), diags.Errors()[0].Detail())
	}

	return p
}
