// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-google/google/sweeper"
)

func init() {
	sweeper.AddTestSweepers("gcp_container_cluster", testSweepContainerClusters)
}

func testSweepContainerClusters(region string) error {
	config, err := sweeper.SharedConfigForRegion(region)
	if err != nil {
		log.Fatalf("error getting shared config for region: %s", err)
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	// List clusters for all zones by using "-" as the zone name
	found, err := config.NewContainerClient(config.UserAgent).Projects.Zones.Clusters.List(config.Project, "-").Do()
	if err != nil {
		log.Printf("error listing container clusters: %s", err)
		return nil
	}

	if len(found.Clusters) == 0 {
		log.Printf("No container clusters found.")
		return nil
	}

	for _, cluster := range found.Clusters {
		if sweeper.IsSweepableTestResource(cluster.Name) {
			log.Printf("Sweeping Container Cluster: %s", cluster.Name)
			clusterURL := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", config.Project, cluster.Location, cluster.Name)
			_, err := config.NewContainerClient(config.UserAgent).Projects.Locations.Clusters.Delete(clusterURL).Do()

			if err != nil {
				log.Printf("Error, failed to delete cluster %s: %s", cluster.Name, err)
				return nil
			}
		}
	}

	return nil
}
