package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	batchKeyTmplServiceUsageEnableServices = "project/%s/services:batchEnable"
)

// BatchRequestEnableServices can be used to batch requests to enable services
// across resource nodes, i.e. to batch creation of several
// google_project_service(s) resources.
func BatchRequestEnableServices(services map[string]struct{}, project string, d *schema.ResourceData, config *Config) error {
	// renamed service create calls are relatively likely to fail, so break out
	// of the batched call to avoid failing that as well
	for k := range services {
		if v, ok := renamedServicesByOldAndNewServiceNames[k]; ok {
			log.Printf("[DEBUG] found renamed service %s (with alternate name %s)", k, v)
			delete(services, k)
			// also remove the other name, so we don't enable it 2x in a row
			delete(services, v)

			// use a short timeout- failures are likely
			log.Printf("[DEBUG] attempting user-specified name %s", k)
			err := enableServiceUsageProjectServices([]string{k}, project, config, 1*time.Minute)
			if err != nil {
				log.Printf("[DEBUG] saw error %s. attempting alternate name %v", err, v)
				err2 := enableServiceUsageProjectServices([]string{v}, project, config, 1*time.Minute)
				if err2 != nil {
					return fmt.Errorf("Saw 2 subsequent errors attempting to enable a renamed service: %s / %s", err, err2)
				}
			}
		}
	}

	req := &BatchRequest{
		ResourceName: project,
		Body:         stringSliceFromGolangSet(services),
		CombineF:     combineServiceUsageServicesBatches,
		SendF:        sendBatchFuncEnableServices(config, d.Timeout(schema.TimeoutCreate)),
		DebugId:      fmt.Sprintf("Enable Project Services %s: %+v", project, services),
	}

	_, err := config.requestBatcherServiceUsage.SendRequestWithTimeout(
		fmt.Sprintf(batchKeyTmplServiceUsageEnableServices, project),
		req,
		d.Timeout(schema.TimeoutCreate))
	return err
}

func combineServiceUsageServicesBatches(srvsRaw interface{}, toAddRaw interface{}) (interface{}, error) {
	srvs, ok := srvsRaw.([]string)
	if !ok {
		return nil, fmt.Errorf("Expected batch body type to be []string, got %v. This is a provider error.", srvsRaw)
	}
	toAdd, ok := toAddRaw.([]string)
	if !ok {
		return nil, fmt.Errorf("Expected new request body type to be []string, got %v. This is a provider error.", toAdd)
	}

	return append(srvs, toAdd...), nil
}

func sendBatchFuncEnableServices(config *Config, timeout time.Duration) batcherSendFunc {
	return func(project string, toEnableRaw interface{}) (interface{}, error) {
		toEnable, ok := toEnableRaw.([]string)
		if !ok {
			return nil, fmt.Errorf("Expected batch body type to be []string, got %v. This is a provider error.", toEnableRaw)
		}
		return nil, enableServiceUsageProjectServices(toEnable, project, config, timeout)
	}
}
