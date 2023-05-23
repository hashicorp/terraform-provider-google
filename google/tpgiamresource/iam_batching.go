package tpgiamresource

import (
	"fmt"
	"time"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
)

const (
	batchKeyTmplModifyIamPolicy = "%s modifyIamPolicy"
)

func BatchRequestModifyIamPolicy(updater ResourceIamUpdater, modify iamPolicyModifyFunc, config *transport_tpg.Config, reqDesc string) error {
	batchKey := fmt.Sprintf(batchKeyTmplModifyIamPolicy, updater.GetMutexKey())

	request := &transport_tpg.BatchRequest{
		ResourceName: updater.GetResourceId(),
		Body:         []iamPolicyModifyFunc{modify},
		CombineF:     combineBatchIamPolicyModifiers,
		SendF:        sendBatchModifyIamPolicy(updater),
		DebugId:      reqDesc,
	}

	_, err := config.RequestBatcherIam.SendRequestWithTimeout(batchKey, request, time.Minute*30)
	return err
}

func combineBatchIamPolicyModifiers(currV interface{}, toAddV interface{}) (interface{}, error) {
	currModifiers, ok := currV.([]iamPolicyModifyFunc)
	if !ok {
		return nil, fmt.Errorf("provider error in batch combiner: expected data to be type []iamPolicyModifyFunc, got %v with type %T", currV, currV)
	}

	newModifiers, ok := toAddV.([]iamPolicyModifyFunc)
	if !ok {
		return nil, fmt.Errorf("provider error in batch combiner: expected data to be type []iamPolicyModifyFunc, got %v with type %T", currV, currV)
	}

	return append(currModifiers, newModifiers...), nil
}

func sendBatchModifyIamPolicy(updater ResourceIamUpdater) transport_tpg.BatcherSendFunc {
	return func(resourceName string, body interface{}) (interface{}, error) {
		modifiers, ok := body.([]iamPolicyModifyFunc)
		if !ok {
			return nil, fmt.Errorf("provider error: expected data to be type []iamPolicyModifyFunc, got %v with type %T", body, body)
		}
		return nil, iamPolicyReadModifyWrite(updater, func(policy *cloudresourcemanager.Policy) error {
			for _, modifyF := range modifiers {
				if err := modifyF(policy); err != nil {
					return err
				}
			}
			return nil
		})
	}
}
