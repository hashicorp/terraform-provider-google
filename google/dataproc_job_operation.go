package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/googleapi"
	"net/http"
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

	// State: Output-only. A state message specifying the overall job state.
	//
	// Possible values:
	//   "STATE_UNSPECIFIED" - The job state is unknown.
	//   "PENDING" - The job is pending; it has been submitted, but is not
	// yet running.
	//   "SETUP_DONE" - Job has been received by the service and completed
	// initial setup; it will soon be submitted to the cluster.
	//   "RUNNING" - The job is running on the cluster.
	//   "CANCEL_PENDING" - A CancelJob request has been received, but is
	// pending.
	//   "CANCEL_STARTED" - Transient in-flight resources have been
	// canceled, and the request to cancel the running job has been issued
	// to the cluster.
	//   "CANCELLED" - The job cancellation was successful.
	//   "DONE" - The job has completed successfully.
	//   "ERROR" - The job has completed, but encountered an error.
	//   "ATTEMPT_FAILURE" - Job attempt has failed. The detail field
	// contains failure details for this attempt.Applies to restartable jobs
	// only.

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
