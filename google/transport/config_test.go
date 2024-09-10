// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package transport_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/provider"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	googleoauth "golang.org/x/oauth2/google"
)

const testOauthScope = "https://www.googleapis.com/auth/compute"

func TestHandleSDKDefaults_BillingProject(t *testing.T) {
	cases := map[string]struct {
		ConfigValue      string
		EnvVariables     map[string]string
		ExpectedValue    string
		ValueNotProvided bool
		ExpectError      bool
	}{
		"billing project value set in the provider config is not overridden by ENVs": {
			ConfigValue: "my-billing-project-from-config",
			EnvVariables: map[string]string{
				"GOOGLE_BILLING_PROJECT": "my-billing-project-from-env",
			},
			ExpectedValue: "my-billing-project-from-config",
		},
		"billing project can be set by environment variable, when no value supplied via the config": {
			EnvVariables: map[string]string{
				"GOOGLE_BILLING_PROJECT": "my-billing-project-from-env",
			},
			ExpectedValue: "my-billing-project-from-env",
		},
		"when no values are provided via config or environment variables, the field remains unset without error": {
			EnvVariables: map[string]string{
				"GOOGLE_BILLING_PROJECT": "", // GOOGLE_BILLING_PROJECT unset
			},
			ValueNotProvided: true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			// Create empty schema.ResourceData using the SDK Provider schema
			emptyConfigMap := map[string]interface{}{}
			d := schema.TestResourceDataRaw(t, provider.Provider().Schema, emptyConfigMap)

			// Set config value(s)
			if tc.ConfigValue != "" {
				d.Set("billing_project", tc.ConfigValue)
			}

			// Set ENVs
			if len(tc.EnvVariables) > 0 {
				for k, v := range tc.EnvVariables {
					t.Setenv(k, v)
				}
			}

			// Act
			err := transport_tpg.HandleSDKDefaults(d)

			// Assert
			if err != nil {
				if !tc.ExpectError {
					t.Fatalf("error: %v", err)
				}
				return
			}

			// Assert
			v, ok := d.GetOk("billing_project")
			if !ok && !tc.ValueNotProvided {
				t.Fatal("expected billing_project to be set in the provider data")
			}
			if ok && tc.ValueNotProvided {
				t.Fatal("expected billing_project to not be set in the provider data")
			}

			if v != tc.ExpectedValue {
				t.Fatalf("unexpected value: wanted %v, got, %v", tc.ExpectedValue, v)
			}
		})
	}
}

// The `user_project_override` field is an odd one out, as other provider schema fields tend to be strings
// and `user_project_override` is a boolean
func TestHandleSDKDefaults_UserProjectOverride(t *testing.T) {
	cases := map[string]struct {
		SetViaConfig     bool // Awkward, but necessary as zero value of ConfigValue could be intended
		ConfigValue      bool
		ValueNotProvided bool
		EnvVariables     map[string]string
		ExpectedValue    bool
		ExpectError      bool
	}{
		"user_project_override value set in the provider schema is not overridden by ENVs": {
			SetViaConfig: true,
			ConfigValue:  false,
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "true",
			},
			ExpectedValue: false,
		},
		"user_project_override can be set by environment variable: true": {
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "true",
			},
			ExpectedValue: true,
		},
		"user_project_override can be set by environment variable: false": {
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "false",
			},
			ExpectedValue: false,
		},
		"user_project_override can be set by environment variable: 1": {
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "1",
			},
			ExpectedValue: true,
		},
		"user_project_override can be set by environment variable: 0": {
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "0",
			},
			ExpectedValue: false,
		},
		"error returned due to non-boolean environment variables": {
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "I'm not a boolean",
			},
			ExpectError: true,
		},
		"when no values are provided via config or environment variables, the field remains unset without error": {
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "", // USER_PROJECT_OVERRIDE unset
			},
			ValueNotProvided: true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			// Create empty schema.ResourceData using the SDK Provider schema
			emptyConfigMap := map[string]interface{}{}
			d := schema.TestResourceDataRaw(t, provider.Provider().Schema, emptyConfigMap)

			// Set config value(s)
			if tc.SetViaConfig {
				d.Set("user_project_override", tc.ConfigValue)
			}

			// Set ENVs
			if len(tc.EnvVariables) > 0 {
				for k, v := range tc.EnvVariables {
					t.Setenv(k, v)
				}
			}

			// Act
			err := transport_tpg.HandleSDKDefaults(d)

			// Assert
			if err != nil {
				if !tc.ExpectError {
					t.Fatalf("error: %v", err)
				}
				return
			}

			v, ok := d.GetOkExists("user_project_override")
			if !ok && !tc.ValueNotProvided {
				t.Fatal("expected user_project_override to be set in the provider data")
			}
			if ok && tc.ValueNotProvided {
				t.Fatal("expected user_project_override to not be set in the provider data")
			}

			if v != tc.ExpectedValue {
				t.Fatalf("unexpected value: wanted %v, got, %v", tc.ExpectedValue, v)
			}
		})
	}
}

func TestConfigLoadAndValidate_accountFilePath(t *testing.T) {
	config := &transport_tpg.Config{
		Credentials: transport_tpg.TestFakeCredentialsPath,
		Project:     "my-gce-project",
		Region:      "us-central1",
	}

	transport_tpg.ConfigureBasePaths(config)

	err := config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestConfigLoadAndValidate_accountFileJSON(t *testing.T) {
	contents, err := ioutil.ReadFile(transport_tpg.TestFakeCredentialsPath)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	config := &transport_tpg.Config{
		Credentials: string(contents),
		Project:     "my-gce-project",
		Region:      "us-central1",
	}

	transport_tpg.ConfigureBasePaths(config)

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestConfigLoadAndValidate_accountFileJSONInvalid(t *testing.T) {
	config := &transport_tpg.Config{
		Credentials: "{this is not json}",
		Project:     "my-gce-project",
		Region:      "us-central1",
	}

	transport_tpg.ConfigureBasePaths(config)

	if config.LoadAndValidate(context.Background()) == nil {
		t.Fatalf("expected error, but got nil")
	}
}

func TestAccConfigLoadValidate_credentials(t *testing.T) {
	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	acctest.AccTestPreCheck(t)

	creds := envvar.GetTestCredsFromEnv()
	proj := envvar.GetTestProjectFromEnv()

	config := &transport_tpg.Config{
		Credentials: creds,
		Project:     proj,
		Region:      "us-central1",
	}

	transport_tpg.ConfigureBasePaths(config)

	err := config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = config.NewComputeClient(config.UserAgent).Zones.Get(proj, "us-central1-a").Do()
	if err != nil {
		t.Fatalf("expected call with loaded config client to work, got error: %s", err)
	}
}

func TestAccConfigLoadValidate_impersonated(t *testing.T) {
	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	acctest.AccTestPreCheck(t)

	serviceaccount := transport_tpg.MultiEnvSearch([]string{"IMPERSONATE_SERVICE_ACCOUNT_ACCTEST"})
	creds := envvar.GetTestCredsFromEnv()
	proj := envvar.GetTestProjectFromEnv()

	config := &transport_tpg.Config{
		Credentials:               creds,
		ImpersonateServiceAccount: serviceaccount,
		Project:                   proj,
		Region:                    "us-central1",
	}

	transport_tpg.ConfigureBasePaths(config)

	err := config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = config.NewComputeClient(config.UserAgent).Zones.Get(proj, "us-central1-a").Do()
	if err != nil {
		t.Fatalf("expected API call with loaded config to work, got error: %s", err)
	}
}

func TestAccConfigLoadValidate_accessTokenImpersonated(t *testing.T) {
	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	acctest.AccTestPreCheck(t)

	creds := envvar.GetTestCredsFromEnv()
	proj := envvar.GetTestProjectFromEnv()
	serviceaccount := transport_tpg.MultiEnvSearch([]string{"IMPERSONATE_SERVICE_ACCOUNT_ACCTEST"})

	c, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(creds), transport_tpg.DefaultClientScopes...)
	if err != nil {
		t.Fatalf("invalid test credentials: %s", err)
	}

	token, err := c.TokenSource.Token()
	if err != nil {
		t.Fatalf("Unable to generate test access token: %s", err)
	}

	config := &transport_tpg.Config{
		AccessToken:               token.AccessToken,
		ImpersonateServiceAccount: serviceaccount,
		Project:                   proj,
		Region:                    "us-central1",
	}

	transport_tpg.ConfigureBasePaths(config)

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = config.NewComputeClient(config.UserAgent).Zones.Get(proj, "us-central1-a").Do()
	if err != nil {
		t.Fatalf("expected API call with loaded config to work, got error: %s", err)
	}
}

func TestAccConfigLoadValidate_accessToken(t *testing.T) {
	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	acctest.AccTestPreCheck(t)

	creds := envvar.GetTestCredsFromEnv()
	proj := envvar.GetTestProjectFromEnv()

	c, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(creds), testOauthScope)
	if err != nil {
		t.Fatalf("invalid test credentials: %s", err)
	}

	token, err := c.TokenSource.Token()
	if err != nil {
		t.Fatalf("Unable to generate test access token: %s", err)
	}

	config := &transport_tpg.Config{
		AccessToken: token.AccessToken,
		Project:     proj,
		Region:      "us-central1",
	}

	transport_tpg.ConfigureBasePaths(config)

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = config.NewComputeClient(config.UserAgent).Zones.Get(proj, "us-central1-a").Do()
	if err != nil {
		t.Fatalf("expected API call with loaded config to work, got error: %s", err)
	}
}

func TestConfigLoadAndValidate_customScopes(t *testing.T) {
	config := &transport_tpg.Config{
		Credentials: transport_tpg.TestFakeCredentialsPath,
		Project:     "my-gce-project",
		Region:      "us-central1",
		Scopes:      []string{"https://www.googleapis.com/auth/compute"},
	}

	transport_tpg.ConfigureBasePaths(config)

	err := config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(config.Scopes) != 1 {
		t.Fatalf("expected 1 scope, got %d scopes: %v", len(config.Scopes), config.Scopes)
	}
	if config.Scopes[0] != "https://www.googleapis.com/auth/compute" {
		t.Fatalf("expected scope to be %q, got %q", "https://www.googleapis.com/auth/compute", config.Scopes[0])
	}
}

func TestConfigLoadAndValidate_defaultBatchingConfig(t *testing.T) {
	// Use default batching config
	batchCfg, err := transport_tpg.ExpandProviderBatchingConfig(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	config := &transport_tpg.Config{
		Credentials:    transport_tpg.TestFakeCredentialsPath,
		Project:        "my-gce-project",
		Region:         "us-central1",
		BatchingConfig: batchCfg,
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedDur := time.Second * transport_tpg.DefaultBatchSendIntervalSec
	if config.RequestBatcherServiceUsage.SendAfter != expectedDur {
		t.Fatalf("expected SendAfter to be %d seconds, got %v",
			transport_tpg.DefaultBatchSendIntervalSec,
			config.RequestBatcherServiceUsage.SendAfter)
	}
}

func TestConfigLoadAndValidate_customBatchingConfig(t *testing.T) {
	batchCfg, err := transport_tpg.ExpandProviderBatchingConfig([]interface{}{
		map[string]interface{}{
			"send_after":      "1s",
			"enable_batching": false,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if batchCfg.SendAfter != time.Second {
		t.Fatalf("expected batchCfg SendAfter to be 1 second, got %v", batchCfg.SendAfter)
	}
	if batchCfg.EnableBatching {
		t.Fatalf("expected EnableBatching to be false")
	}

	config := &transport_tpg.Config{
		Credentials:    transport_tpg.TestFakeCredentialsPath,
		Project:        "my-gce-project",
		Region:         "us-central1",
		BatchingConfig: batchCfg,
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedDur := time.Second * 1
	if config.RequestBatcherServiceUsage.SendAfter != expectedDur {
		t.Fatalf("expected SendAfter to be %d seconds, got %v",
			1,
			config.RequestBatcherServiceUsage.SendAfter)
	}

	if config.RequestBatcherServiceUsage.EnableBatching {
		t.Fatalf("expected EnableBatching to be false")
	}
}

func TestRemoveBasePathVersion(t *testing.T) {
	cases := []struct {
		BaseURL  string
		Expected string
	}{
		{"https://www.googleapis.com/compute/version_v1/", "https://www.googleapis.com/compute/"},
		{"https://runtimeconfig.googleapis.com/v1beta1/", "https://runtimeconfig.googleapis.com/"},
		{"https://www.googleapis.com/compute/v1/", "https://www.googleapis.com/compute/"},
		{"https://staging-version.googleapis.com/", "https://staging-version.googleapis.com/"},
		// For URLs with any parts, the last part is always removed- it's assumed to be the version.
		{"https://runtimeconfig.googleapis.com/runtimeconfig/", "https://runtimeconfig.googleapis.com/"},
	}

	for _, c := range cases {
		if c.Expected != transport_tpg.RemoveBasePathVersion(c.BaseURL) {
			t.Errorf("replace url failed: got %s wanted %s", transport_tpg.RemoveBasePathVersion(c.BaseURL), c.Expected)
		}
	}
}

func TestGetRegionFromRegionSelfLink(t *testing.T) {
	cases := map[string]struct {
		Input          string
		ExpectedOutput string
	}{
		"A short region name is returned unchanged": {
			Input:          "us-central1",
			ExpectedOutput: "us-central1",
		},
		"A selflink is shortened to a region name": {
			Input:          "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1",
			ExpectedOutput: "us-central1",
		},
		"Logic is specific to region selflinks; zone selflinks are not shortened": {
			Input:          "https://www.googleapis.com/compute/v1/projects/my-project/zones/asia-east1-a",
			ExpectedOutput: "https://www.googleapis.com/compute/v1/projects/my-project/zones/asia-east1-a",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			region := transport_tpg.GetRegionFromRegionSelfLink(tc.Input)

			if region != tc.ExpectedOutput {
				t.Fatalf("want %s,  got %s", region, tc.ExpectedOutput)
			}
		})
	}
}
