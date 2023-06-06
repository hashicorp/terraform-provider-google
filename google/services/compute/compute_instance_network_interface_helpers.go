// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/compute/v1"
)

func computeInstanceDeleteAccessConfigs(d *schema.ResourceData, config *transport_tpg.Config, instNetworkInterface *compute.NetworkInterface, project, zone, userAgent, instanceName string) error {
	// Delete any accessConfig that currently exists in instNetworkInterface
	for _, ac := range instNetworkInterface.AccessConfigs {
		op, err := config.NewComputeClient(userAgent).Instances.DeleteAccessConfig(
			project, zone, instanceName, ac.Name, instNetworkInterface.Name).Do()
		if err != nil {
			return fmt.Errorf("Error deleting old access_config: %s", err)
		}
		opErr := ComputeOperationWaitTime(config, op, project, "old access_config to delete", userAgent, d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
	}
	return nil
}

func computeInstanceAddAccessConfigs(d *schema.ResourceData, config *transport_tpg.Config, instNetworkInterface *compute.NetworkInterface, accessConfigs []*compute.AccessConfig, project, zone, userAgent, instanceName string) error {
	// Create new ones
	for _, ac := range accessConfigs {
		op, err := config.NewComputeClient(userAgent).Instances.AddAccessConfig(project, zone, instanceName, instNetworkInterface.Name, ac).Do()
		if err != nil {
			return fmt.Errorf("Error adding new access_config: %s", err)
		}
		opErr := ComputeOperationWaitTime(config, op, project, "new access_config to add", userAgent, d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
	}
	return nil
}

func computeInstanceCreateUpdateWhileStoppedCall(d *schema.ResourceData, config *transport_tpg.Config, networkInterfacePatchObj *compute.NetworkInterface, accessConfigs []*compute.AccessConfig, accessConfigsHaveChanged bool, index int, project, zone, userAgent, instanceName string) func(inst *compute.Instance) error {

	// Access configs' ip changes when the instance stops invalidating our fingerprint
	// expect caller to re-validate instance before calling patch this is why we expect
	// instance to be passed in
	return func(instance *compute.Instance) error {

		instNetworkInterface := instance.NetworkInterfaces[index]
		networkInterfacePatchObj.Fingerprint = instNetworkInterface.Fingerprint

		// Access config can run into some issues since we can't tell the difference between
		// the users declared intent (config within their hcl file) and what we have inferred from the
		// server (terraform state). Access configs contain an ip subproperty that can be incompatible
		// with the subnetwork/network we are transitioning to. Due to this we only change access
		// configs if we notice the configuration (user intent) changes.
		if accessConfigsHaveChanged {
			err := computeInstanceDeleteAccessConfigs(d, config, instNetworkInterface, project, zone, userAgent, instanceName)
			if err != nil {
				return err
			}
		}

		op, err := config.NewComputeClient(userAgent).Instances.UpdateNetworkInterface(project, zone, instanceName, instNetworkInterface.Name, networkInterfacePatchObj).Do()
		if err != nil {
			return errwrap.Wrapf("Error updating network interface: {{err}}", err)
		}
		opErr := ComputeOperationWaitTime(config, op, project, "network interface to update", userAgent, d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}

		if accessConfigsHaveChanged {
			err := computeInstanceAddAccessConfigs(d, config, instNetworkInterface, accessConfigs, project, zone, userAgent, instanceName)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
