// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sweeper

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// List of prefixes used for test resource names
var testResourcePrefixes = []string{
	// tf-test and tf_test are automatically prepended to resource ids in examples that
	// include a "-" or "_" respectively, and they are the preferred prefix for our test resources to use
	"tf-test",
	"tf_test",
	"tfgen",
	"gke-us-central1-tf",  // composer-created disks which are abandoned by design (https://cloud.google.com/composer/pricing)
	"gcs-bucket-tf-test-", // https://github.com/hashicorp/terraform-provider-google/issues/8909
	"df-",                 // https://github.com/hashicorp/terraform-provider-google/issues/8909
	"resourcegroup-",      // https://github.com/hashicorp/terraform-provider-google/issues/8924
	"cluster-",            // https://github.com/hashicorp/terraform-provider-google/issues/8924
	"k8s-fw-",             // firewall rules are getting created and not cleaned up by k8 resources using this prefix
}

// SharedConfigForRegion returns a common config setup needed for the sweeper
// functions for a given region
func SharedConfigForRegion(region string) (*transport_tpg.Config, error) {
	project := envvar.GetTestProjectFromEnv()
	if project == "" {
		return nil, fmt.Errorf("set project using any of these env variables %v", envvar.ProjectEnvVars)
	}

	if v := transport_tpg.MultiEnvSearch(envvar.CredsEnvVars); v == "" {
		return nil, fmt.Errorf("set credentials using any of these env variables %v", envvar.CredsEnvVars)
	}

	conf := &transport_tpg.Config{
		Credentials: envvar.GetTestCredsFromEnv(),
		Region:      region,
		Project:     project,
	}

	transport_tpg.ConfigureBasePaths(conf)

	return conf, nil
}

func IsSweepableTestResource(resourceName string) bool {
	for _, p := range testResourcePrefixes {
		if strings.HasPrefix(resourceName, p) {
			return true
		}
	}
	return false
}

func AddTestSweepers(name string, sweeper func(region string) error) {
	_, filename, _, _ := runtime.Caller(0)
	hash := crc32.NewIEEE()
	hash.Write([]byte(filename))
	hashedFilename := hex.EncodeToString(hash.Sum(nil))
	uniqueName := name + "_" + hashedFilename

	resource.AddTestSweepers(uniqueName, &resource.Sweeper{
		Name: name,
		F:    sweeper,
	})
}
