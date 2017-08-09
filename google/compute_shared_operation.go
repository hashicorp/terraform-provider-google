package google

import (
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

func computeSharedOperationWait(config *Config, op interface{}, project string, activity string) error {
	return computeSharedOperationWaitTime(config, op, project, 4, activity)
}

func computeSharedOperationWaitTime(config *Config, op interface{}, project string, minutes int, activity string) error {
	if op == nil {
		panic("Attempted to wait on an Operation that was nil.")
	}

	switch op.(type) {
	case *compute.Operation:
		return computeOperationWaitTime(config, op.(*compute.Operation), project, activity, minutes)
	case *computeBeta.Operation:
		return computeBetaOperationWaitTime(config, op.(*computeBeta.Operation), project, activity, minutes)
	default:
		panic("Attempted to wait on an Operation of unknown type.")
	}
}

func computeSharedOperationWaitZone(config *Config, op interface{}, project string, zone, activity string) error {
	return computeSharedOperationWaitZoneTime(config, op, project, zone, 4, activity)
}

func computeSharedOperationWaitZoneTime(config *Config, op interface{}, project string, zone string, minutes int, activity string) error {
	switch op.(type) {
	case *compute.Operation:
		return computeOperationWaitTime(config, op.(*compute.Operation), project, activity, minutes)
	case *computeBeta.Operation:
		return computeBetaOperationWaitZoneTime(config, op.(*computeBeta.Operation), project, zone, minutes, activity)
	case nil:
		panic("Attempted to wait on an Operation that was nil.")
	default:
		panic("Attempted to wait on an Operation of unknown type.")
	}
}
