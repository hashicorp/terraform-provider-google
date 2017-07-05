package google

import (
	"encoding/json"
	"google.golang.org/api/compute/v1"
)

func convert(item, out interface{}) error {
	bytes, err := json.Marshal(item)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return err
	}

	return nil
}

func convertInstanceGroupManagerToV1(self interface{}) (*compute.InstanceGroupManager, error) {
	item := &(compute.InstanceGroupManager{})
	err := convert(self, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func convertInstanceGroupManagersSetTargetPoolsRequestToV1(self interface{}) (*compute.InstanceGroupManagersSetTargetPoolsRequest, error) {
	item := &(compute.InstanceGroupManagersSetTargetPoolsRequest{})
	err := convert(self, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func convertInstanceGroupManagersSetInstanceTemplateRequestToV1(self interface{}) (*compute.InstanceGroupManagersSetInstanceTemplateRequest, error) {
	item := &(compute.InstanceGroupManagersSetInstanceTemplateRequest{})
	err := convert(self, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func convertInstanceGroupsSetNamedPortsRequestToV1(self interface{}) (*compute.InstanceGroupsSetNamedPortsRequest, error) {
	item := &(compute.InstanceGroupsSetNamedPortsRequest{})
	err := convert(self, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func convertInstanceGroupManagersListManagedInstancesResponseToV1(self interface{}) (*compute.InstanceGroupManagersListManagedInstancesResponse, error) {
	item := &(compute.InstanceGroupManagersListManagedInstancesResponse{})
	err := convert(self, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func convertInstanceGroupManagersRecreateInstancesRequestToV1(self interface{}) (*compute.InstanceGroupManagersRecreateInstancesRequest, error) {
	item := &(compute.InstanceGroupManagersRecreateInstancesRequest{})
	err := convert(self, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}
