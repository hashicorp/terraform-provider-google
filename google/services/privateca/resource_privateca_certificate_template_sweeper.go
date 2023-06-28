// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package privateca

import (
	"context"
	"log"
	"testing"

	privateca "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/privateca"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func init() {
	sweeper.AddTestSweepers("PrivatecaCertificateTemplate", testSweepPrivatecaCertificateTemplate)
}

func testSweepPrivatecaCertificateTemplate(region string) error {
	log.Print("[INFO][SWEEPER_LOG] Starting sweeper for PrivatecaCertificateTemplate")

	config, err := sweeper.SharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}

	t := &testing.T{}
	billingId := envvar.GetTestBillingAccountFromEnv(t)

	// Setup variables to be used for Delete arguments.
	d := map[string]string{
		"project":         config.Project,
		"region":          region,
		"location":        region,
		"zone":            "-",
		"billing_account": billingId,
	}

	client := transport_tpg.NewDCLPrivatecaClient(config, config.UserAgent, "", 0)
	err = client.DeleteAllCertificateTemplate(context.Background(), d["project"], d["location"], isDeletablePrivatecaCertificateTemplate)
	if err != nil {
		return err
	}
	return nil
}

func isDeletablePrivatecaCertificateTemplate(r *privateca.CertificateTemplate) bool {
	return sweeper.IsSweepableTestResource(*r.Name)
}
