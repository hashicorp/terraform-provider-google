// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwprovider"
	tpgprovider "github.com/hashicorp/terraform-provider-google/google/provider"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	fwDiags "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func IsVcrEnabled() bool {
	envPath := os.Getenv("VCR_PATH")
	vcrMode := os.Getenv("VCR_MODE")
	return envPath != "" && vcrMode != ""
}

var configsLock = sync.RWMutex{}
var sourcesLock = sync.RWMutex{}

var configs map[string]*transport_tpg.Config
var fwProviders map[string]*frameworkTestProvider

var sources map[string]VcrSource

// VcrSource is a source for a given VCR test with the value that seeded it
type VcrSource struct {
	seed   int64
	source rand.Source
}

// Produces a rand.Source for VCR testing based on the given mode.
// In RECORDING mode, generates a new seed and saves it to a file, using the seed for the source
// In REPLAYING mode, reads a seed from a file and creates a source from it
func vcrSource(t *testing.T, path, mode string) (*VcrSource, error) {
	sourcesLock.RLock()
	s, ok := sources[t.Name()]
	sourcesLock.RUnlock()
	if ok {
		return &s, nil
	}
	tflog.Debug(context.Background(), fmt.Sprintf("VCR_MODE: %s", mode))
	switch mode {
	case "RECORDING":
		seed := rand.Int63()
		s := rand.NewSource(seed)
		vcrSource := VcrSource{seed: seed, source: s}
		sourcesLock.Lock()
		sources[t.Name()] = vcrSource
		sourcesLock.Unlock()
		return &vcrSource, nil
	case "REPLAYING":
		seed, err := readSeedFromFile(vcrSeedFile(path, t.Name()))
		if err != nil {
			return nil, fmt.Errorf("no cassette found on disk for %s, please replay this testcase in recording mode - %w", t.Name(), err)
		}
		s := rand.NewSource(seed)
		vcrSource := VcrSource{seed: seed, source: s}
		sourcesLock.Lock()
		sources[t.Name()] = vcrSource
		sourcesLock.Unlock()
		return &vcrSource, nil
	default:
		log.Printf("[DEBUG] No valid environment var set for VCR_MODE, expected RECORDING or REPLAYING, skipping VCR. VCR_MODE: %s", mode)
		return nil, errors.New("No valid VCR_MODE set")
	}
}

func readSeedFromFile(fileName string) (int64, error) {
	// Max number of digits for int64 is 19
	data := make([]byte, 19)
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	_, err = f.Read(data)
	if err != nil {
		return 0, err
	}
	// Remove NULL characters from seed
	data = bytes.Trim(data, "\x00")
	seed := string(data)
	return tpgresource.StringToFixed64(seed)
}

func writeSeedToFile(seed int64, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(strconv.FormatInt(seed, 10))
	if err != nil {
		return err
	}
	return nil
}

// Retrieves a unique test name used for writing files
// replaces all `/` characters that would cause filepath issues
// This matters during tests that dispatch multiple tests, for example TestAccLoggingFolderExclusion
func vcrSeedFile(path, name string) string {
	return filepath.Join(path, fmt.Sprintf("%s.seed", vcrFileName(name)))
}

func vcrFileName(name string) string {
	return strings.ReplaceAll(name, "/", "_")
}

// VcrTest is a wrapper for resource.Test to swap out providers for VCR providers and handle VCR specific things
// Can be called when VCR is not enabled, and it will behave as normal
func VcrTest(t *testing.T, c resource.TestCase) {
	if IsVcrEnabled() {
		defer closeRecorder(t)
	} else if isReleaseDiffEnabled() {
		c = initializeReleaseDiffTest(c, t.Name())
	}
	resource.Test(t, c)
}

// We need to explicitly close the VCR recorder to save the cassette
func closeRecorder(t *testing.T) {
	configsLock.RLock()
	config, ok := configs[t.Name()]
	configsLock.RUnlock()
	if ok {
		// We did not cache the config if it does not use VCR
		if !t.Failed() && IsVcrEnabled() {
			// If a test succeeds, write new seed/yaml to files
			err := config.Client.Transport.(*recorder.Recorder).Stop()
			if err != nil {
				t.Error(err)
			}
			envPath := os.Getenv("VCR_PATH")

			sourcesLock.RLock()
			vcrSource, ok := sources[t.Name()]
			sourcesLock.RUnlock()
			if ok {
				err = writeSeedToFile(vcrSource.seed, vcrSeedFile(envPath, t.Name()))
				if err != nil {
					t.Error(err)
				}
			}
		}
		// Clean up test config
		configsLock.Lock()
		delete(configs, t.Name())
		configsLock.Unlock()

		sourcesLock.Lock()
		delete(sources, t.Name())
		sourcesLock.Unlock()
	}

	configsLock.RLock()
	fwProvider, fwOk := fwProviders[t.Name()]
	configsLock.RUnlock()
	if fwOk {
		// We did not cache the config if it does not use VCR
		if !t.Failed() && IsVcrEnabled() {
			// If a test succeeds, write new seed/yaml to files
			err := fwProvider.Client.Transport.(*recorder.Recorder).Stop()
			if err != nil {
				t.Error(err)
			}
			envPath := os.Getenv("VCR_PATH")

			sourcesLock.RLock()
			vcrSource, ok := sources[t.Name()]
			sourcesLock.RUnlock()
			if ok {
				err = writeSeedToFile(vcrSource.seed, vcrSeedFile(envPath, t.Name()))
				if err != nil {
					t.Error(err)
				}
			}
		}
		// Clean up test config
		configsLock.Lock()
		delete(fwProviders, t.Name())
		configsLock.Unlock()

		sourcesLock.Lock()
		delete(sources, t.Name())
		sourcesLock.Unlock()
	}
}

func isReleaseDiffEnabled() bool {
	releaseDiff := os.Getenv("RELEASE_DIFF")
	return releaseDiff != ""
}

func initializeReleaseDiffTest(c resource.TestCase, testName string) resource.TestCase {
	var releaseProvider string
	packagePath := fmt.Sprint(reflect.TypeOf(transport_tpg.Config{}).PkgPath())
	if strings.Contains(packagePath, "google-beta") {
		releaseProvider = "google-beta"
	} else {
		releaseProvider = "google"
	}

	if c.ExternalProviders != nil {
		c.ExternalProviders[releaseProvider] = resource.ExternalProvider{}
	} else {
		c.ExternalProviders = map[string]resource.ExternalProvider{
			releaseProvider: {},
		}
	}

	localProviderName := "google-local"
	if c.Providers != nil {
		c.Providers = map[string]*schema.Provider{
			localProviderName: GetSDKProvider(testName),
		}
		c.ProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
			localProviderName: func() (tfprotov5.ProviderServer, error) {
				return nil, nil
			},
		}
	} else {
		c.ProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
			localProviderName: func() (tfprotov5.ProviderServer, error) {
				provider, err := MuxedProviders(testName)
				return provider(), err
			},
		}
	}

	var replacementSteps []resource.TestStep
	for _, testStep := range c.Steps {
		if testStep.Config != "" {
			ogConfig := testStep.Config
			testStep.Config = reformConfigWithProvider(ogConfig, localProviderName)
			if testStep.ExpectError == nil && testStep.PlanOnly == false {
				newStep := resource.TestStep{
					Config: reformConfigWithProvider(ogConfig, releaseProvider),
				}
				testStep.PlanOnly = true
				testStep.ExpectNonEmptyPlan = false
				replacementSteps = append(replacementSteps, newStep)
			}
			replacementSteps = append(replacementSteps, testStep)
		} else {
			replacementSteps = append(replacementSteps, testStep)
		}
	}

	c.Steps = replacementSteps

	return c
}

func reformConfigWithProvider(config, provider string) string {
	configBytes := []byte(config)
	providerReplacement := fmt.Sprintf("provider = %s", provider)
	providerReplacementBytes := []byte(providerReplacement)
	providerBlock := regexp.MustCompile(`provider *=.*google-beta.*`)

	if providerBlock.Match(configBytes) {
		return string(providerBlock.ReplaceAll(configBytes, providerReplacementBytes))
	}

	providerReplacement = fmt.Sprintf("${1}\n\t%s", providerReplacement)
	providerReplacementBytes = []byte(providerReplacement)
	resourceHeader := regexp.MustCompile(`(resource .*google_.* .*\w+.*\{.*)`)
	return string(resourceHeader.ReplaceAll(configBytes, providerReplacementBytes))
}

func HandleVCRConfiguration(ctx context.Context, testName string, rndTripper http.RoundTripper, pollInterval time.Duration) (time.Duration, http.RoundTripper, fwDiags.Diagnostics) {
	var diags fwDiags.Diagnostics
	var vcrMode recorder.Mode
	switch vcrEnv := os.Getenv("VCR_MODE"); vcrEnv {
	case "RECORDING":
		vcrMode = recorder.ModeRecording
	case "REPLAYING":
		vcrMode = recorder.ModeReplaying
		// When replaying, set the poll interval low to speed up tests
		pollInterval = 10 * time.Millisecond
	default:
		tflog.Debug(ctx, fmt.Sprintf("No valid environment var set for VCR_MODE, expected RECORDING or REPLAYING, skipping VCR. VCR_MODE: %s", vcrEnv))
		return pollInterval, rndTripper, diags
	}

	envPath := os.Getenv("VCR_PATH")
	if envPath == "" {
		tflog.Debug(ctx, "No environment var set for VCR_PATH, skipping VCR")
		return pollInterval, rndTripper, diags
	}
	path := filepath.Join(envPath, vcrFileName(testName))

	rec, err := recorder.NewAsMode(path, vcrMode, rndTripper)
	if err != nil {
		diags.AddError("error creating record as new mode", err.Error())
		return pollInterval, rndTripper, diags
	}
	// Defines how VCR will match requests to responses.
	rec.SetMatcher(func(r *http.Request, i cassette.Request) bool {
		// Default matcher compares method and URL only
		if !cassette.DefaultMatcher(r, i) {
			return false
		}
		if r.Body == nil {
			return true
		}
		contentType := r.Header.Get("Content-Type")
		// If body contains media, don't try to compare
		if strings.Contains(contentType, "multipart/related") {
			return true
		}

		var b bytes.Buffer
		if _, err := b.ReadFrom(r.Body); err != nil {
			tflog.Debug(ctx, fmt.Sprintf("Failed to read request body from cassette: %v", err))
			return false
		}
		r.Body = ioutil.NopCloser(&b)
		reqBody := b.String()
		// If body matches identically, we are done
		if reqBody == i.Body {
			return true
		}

		// JSON might be the same, but reordered. Try parsing json and comparing
		if strings.Contains(contentType, "application/json") {
			var reqJson, cassetteJson interface{}
			if err := json.Unmarshal([]byte(reqBody), &reqJson); err != nil {
				tflog.Debug(ctx, fmt.Sprintf("Failed to unmarshall request json: %v", err))
				return false
			}
			if err := json.Unmarshal([]byte(i.Body), &cassetteJson); err != nil {
				tflog.Debug(ctx, fmt.Sprintf("Failed to unmarshall cassette json: %v", err))
				return false
			}
			return reflect.DeepEqual(reqJson, cassetteJson)
		}
		return false
	})

	return pollInterval, rec, diags
}

// MuxedProviders configures the providers, thus, if we want the providers to be configured
// to use VCR, the configure functions need to be altered. The only way to do this is to create
// test versions of the provider that will call the same configure function, only append the VCR
// configuration to it.

func NewFrameworkTestProvider(testName string) *frameworkTestProvider {
	return &frameworkTestProvider{
		FrameworkProvider: fwprovider.FrameworkProvider{
			Version: "test",
		},
		TestName: testName,
	}
}

// frameworkTestProvider is a test version of the plugin-framework version of the provider
// that embeds FrameworkProvider whose configure function we can use
// the Configure function is overwritten in the framework_provider_test file
type frameworkTestProvider struct {
	fwprovider.FrameworkProvider
	TestName string
}

// Configure is here to overwrite the FrameworkProvider configure function for VCR testing
func (p *frameworkTestProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	p.FrameworkProvider.Configure(ctx, req, resp)
	if IsVcrEnabled() {
		if resp.Diagnostics.HasError() {
			return
		}

		var diags fwDiags.Diagnostics
		p.PollInterval, p.Client.Transport, diags = HandleVCRConfiguration(ctx, p.TestName, p.Client.Transport, p.PollInterval)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		configsLock.Lock()
		fwProviders[p.TestName] = p
		configsLock.Unlock()
		return
	} else {
		tflog.Debug(ctx, "VCR_PATH or VCR_MODE not set, skipping VCR")
	}
}

func configureApiClient(ctx context.Context, p *fwprovider.FrameworkProvider, diags *fwDiags.Diagnostics) {
	var data fwmodels.ProviderModel
	var d fwDiags.Diagnostics

	// Set defaults if needed - the only attribute without a default is ImpersonateServiceAccountDelegates
	// this is a bit of a hack, but we'll just initialize it here so that it's been initialized at least
	data.ImpersonateServiceAccountDelegates, d = types.ListValue(types.StringType, []attr.Value{})
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	p.LoadAndValidateFramework(ctx, &data, "test", diags, p.Version)
}

// GetSDKProvider gets the SDK provider with an overwritten configure function to be called by MuxedProviders
func GetSDKProvider(testName string) *schema.Provider {
	prov := tpgprovider.Provider()
	if IsVcrEnabled() {
		old := prov.ConfigureContextFunc
		prov.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return getCachedConfig(ctx, d, old, testName)
		}
	} else {
		log.Print("[DEBUG] VCR_PATH or VCR_MODE not set, skipping VCR")
	}
	return prov
}

// Returns a cached config if VCR testing is enabled. This enables us to use a single HTTP transport
// for a given test, allowing for recording of HTTP interactions.
// Why this exists: schema.Provider.ConfigureFunc is called multiple times for a given test
// ConfigureFunc on our provider creates a new HTTP client and sets base paths (config.go LoadAndValidate)
// VCR requires a single HTTP client to handle all interactions so it can record and replay responses so
// this caches HTTP clients per test by replacing ConfigureFunc
func getCachedConfig(ctx context.Context, d *schema.ResourceData, configureFunc schema.ConfigureContextFunc, testName string) (*transport_tpg.Config, diag.Diagnostics) {
	configsLock.RLock()
	v, ok := configs[testName]
	configsLock.RUnlock()
	if ok {
		return v, nil
	}
	c, diags := configureFunc(ctx, d)
	if diags.HasError() {
		return nil, diags
	}

	var fwD fwDiags.Diagnostics
	config := c.(*transport_tpg.Config)
	config.PollInterval, config.Client.Transport, fwD = HandleVCRConfiguration(ctx, testName, config.Client.Transport, config.PollInterval)
	if fwD.HasError() {
		diags = append(diags, *tpgresource.FrameworkDiagsToSdkDiags(fwD)...)
		return nil, diags
	}

	configsLock.Lock()
	configs[testName] = config
	configsLock.Unlock()
	return config, nil
}
