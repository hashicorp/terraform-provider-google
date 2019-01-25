package google

import (
	"fmt"
	"regexp"

	"google.golang.org/api/appengine/v1"
)

var (
	appEngineOperationIdRegexp = regexp.MustCompile(fmt.Sprintf("apps/%s/operations/(.*)", ProjectRegex))
)

type AppEngineOperationWaiter struct {
	Service *appengine.APIService
	AppId   string
	CommonOperationWaiter
}

func (w *AppEngineOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	matches := appEngineOperationIdRegexp.FindStringSubmatch(w.Op.Name)
	if len(matches) != 2 {
		return nil, fmt.Errorf("Expected %d results of parsing operation name, got %d from %s", 2, len(matches), w.Op.Name)
	}
	return w.Service.Apps.Operations.Get(w.AppId, matches[1]).Do()
}

func appEngineOperationWait(client *appengine.APIService, op *appengine.Operation, appId, activity string) error {
	return appEngineOperationWaitTime(client, op, appId, activity, 4)
}

func appEngineOperationWaitTime(client *appengine.APIService, op *appengine.Operation, appId, activity string, timeoutMinutes int) error {
	w := &AppEngineOperationWaiter{
		Service: client,
		AppId:   appId,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}
