// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/sweeper"
)

var (
	TestPrefix = "tf-test"
)

func init() {
	// SKIP_PROJECT_SWEEPER can be set for a sweeper run to prevent it from
	// sweeping projects. This can be useful when running sweepers in
	// organizations where acceptance tests intiated by another project may
	// already be in-progress.
	// Example: SKIP_PROJECT_SWEEPER=1 go test ./google -v -sweep=us-central1 -sweep-run=
	if os.Getenv("SKIP_PROJECT_SWEEPER") != "" {
		return
	}

	sweeper.AddTestSweepers("GoogleProject", testSweepProject)
}

func testSweepProject(region string) error {
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

	org := envvar.UnsafeGetTestOrgFromEnv()
	if org == "" {
		log.Printf("[INFO][SWEEPER_LOG] no organization set, failing project sweeper")
		return fmt.Errorf("no organization set")
	}

	token := ""
	for paginate := true; paginate; {
		// Filter for projects with test prefix
		filter := fmt.Sprintf("id:\"%s*\" -lifecycleState:DELETE_REQUESTED parent.id:%v", TestPrefix, org)
		found, err := config.NewResourceManagerClient(config.UserAgent).Projects.List().Filter(filter).PageToken(token).Do()
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error listing projects: %s", err)
			return nil
		}

		for _, project := range found.Projects {
			log.Printf("[INFO][SWEEPER_LOG] Sweeping Project id: %s", project.ProjectId)
			_, err := config.NewResourceManagerClient(config.UserAgent).Projects.Delete(project.ProjectId).Do()
			if err != nil {
				log.Printf("[INFO][SWEEPER_LOG] Error, failed to delete project %s: %s", project.Name, err)
				continue
			}
		}
		token = found.NextPageToken
		paginate = token != ""
	}

	return nil
}
