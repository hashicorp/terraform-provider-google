// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/compute/v1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeInstanceMigrateState(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	cases := map[string]struct {
		StateVersion int
		Attributes   map[string]string
		Expected     map[string]string
	}{
		"v0.4.2 and earlier": {
			StateVersion: 0,
			Attributes: map[string]string{
				"disk.#":               "0",
				"metadata.#":           "2",
				"metadata.0.foo":       "bar",
				"metadata.1.baz":       "qux",
				"metadata.2.with.dots": "should.work",
			},
			Expected: map[string]string{
				"create_timeout":     "4",
				"metadata.foo":       "bar",
				"metadata.baz":       "qux",
				"metadata.with.dots": "should.work",
			},
		},
		"change scope from list to set": {
			StateVersion: 1,
			Attributes: map[string]string{
				"service_account.#":          "1",
				"service_account.0.email":    "xxxxxx-compute@developer.gserviceaccount.com",
				"service_account.0.scopes.#": "4",
				"service_account.0.scopes.0": "https://www.googleapis.com/auth/compute",
				"service_account.0.scopes.1": "https://www.googleapis.com/auth/datastore",
				"service_account.0.scopes.2": "https://www.googleapis.com/auth/devstorage.full_control",
				"service_account.0.scopes.3": "https://www.googleapis.com/auth/logging.write",
			},
			Expected: map[string]string{
				"create_timeout":                      "4",
				"service_account.#":                   "1",
				"service_account.0.email":             "xxxxxx-compute@developer.gserviceaccount.com",
				"service_account.0.scopes.#":          "4",
				"service_account.0.scopes.1693978638": "https://www.googleapis.com/auth/devstorage.full_control",
				"service_account.0.scopes.172152165":  "https://www.googleapis.com/auth/logging.write",
				"service_account.0.scopes.299962681":  "https://www.googleapis.com/auth/compute",
				"service_account.0.scopes.3435931483": "https://www.googleapis.com/auth/datastore",
			},
		},
		"add new create_timeout attribute": {
			StateVersion: 2,
			Attributes:   map[string]string{},
			Expected: map[string]string{
				"create_timeout": "4",
			},
		},
		"remove empty initialize_params": {
			StateVersion: 5,
			Attributes: map[string]string{
				"boot_disk.0.initialize_params.#":      "1",
				"boot_disk.0.initialize_params.0.size": "0",
			},
			Expected: map[string]string{
				"boot_disk.0.initialize_params.#": "0",
			},
		},
	}

	config := getInitializedConfig(t)

	instanceName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	instance := &compute.Instance{
		Name: instanceName,
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
		},
		MachineType: "zones/" + config.Zone + "/machineTypes/e2-medium",
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
	}
	op, err := config.NewComputeClient(config.UserAgent).Instances.Insert(config.Project, config.Zone, instance).Do()
	if err != nil {
		t.Fatalf("Error creating instance: %s", err)
	}
	waitErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "instance to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpInstance(config, instanceName, config.Zone)

	for tn, tc := range cases {
		runInstanceMigrateTest(t, instanceName, tn, tc.StateVersion, tc.Attributes, tc.Expected, config)
	}
}

func TestAccComputeInstanceMigrateState_empty(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	var is *terraform.InstanceState
	var meta interface{}

	// should handle nil
	is, err := tpgcompute.ResourceComputeInstanceMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
	if is != nil {
		t.Fatalf("expected nil instancestate, got: %#v", is)
	}

	// should handle non-nil but empty
	is = &terraform.InstanceState{}
	_, err = tpgcompute.ResourceComputeInstanceMigrateState(0, is, meta)

	if err != nil {
		t.Fatalf("err: %#v", err)
	}
}

func TestAccComputeInstanceMigrateState_bootDisk(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	config := getInitializedConfig(t)
	zone := "us-central1-f"

	// Seed test data
	instanceName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	instance := &compute.Instance{
		Name: instanceName,
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
		},
		MachineType: "zones/" + zone + "/machineTypes/e2-medium",
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
	}
	op, err := config.NewComputeClient(config.UserAgent).Instances.Insert(config.Project, zone, instance).Do()

	if err != nil {
		t.Fatalf("Error creating instance: %s", err)
	}
	waitErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "instance to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpInstance(config, instanceName, zone)

	attributes := map[string]string{
		"disk.#":                            "1",
		"disk.0.disk":                       "disk-1",
		"disk.0.type":                       "pd-ssd",
		"disk.0.auto_delete":                "true",
		"disk.0.size":                       "12",
		"disk.0.device_name":                "persistent-disk-0",
		"disk.0.disk_encryption_key_raw":    "encrypt-key",
		"disk.0.disk_encryption_key_sha256": "encrypt-key-sha",
		"zone":                              zone,
	}
	expected := map[string]string{
		"boot_disk.#":                            "1",
		"boot_disk.0.auto_delete":                "true",
		"boot_disk.0.device_name":                "persistent-disk-0",
		"boot_disk.0.disk_encryption_key_raw":    "encrypt-key",
		"boot_disk.0.disk_encryption_key_sha256": "encrypt-key-sha",
		"boot_disk.0.initialize_params.#":        "1",
		"boot_disk.0.initialize_params.0.size":   "12",
		"boot_disk.0.initialize_params.0.type":   "pd-ssd",
		"boot_disk.0.source":                     instanceName,
		"zone":                                   zone,
		"create_timeout":                         "4",
	}

	runInstanceMigrateTest(t, instanceName, "migrate disk to boot disk", 2 /* state version */, attributes, expected, config)
}

func TestAccComputeInstanceMigrateState_v4FixBootDisk(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	config := getInitializedConfig(t)
	zone := "us-central1-f"

	// Seed test data
	instanceName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	instance := &compute.Instance{
		Name: instanceName,
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
		},
		MachineType: "zones/" + zone + "/machineTypes/e2-medium",
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
	}
	op, err := config.NewComputeClient(config.UserAgent).Instances.Insert(config.Project, zone, instance).Do()

	if err != nil {
		t.Fatalf("Error creating instance: %s", err)
	}
	waitErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "instance to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpInstance(config, instanceName, zone)

	attributes := map[string]string{
		"disk.#":                            "1",
		"disk.0.disk":                       "disk-1",
		"disk.0.type":                       "pd-ssd",
		"disk.0.auto_delete":                "true",
		"disk.0.size":                       "12",
		"disk.0.device_name":                "persistent-disk-0",
		"disk.0.disk_encryption_key_raw":    "encrypt-key",
		"disk.0.disk_encryption_key_sha256": "encrypt-key-sha",
		"zone":                              zone,
	}
	expected := map[string]string{
		"boot_disk.#":                            "1",
		"boot_disk.0.auto_delete":                "true",
		"boot_disk.0.device_name":                "persistent-disk-0",
		"boot_disk.0.disk_encryption_key_raw":    "encrypt-key",
		"boot_disk.0.disk_encryption_key_sha256": "encrypt-key-sha",
		"boot_disk.0.initialize_params.#":        "1",
		"boot_disk.0.initialize_params.0.size":   "12",
		"boot_disk.0.initialize_params.0.type":   "pd-ssd",
		"boot_disk.0.source":                     instanceName,
		"zone":                                   zone,
	}

	runInstanceMigrateTest(t, instanceName, "migrate disk to boot disk", 4 /* state version */, attributes, expected, config)
}

func TestAccComputeInstanceMigrateState_attachedDiskFromSource(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	config := getInitializedConfig(t)
	zone := "us-central1-f"

	// Seed test data
	diskName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	disk := &compute.Disk{
		Name:        diskName,
		SourceImage: "projects/debian-cloud/global/images/family/debian-11",
		Zone:        zone,
	}
	op, err := config.NewComputeClient(config.UserAgent).Disks.Insert(config.Project, zone, disk).Do()
	if err != nil {
		t.Fatalf("Error creating disk: %s", err)
	}
	waitErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "disk to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpDisk(config, diskName, zone)

	instanceName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	instance := &compute.Instance{
		Name: instanceName,
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
			{
				Source: "projects/" + config.Project + "/zones/" + zone + "/disks/" + diskName,
			},
		},
		MachineType: "zones/" + zone + "/machineTypes/e2-medium",
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
	}
	op, err = config.NewComputeClient(config.UserAgent).Instances.Insert(config.Project, zone, instance).Do()
	if err != nil {
		t.Fatalf("Error creating instance: %s", err)
	}
	waitErr = tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "instance to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpInstance(config, instanceName, zone)

	attributes := map[string]string{
		"boot_disk.#":                       "1",
		"disk.#":                            "1",
		"disk.0.disk":                       diskName,
		"disk.0.device_name":                "persistent-disk-1",
		"disk.0.disk_encryption_key_raw":    "encrypt-key",
		"disk.0.disk_encryption_key_sha256": "encrypt-key-sha",
		"zone":                              zone,
	}
	expected := map[string]string{
		"boot_disk.#":                                "1",
		"attached_disk.#":                            "1",
		"attached_disk.0.source":                     "https://www.googleapis.com/compute/v1/projects/" + config.Project + "/zones/" + zone + "/disks/" + diskName,
		"attached_disk.0.device_name":                "persistent-disk-1",
		"attached_disk.0.disk_encryption_key_raw":    "encrypt-key",
		"attached_disk.0.disk_encryption_key_sha256": "encrypt-key-sha",
		"zone":           zone,
		"create_timeout": "4",
	}

	runInstanceMigrateTest(t, instanceName, "migrate disk to attached disk", 2 /* state version */, attributes, expected, config)
}

func TestAccComputeInstanceMigrateState_v4FixAttachedDiskFromSource(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	config := getInitializedConfig(t)
	zone := "us-central1-f"

	// Seed test data
	diskName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	disk := &compute.Disk{
		Name:        diskName,
		SourceImage: "projects/debian-cloud/global/images/family/debian-11",
		Zone:        zone,
	}
	op, err := config.NewComputeClient(config.UserAgent).Disks.Insert(config.Project, zone, disk).Do()
	if err != nil {
		t.Fatalf("Error creating disk: %s", err)
	}
	waitErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "disk to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpDisk(config, diskName, zone)

	instanceName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	instance := &compute.Instance{
		Name: instanceName,
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
			{
				Source: "projects/" + config.Project + "/zones/" + zone + "/disks/" + diskName,
			},
		},
		MachineType: "zones/" + zone + "/machineTypes/e2-medium",
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
	}
	op, err = config.NewComputeClient(config.UserAgent).Instances.Insert(config.Project, zone, instance).Do()
	if err != nil {
		t.Fatalf("Error creating instance: %s", err)
	}
	waitErr = tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "instance to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpInstance(config, instanceName, zone)

	attributes := map[string]string{
		"boot_disk.#":                       "1",
		"disk.#":                            "1",
		"disk.0.disk":                       diskName,
		"disk.0.device_name":                "persistent-disk-1",
		"disk.0.disk_encryption_key_raw":    "encrypt-key",
		"disk.0.disk_encryption_key_sha256": "encrypt-key-sha",
		"zone":                              zone,
	}
	expected := map[string]string{
		"boot_disk.#":                                "1",
		"attached_disk.#":                            "1",
		"attached_disk.0.source":                     "https://www.googleapis.com/compute/v1/projects/" + config.Project + "/zones/" + zone + "/disks/" + diskName,
		"attached_disk.0.device_name":                "persistent-disk-1",
		"attached_disk.0.disk_encryption_key_raw":    "encrypt-key",
		"attached_disk.0.disk_encryption_key_sha256": "encrypt-key-sha",
		"zone": zone,
	}

	runInstanceMigrateTest(t, instanceName, "migrate disk to attached disk", 4 /* state version */, attributes, expected, config)
}

func TestAccComputeInstanceMigrateState_attachedDiskFromEncryptionKey(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	config := getInitializedConfig(t)
	zone := "us-central1-f"

	instanceName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	instance := &compute.Instance{
		Name: instanceName,
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
			{
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
				DiskEncryptionKey: &compute.CustomerEncryptionKey{
					RawKey: "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0=",
				},
			},
		},
		MachineType: "zones/" + zone + "/machineTypes/e2-medium",
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
	}
	op, err := config.NewComputeClient(config.UserAgent).Instances.Insert(config.Project, zone, instance).Do()
	if err != nil {
		t.Fatalf("Error creating instance: %s", err)
	}
	waitErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "instance to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpInstance(config, instanceName, zone)

	attributes := map[string]string{
		"boot_disk.#":                       "1",
		"disk.#":                            "1",
		"disk.0.image":                      "projects/debian-cloud/global/images/family/debian-11",
		"disk.0.disk_encryption_key_raw":    "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0=",
		"disk.0.disk_encryption_key_sha256": "esTuF7d4eatX4cnc4JsiEiaI+Rff78JgPhA/v1zxX9E=",
		"zone":                              zone,
	}
	expected := map[string]string{
		"boot_disk.#":                                "1",
		"attached_disk.#":                            "1",
		"attached_disk.0.source":                     "https://www.googleapis.com/compute/v1/projects/" + config.Project + "/zones/" + zone + "/disks/" + instanceName + "-1",
		"attached_disk.0.device_name":                "persistent-disk-1",
		"attached_disk.0.disk_encryption_key_raw":    "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0=",
		"attached_disk.0.disk_encryption_key_sha256": "esTuF7d4eatX4cnc4JsiEiaI+Rff78JgPhA/v1zxX9E=",
		"zone":           zone,
		"create_timeout": "4",
	}

	runInstanceMigrateTest(t, instanceName, "migrate disk to attached disk", 2 /* state version */, attributes, expected, config)
}

func TestAccComputeInstanceMigrateState_v4FixAttachedDiskFromEncryptionKey(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	config := getInitializedConfig(t)
	zone := "us-central1-f"

	instanceName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	instance := &compute.Instance{
		Name: instanceName,
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
			{
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
				DiskEncryptionKey: &compute.CustomerEncryptionKey{
					RawKey: "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0=",
				},
			},
		},
		MachineType: "zones/" + zone + "/machineTypes/e2-medium",
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
	}
	op, err := config.NewComputeClient(config.UserAgent).Instances.Insert(config.Project, zone, instance).Do()
	if err != nil {
		t.Fatalf("Error creating instance: %s", err)
	}
	waitErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "instance to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpInstance(config, instanceName, zone)

	attributes := map[string]string{
		"boot_disk.#":                       "1",
		"disk.#":                            "1",
		"disk.0.image":                      "projects/debian-cloud/global/images/family/debian-11",
		"disk.0.disk_encryption_key_raw":    "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0=",
		"disk.0.disk_encryption_key_sha256": "esTuF7d4eatX4cnc4JsiEiaI+Rff78JgPhA/v1zxX9E=",
		"zone":                              zone,
	}
	expected := map[string]string{
		"boot_disk.#":                                "1",
		"attached_disk.#":                            "1",
		"attached_disk.0.source":                     "https://www.googleapis.com/compute/v1/projects/" + config.Project + "/zones/" + zone + "/disks/" + instanceName + "-1",
		"attached_disk.0.device_name":                "persistent-disk-1",
		"attached_disk.0.disk_encryption_key_raw":    "SGVsbG8gZnJvbSBHb29nbGUgQ2xvdWQgUGxhdGZvcm0=",
		"attached_disk.0.disk_encryption_key_sha256": "esTuF7d4eatX4cnc4JsiEiaI+Rff78JgPhA/v1zxX9E=",
		"zone": zone,
	}

	runInstanceMigrateTest(t, instanceName, "migrate disk to attached disk", 4 /* state version */, attributes, expected, config)
}

func TestAccComputeInstanceMigrateState_attachedDiskFromAutoDeleteAndImage(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	config := getInitializedConfig(t)
	zone := "us-central1-f"

	instanceName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	instance := &compute.Instance{
		Name: instanceName,
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
			{
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
			{
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/debian-11-bullseye-v20220719",
				},
			},
		},
		MachineType: "zones/" + zone + "/machineTypes/e2-medium",
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
	}
	op, err := config.NewComputeClient(config.UserAgent).Instances.Insert(config.Project, zone, instance).Do()
	if err != nil {
		t.Fatalf("Error creating instance: %s", err)
	}
	waitErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "instance to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpInstance(config, instanceName, zone)

	attributes := map[string]string{
		"boot_disk.#":        "1",
		"disk.#":             "2",
		"disk.0.image":       "projects/debian-cloud/global/images/debian-11-bullseye-v20220719",
		"disk.0.auto_delete": "true",
		"disk.1.image":       "global/images/family/debian-11",
		"disk.1.auto_delete": "true",
		"zone":               zone,
	}
	expected := map[string]string{
		"boot_disk.#":                 "1",
		"attached_disk.#":             "2",
		"attached_disk.0.source":      "https://www.googleapis.com/compute/v1/projects/" + config.Project + "/zones/" + zone + "/disks/" + instanceName + "-2",
		"attached_disk.0.device_name": "persistent-disk-2",
		"attached_disk.1.source":      "https://www.googleapis.com/compute/v1/projects/" + config.Project + "/zones/" + zone + "/disks/" + instanceName + "-1",
		"attached_disk.1.device_name": "persistent-disk-1",
		"zone":                        zone,
		"create_timeout":              "4",
	}

	runInstanceMigrateTest(t, instanceName, "migrate disk to attached disk", 2 /* state version */, attributes, expected, config)
}

func TestAccComputeInstanceMigrateState_v4FixAttachedDiskFromAutoDeleteAndImage(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	config := getInitializedConfig(t)
	zone := "us-central1-f"

	instanceName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	instance := &compute.Instance{
		Name: instanceName,
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
			{
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
			{
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/debian-11-bullseye-v20220719",
				},
			},
		},
		MachineType: "zones/" + zone + "/machineTypes/e2-medium",
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
	}
	op, err := config.NewComputeClient(config.UserAgent).Instances.Insert(config.Project, zone, instance).Do()
	if err != nil {
		t.Fatalf("Error creating instance: %s", err)
	}
	waitErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "instance to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpInstance(config, instanceName, zone)

	attributes := map[string]string{
		"boot_disk.#":        "1",
		"disk.#":             "2",
		"disk.0.image":       "projects/debian-cloud/global/images/debian-11-bullseye-v20220719",
		"disk.0.auto_delete": "true",
		"disk.1.image":       "global/images/family/debian-11",
		"disk.1.auto_delete": "true",
		"zone":               zone,
	}
	expected := map[string]string{
		"boot_disk.#":                 "1",
		"attached_disk.#":             "2",
		"attached_disk.0.source":      "https://www.googleapis.com/compute/v1/projects/" + config.Project + "/zones/" + zone + "/disks/" + instanceName + "-2",
		"attached_disk.0.device_name": "persistent-disk-2",
		"attached_disk.1.source":      "https://www.googleapis.com/compute/v1/projects/" + config.Project + "/zones/" + zone + "/disks/" + instanceName + "-1",
		"attached_disk.1.device_name": "persistent-disk-1",
		"zone":                        zone,
	}

	runInstanceMigrateTest(t, instanceName, "migrate disk to attached disk", 4 /* state version */, attributes, expected, config)
}

func TestAccComputeInstanceMigrateState_scratchDisk(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	config := getInitializedConfig(t)
	zone := "us-central1-f"

	// Seed test data
	instanceName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	instance := &compute.Instance{
		Name: instanceName,
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
			{
				AutoDelete: true,
				Type:       "SCRATCH",
				InitializeParams: &compute.AttachedDiskInitializeParams{
					DiskType: "zones/" + zone + "/diskTypes/local-ssd",
				},
			},
		},
		// can't be e2 because of local-ssd
		MachineType: "zones/" + zone + "/machineTypes/n1-standard-1",
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
	}
	op, err := config.NewComputeClient(config.UserAgent).Instances.Insert(config.Project, zone, instance).Do()
	if err != nil {
		t.Fatalf("Error creating instance: %s", err)
	}
	waitErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "instance to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpInstance(config, instanceName, zone)

	attributes := map[string]string{
		"boot_disk.#":        "1",
		"disk.#":             "1",
		"disk.0.auto_delete": "true",
		"disk.0.type":        "local-ssd",
		"disk.0.scratch":     "true",
		"zone":               zone,
	}
	expected := map[string]string{
		"boot_disk.#":              "1",
		"scratch_disk.#":           "1",
		"scratch_disk.0.interface": "SCSI",
		"zone":                     zone,
		"create_timeout":           "4",
	}

	runInstanceMigrateTest(t, instanceName, "migrate disk to scratch disk", 2 /* state version */, attributes, expected, config)
}

func TestAccComputeInstanceMigrateState_v4FixScratchDisk(t *testing.T) {
	t.Parallel()

	if os.Getenv(envvar.TestEnvVar) == "" {
		t.Skipf("Network access not allowed; use %s=1 to enable", envvar.TestEnvVar)
	}
	config := getInitializedConfig(t)
	zone := "us-central1-f"

	// Seed test data
	instanceName := fmt.Sprintf("instance-test-%s", acctest.RandString(t, 10))
	instance := &compute.Instance{
		Name: instanceName,
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/debian-cloud/global/images/family/debian-11",
				},
			},
			{
				AutoDelete: true,
				Type:       "SCRATCH",
				InitializeParams: &compute.AttachedDiskInitializeParams{
					DiskType: "zones/" + zone + "/diskTypes/local-ssd",
				},
			},
		},
		MachineType: "zones/" + zone + "/machineTypes/n1-standard-1", // can't be e2 because of local-ssd
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: "global/networks/default",
			},
		},
	}
	op, err := config.NewComputeClient(config.UserAgent).Instances.Insert(config.Project, zone, instance).Do()
	if err != nil {
		t.Fatalf("Error creating instance: %s", err)
	}
	waitErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "instance to create", config.UserAgent, 4*time.Minute)
	if waitErr != nil {
		t.Fatal(waitErr)
	}
	defer cleanUpInstance(config, instanceName, zone)

	attributes := map[string]string{
		"boot_disk.#":        "1",
		"disk.#":             "1",
		"disk.0.auto_delete": "true",
		"disk.0.type":        "local-ssd",
		"disk.0.scratch":     "true",
		"zone":               zone,
	}
	expected := map[string]string{
		"boot_disk.#":              "1",
		"scratch_disk.#":           "1",
		"scratch_disk.0.interface": "SCSI",
		"zone":                     zone,
	}

	runInstanceMigrateTest(t, instanceName, "migrate disk to scratch disk", 4 /* state version */, attributes, expected, config)
}

func runInstanceMigrateTest(t *testing.T, id, testName string, version int, attributes, expected map[string]string, meta interface{}) {
	is := &terraform.InstanceState{
		ID:         id,
		Attributes: attributes,
	}
	_, err := tpgcompute.ResourceComputeInstanceMigrateState(version, is, meta)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range expected {
		// source is the only self link, so compare by relpaths if source is being
		// compared
		if strings.HasSuffix(k, "source") {
			if !tpgresource.CompareSelfLinkOrResourceName("", attributes[k], v, nil) && attributes[k] != v {
				t.Fatalf(
					"bad uri: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					testName, k, expected[k], k, attributes[k], attributes)
			}
		} else {
			if attributes[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					testName, k, expected[k], k, attributes[k], attributes)
			}
		}
	}

	for k, v := range attributes {
		// source is the only self link, so compare by relpaths if source is being
		// compared
		if strings.HasSuffix(k, "source") {
			if !tpgresource.CompareSelfLinkOrResourceName("", expected[k], v, nil) && expected[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					testName, k, expected[k], k, attributes[k], expected)
			}
		} else {
			if expected[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					testName, k, expected[k], k, attributes[k], expected)
			}
		}
	}
}

func cleanUpInstance(config *transport_tpg.Config, instanceName, zone string) {
	op, err := config.NewComputeClient(config.UserAgent).Instances.Delete(config.Project, zone, instanceName).Do()
	if err != nil {
		log.Printf("[WARNING] Error deleting instance %q, dangling resources may exist: %s", instanceName, err)
		return
	}

	// Wait for the operation to complete
	opErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "instance to delete", config.UserAgent, 4*time.Minute)
	if opErr != nil {
		log.Printf("[WARNING] Error deleting instance %q, dangling resources may exist: %s", instanceName, opErr)
	}
}

func cleanUpDisk(config *transport_tpg.Config, diskName, zone string) {
	op, err := config.NewComputeClient(config.UserAgent).Disks.Delete(config.Project, zone, diskName).Do()
	if err != nil {
		log.Printf("[WARNING] Error deleting disk %q, dangling resources may exist: %s", diskName, err)
		return
	}

	// Wait for the operation to complete
	opErr := tpgcompute.ComputeOperationWaitTime(config, op, config.Project, "disk to delete", config.UserAgent, 4*time.Minute)
	if opErr != nil {
		log.Printf("[WARNING] Error deleting disk %q, dangling resources may exist: %s", diskName, opErr)
	}
}

func getInitializedConfig(t *testing.T) *transport_tpg.Config {
	// Migrate tests are non standard and handle the config directly
	acctest.SkipIfVcr(t)
	// Check that all required environment variables are set
	acctest.AccTestPreCheck(t)

	config := &transport_tpg.Config{
		Project:     envvar.GetTestProjectFromEnv(),
		Credentials: envvar.GetTestCredsFromEnv(),
		Region:      envvar.GetTestRegionFromEnv(),
		Zone:        envvar.GetTestZoneFromEnv(),
	}

	transport_tpg.ConfigureBasePaths(config)

	err := config.LoadAndValidate(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return config
}
