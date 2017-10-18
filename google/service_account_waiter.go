package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
)

type ServiceAccountKeyWaiter struct {
	Service       *iam.ProjectsServiceAccountsKeysService
	PublicKeyType string
	KeyName       string
}

func (w *ServiceAccountKeyWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var err error
		var sak *iam.ServiceAccountKey
		sak, err = w.Service.Get(w.KeyName).PublicKeyType(w.PublicKeyType).Do()

		if err != nil {
			if err.(*googleapi.Error).Code == 404 {
				return nil, "PENDING", nil
			} else {
				return nil, "", err
			}
		} else {
			return sak, "DONE", nil
		}
	}
}

func (w *ServiceAccountKeyWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"DONE"},
		Refresh: w.RefreshFunc(),
	}
}

func serviceAccountKeyWaitTime(client *iam.ProjectsServiceAccountsKeysService, keyName, publicKeyType, activity string, timeoutMin int) error {
	w := &ServiceAccountKeyWaiter{
		Service:       client,
		PublicKeyType: publicKeyType,
		KeyName:       keyName,
	}

	state := w.Conf()
	state.Delay = 10 * time.Second
	state.Timeout = time.Duration(timeoutMin) * time.Minute
	state.MinTimeout = 2 * time.Second
	_, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	return nil
}
