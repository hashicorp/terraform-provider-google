package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	batchKeyTmplServiceUsageEnableServices = "project/%s/services:batchEnable"
	batchKeyTmplServiceUsageListServices   = "project/%s/services"
)

// BatchRequestEnableServices can be used to batch requests to enable services
// across resource nodes, i.e. to batch creation of several
// google_project_service(s) resources.
func BatchRequestEnableService(service string, project string, d *schema.ResourceData, config *Config) error {
	// Renamed service create calls are relatively likely to fail, so don't try to batch the call.
	if altName, ok := renamedServicesByOldAndNewServiceNames[service]; ok {
		return tryEnableRenamedService(service, altName, project, d, config)
	}

	req := &BatchRequest{
		ResourceName: project,
		Body:         []string{service},
		CombineF:     combineServiceUsageServicesBatches,
		SendF:        sendBatchFuncEnableServices(config, d.Timeout(schema.TimeoutCreate)),
		DebugId:      fmt.Sprintf("Enable Project Service %q for project %q", service, project),
	}

	_, err := config.requestBatcherServiceUsage.SendRequestWithTimeout(
		fmt.Sprintf(batchKeyTmplServiceUsageEnableServices, project),
		req,
		d.Timeout(schema.TimeoutCreate))
	return err
}

func tryEnableRenamedService(service, altName string, project string, d *schema.ResourceData, config *Config) error {
	log.Printf("[DEBUG] found renamed service %s (with alternate name %s)", service, altName)
	// use a short timeout- failures are likely

	log.Printf("[DEBUG] attempting enabling service with user-specified name %s", service)
	err := enableServiceUsageProjectServices([]string{altName}, project, config, 1*time.Minute)
	if err != nil {
		log.Printf("[DEBUG] saw error %s. attempting alternate name %v", err, altName)
		err2 := enableServiceUsageProjectServices([]string{altName}, project, config, 1*time.Minute)
		if err2 != nil {
			return fmt.Errorf("Saw 2 subsequent errors attempting to enable a renamed service: %s / %s", err, err2)
		}
	}
	return nil
}

func BatchRequestReadServices(project string, d *schema.ResourceData, config *Config) (interface{}, error) {
	req := &BatchRequest{
		ResourceName: project,
		Body:         nil,
		// Use empty CombineF since the request is exactly the same no matter how many services we read.
		CombineF: func(body interface{}, toAdd interface{}) (interface{}, error) { return nil, nil },
		SendF:    sendListServices(config, d.Timeout(schema.TimeoutRead)),
		DebugId:  fmt.Sprintf("List Project Services %s", project),
	}

	return config.requestBatcherServiceUsage.SendRequestWithTimeout(
		fmt.Sprintf(batchKeyTmplServiceUsageListServices, project),
		req,
		d.Timeout(schema.TimeoutRead))
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

func sendBatchFuncEnableServices(config *Config, timeout time.Duration) BatcherSendFunc {
	return func(project string, toEnableRaw interface{}) (interface{}, error) {
		toEnable, ok := toEnableRaw.([]string)
		if !ok {
			return nil, fmt.Errorf("Expected batch body type to be []string, got %v. This is a provider error.", toEnableRaw)
		}
		return nil, enableServiceUsageProjectServices(toEnable, project, config, timeout)
	}
}

func sendListServices(config *Config, timeout time.Duration) BatcherSendFunc {
	return func(project string, _ interface{}) (interface{}, error) {
		return listCurrentlyEnabledServices(project, config, timeout)
	}
}
