package google

import (
	"google.golang.org/api/cloudfunctions/v1"
)

type CloudFunctionsOperationWaiter struct {
	Service *cloudfunctions.Service
	CommonOperationWaiter
}

func (w *CloudFunctionsOperationWaiter) QueryOp() (interface{}, error) {
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func cloudFunctionsOperationWait(service *cloudfunctions.Service, op *cloudfunctions.Operation, activity string) error {
	w := &CloudFunctionsOperationWaiter{
		Service: service,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, 4)
}
