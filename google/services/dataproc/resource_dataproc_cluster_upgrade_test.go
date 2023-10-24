// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataproc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"google.golang.org/api/dataproc/v1"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// Tests schema version migration by creating a cluster with an old version of the provider (4.65.0)
// and then updating it with the current version the provider.
func TestAccDataprocClusterLabelsMigration_withoutLabels_withoutChanges(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	rnd := acctest.RandString(t, 10)
	var cluster dataproc.Cluster
	networkName := acctest.BootstrapSharedTestNetwork(t, "dataproc-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "dataproc-cluster", networkName)
	acctest.BootstrapFirewallForDataprocSharedNetwork(t, "dataproc-cluster", networkName)

	oldVersion := map[string]resource.ExternalProvider{
		"google": {
			VersionConstraint: "4.65.0", // a version that doesn't separate user defined labels and system labels
			Source:            "registry.terraform.io/hashicorp/google",
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config:            testAccDataprocCluster_withoutLabels(rnd, subnetworkName),
				ExternalProviders: oldVersion,
			},
			{
				Config:                   testAccDataprocCluster_withoutLabels(rnd, subnetworkName),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_labels", &cluster),

					resource.TestCheckNoResourceAttr("google_dataproc_cluster.with_labels", "labels.%"),
					// GCP adds three and goog-dataproc-autozone is added internally, so expect 4.
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "effective_labels.%", "4"),
				),
			},
			{
				Config:                   testAccDataprocCluster_withLabels(rnd, subnetworkName),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_labels", &cluster),

					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.%", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.key1", "value1"),
					// We only provide one, but GCP adds three and goog-dataproc-autozone is added internally, so expect 5.
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "effective_labels.%", "5"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "effective_labels.key1", "value1"),
				),
			},
		},
	})
}

func TestAccDataprocClusterLabelsMigration_withLabels_withoutChanges(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	rnd := acctest.RandString(t, 10)
	var cluster dataproc.Cluster
	networkName := acctest.BootstrapSharedTestNetwork(t, "dataproc-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "dataproc-cluster", networkName)
	acctest.BootstrapFirewallForDataprocSharedNetwork(t, "dataproc-cluster", networkName)

	oldVersion := map[string]resource.ExternalProvider{
		"google": {
			VersionConstraint: "4.65.0", // a version that doesn't separate user defined labels and system labels
			Source:            "registry.terraform.io/hashicorp/google",
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config:            testAccDataprocCluster_withLabels(rnd, subnetworkName),
				ExternalProviders: oldVersion,
			},
			{
				Config:                   testAccDataprocCluster_withLabels(rnd, subnetworkName),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_labels", &cluster),

					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.%", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.key1", "value1"),
					// We only provide one, but GCP adds three and goog-dataproc-autozone is added internally, so expect 5.
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "effective_labels.%", "5"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "effective_labels.key1", "value1"),
				),
			},
			{
				Config:                   testAccDataprocCluster_withLabelsUpdate(rnd, subnetworkName),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_labels", &cluster),

					// We only provide one, so expect 1.
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.%", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.key2", "value2"),
					// We only provide one, but GCP adds three and goog-dataproc-autozone is added internally, so expect 5.
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "effective_labels.%", "5"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "effective_labels.key2", "value2"),
				),
			},
		},
	})
}

func TestAccDataprocClusterLabelsMigration_withUpdate(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	rnd := acctest.RandString(t, 10)
	var cluster dataproc.Cluster
	networkName := acctest.BootstrapSharedTestNetwork(t, "dataproc-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "dataproc-cluster", networkName)
	acctest.BootstrapFirewallForDataprocSharedNetwork(t, "dataproc-cluster", networkName)

	oldVersion := map[string]resource.ExternalProvider{
		"google": {
			VersionConstraint: "4.65.0", // a version that doesn't separate user defined labels and system labels
			Source:            "registry.terraform.io/hashicorp/google",
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckDataprocClusterDestroy(t),
		Steps: []resource.TestStep{
			{
				Config:            testAccDataprocCluster_withoutLabels(rnd, subnetworkName),
				ExternalProviders: oldVersion,
			},
			{
				Config:                   testAccDataprocCluster_withLabels(rnd, subnetworkName),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_labels", &cluster),

					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.%", "1"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "labels.key1", "value1"),
					// We only provide one, but GCP adds three and goog-dataproc-autozone is added internally, so expect 5.
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "effective_labels.%", "5"),
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "effective_labels.key1", "value1"),
				),
			},
			{
				Config:                   testAccDataprocCluster_withoutLabels(rnd, subnetworkName),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataprocClusterExists(t, "google_dataproc_cluster.with_labels", &cluster),

					resource.TestCheckNoResourceAttr("google_dataproc_cluster.with_labels", "labels.%"),
					// We only provide one, but GCP adds three and goog-dataproc-autozone is added internally, so expect 4.
					resource.TestCheckResourceAttr("google_dataproc_cluster.with_labels", "effective_labels.%", "4"),
				),
			},
		},
	})
}
