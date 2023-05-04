package google

import (
	"context"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acctest_tpg "github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

type ResourceDataMock struct {
	FieldsInSchema      map[string]interface{}
	FieldsWithHasChange []string
	id                  string
}

func (d *ResourceDataMock) HasChange(key string) bool {
	exists := false
	for _, val := range d.FieldsWithHasChange {
		if key == val {
			exists = true
		}
	}

	return exists
}

func (d *ResourceDataMock) Get(key string) interface{} {
	v, _ := d.GetOk(key)
	return v
}

func (d *ResourceDataMock) GetOk(key string) (interface{}, bool) {
	v, ok := d.GetOkExists(key)
	if ok && !tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
		return v, true
	} else {
		return v, false
	}
}

func (d *ResourceDataMock) GetOkExists(key string) (interface{}, bool) {
	for k, v := range d.FieldsInSchema {
		if key == k {
			return v, true
		}
	}

	return nil, false
}

func (d *ResourceDataMock) Set(key string, value interface{}) error {
	d.FieldsInSchema[key] = value
	return nil
}

func (d *ResourceDataMock) SetId(v string) {
	d.id = v
}

func (d *ResourceDataMock) Id() string {
	return d.id
}

func (d *ResourceDataMock) GetProviderMeta(dst interface{}) error {
	return nil
}

func (d *ResourceDataMock) Timeout(key string) time.Duration {
	return time.Duration(1)
}

type ResourceDiffMock struct {
	Before     map[string]interface{}
	After      map[string]interface{}
	Cleared    map[string]interface{}
	IsForceNew bool
}

func (d *ResourceDiffMock) GetChange(key string) (interface{}, interface{}) {
	return d.Before[key], d.After[key]
}

func (d *ResourceDiffMock) HasChange(key string) bool {
	old, new := d.GetChange(key)
	return old != new
}

func (d *ResourceDiffMock) Get(key string) interface{} {
	return d.After[key]
}

func (d *ResourceDiffMock) GetOk(key string) (interface{}, bool) {
	v, ok := d.After[key]
	return v, ok
}

func (d *ResourceDiffMock) Clear(key string) error {
	if d.Cleared == nil {
		d.Cleared = map[string]interface{}{}
	}
	d.Cleared[key] = true
	return nil
}

func (d *ResourceDiffMock) ForceNew(key string) error {
	d.IsForceNew = true
	return nil
}

func CheckDataSourceStateMatchesResourceState(dataSourceName, resourceName string) func(*terraform.State) error {
	return acctest_tpg.CheckDataSourceStateMatchesResourceState(dataSourceName, resourceName)
}

func CheckDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName string, ignoreFields map[string]struct{}) func(*terraform.State) error {
	return acctest_tpg.CheckDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName, ignoreFields)
}

// General test utils

// MuxedProviders returns the correct test provider (between the sdk version or the framework version)
func MuxedProviders(testName string) (func() tfprotov5.ProviderServer, error) {
	ctx := context.Background()

	providers := []func() tfprotov5.ProviderServer{
		providerserver.NewProtocol5(NewFrameworkTestProvider(testName)), // framework provider
		GetSDKProvider(testName).GRPCProvider,                           // sdk provider
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		return nil, err
	}

	return muxServer.ProviderServer, nil
}

func RandString(t *testing.T, length int) string {
	if !acctest_tpg.IsVcrEnabled() {
		return acctest.RandString(length)
	}
	envPath := os.Getenv("VCR_PATH")
	vcrMode := os.Getenv("VCR_MODE")
	s, err := vcrSource(t, envPath, vcrMode)
	if err != nil {
		// At this point we haven't created any resources, so fail fast
		t.Fatal(err)
	}

	r := rand.New(s.source)
	result := make([]byte, length)
	set := "abcdefghijklmnopqrstuvwxyz012346789"
	for i := 0; i < length; i++ {
		result[i] = set[r.Intn(len(set))]
	}
	return string(result)
}

func RandInt(t *testing.T) int {
	if !acctest_tpg.IsVcrEnabled() {
		return acctest.RandInt()
	}
	envPath := os.Getenv("VCR_PATH")
	vcrMode := os.Getenv("VCR_MODE")
	s, err := vcrSource(t, envPath, vcrMode)
	if err != nil {
		// At this point we haven't created any resources, so fail fast
		t.Fatal(err)
	}

	return rand.New(s.source).Int()
}

// ProtoV5ProviderFactories returns a muxed ProviderServer that uses the provider code from this repo (SDK and plugin-framework).
// Used to set ProtoV5ProviderFactories in a resource.TestStep within an acceptance test.
func ProtoV5ProviderFactories(t *testing.T) map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"google": func() (tfprotov5.ProviderServer, error) {
			provider, err := MuxedProviders(t.Name())
			return provider(), err
		},
	}
}

// ProtoV5ProviderBetaFactories returns the same as ProtoV5ProviderFactories only the provider is mapped with
// "google-beta" to ensure that registry examples use `google-beta` if the example is versioned as beta;
// normal beta tests should continue to use ProtoV5ProviderFactories
func ProtoV5ProviderBetaFactories(t *testing.T) map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){}
}

type TimeoutError struct {
	timeout bool
}

func (e *TimeoutError) Timeout() bool {
	return e.timeout
}

func (e *TimeoutError) Error() string {
	return "timeout error"
}
