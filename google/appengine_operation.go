package google

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"google.golang.org/api/appengine/v1"
)

var (
	appEngineOperationIdRegexp = regexp.MustCompile(fmt.Sprintf("apps/%s/operations/(.*)", ProjectRegex))
)

type AppEngineOperationWaiter struct {
	Service *appengine.APIService
	Op      *appengine.Operation
	AppId   string
}

func (w *AppEngineOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		matches := appEngineOperationIdRegexp.FindStringSubmatch(w.Op.Name)
		if len(matches) != 2 {
			return nil, "", fmt.Errorf("Expected %d results of parsing operation name, got %d from %s", 2, len(matches), w.Op.Name)
		}
		op, err := w.Service.Apps.Operations.Get(w.AppId, matches[1]).Do()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Got %v when asking for operation %q", op.Done, w.Op.Name)
		return op, strconv.FormatBool(op.Done), nil
	}
}

func (w *AppEngineOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"false"},
		Target:  []string{"true"},
		Refresh: w.RefreshFunc(),
	}
}

// AppEngineOperationError wraps appengine.Status and implements the
// error interface so it can be returned.
type AppEngineOperationError appengine.Status

func (e AppEngineOperationError) Error() string {
	return e.Message
}

func appEngineOperationWait(client *appengine.APIService, op *appengine.Operation, appId, activity string) error {
	return appEngineOperationWaitTime(client, op, appId, activity, 4)
}

func appEngineOperationWaitTime(client *appengine.APIService, op *appengine.Operation, appId, activity string, timeoutMin int) error {
	if op.Done {
		if op.Error != nil {
			return AppEngineOperationError(*op.Error)
		}
		return nil
	}

	w := &AppEngineOperationWaiter{
		Service: client,
		Op:      op,
		AppId:   appId,
	}

	state := w.Conf()
	state.Delay = 10 * time.Second
	state.Timeout = time.Duration(timeoutMin) * time.Minute
	state.MinTimeout = 2 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	resultOp := opRaw.(*appengine.Operation)
	if resultOp.Error != nil {
		return AppEngineOperationError(*resultOp.Error)
	}

	return nil
}
