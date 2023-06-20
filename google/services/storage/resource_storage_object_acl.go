// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/api/storage/v1"
)

func ResourceStorageObjectAcl() *schema.Resource {
	return &schema.Resource{
		Create:        resourceStorageObjectAclCreate,
		Read:          resourceStorageObjectAclRead,
		Update:        resourceStorageObjectAclUpdate,
		Delete:        resourceStorageObjectAclDelete,
		CustomizeDiff: resourceStorageObjectAclDiff,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"object": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"predefined_acl": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"role_entity"},
				ValidateFunc:  validation.StringInSlice([]string{"private", "bucketOwnerRead", "bucketOwnerFullControl", "projectPrivate", "authenticatedRead", "publicRead", ""}, false),
			},

			"role_entity": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateRoleEntityPair,
				},
				ConflictsWith: []string{"predefined_acl"},
			},
		},
		UseJSONNumber: true,
	}
}

// We can't edit the object owner (at risk of 403 errors), and users will always see a diff if they
// don't explicitly specify that it has OWNER permissions.
// Suppressing it means their configs won't be *strictly* correct as they will be missing the object
// owner having OWNER. It's impossible to remove that permission though, so this custom diff
// makes configs with or without that line indistinguishable.
func resourceStorageObjectAclDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	bucket, ok := diff.GetOk("bucket")
	if !ok {
		// During `plan` when this is interpolated from a resource that hasn't been created yet
		// required fields may not be present yet
		return nil
	}
	object, ok := diff.GetOk("object")
	if !ok {
		// During `plan` when this is interpolated from a resource that hasn't been created yet
		// required fields may not be present yet
		return nil
	}

	sObject, err := config.NewStorageClient(config.UserAgent).Objects.Get(bucket.(string), object.(string)).Projection("full").Do()
	if err != nil {
		// Failing here is OK! Generally, it means we are at Create although it could mean the resource is gone.
		// Create won't show the object owner being given
		return nil
	}

	var objectOwner string
	if sObject.Owner != nil {
		objectOwner = sObject.Owner.Entity
	}
	ownerRole := fmt.Sprintf("%s:%s", "OWNER", objectOwner)
	oldRoleEntity, newRoleEntity := diff.GetChange("role_entity")

	// We can fail at plan time if the object owner/creator is being set to
	// a reader
	for _, v := range newRoleEntity.(*schema.Set).List() {
		res := getValidatedRoleEntityPair(v.(string))

		if res.Entity == objectOwner && res.Role != "OWNER" {
			return fmt.Errorf("New state tried setting object owner entity (%s) to non-'OWNER' role", objectOwner)
		}
	}

	// Diffs won't match in Plan and Apply pre-create if we naively add the RE
	// every time. So instead, we check to see if the old state (upstream/gcp
	// because we will have just done a refresh) contains it first.
	if oldRoleEntity.(*schema.Set).Contains(ownerRole) &&
		!newRoleEntity.(*schema.Set).Contains(ownerRole) {
		newRoleEntity.(*schema.Set).Add(ownerRole)
		return diff.SetNew("role_entity", newRoleEntity)
	}

	return nil
}

func getObjectAclId(object string) string {
	return object + "-acl"
}

func resourceStorageObjectAclCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)
	object := d.Get("object").(string)

	lockName, err := tpgresource.ReplaceVars(d, config, "storage/buckets/{{bucket}}/objects/{{object}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	// If we're using a predefined acl we just use the canned api.
	if predefinedAcl, ok := d.GetOk("predefined_acl"); ok {
		res, err := config.NewStorageClient(userAgent).Objects.Get(bucket, object).Do()
		if err != nil {
			return fmt.Errorf("Error reading object %s in %s: %v", object, bucket, err)
		}

		_, err = config.NewStorageClient(userAgent).Objects.Update(bucket, object, res).PredefinedAcl(predefinedAcl.(string)).Do()
		if err != nil {
			return fmt.Errorf("Error updating object %s in %s: %v", object, bucket, err)
		}

		return resourceStorageObjectAclRead(d, meta)
	} else if reMap := d.Get("role_entity").(*schema.Set); reMap.Len() > 0 {
		sObject, err := config.NewStorageClient(userAgent).Objects.Get(bucket, object).Projection("full").Do()
		if err != nil {
			return fmt.Errorf("Error reading object %s in %s: %v", object, bucket, err)
		}

		var objectOwner string
		if sObject.Owner != nil {
			objectOwner = sObject.Owner.Entity
		}
		roleEntitiesUpstream, err := getRoleEntitiesAsStringsFromApi(config, bucket, object, userAgent)
		if err != nil {
			return fmt.Errorf("Error reading object %s in %s: %v", object, bucket, err)
		}

		create, update, remove, err := getRoleEntityChange(roleEntitiesUpstream, tpgresource.ConvertStringSet(reMap), objectOwner)
		if err != nil {
			return fmt.Errorf("Error reading object %s in %s. Invalid schema: %v", object, bucket, err)
		}

		err = performStorageObjectRoleEntityOperations(create, update, remove, config, bucket, object, userAgent)
		if err != nil {
			return fmt.Errorf("Error creating object %s in %s: %v", object, bucket, err)
		}

		return resourceStorageObjectAclRead(d, meta)
	}

	return fmt.Errorf("Error, you must specify either \"predefined_acl\" or \"role_entity\"")
}

func resourceStorageObjectAclRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)
	object := d.Get("object").(string)

	roleEntities, err := getRoleEntitiesAsStringsFromApi(config, bucket, object, userAgent)
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Storage Object ACL for Bucket %q", d.Get("bucket").(string)))
	}

	err = d.Set("role_entity", roleEntities)
	if err != nil {
		return err
	}

	d.SetId(getObjectAclId(object))
	return nil
}

func resourceStorageObjectAclUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)
	object := d.Get("object").(string)

	lockName, err := tpgresource.ReplaceVars(d, config, "storage/buckets/{{bucket}}/objects/{{object}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	if _, ok := d.GetOk("predefined_acl"); d.HasChange("predefined_acl") && ok {
		res, err := config.NewStorageClient(userAgent).Objects.Get(bucket, object).Do()
		if err != nil {
			return fmt.Errorf("Error reading object %s in %s: %v", object, bucket, err)
		}

		_, err = config.NewStorageClient(userAgent).Objects.Update(bucket, object, res).PredefinedAcl(d.Get("predefined_acl").(string)).Do()
		if err != nil {
			return fmt.Errorf("Error updating object %s in %s: %v", object, bucket, err)
		}

		return resourceStorageObjectAclRead(d, meta)
	} else {
		sObject, err := config.NewStorageClient(userAgent).Objects.Get(bucket, object).Projection("full").Do()
		if err != nil {
			return fmt.Errorf("Error reading object %s in %s: %v", object, bucket, err)
		}

		var objectOwner string
		if sObject.Owner != nil {
			objectOwner = sObject.Owner.Entity
		}

		o, n := d.GetChange("role_entity")
		create, update, remove, err := getRoleEntityChange(
			tpgresource.ConvertStringSet(o.(*schema.Set)),
			tpgresource.ConvertStringSet(n.(*schema.Set)),
			objectOwner)
		if err != nil {
			return fmt.Errorf("Error reading object %s in %s. Invalid schema: %v", object, bucket, err)
		}

		err = performStorageObjectRoleEntityOperations(create, update, remove, config, bucket, object, userAgent)
		if err != nil {
			return fmt.Errorf("Error updating object %s in %s: %v", object, bucket, err)
		}

		return resourceStorageObjectAclRead(d, meta)
	}
}

func resourceStorageObjectAclDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)
	object := d.Get("object").(string)

	lockName, err := tpgresource.ReplaceVars(d, config, "storage/buckets/{{bucket}}/objects/{{object}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	res, err := config.NewStorageClient(userAgent).Objects.Get(bucket, object).Do()
	if err != nil {
		return fmt.Errorf("Error reading object %s in %s: %v", object, bucket, err)
	}

	_, err = config.NewStorageClient(userAgent).Objects.Update(bucket, object, res).PredefinedAcl("private").Do()
	if err != nil {
		return fmt.Errorf("Error updating object %s in %s: %v", object, bucket, err)
	}

	return nil
}

func getRoleEntitiesAsStringsFromApi(config *transport_tpg.Config, bucket, object, userAgent string) ([]string, error) {
	var roleEntities []string
	res, err := config.NewStorageClient(userAgent).ObjectAccessControls.List(bucket, object).Do()
	if err != nil {
		return nil, err
	}

	for _, roleEntity := range res.Items {
		role := roleEntity.Role
		entity := roleEntity.Entity
		roleEntities = append(roleEntities, fmt.Sprintf("%s:%s", role, entity))
	}

	return roleEntities, nil
}

// Creates 3 lists of changes we need to make to go from one set of entities to another- which entities need to be created, update, and deleted
// Not resource specific
func getRoleEntityChange(old []string, new []string, owner string) (create, update, remove []*RoleEntity, err error) {
	newEntitiesUsed := make(map[string]struct{})
	for _, v := range new {
		res := getValidatedRoleEntityPair(v)
		if _, ok := newEntitiesUsed[res.Entity]; ok {
			return nil, nil, nil, fmt.Errorf("New state has duplicate Entity: %v", res.Entity)
		}

		newEntitiesUsed[res.Entity] = struct{}{}
	}

	oldEntitiesUsed := make(map[string]string)
	for _, v := range old {
		res := getValidatedRoleEntityPair(v)

		// Updating the owner will error out, so let's avoid it.
		if res.Entity == owner {
			continue
		}

		oldEntitiesUsed[res.Entity] = res.Role
	}

	for _, re := range new {
		res := getValidatedRoleEntityPair(re)

		// Updating the owner will error out, so let's never do it.
		if res.Entity == owner {
			continue
		}

		if v, ok := oldEntitiesUsed[res.Entity]; ok {
			if res.Role != v {
				update = append(update, res)
			}

			delete(oldEntitiesUsed, res.Entity)
		} else {
			create = append(create, res)
		}
	}

	for _, v := range old {
		res := getValidatedRoleEntityPair(v)

		if _, ok := oldEntitiesUsed[res.Entity]; ok {
			remove = append(remove, res)
		}
	}

	return create, update, remove, nil
}

// Takes in lists of changes to make to a Storage Object's ACL and makes those changes
func performStorageObjectRoleEntityOperations(create []*RoleEntity, update []*RoleEntity, remove []*RoleEntity, config *transport_tpg.Config, bucket, object, userAgent string) error {
	for _, v := range create {
		objectAccessControl := &storage.ObjectAccessControl{
			Role:   v.Role,
			Entity: v.Entity,
		}
		_, err := config.NewStorageClient(userAgent).ObjectAccessControls.Insert(bucket, object, objectAccessControl).Do()
		if err != nil {
			return fmt.Errorf("Error inserting ACL item %s for object %s in %s: %v", v.Entity, object, bucket, err)
		}
	}

	for _, v := range update {
		objectAccessControl := &storage.ObjectAccessControl{
			Role:   v.Role,
			Entity: v.Entity,
		}
		_, err := config.NewStorageClient(userAgent).ObjectAccessControls.Update(bucket, object, v.Entity, objectAccessControl).Do()
		if err != nil {
			return fmt.Errorf("Error updating ACL item %s for object %s in %s: %v", v.Entity, object, bucket, err)
		}
	}

	for _, v := range remove {
		err := config.NewStorageClient(userAgent).ObjectAccessControls.Delete(bucket, object, v.Entity).Do()
		if err != nil {
			return fmt.Errorf("Error deleting ACL item %s for object %s in %s: %v", v.Entity, object, bucket, err)
		}
	}

	return nil
}

func validateRoleEntityPair(v interface{}, k string) (ws []string, errors []error) {
	split := strings.Split(v.(string), ":")
	if len(split) != 2 {
		errors = append(errors, fmt.Errorf("Role entity pairs must be formatted as 'ROLE:entity': %s", v))
	}

	return
}

func getValidatedRoleEntityPair(roleEntity string) *RoleEntity {
	split := strings.Split(roleEntity, ":")
	return &RoleEntity{Role: split[0], Entity: split[1]}
}
