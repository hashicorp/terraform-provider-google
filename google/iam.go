package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

const maxBackoffSeconds = 30

// The ResourceIamUpdater interface is implemented for each GCP resource supporting IAM policy.
//
// Implementations should keep track of the resource identifier.
type ResourceIamUpdater interface {
	// Fetch the existing IAM policy attached to a resource.
	GetResourceIamPolicy() (*cloudresourcemanager.Policy, error)

	// Replaces the existing IAM Policy attached to a resource.
	SetResourceIamPolicy(policy *cloudresourcemanager.Policy) error

	// A mutex guards against concurrent to call to the SetResourceIamPolicy method.
	// The mutex key should be made of the resource type and resource id.
	// For example: `iam-project-{id}`.
	GetMutexKey() string

	// Returns the unique resource identifier.
	GetResourceId() string

	// Textual description of this resource to be used in error message.
	// The description should include the unique resource identifier.
	DescribeResource() string
}

type newResourceIamUpdaterFunc func(d *schema.ResourceData, config *Config) (ResourceIamUpdater, error)
type iamPolicyModifyFunc func(p *cloudresourcemanager.Policy) error

// This method parses identifiers specific to the resource (d.GetId()) into the ResourceData
// object, so that it can be given to the resource's Read method.  Externally, this is wrapped
// into schema.StateFunc functions - one each for a _member, a _binding, and a _policy.  Any
// GCP resource supporting IAM policy might support one, two, or all of these.  Any GCP resource
// for which an implementation of this interface exists could support any of the three.

type resourceIdParserFunc func(d *schema.ResourceData, config *Config) error

// Wrapper around updater.GetResourceIamPolicy() to handle retry/backoff
// for just reading policies from IAM
func iamPolicyReadWithRetry(updater ResourceIamUpdater) (*cloudresourcemanager.Policy, error) {
	mutexKey := updater.GetMutexKey()
	mutexKV.Lock(mutexKey)
	defer mutexKV.Unlock(mutexKey)

	log.Printf("[DEBUG] Retrieving policy for %s\n", updater.DescribeResource())
	var policy *cloudresourcemanager.Policy
	err := retryTime(func() (perr error) {
		policy, perr = updater.GetResourceIamPolicy()
		return perr
	}, 10)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Retrieved policy for %s: %+v\n", updater.DescribeResource(), policy)
	return policy, nil
}

func iamPolicyReadModifyWrite(updater ResourceIamUpdater, modify iamPolicyModifyFunc) error {
	mutexKey := updater.GetMutexKey()
	mutexKV.Lock(mutexKey)
	defer mutexKV.Unlock(mutexKey)

	backoff := time.Second
	for {
		log.Printf("[DEBUG]: Retrieving policy for %s\n", updater.DescribeResource())
		p, err := updater.GetResourceIamPolicy()
		if isGoogleApiErrorWithCode(err, 429) {
			log.Printf("[DEBUG] 429 while attempting to read policy for %s, waiting %v before attempting again", updater.DescribeResource(), backoff)
			time.Sleep(backoff)
			continue
		} else if err != nil {
			return err
		}
		log.Printf("[DEBUG]: Retrieved policy for %s: %+v\n", updater.DescribeResource(), p)

		err = modify(p)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG]: Setting policy for %s to %+v\n", updater.DescribeResource(), p)
		err = updater.SetResourceIamPolicy(p)
		if err == nil {
			fetchBackoff := 1 * time.Second
			for successfulFetches := 0; successfulFetches < 3; {
				if fetchBackoff > maxBackoffSeconds*time.Second {
					return fmt.Errorf("Error applying IAM policy to %s: Waited too long for propagation.\n", updater.DescribeResource())
				}
				time.Sleep(fetchBackoff)
				log.Printf("[DEBUG]: Retrieving policy for %s\n", updater.DescribeResource())
				new_p, err := updater.GetResourceIamPolicy()
				if err != nil {
					// Quota for Read is pretty limited, so watch out for running out of quota.
					if isGoogleApiErrorWithCode(err, 429) {
						fetchBackoff = fetchBackoff * 2
					} else {
						return err
					}
				}
				log.Printf("[DEBUG]: Retrieved policy for %s: %+v\n", updater.DescribeResource(), p)
				if new_p == nil {
					// https://github.com/terraform-providers/terraform-provider-google/issues/2625
					fetchBackoff = fetchBackoff * 2
					continue
				}
				modified_p := new_p
				// This relies on the fact that `modify` is idempotent: since other changes might have
				// happened between the call to set the policy and now, we just need to make sure that
				// our change has been made.  'modify(p) == p' is our check for whether this has been
				// correctly applied.
				err = modify(modified_p)
				if err != nil {
					return err
				}
				if modified_p == new_p {
					successfulFetches += 1
				} else {
					fetchBackoff = fetchBackoff * 2
				}
			}
			break
		}
		if isConflictError(err) {
			log.Printf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s\n", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				return errwrap.Wrapf(fmt.Sprintf("Error applying IAM policy to %s: Too many conflicts.  Latest error: {{err}}", updater.DescribeResource()), err)
			}
			continue
		}
		return errwrap.Wrapf(fmt.Sprintf("Error applying IAM policy for %s: {{err}}", updater.DescribeResource()), err)
	}
	log.Printf("[DEBUG]: Set policy for %s", updater.DescribeResource())
	return nil
}

// Takes a single binding and will either overwrite the same role in a list or append it to the end
func overwriteBinding(bindings []*cloudresourcemanager.Binding, overwrite *cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	var found bool

	for i, b := range bindings {
		if b.Role == overwrite.Role {
			bindings[i] = overwrite
			found = true
			break
		}
	}

	if !found {
		bindings = append(bindings, overwrite)
	}

	return bindings
}

// Merge multiple Bindings such that Bindings with the same Role result in
// a single Binding with combined Members
func mergeBindings(bindings []*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	bm := rolesToMembersMap(bindings)
	rb := make([]*cloudresourcemanager.Binding, 0)

	for role, members := range bm {
		var b cloudresourcemanager.Binding
		b.Role = role
		b.Members = make([]string, 0)
		for m := range members {
			b.Members = append(b.Members, m)
		}
		if len(b.Members) > 0 {
			rb = append(rb, &b)
		}
	}

	return rb
}

// Map a role to a map of members, allowing easy merging of multiple bindings.
func rolesToMembersMap(bindings []*cloudresourcemanager.Binding) map[string]map[string]bool {
	bm := make(map[string]map[string]bool)
	// Get each binding
	for _, b := range bindings {
		// Initialize members map
		if _, ok := bm[b.Role]; !ok {
			bm[b.Role] = make(map[string]bool)
		}
		// Get each member (user/principal) for the binding
		for _, m := range b.Members {
			// Add the member
			bm[b.Role][m] = true
		}
	}
	return bm
}

// Merge multiple Audit Configs such that configs with the same service result in
// a single exemption list with combined members
func mergeAuditConfigs(auditConfigs []*cloudresourcemanager.AuditConfig) []*cloudresourcemanager.AuditConfig {
	am := auditConfigToServiceMap(auditConfigs)
	var ac []*cloudresourcemanager.AuditConfig
	for service, auditLogConfigs := range am {
		var a cloudresourcemanager.AuditConfig
		a.Service = service
		a.AuditLogConfigs = make([]*cloudresourcemanager.AuditLogConfig, 0, len(auditLogConfigs))
		for k, v := range auditLogConfigs {
			var alc cloudresourcemanager.AuditLogConfig
			alc.LogType = k
			for member := range v {
				alc.ExemptedMembers = append(alc.ExemptedMembers, member)
			}
			a.AuditLogConfigs = append(a.AuditLogConfigs, &alc)
		}
		if len(a.AuditLogConfigs) > 0 {
			ac = append(ac, &a)
		}
	}
	return ac
}

// Build a service map with the log_type and bindings below it
func auditConfigToServiceMap(auditConfig []*cloudresourcemanager.AuditConfig) map[string]map[string]map[string]bool {
	ac := make(map[string]map[string]map[string]bool)
	// Get each config
	for _, c := range auditConfig {
		// Initialize service map
		if _, ok := ac[c.Service]; !ok {
			ac[c.Service] = map[string]map[string]bool{}
		}
		// loop through audit log configs
		for _, lc := range c.AuditLogConfigs {
			// Initialize service map
			if _, ok := ac[c.Service][lc.LogType]; !ok {
				ac[c.Service][lc.LogType] = map[string]bool{}
			}
			// Get each member (user/principal) for the binding
			for _, m := range lc.ExemptedMembers {
				// Add the member
				if _, ok := ac[c.Service][lc.LogType][m]; !ok {
					ac[c.Service][lc.LogType][m] = true
				}
			}
		}
	}
	return ac
}
