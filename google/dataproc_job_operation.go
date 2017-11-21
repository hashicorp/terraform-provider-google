package google

import (
	"fmt"
	"time"

	"net/http"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/googleapi"
)

type DataprocJobOperationWaiter struct {
	Service   *dataproc.Service
	Region    string
	ProjectId string
	JobId     string
}

func (w *DataprocJobOperationWaiter) ConfForDelete() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"EXISTS"},
		Target:  []string{"DELETED"},
		Refresh: w.RefreshFuncForDelete(),
	}
}

func (w *DataprocJobOperationWaiter) Conf() *resource.StateChangeConf {
	// For more info on each of the states please see
	// https://cloud.google.com/dataproc/docs/reference/rest/v1/projects.regions.jobs#JobStatus
	return &resource.StateChangeConf{
		Pending: []string{"PENDING", "CANCEL_PENDING", "CANCEL_STARTED", "SETUP_DONE", "RUNNING"},
		Target:  []string{"CANCELLED", "DONE", "ATTEMPT_FAILURE", "ERROR"},
		Refresh: w.RefreshFunc(),
	}
}

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	ae, ok := err.(*googleapi.Error)
	return ok && ae.Code == http.StatusNotFound
}

func (w *DataprocJobOperationWaiter) RefreshFuncForDelete() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		_, err := w.Service.Projects.Regions.Jobs.Get(w.ProjectId, w.Region, w.JobId).Do()

		if err != nil {
			if isNotFound(err) {
				return "NA", "DELETED", nil
			}
			return nil, "", err
		}

		return "JOB", "EXISTS", err
	}
}

func (w *DataprocJobOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		job, err := w.Service.Projects.Regions.Jobs.Get(w.ProjectId, w.Region, w.JobId).Do()

		if err != nil {
			return nil, "", err
		}

		return job, job.Status.State, err
	}
}

func dataprocDeleteOperationWait(config *Config, region, projectId, jobId string, activity string, timeoutMinutes, minTimeoutSeconds int) error {
	w := &DataprocJobOperationWaiter{
		Service:   config.clientDataproc,
		Region:    region,
		ProjectId: projectId,
		JobId:     jobId,
	}

	state := w.ConfForDelete()
	state.Timeout = time.Duration(timeoutMinutes) * time.Minute
	state.MinTimeout = time.Duration(minTimeoutSeconds) * time.Second
	_, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	return nil
}

func dataprocJobOperationWait(config *Config, region, projectId, jobId string, activity string, timeoutMinutes, minTimeoutSeconds int) error {
	w := &DataprocJobOperationWaiter{
		Service:   config.clientDataproc,
		Region:    region,
		ProjectId: projectId,
		JobId:     jobId,
	}

	state := w.Conf()
	state.Timeout = time.Duration(timeoutMinutes) * time.Minute
	state.MinTimeout = time.Duration(minTimeoutSeconds) * time.Second
	_, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for operation %s: %s", activity, err)
	}

	return nil
}
