package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/dns/v1"
)

type DnsOperationWaiter struct {
	Service *dns.ManagedZoneOperationsService
	Op      *dns.Operation
	Project string
}

func (w *DnsOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var op *dns.Operation
		var err error

		if w.Op.ZoneContext != nil {
			op, err = w.Service.Get(w.Project, w.Op.ZoneContext.NewValue.Name, w.Op.Id).Do()
		} else {
			return nil, "", fmt.Errorf("unsupported DNS operation %q", w.Op.Id)
		}

		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Got %q when asking for operation %q", op.Status, w.Op.Id)

		return op, op.Status, nil
	}
}

func (w *DnsOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"done"},
		Refresh: w.RefreshFunc(),
	}
}

func dnsOperationWait(service *dns.Service, op *dns.Operation, project, activity string) error {
	return dnsOperationWaitTime(service, op, project, activity, 4)
}

func dnsOperationWaitTime(service *dns.Service, op *dns.Operation, project, activity string, timeoutMin int) error {
	if op.Status == "done" {
		return nil
	}

	w := &DnsOperationWaiter{
		Service: service.ManagedZoneOperations,
		Op:      op,
		Project: project,
	}

	state := w.Conf()
	state.Delay = 10 * time.Second
	state.Timeout = time.Duration(timeoutMin) * time.Minute
	state.MinTimeout = 2 * time.Second
	_, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	return nil
}
