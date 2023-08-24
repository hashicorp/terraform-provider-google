// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func CheckDataSourceStateMatchesResourceState(dataSourceName, resourceName string) func(*terraform.State) error {
	return CheckDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName, map[string]struct{}{})
}

func CheckDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName string, ignoreFields map[string]struct{}) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		errMsg := ""
		// Data sources are often derived from resources, so iterate over the resource fields to
		// make sure all fields are accounted for in the data source.
		// If a field exists in the data source but not in the resource, its expected value should
		// be checked separately.
		for k := range rsAttr {
			if _, ok := ignoreFields[k]; ok {
				continue
			}
			if k == "%" {
				continue
			}
			if dsAttr[k] != rsAttr[k] {
				// ignore data sources where an empty list is being compared against a null list.
				if k[len(k)-1:] == "#" && (dsAttr[k] == "" || dsAttr[k] == "0") && (rsAttr[k] == "" || rsAttr[k] == "0") {
					continue
				}
				errMsg += fmt.Sprintf("%s is %s; want %s\n", k, dsAttr[k], rsAttr[k])
			}
		}

		if errMsg != "" {
			return errors.New(errMsg)
		}

		return nil
	}
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
	if !IsVcrEnabled() {
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
	if !IsVcrEnabled() {
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

// This is a Printf sibling (Nprintf; Named Printf), which handles strings like
// Nprintf("Hello %{target}!", map[string]interface{}{"target":"world"}) == "Hello world!".
// This is particularly useful for generated tests, where we don't want to use Printf,
// since that would require us to generate a very particular ordering of arguments.
func Nprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		format = strings.Replace(format, "%{"+key+"}", fmt.Sprintf("%v", val), -1)
	}
	return format
}

func TestBucketName(t *testing.T) string {
	return fmt.Sprintf("%s-%d", "tf-test-bucket", RandInt(t))
}

func CreateZIPArchiveForCloudFunctionSource(t *testing.T, sourcePath string) string {
	source, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err.Error())
	}
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	f, err := w.Create("index.js")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = f.Write(source)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		t.Fatal(err.Error())
	}
	// Create temp file to write zip to
	tmpfile, err := ioutil.TempFile("", "sourceArchivePrefix")
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, err := tmpfile.Write(buf.Bytes()); err != nil {
		t.Fatal(err.Error())
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err.Error())
	}
	return tmpfile.Name()
}
