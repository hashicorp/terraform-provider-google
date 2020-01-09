package google

import (
	"fmt"
	"google.golang.org/api/cloudresourcemanager/v1"
	"time"
)

const (
	batchKeyTmplModifyIamPolicy = "%s modifyIamPolicy"

	IamBatchingEnabled  = true
	IamBatchingDisabled = false
)

func BatchRequestModifyIamPolicy(updater ResourceIamUpdater, modify iamPolicyModifyFunc, config *Config, reqDesc string) error {
	batchKey := fmt.Sprintf(batchKeyTmplModifyIamPolicy, updater.GetMutexKey())

	request := &BatchRequest{
		ResourceName: updater.GetResourceId(),
		Body:         []iamPolicyModifyFunc{modify},
		CombineF:     combineBatchIamPolicyModifiers,
		SendF:        sendBatchModifyIamPolicy(updater),
		DebugId:      reqDesc,
	}

	_, err := config.requestBatcherIam.SendRequestWithTimeout(batchKey, request, time.Minute*30)
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

func sendBatchModifyIamPolicy(updater ResourceIamUpdater) batcherSendFunc {
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
