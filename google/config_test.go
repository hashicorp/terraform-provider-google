package google

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"golang.org/x/oauth2/google"
)

const testFakeCredentialsPath = "./test-fixtures/fake_account.json"
const testOauthScope = "https://www.googleapis.com/auth/compute"

func TestConfigLoadAndValidate_accountFilePath(t *testing.T) {
	config := &Config{
		Credentials: testFakeCredentialsPath,
		Project:     "my-gce-project",
		Region:      "us-central1",
	}

	ConfigureBasePaths(config)

	err := config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestConfigLoadAndValidate_accountFileJSON(t *testing.T) {
	contents, err := ioutil.ReadFile(testFakeCredentialsPath)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	config := &Config{
		Credentials: string(contents),
		Project:     "my-gce-project",
		Region:      "us-central1",
	}

	ConfigureBasePaths(config)

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestConfigLoadAndValidate_accountFileJSONInvalid(t *testing.T) {
	config := &Config{
		Credentials: "{this is not json}",
		Project:     "my-gce-project",
		Region:      "us-central1",
	}

	ConfigureBasePaths(config)

	if config.LoadAndValidate(context.Background()) == nil {
		t.Fatalf("expected error, but got nil")
	}
}

func TestAccConfigLoadValidate_credentials(t *testing.T) {
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Network access not allowed; use %s=1 to enable", resource.TestEnvVar))
	}
	testAccPreCheck(t)

	creds := getTestCredsFromEnv()
	proj := getTestProjectFromEnv()

	config := &Config{
		Credentials: creds,
		Project:     proj,
		Region:      "us-central1",
	}

	ConfigureBasePaths(config)

	err := config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = config.clientCompute.Zones.Get(proj, "us-central1-a").Do()
	if err != nil {
		t.Fatalf("expected call with loaded config client to work, got error: %s", err)
	}
}

func TestAccConfigLoadValidate_accessToken(t *testing.T) {
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Network access not allowed; use %s=1 to enable", resource.TestEnvVar))
	}
	testAccPreCheck(t)

	creds := getTestCredsFromEnv()
	proj := getTestProjectFromEnv()

	c, err := google.CredentialsFromJSON(context.Background(), []byte(creds), testOauthScope)
	if err != nil {
		t.Fatalf("invalid test credentials: %s", err)
	}

	token, err := c.TokenSource.Token()
	if err != nil {
		t.Fatalf("Unable to generate test access token: %s", err)
	}

	config := &Config{
		AccessToken: token.AccessToken,
		Project:     proj,
		Region:      "us-central1",
	}

	ConfigureBasePaths(config)

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = config.clientCompute.Zones.Get(proj, "us-central1-a").Do()
	if err != nil {
		t.Fatalf("expected API call with loaded config to work, got error: %s", err)
	}
}

func TestConfigLoadAndValidate_customScopes(t *testing.T) {
	config := &Config{
		Credentials: testFakeCredentialsPath,
		Project:     "my-gce-project",
		Region:      "us-central1",
		Scopes:      []string{"https://www.googleapis.com/auth/compute"},
	}

	ConfigureBasePaths(config)

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
	batchCfg, err := expandProviderBatchingConfig(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	config := &Config{
		Credentials:    testFakeCredentialsPath,
		Project:        "my-gce-project",
		Region:         "us-central1",
		BatchingConfig: batchCfg,
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedDur := time.Second * defaultBatchSendIntervalSec
	if config.requestBatcherServiceUsage.sendAfter != expectedDur {
		t.Fatalf("expected sendAfter to be %d seconds, got %v",
			defaultBatchSendIntervalSec,
			config.requestBatcherServiceUsage.sendAfter)
	}
}

func TestConfigLoadAndValidate_customBatchingConfig(t *testing.T) {
	batchCfg, err := expandProviderBatchingConfig([]interface{}{
		map[string]interface{}{
			"send_after":      "1s",
			"enable_batching": false,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if batchCfg.sendAfter != time.Second {
		t.Fatalf("expected batchCfg sendAfter to be 1 second, got %v", batchCfg.sendAfter)
	}
	if batchCfg.enableBatching {
		t.Fatalf("expected enableBatching to be false")
	}

	config := &Config{
		Credentials:    testFakeCredentialsPath,
		Project:        "my-gce-project",
		Region:         "us-central1",
		BatchingConfig: batchCfg,
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedDur := time.Second * 1
	if config.requestBatcherServiceUsage.sendAfter != expectedDur {
		t.Fatalf("expected sendAfter to be %d seconds, got %v",
			1,
			config.requestBatcherServiceUsage.sendAfter)
	}

	if config.requestBatcherServiceUsage.enableBatching {
		t.Fatalf("expected enableBatching to be false")
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
		if c.Expected != removeBasePathVersion(c.BaseURL) {
			t.Errorf("replace url failed: got %s wanted %s", removeBasePathVersion(c.BaseURL), c.Expected)
		}
	}
}
