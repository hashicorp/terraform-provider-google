package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"log"
	"time"
)

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

func iamPolicyReadModifyWrite(updater ResourceIamUpdater, modify iamPolicyModifyFunc) error {
	mutexKey := updater.GetMutexKey()
	mutexKV.Lock(mutexKey)
	defer mutexKV.Unlock(mutexKey)

	for {
		backoff := time.Second
		log.Printf("[DEBUG]: Retrieving policy for %s\n", updater.DescribeResource())
		p, err := updater.GetResourceIamPolicy()
		if err != nil {
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
			break
		}
		if isConflictError(err) {
			log.Printf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s\n", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				return fmt.Errorf("Error applying IAM policy to %s: too many concurrent policy changes.\n", updater.DescribeResource())
			}
			continue
		}
		return fmt.Errorf("Error applying IAM policy for %s: %v", updater.DescribeResource(), err)
	}
	log.Printf("[DEBUG]: Set policy for %s", updater.DescribeResource())
	return nil
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
		rb = append(rb, &b)
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
