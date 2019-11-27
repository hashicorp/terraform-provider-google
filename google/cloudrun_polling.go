package google

import (
	"errors"
	"fmt"

	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const readyStatus string = "Ready"

type Condition struct {
	Type    string
	Status  string
	Reason  string
	Message string
}

// KnativeStatus is a struct that can contain a Knative style resource's Status block. It is not
// intended to be used for anything other than polling for the success of the given resource.
type KnativeStatus struct {
	Metadata struct {
		Name      string
		Namespace string
		SelfLink  string
	}
	Status struct {
		Conditions []Condition
	}
}

// ConditionByType is a helper method for extracting a given condition
func (s KnativeStatus) ConditionByType(typ string) *Condition {
	for _, condition := range s.Status.Conditions {
		if condition.Type == typ {
			c := condition
			return &c
		}
	}
	return nil
}

// LatestMessage will return a human consumable status of the resource. This can
// be used to determine the human actionable error the GET doesn't return an explicit
// error but the resource is in an error state.
func (s KnativeStatus) LatestMessage() string {
	c := s.ConditionByType(readyStatus)
	if c != nil {
		return fmt.Sprintf("%s - %s", c.Reason, c.Message)
	}

	return ""
}

// State will return a string representing the status of the Ready condition.
// No other conditions are currently returned as part of the state.
func (s KnativeStatus) State(res interface{}) string {
	for _, condition := range s.Status.Conditions {
		if condition.Type == "Ready" {
			return fmt.Sprintf("%s:%s", condition.Type, condition.Status)
		}
	}
	return "Empty"
}

// CloudRunPolling allows for polling against a cloud run resource that implements the
// Kubernetes style status schema.
type CloudRunPolling struct {
	Config  *Config
	WaitURL string
}

func (p *CloudRunPolling) PendingStates() []string {
	return []string{"Ready:Unknown", "Empty"}
}
func (p *CloudRunPolling) TargetStates() []string {
	return []string{"Ready:True"}
}
func (p *CloudRunPolling) ErrorStates() []string {
	return []string{"Ready:False"}
}

func cloudRunPollingWaitTime(config *Config, res map[string]interface{}, project, url, activity string, timeoutMinutes int) error {
	w := &CloudRunPolling{}

	scc := &resource.StateChangeConf{
		Pending: w.PendingStates(),
		Target:  w.TargetStates(),
		Refresh: func() (interface{}, string, error) {
			res, err := sendRequest(config, "GET", project, url, nil)
			if err != nil {
				return res, "", err
			}

			status := KnativeStatus{}
			err = Convert(res, &status)
			if err != nil {
				return res, "", err
			}

			for _, errState := range w.ErrorStates() {
				if status.State(res) == errState {
					err = errors.New(status.LatestMessage())
				}
			}

			return res, status.State(res), err
		},
		Timeout: time.Duration(timeoutMinutes) * time.Minute,
	}

	_, err := scc.WaitForState()
	return err
}
