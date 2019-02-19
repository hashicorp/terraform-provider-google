package google

import (
	"fmt"

	"google.golang.org/api/cloudfunctions/v1"
)

type CloudFunctionsOperationWaiter struct {
	Service *cloudfunctions.Service
	CommonOperationWaiter
}

func (w *CloudFunctionsOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func cloudFunctionsOperationWait(service *cloudfunctions.Service, op *cloudfunctions.Operation, activity string, timeoutMin int) error {
	w := &CloudFunctionsOperationWaiter{
		Service: service,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMin)
}
