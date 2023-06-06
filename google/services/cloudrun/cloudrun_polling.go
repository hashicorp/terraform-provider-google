// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudrun

import (
	"fmt"
	"log"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
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
		Conditions         []Condition
		ObservedGeneration float64
	}
}

func getGeneration(res map[string]interface{}) (int, error) {
	metadata, ok := res["metadata"]
	if !ok {
		return 0, fmt.Errorf("Unable to find knative metadata")
	}
	m, ok := metadata.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("Unable to find generation in knative metadata")
	}
	gen, ok := m["generation"]
	if !ok {
		return 0, fmt.Errorf("Unable to find generation in knative metadata")
	}
	return int(gen.(float64)), nil
}

func PollCheckKnativeStatusFunc(knativeRestResponse map[string]interface{}) func(resp map[string]interface{}, respErr error) transport_tpg.PollResult {
	return func(resp map[string]interface{}, respErr error) transport_tpg.PollResult {
		if respErr != nil {
			return transport_tpg.ErrorPollResult(respErr)
		}
		s := KnativeStatus{}
		if err := tpgresource.Convert(resp, &s); err != nil {
			return transport_tpg.ErrorPollResult(errwrap.Wrapf("unable to get KnativeStatus: {{err}}", err))
		}

		gen, err := getGeneration(knativeRestResponse)
		if err != nil {
			return transport_tpg.ErrorPollResult(errwrap.Wrapf("unable to find Knative generation: {{err}}", err))
		}
		if int(s.Status.ObservedGeneration) < gen {
			return transport_tpg.PendingStatusPollResult("waiting for observed generation to match")
		}
		for _, condition := range s.Status.Conditions {
			if condition.Type == readyStatusType {
				log.Printf("[DEBUG] checking KnativeStatus Ready condition %s: %s", condition.Status, condition.Message)
				switch condition.Status {
				case "True":
					// Resource is ready
					return transport_tpg.SuccessPollResult()
				case "Unknown":
					// DomainMapping can enter a 'terminal' state where "Ready" status is "Unknown"
					// but the resource is waiting for external verification of DNS records.
					if condition.Reason == pendingCertificateReason {
						return transport_tpg.SuccessPollResult()
					}
					return transport_tpg.PendingStatusPollResult(fmt.Sprintf("%s:%s", condition.Status, condition.Message))
				case "False":
					return transport_tpg.ErrorPollResult(fmt.Errorf(`resource is in failed state "Ready:False", message: %s`, condition.Message))
				}
			}
		}
		return transport_tpg.PendingStatusPollResult("no status yet")
	}
}
