package google

import (
	"encoding/json"
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

func appEngineOperationWait(config *Config, res interface{}, appId, activity string) error {
	return appEngineOperationWaitTime(config, res, appId, activity, 4)
}

func appEngineOperationWaitTimeWithResponse(config *Config, res interface{}, response *map[string]interface{}, appId, activity string, timeoutMinutes int) error {
	op := &appengine.Operation{}
	err := Convert(res, op)
	if err != nil {
		return err
	}

	w := &AppEngineOperationWaiter{
		Service: config.clientAppEngine,
		AppId:   appId,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	if err := OperationWait(w, activity, timeoutMinutes); err != nil {
		return err
	}
	return json.Unmarshal([]byte(w.CommonOperationWaiter.Op.Response), response)
}

func appEngineOperationWaitTime(config *Config, res interface{}, appId, activity string, timeoutMinutes int) error {
	op := &appengine.Operation{}
	err := Convert(res, op)
	if err != nil {
		return err
	}

	w := &AppEngineOperationWaiter{
		Service: config.clientAppEngine,
		AppId:   appId,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}
