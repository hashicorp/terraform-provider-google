// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Utils for modifying IAM policies for resources across GCP
package tpgiamresource

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/cloudresourcemanager/v1"
)

const maxBackoffSeconds = 30
const IamPolicyVersion = 3

// These types are implemented per GCP resource type and specify how to do per-resource IAM operations.
// They are used in the generic Terraform IAM resource definitions
// (e.g. _member/_binding/_policy/_audit_config)
type (
	// The ResourceIamUpdater interface is implemented for each GCP resource supporting IAM policy.
	// Implementations should be created per resource and should keep track of the resource identifier.
	ResourceIamUpdater interface {
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

	// Factory for generating ResourceIamUpdater for given ResourceData resource
	NewResourceIamUpdaterFunc func(d tpgresource.TerraformResourceData, config *transport_tpg.Config) (ResourceIamUpdater, error)

	// Describes how to modify a policy for a given Terraform IAM (_policy/_member/_binding/_audit_config) resource
	iamPolicyModifyFunc func(p *cloudresourcemanager.Policy) error

	// Parser for Terraform resource identifier (d.Id) for resource whose IAM policy is being changed
	ResourceIdParserFunc func(d *schema.ResourceData, config *transport_tpg.Config) error
)

// Locking wrapper around read-only operation with retries.
func iamPolicyReadWithRetry(updater ResourceIamUpdater) (*cloudresourcemanager.Policy, error) {
	mutexKey := updater.GetMutexKey()
	transport_tpg.MutexStore.Lock(mutexKey)
	defer transport_tpg.MutexStore.Unlock(mutexKey)

	log.Printf("[DEBUG] Retrieving policy for %s\n", updater.DescribeResource())
	var policy *cloudresourcemanager.Policy
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (perr error) {
			policy, perr = updater.GetResourceIamPolicy()
			return perr
		},
		Timeout: 10 * time.Minute,
	})
	if err != nil {
		return nil, err
	}
	log.Print(spew.Sprintf("[DEBUG] Retrieved policy for %s: %#v\n", updater.DescribeResource(), policy))
	return policy, nil
}

// Locking wrapper around read-modify-write cycle for IAM policy.
func iamPolicyReadModifyWrite(updater ResourceIamUpdater, modify iamPolicyModifyFunc) error {
	mutexKey := updater.GetMutexKey()
	transport_tpg.MutexStore.Lock(mutexKey)
	defer transport_tpg.MutexStore.Unlock(mutexKey)

	backoff := time.Second
	for {
		log.Printf("[DEBUG]: Retrieving policy for %s\n", updater.DescribeResource())
		p, err := updater.GetResourceIamPolicy()
		if transport_tpg.IsGoogleApiErrorWithCode(err, 429) {
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
					if transport_tpg.IsGoogleApiErrorWithCode(err, 429) {
						fetchBackoff = fetchBackoff * 2
					} else {
						return err
					}
				}
				log.Printf("[DEBUG]: Retrieved policy for %s: %+v\n", updater.DescribeResource(), p)
				if new_p == nil {
					// https://github.com/hashicorp/terraform-provider-google/issues/2625
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
		if tpgresource.IsConflictError(err) {
			log.Printf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s\n", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				return errwrap.Wrapf(fmt.Sprintf("Error applying IAM policy to %s: Too many conflicts.  Latest error: {{err}}", updater.DescribeResource()), err)
			}
			continue
		}

		// retry in the case that a service account is not found. This can happen when a service account is deleted
		// out of band.
		if isServiceAccountNotFoundError, _ := transport_tpg.IamServiceAccountNotFound(err); isServiceAccountNotFoundError {
			// calling a retryable function within a retry loop is not
			// strictly the _best_ idea, but this error only happens in
			// high-traffic projects anyways
			currentPolicy, rerr := iamPolicyReadWithRetry(updater)
			if rerr != nil {
				if p.Etag != currentPolicy.Etag {
					// not matching indicates that there is a new state to attempt to apply
					log.Printf("current and old etag did not match for %s, retrying", updater.DescribeResource())
					time.Sleep(backoff)
					backoff = backoff * 2
					continue
				}

				log.Printf("current and old etag matched for %s, not retrying", updater.DescribeResource())
			} else {
				// if the error is non-nil, just fall through and return the base error
				log.Printf("[DEBUG]: error checking etag for policy %s. error: %v", updater.DescribeResource(), rerr)
			}
		}

		log.Printf("[DEBUG]: not retrying IAM policy for %s. error: %v", updater.DescribeResource(), err)
		return errwrap.Wrapf(fmt.Sprintf("Error applying IAM policy for %s: {{err}}", updater.DescribeResource()), err)
	}
	log.Printf("[DEBUG]: Set policy for %s", updater.DescribeResource())
	return nil
}

// Flattens a list of Bindings so each role+condition has a single Binding with combined members
func MergeBindings(bindings []*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	bm := createIamBindingsMap(bindings)
	return listFromIamBindingMap(bm)
}

type conditionKey struct {
	Description string
	Expression  string
	Title       string
}

func conditionKeyFromCondition(condition *cloudresourcemanager.Expr) conditionKey {
	if condition == nil {
		return conditionKey{}
	}
	return conditionKey{condition.Description, condition.Expression, condition.Title}
}

func (k conditionKey) Empty() bool {
	return k == conditionKey{}
}

func (k conditionKey) String() string {
	return fmt.Sprintf("%s/%s/%s", k.Title, k.Description, k.Expression)
}

type iamBindingKey struct {
	Role      string
	Condition conditionKey
}

// Removes a single role+condition binding from a list of Bindings
func filterBindingsWithRoleAndCondition(b []*cloudresourcemanager.Binding, role string, condition *cloudresourcemanager.Expr) []*cloudresourcemanager.Binding {
	bMap := createIamBindingsMap(b)
	key := iamBindingKey{role, conditionKeyFromCondition(condition)}
	delete(bMap, key)
	return listFromIamBindingMap(bMap)
}

// Removes given role+condition/bound-member pairs from the given Bindings (i.e subtraction).
func subtractFromBindings(bindings []*cloudresourcemanager.Binding, toRemove ...*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	currMap := createIamBindingsMap(bindings)
	toRemoveMap := createIamBindingsMap(toRemove)

	for key, removeSet := range toRemoveMap {
		members, ok := currMap[key]
		if !ok {
			continue
		}
		// Remove all removed members
		for m := range removeSet {
			delete(members, m)
		}
		// Remove role+condition from bindings
		if len(members) == 0 {
			delete(currMap, key)
		}
	}

	return listFromIamBindingMap(currMap)
}

func iamMemberIsCaseSensitive(member string) bool {
	// allAuthenticatedUsers and allUsers are special identifiers that are case sensitive. See:
	// https://cloud.google.com/iam/docs/overview#all-authenticated-users
	return strings.Contains(member, "allAuthenticatedUsers") || strings.Contains(member, "allUsers") ||
		strings.HasPrefix(member, "principalSet:") || strings.HasPrefix(member, "principal:") ||
		strings.HasPrefix(member, "principalHierarchy:")
}

// normalizeIamMemberCasing returns the case adjusted value of an iamMember
// this is important as iam will ignore casing unless it is one of the following
// member types: principalSet, principal, principalHierarchy
// members are in <type>:<value> format
// <type> is case sensitive
// <value> isn't in most cases
// so lowercase the value unless iamMemberIsCaseSensitive and leave the type alone
// since Dec '19 members can be prefixed with "deleted:" to indicate the principal
// has been deleted
func normalizeIamMemberCasing(member string) string {
	var pieces []string
	if strings.HasPrefix(member, "deleted:") {
		pieces = strings.SplitN(member, ":", 3)
		if len(pieces) > 2 && !iamMemberIsCaseSensitive(strings.TrimPrefix(member, "deleted:")) {
			pieces[2] = strings.ToLower(pieces[2])
		}
	} else if !iamMemberIsCaseSensitive(member) {
		pieces = strings.SplitN(member, ":", 2)
		if len(pieces) > 1 {
			pieces[1] = strings.ToLower(pieces[1])
		}
	}

	if len(pieces) > 0 {
		member = strings.Join(pieces, ":")
	}
	return member
}

// Construct map of role to set of members from list of bindings.
func createIamBindingsMap(bindings []*cloudresourcemanager.Binding) map[iamBindingKey]map[string]struct{} {
	bm := make(map[iamBindingKey]map[string]struct{})
	// Get each binding
	for _, b := range bindings {
		members := make(map[string]struct{})
		key := iamBindingKey{b.Role, conditionKeyFromCondition(b.Condition)}
		// Initialize members map
		if _, ok := bm[key]; ok {
			members = bm[key]
		}
		// Get each member (user/principal) for the binding
		for _, m := range b.Members {
			m = normalizeIamMemberCasing(m)
			// Add the member
			members[m] = struct{}{}
		}
		if len(members) > 0 {
			bm[key] = members
		} else {
			delete(bm, key)
		}
	}
	return bm
}

// Return list of Bindings for a map of role to member sets
func listFromIamBindingMap(bm map[iamBindingKey]map[string]struct{}) []*cloudresourcemanager.Binding {
	rb := make([]*cloudresourcemanager.Binding, 0, len(bm))
	var keys []iamBindingKey
	for k := range bm {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		keyI := keys[i]
		keyJ := keys[j]
		return fmt.Sprintf("%s%s", keyI.Role, keyI.Condition.String()) < fmt.Sprintf("%s%s", keyJ.Role, keyJ.Condition.String())
	})
	for _, key := range keys {
		members := bm[key]
		if len(members) == 0 {
			continue
		}
		b := &cloudresourcemanager.Binding{
			Role:    key.Role,
			Members: tpgresource.StringSliceFromGolangSet(members),
		}
		if !key.Condition.Empty() {
			b.Condition = &cloudresourcemanager.Expr{
				Description: key.Condition.Description,
				Expression:  key.Condition.Expression,
				Title:       key.Condition.Title,
			}
		}
		rb = append(rb, b)
	}
	return rb
}

// Flattens AuditConfigs so each role has a single Binding with combined members\
func removeAllAuditConfigsWithService(ac []*cloudresourcemanager.AuditConfig, service string) []*cloudresourcemanager.AuditConfig {
	acMap := createIamAuditConfigsMap(ac)
	delete(acMap, service)
	return listFromIamAuditConfigMap(acMap)
}

// Build a AuditConfig service to audit log config map
func createIamAuditConfigsMap(auditConfigs []*cloudresourcemanager.AuditConfig) map[string]map[string]map[string]struct{} {
	acMap := make(map[string]map[string]map[string]struct{})

	for _, ac := range auditConfigs {
		if _, ok := acMap[ac.Service]; !ok {
			acMap[ac.Service] = make(map[string]map[string]struct{})
		}
		alcMap := acMap[ac.Service]
		for _, alc := range ac.AuditLogConfigs {
			if _, ok := alcMap[alc.LogType]; !ok {
				alcMap[alc.LogType] = make(map[string]struct{})
			}
			memberMap := alcMap[alc.LogType]
			// Add members to map for log type.
			for _, m := range alc.ExemptedMembers {
				memberMap[m] = struct{}{}
			}
		}
	}

	return acMap
}

// Construct list of AuditConfigs from audit config maps.
func listFromIamAuditConfigMap(acMap map[string]map[string]map[string]struct{}) []*cloudresourcemanager.AuditConfig {
	ac := make([]*cloudresourcemanager.AuditConfig, 0, len(acMap))

	for service, logConfigMap := range acMap {
		if len(logConfigMap) == 0 {
			continue
		}

		logConfigs := make([]*cloudresourcemanager.AuditLogConfig, 0, len(logConfigMap))
		for logType, memberSet := range logConfigMap {
			alc := &cloudresourcemanager.AuditLogConfig{
				LogType:         logType,
				ForceSendFields: []string{"exemptedMembers"},
			}
			if len(memberSet) > 0 {
				alc.ExemptedMembers = tpgresource.StringSliceFromGolangSet(memberSet)
			}
			logConfigs = append(logConfigs, alc)
		}

		ac = append(ac, &cloudresourcemanager.AuditConfig{
			Service:         service,
			AuditLogConfigs: logConfigs,
		})
	}
	return ac
}

func jsonPolicyDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if old == "" && new == "" {
		return true
	}

	var oldPolicy, newPolicy cloudresourcemanager.Policy
	if old != "" && new != "" {
		if err := json.Unmarshal([]byte(old), &oldPolicy); err != nil {
			log.Printf("[ERROR] Could not unmarshal old policy %s: %v", old, err)
			return false
		}
		if err := json.Unmarshal([]byte(new), &newPolicy); err != nil {
			log.Printf("[ERROR] Could not unmarshal new policy %s: %v", new, err)
			return false
		}

		return compareIamPolicies(&newPolicy, &oldPolicy)
	}

	return false
}

func compareIamPolicies(a, b *cloudresourcemanager.Policy) bool {
	if a.Etag != b.Etag {
		log.Printf("[DEBUG] policies etag differ: %q vs %q", a.Etag, b.Etag)
		return false
	}
	if a.Version != b.Version {
		log.Printf("[DEBUG] policies version differ: %q vs %q", a.Version, b.Version)
		return false
	}
	if !CompareBindings(a.Bindings, b.Bindings) {
		log.Printf("[DEBUG] policies bindings differ: %#v vs %#v", a.Bindings, b.Bindings)
		return false
	}
	if !CompareAuditConfigs(a.AuditConfigs, b.AuditConfigs) {
		log.Printf("[DEBUG] policies audit configs differ: %#v vs %#v", a.AuditConfigs, b.AuditConfigs)
		return false
	}
	return true
}

func CompareBindings(a, b []*cloudresourcemanager.Binding) bool {
	aMap := createIamBindingsMap(a)
	bMap := createIamBindingsMap(b)
	return reflect.DeepEqual(aMap, bMap)
}

func CompareAuditConfigs(a, b []*cloudresourcemanager.AuditConfig) bool {
	aMap := createIamAuditConfigsMap(a)
	bMap := createIamAuditConfigsMap(b)
	return reflect.DeepEqual(aMap, bMap)
}

type IamSettings struct {
	DeprecationMessage string
	EnableBatching     bool
}

func NewIamSettings(options ...func(*IamSettings)) *IamSettings {
	settings := &IamSettings{}
	for _, o := range options {
		o(settings)
	}
	return settings
}

func IamWithDeprecationMessage(message string) func(s *IamSettings) {
	return func(s *IamSettings) {
		s.DeprecationMessage = message
	}
}

func IamWithGAResourceDeprecation() func(s *IamSettings) {
	return IamWithDeprecationMessage("This resource has been deprecated in the google (GA) provider, and will only be available in the google-beta provider in a future release.")
}

func IamWithBatching(s *IamSettings) {
	s.EnableBatching = true
}

// Util to deref and print auditConfigs
func DebugPrintAuditConfigs(bs []*cloudresourcemanager.AuditConfig) string {
	v, _ := json.MarshalIndent(bs, "", "\t")
	return string(v)
}

// Util to deref and print bindings
func DebugPrintBindings(bs []*cloudresourcemanager.Binding) string {
	v, _ := json.MarshalIndent(bs, "", "\t")
	return string(v)
}

// Returns a map representing iam bindings that are in the first map but not the second.
func missingBindingsMap(aMap, bMap map[iamBindingKey]map[string]struct{}) map[iamBindingKey]map[string]struct{} {
	results := make(map[iamBindingKey]map[string]struct{})
	for key, aMembers := range aMap {
		if bMembers, ok := bMap[key]; ok {
			// The key is in both maps.
			resultMembers := make(map[string]struct{})

			for aMember := range aMembers {
				if _, ok := bMembers[aMember]; !ok {
					// The member is in a but not in b.
					resultMembers[aMember] = struct{}{}
				}
			}
			for bMember := range bMembers {
				if _, ok := aMembers[bMember]; !ok {
					// The member is in b but not in a.
					resultMembers[bMember] = struct{}{}
				}
			}

			if len(resultMembers) > 0 {
				results[key] = resultMembers
			}
		} else {
			// The key is in map a but not map b.
			results[key] = aMembers
		}
	}

	for key, bMembers := range bMap {
		if _, ok := aMap[key]; !ok {
			// The key is in map b but not map a.
			results[key] = bMembers
		}
	}

	return results
}

// Returns the bindings that are in the first set of bindings but not the second.
func MissingBindings(a, b []*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	aMap := createIamBindingsMap(a)
	bMap := createIamBindingsMap(b)

	var results []*cloudresourcemanager.Binding
	for key, membersSet := range missingBindingsMap(aMap, bMap) {
		members := make([]string, 0, len(membersSet))
		for member := range membersSet {
			members = append(members, member)
		}
		results = append(results, &cloudresourcemanager.Binding{
			Role:    key.Role,
			Members: members,
		})
	}
	return results
}
