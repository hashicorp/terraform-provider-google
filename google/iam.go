package google

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
	"log"
	"strings"
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

func iamPolicyImport(resourceIdParser resourceIdParserFunc) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		if resourceIdParser == nil {
			return nil, errors.New("Import not supported for this IAM resource.")
		}
		config := m.(*Config)
		err := resourceIdParser(d, config)
		if err != nil {
			return nil, err
		}
		return []*schema.ResourceData{d}, nil
	}
}

func iamBindingImport(resourceIdParser resourceIdParserFunc) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		if resourceIdParser == nil {
			return nil, errors.New("Import not supported for this IAM resource.")
		}
		config := m.(*Config)
		s := strings.Split(d.Id(), " ")
		if len(s) != 2 {
			d.SetId("")
			return nil, fmt.Errorf("Wrong number of parts to Binding id %s; expected 'resource_name role'.", s)
		}
		id, role := s[0], s[1]
		d.SetId(id)
		d.Set("role", role)
		err := resourceIdParser(d, config)
		if err != nil {
			return nil, err
		}
		// It is possible to return multiple bindings, since we can learn about all the bindings
		// for this resource here.  Unfortunately, `terraform import` has some messy behavior here -
		// there's no way to know at this point which resource is being imported, so it's not possible
		// to order this list in a useful way.  In the event of a complex set of bindings, the user
		// will have a terribly confusing set of imported resources and no way to know what matches
		// up to what.  And since the only users who will do a terraform import on their IAM bindings
		// are users who aren't too familiar with Google Cloud IAM (because a "create" for bindings or
		// members is idempotent), it's reasonable to expect that the user will be very alarmed by the
		// plan that terraform will output which mentions destroying a dozen-plus IAM bindings.  With
		// that in mind, we return only the binding that matters.
		return []*schema.ResourceData{d}, nil
	}
}

func iamMemberImport(resourceIdParser resourceIdParserFunc) schema.StateFunc {
	return func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		if resourceIdParser == nil {
			return nil, errors.New("Import not supported for this IAM resource.")
		}
		config := m.(*Config)
		s := strings.Split(d.Id(), " ")
		if len(s) != 3 {
			d.SetId("")
			return nil, fmt.Errorf("Wrong number of parts to Member id %s; expected 'resource_name role username'.", s)
		}
		id, role, member := s[0], s[1], s[2]
		d.SetId(id)
		d.Set("role", role)
		d.Set("member", member)
		err := resourceIdParser(d, config)
		if err != nil {
			return nil, err
		}
		return []*schema.ResourceData{d}, nil
	}
}

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
