package google

import (
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/dataproc/v1"
)

type DataprocJobOperationWaiter struct {
	Service   *dataproc.Service
	Region    string
	ProjectId string
	JobId     string
	Status    string
}

func (w *DataprocJobOperationWaiter) State() string {
	if w == nil {
		return "<nil>"
	}
	return w.Status
}

func (w *DataprocJobOperationWaiter) Error() error {
	// The "operation" is just the job, which has no special error field that we
	// want to expose.
	return nil
}

func (w *DataprocJobOperationWaiter) IsRetryable(error) bool {
	return false
}

func (w *DataprocJobOperationWaiter) SetOp(job interface{}) error {
	// The "operation" is just the job. Instead of holding onto the whole job
	// object, we only care about the state, which gets set in QueryOp, so this
	// doesn't have to do anything.
	return nil
}

func (w *DataprocJobOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	job, err := w.Service.Projects.Regions.Jobs.Get(w.ProjectId, w.Region, w.JobId).Do()
	if job != nil {
		w.Status = job.Status.State
	}
	return job, err
}

func (w *DataprocJobOperationWaiter) OpName() string {
	if w == nil {
		return "<nil>"
	}
	return w.JobId
}

func (w *DataprocJobOperationWaiter) PendingStates() []string {
	return []string{"PENDING", "CANCEL_PENDING", "CANCEL_STARTED", "SETUP_DONE", "RUNNING"}
}

func (w *DataprocJobOperationWaiter) TargetStates() []string {
	return []string{"CANCELLED", "DONE", "ATTEMPT_FAILURE", "ERROR"}
}

func dataprocJobOperationWait(config *Config, region, projectId, jobId string, activity string, timeout time.Duration) error {
	w := &DataprocJobOperationWaiter{
		Service:   config.clientDataproc,
		Region:    region,
		ProjectId: projectId,
		JobId:     jobId,
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}

type DataprocDeleteJobOperationWaiter struct {
	DataprocJobOperationWaiter
}

func (w *DataprocDeleteJobOperationWaiter) PendingStates() []string {
	return []string{"EXISTS", "ERROR"}
}

func (w *DataprocDeleteJobOperationWaiter) TargetStates() []string {
	return []string{"DELETED"}
}

func (w *DataprocDeleteJobOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	job, err := w.Service.Projects.Regions.Jobs.Get(w.ProjectId, w.Region, w.JobId).Do()
	if err != nil {
		if isGoogleApiErrorWithCode(err, http.StatusNotFound) {
			w.Status = "DELETED"
			return job, nil
		}
		w.Status = "ERROR"
	}
	w.Status = "EXISTS"
	return job, err
}

func dataprocDeleteOperationWait(config *Config, region, projectId, jobId string, activity string, timeout time.Duration) error {
	w := &DataprocDeleteJobOperationWaiter{
		DataprocJobOperationWaiter{
			Service:   config.clientDataproc,
			Region:    region,
			ProjectId: projectId,
			JobId:     jobId,
		},
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}
