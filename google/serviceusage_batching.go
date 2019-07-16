package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

const (
	batchKeyTmplServiceUsageEnableServices = "project/%s/services:batchEnable"
)

// globalBatchEnableServices can be used to batch requests to enable services
// across resource nodes, i.e. to batch creation of several
// google_project_service(s) resources.
func globalBatchEnableServices(services []string, project string, d *schema.ResourceData, config *Config) error {
	req := &BatchRequest{
		ResourceName: project,
		Body:         services,
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
