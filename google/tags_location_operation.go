package google

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type TagsLocationOperationWaiter struct {
	Config    *Config
	UserAgent string
	CommonOperationWaiter
}

func (w *TagsLocationOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	location := GetLocationFromOpName(w.CommonOperationWaiter.Op.Name)
	if location != w.CommonOperationWaiter.Op.Name {
		// Found location in Op.Name, fill it in TagsLocationBasePath and rewrite URL
		url := fmt.Sprintf("%s%s", strings.Replace(w.Config.TagsLocationBasePath, "{{location}}", location, 1), w.CommonOperationWaiter.Op.Name)
		return SendRequest(w.Config, "GET", "", url, w.UserAgent, nil)
	} else {
		url := fmt.Sprintf("%s%s", w.Config.TagsBasePath, w.CommonOperationWaiter.Op.Name)
		return SendRequest(w.Config, "GET", "", url, w.UserAgent, nil)
	}
}

func createTagsLocationWaiter(config *Config, op map[string]interface{}, activity, userAgent string) (*TagsLocationOperationWaiter, error) {
	w := &TagsLocationOperationWaiter{
		Config:    config,
		UserAgent: userAgent,
	}
	if err := w.CommonOperationWaiter.SetOp(op); err != nil {
		return nil, err
	}
	return w, nil
}

func TagsLocationOperationWaitTimeWithResponse(config *Config, op map[string]interface{}, response *map[string]interface{}, activity, userAgent string, timeout time.Duration) error {
	w, err := createTagsLocationWaiter(config, op, activity, userAgent)
	if err != nil {
		return err
	}
	if err := OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return err
	}
	return json.Unmarshal([]byte(w.CommonOperationWaiter.Op.Response), response)
}

func TagsLocationOperationWaitTime(config *Config, op map[string]interface{}, activity, userAgent string, timeout time.Duration) error {
	if val, ok := op["name"]; !ok || val == "" {
		// This was a synchronous call - there is no operation to wait for.
		return nil
	}
	w, err := createTagsLocationWaiter(config, op, activity, userAgent)
	if err != nil {
		// If w is nil, the op was synchronous.
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}

func GetLocationFromOpName(opName string) string {
	re := regexp.MustCompile("operations/(?:rctb|rdtb)\\.([a-zA-Z0-9-]*)\\.([0-9]*)")
	switch {
	case re.MatchString(opName):
		if res := re.FindStringSubmatch(opName); len(res) == 3 && res[1] != "" {
			return res[1]
		}
	}
	return opName
}
