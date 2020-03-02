package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/errwrap"
)

const readyStatusType string = "Ready"
const pendingCertificateReason string = "CertificatePending"

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

func PollCheckKnativeStatus(resp map[string]interface{}, respErr error) PollResult {
	if respErr != nil {
		return ErrorPollResult(respErr)
	}
	s := KnativeStatus{}
	if err := Convert(resp, &s); err != nil {
		return ErrorPollResult(errwrap.Wrapf("unable to get KnativeStatus: {{err}}", err))
	}

	for _, condition := range s.Status.Conditions {
		if condition.Type == readyStatusType {
			log.Printf("[DEBUG] checking KnativeStatus Ready condition %s: %s", condition.Status, condition.Message)
			switch condition.Status {
			case "True":
				// Resource is ready
				return SuccessPollResult()
			case "Unknown":
				// DomainMapping can enter a 'terminal' state where "Ready" status is "Unknown"
				// but the resource is waiting for external verification of DNS records.
				if condition.Reason == pendingCertificateReason {
					return SuccessPollResult()
				}
				return PendingStatusPollResult(fmt.Sprintf("%s:%s", condition.Status, condition.Message))
			case "False":
				return ErrorPollResult(fmt.Errorf(`resource is in failed state "Ready:False", message: %s`, condition.Message))
			}
		}
	}
	return PendingStatusPollResult("no status yet")
}
