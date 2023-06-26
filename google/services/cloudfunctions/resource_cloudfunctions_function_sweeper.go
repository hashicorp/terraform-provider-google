// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudfunctions

import (
	"fmt"
	"log"
	"os"
	"strings"

	"io/ioutil"

	"github.com/hashicorp/terraform-provider-google/google/sweeper"
)

const testFunctionsSourceArchivePrefix = "cloudfunczip"

func init() {
	sweeper.AddTestSweepers("gcp_cloud_function_source_archive", sweepCloudFunctionSourceZipArchives)
}

func sweepCloudFunctionSourceZipArchives(_ string) error {
	files, err := ioutil.ReadDir(os.TempDir())
	if err != nil {
		log.Printf("Error reading files: %s", err)
		return nil
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name(), testFunctionsSourceArchivePrefix) {
			filepath := fmt.Sprintf("%s/%s", os.TempDir(), f.Name())
			if err := os.Remove(filepath); err != nil {
				log.Printf("Error removing files: %s", err)
				return nil
			}
			log.Printf("[INFO] cloud functions sweeper removed old file %s", filepath)
		}
	}
	return nil
}
