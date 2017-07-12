package google

import (
	"fmt"

	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

type ScopeType uint8

const (
	GLOBAL ScopeType = iota
	REGION
	ZONE
)

// A ComputeService that delegates requests to the appropriate API
// Takes in and returns models of the highest API level supported.
type ComputeMultiversionService struct {
	v1     *compute.Service
	v0beta *computeBeta.Service
}

func (s *ComputeMultiversionService) WaitOperation(project string, operationName string, scopeType ScopeType, scope string) (*computeBeta.Operation, error) {
	var operation *compute.Operation
	var err error
	switch scopeType {
	case GLOBAL:
		operation, err = s.v1.GlobalOperations.Get(project, operationName).Do()
	case REGION:
		operation, err = s.v1.RegionOperations.Get(project, scope, operationName).Do()
	case ZONE:
		operation, err = s.v1.ZoneOperations.Get(project, scope, operationName).Do()
	default:
		operation, err = nil, fmt.Errorf("Awaited operation with unknown scope. %v %s", scopeType, scope)
	}

	if err != nil {
		return nil, err
	}

	v0BetaOperation := &computeBeta.Operation{}
	err = Convert(operation, v0BetaOperation)
	if err != nil {
		return nil, err
	}

	return v0BetaOperation, nil

}

func (s *ComputeMultiversionService) InsertInstanceGroupManager(project string, zone string, manager *computeBeta.InstanceGroupManager, version ComputeApiVersion) (*computeBeta.Operation, error) {
	op := &computeBeta.Operation{}
	switch version {
	case v1:
		v1Manager := &compute.InstanceGroupManager{}
		err := Convert(manager, v1Manager)
		if err != nil {
			return nil, err
		}

		v1Op, err := s.v1.InstanceGroupManagers.Insert(project, zone, v1Manager).Do()
		err = Convert(v1Op, op)
		if err != nil {
			return nil, err
		}

		return op, nil
	case v0beta:
		v0BetaManager := &computeBeta.InstanceGroupManager{}
		err := Convert(manager, v0BetaManager)
		if err != nil {
			return nil, err
		}

		v0BetaOp, err := s.v0beta.InstanceGroupManagers.Insert(project, zone, v0BetaManager).Do()
		err = Convert(v0BetaOp, op)
		if err != nil {
			return nil, err
		}

		return op, nil
	}

	return nil, fmt.Errorf("Unknown API version.")
}
