// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package composer

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-google/google/sweeper"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"google.golang.org/api/storage/v1"
)

func init() {
	sweeper.AddTestSweepers("gcp_composer_environment", testSweepComposerResources)
}

/**
 * CLEAN UP HELPER FUNCTIONS
 * Because the environments are flaky and bucket deletion rates can be
 * rate-limited, for now just warn instead of returning actual errors.
 */
func testSweepComposerResources(region string) error {
	config, err := sweeper.SharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting shared config for region: %s", err)
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	// us-central is passed as the region for our sweepers, but there are also
	// many tests that use the us-east1 region
	regions := []string{"us-central1", "us-east1"}
	for _, r := range regions {
		// Environments need to be cleaned up because the service is flaky.
		if err := testSweepComposerEnvironments(config, r); err != nil {
			log.Printf("[WARNING] unable to clean up all environments: %s", err)
		}

		// Buckets need to be cleaned up because they just don't get deleted on purpose.
		if err := testSweepComposerEnvironmentBuckets(config, r); err != nil {
			log.Printf("[WARNING] unable to clean up all environment storage buckets: %s", err)
		}
	}

	return nil
}

func testSweepComposerEnvironments(config *transport_tpg.Config, region string) error {
	found, err := config.NewComposerClient(config.UserAgent).Projects.Locations.Environments.List(
		fmt.Sprintf("projects/%s/locations/%s", config.Project, region)).Do()
	if err != nil {
		return fmt.Errorf("error listing storage buckets for composer environment: %s", err)
	}

	if len(found.Environments) == 0 {
		log.Printf("composer: no environments need to be cleaned up")
		return nil
	}

	log.Printf("composer: %d environments need to be cleaned up", len(found.Environments))

	var allErrors error
	for _, e := range found.Environments {
		createdAt, err := time.Parse(time.RFC3339Nano, e.CreateTime)
		if err != nil {
			return fmt.Errorf("composer: environment %q has invalid create time %q", e.Name, e.CreateTime)
		}
		// Skip environments that were created in same day
		// This sweeper should really only clean out very old environments.
		if time.Since(createdAt) < time.Hour*24 {
			log.Printf("composer: skipped environment %q, it was created today", e.Name)
			continue
		}

		switch e.State {
		case "CREATING":
			fallthrough
		case "UPDATING":
			log.Printf("composer: skipping pending Environment %q with state %q", e.Name, e.State)
		case "DELETING":
			log.Printf("composer: skipping pending Environment %q that is currently deleting", e.Name)
		case "RUNNING":
			fallthrough
		case "ERROR":
			fallthrough
		default:
			op, deleteErr := config.NewComposerClient(config.UserAgent).Projects.Locations.Environments.Delete(e.Name).Do()
			if deleteErr != nil {
				allErrors = multierror.Append(allErrors, fmt.Errorf("composer: unable to delete environment %q: %s", e.Name, deleteErr))
				continue
			}
			waitErr := ComposerOperationWaitTime(config, op, config.Project, "Sweeping old test environments", config.UserAgent, 10*time.Minute)
			if waitErr != nil {
				allErrors = multierror.Append(allErrors, fmt.Errorf("composer: unable to delete environment %q: %s", e.Name, waitErr))
			}
		}
	}
	return allErrors
}

func testSweepComposerEnvironmentBuckets(config *transport_tpg.Config, region string) error {
	artifactsBName := fmt.Sprintf("artifacts.%s.appspot.com", config.Project)
	artifactBucket, err := config.NewStorageClient(config.UserAgent).Buckets.Get(artifactsBName).Do()
	if err != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			log.Printf("composer environment bucket %q not found, doesn't need to be cleaned up", artifactsBName)
		} else {
			return err
		}
	} else if err = testSweepComposerEnvironmentCleanUpBucket(config, artifactBucket); err != nil {
		return err
	}

	found, err := config.NewStorageClient(config.UserAgent).Buckets.List(config.Project).Prefix(region).Do()
	if err != nil {
		return fmt.Errorf("error listing storage buckets created when testing composer environment: %s", err)
	}
	if len(found.Items) == 0 {
		log.Printf("No environment-specific buckets need to be cleaned up")
		return nil
	}

	for _, bucket := range found.Items {
		if _, ok := bucket.Labels["goog-composer-environment"]; !ok {
			continue
		}
		if err := testSweepComposerEnvironmentCleanUpBucket(config, bucket); err != nil {
			return err
		}
	}
	return nil
}

func testSweepComposerEnvironmentCleanUpBucket(config *transport_tpg.Config, bucket *storage.Bucket) error {
	var allErrors error
	objList, err := config.NewStorageClient(config.UserAgent).Objects.List(bucket.Name).Do()
	if err != nil {
		allErrors = multierror.Append(allErrors,
			fmt.Errorf("Unable to list objects to delete for bucket %q: %s", bucket.Name, err))
	}

	for _, o := range objList.Items {
		if err := config.NewStorageClient(config.UserAgent).Objects.Delete(bucket.Name, o.Name).Do(); err != nil {
			allErrors = multierror.Append(allErrors,
				fmt.Errorf("Unable to delete object %q from bucket %q: %s", o.Name, bucket.Name, err))
		}
	}

	if err := config.NewStorageClient(config.UserAgent).Buckets.Delete(bucket.Name).Do(); err != nil {
		allErrors = multierror.Append(allErrors, fmt.Errorf("Unable to delete bucket %q: %s", bucket.Name, err))
	}

	if allErrors != nil {
		return fmt.Errorf("Unable to clean up bucket %q: %v", bucket.Name, allErrors)
	}

	log.Printf("Cleaned up bucket %q for composer environment tests", bucket.Name)
	return nil
}
