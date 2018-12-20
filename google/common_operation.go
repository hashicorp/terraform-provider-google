package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/googleapi"
)

type Waiter interface {
	State() string
	Error() error
	SetOp(interface{}) error
	QueryOp() (interface{}, error)
	OpName() string
	PendingStates() []string
	TargetStates() []string
}

type CommonOperationWaiter struct {
	Op CommonOperation
}

func (w *CommonOperationWaiter) State() string {
	return fmt.Sprintf("done: %v", w.Op.Done)
}

func (w *CommonOperationWaiter) Error() error {
	if w.Op.Error != nil {
		return fmt.Errorf("Error code %v, message: %s", w.Op.Error.Code, w.Op.Error.Message)
	}
	return nil
}

func (w *CommonOperationWaiter) SetOp(op interface{}) error {
	err := Convert(op, &w.Op)
	if err != nil {
		return err
	}
	return nil
}

func (w *CommonOperationWaiter) OpName() string {
	return w.Op.Name
}

func (w *CommonOperationWaiter) PendingStates() []string {
	return []string{"done: false"}
}

func (w *CommonOperationWaiter) TargetStates() []string {
	return []string{"done: true"}
}

func OperationDone(w Waiter) bool {
	for _, s := range w.TargetStates() {
		if s == w.State() {
			return true
		}
	}
	return false
}

func CommonRefreshFunc(w Waiter) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		op, err := w.QueryOp()
		if err != nil && !isGoogleApiErrorWithCode(err, 429) && !isGoogleApiErrorWithCode(err, 503) {
			return nil, "", err
		}
		if err = w.SetOp(op); err != nil {
			return nil, "", err
		}
		if err = w.Error(); err != nil {
			return nil, "", err
		}
		log.Printf("[DEBUG] Got %v while polling for operation %s's status", w.State(), w.OpName())

		return op, w.State(), nil
	}
}

func OperationWait(w Waiter, activity string, timeoutMinutes int) error {
	if OperationDone(w) {
		if w.Error() != nil {
			return w.Error()
		}
		return nil
	}

	c := &resource.StateChangeConf{
		Pending:    w.PendingStates(),
		Target:     w.TargetStates(),
		Refresh:    CommonRefreshFunc(w),
		Timeout:    time.Duration(timeoutMinutes) * time.Minute,
		MinTimeout: 2 * time.Second,
	}
	opRaw, err := c.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	err = w.SetOp(opRaw)
	if err != nil {
		return err
	}
	if w.Error() != nil {
		return w.Error()
	}

	return nil
}

type CommonOperation struct {
	// Done: If the value is `false`, it means the operation is still in
	// progress.
	// If `true`, the operation is completed, and either `error` or
	// `response` is
	// available.
	Done bool `json:"done,omitempty"`

	// Error: The error result of the operation in case of failure or
	// cancellation.
	Error *CommonOperationStatus `json:"error,omitempty"`

	// Metadata: Service-specific metadata associated with the operation.
	// It typically
	// contains progress information and common metadata such as create
	// time.
	// Some services might not provide such metadata.  Any method that
	// returns a
	// long-running operation should document the metadata type, if any.
	Metadata googleapi.RawMessage `json:"metadata,omitempty"`

	// Name: The server-assigned name, which is only unique within the same
	// service that
	// originally returns it. If you use the default HTTP mapping,
	// the
	// `name` should have the format of `operations/some/unique/name`.
	Name string `json:"name,omitempty"`

	// Response: The normal response of the operation in case of success.
	// If the original
	// method returns no data on success, such as `Delete`, the response
	// is
	// `google.protobuf.Empty`.  If the original method is
	// standard
	// `Get`/`Create`/`Update`, the response should be the resource.  For
	// other
	// methods, the response should have the type `XxxResponse`, where
	// `Xxx`
	// is the original method name.  For example, if the original method
	// name
	// is `TakeSnapshot()`, the inferred response type
	// is
	// `TakeSnapshotResponse`.
	Response googleapi.RawMessage `json:"response,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Done") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Done") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

type CommonOperationStatus struct {
	// Code: The status code, which should be an enum value of
	// google.rpc.Code.
	Code int64 `json:"code,omitempty"`

	// Details: A list of messages that carry the error details.  There is a
	// common set of
	// message types for APIs to use.
	Details []googleapi.RawMessage `json:"details,omitempty"`

	// Message: A developer-facing error message, which should be in
	// English. Any
	// user-facing error message should be localized and sent in
	// the
	// google.rpc.Status.details field, or localized by the client.
	Message string `json:"message,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Code") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Code") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}
