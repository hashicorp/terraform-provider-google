// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"google.golang.org/api/storage/v1"
)

func ResourceStorageBucketAcl() *schema.Resource {
	return &schema.Resource{
		Create:        resourceStorageBucketAclCreate,
		Read:          resourceStorageBucketAclRead,
		Update:        resourceStorageBucketAclUpdate,
		Delete:        resourceStorageBucketAclDelete,
		CustomizeDiff: resourceStorageRoleEntityCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the bucket it applies to.`,
			},

			"default_acl": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Configure this ACL to be the default ACL.`,
			},

			"predefined_acl": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"role_entity"},
				Description:   `The canned GCS ACL to apply. Must be set if role_entity is not.`,
			},

			"role_entity": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"predefined_acl"},
				Description:   `List of role/entity pairs in the form ROLE:entity. See GCS Bucket ACL documentation  for more details. Must be set if predefined_acl is not.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceStorageRoleEntityCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	keys := diff.GetChangedKeysPrefix("role_entity")
	if len(keys) < 1 {
		return nil
	}
	count := diff.Get("role_entity.#").(int)
	if count < 1 {
		return nil
	}
	state := map[string]struct{}{}
	conf := map[string]struct{}{}
	for i := 0; i < count; i++ {
		old, new := diff.GetChange(fmt.Sprintf("role_entity.%d", i))
		state[old.(string)] = struct{}{}
		conf[new.(string)] = struct{}{}
	}
	if len(state) != len(conf) {
		return nil
	}
	for k := range state {
		if _, ok := conf[k]; !ok {
			// project-owners- is explicitly stripped from the roles that this
			// resource will delete
			if strings.Contains(k, "OWNER:project-owners-") {
				continue
			}
			return nil
		}
	}
	return diff.Clear("role_entity")
}

type RoleEntity struct {
	Role   string
	Entity string
}

func getBucketAclId(bucket string) string {
	return bucket + "-acl"
}

func GetRoleEntityPair(role_entity string) (*RoleEntity, error) {
	split := strings.Split(role_entity, ":")
	if len(split) != 2 {
		return nil, fmt.Errorf("Error, each role entity pair must be " +
			"formatted as ROLE:entity")
	}

	return &RoleEntity{Role: split[0], Entity: split[1]}, nil
}

func resourceStorageBucketAclCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)
	predefined_acl := ""
	default_acl := ""
	role_entity := make([]interface{}, 0)

	if v, ok := d.GetOk("predefined_acl"); ok {
		predefined_acl = v.(string)
	}

	if v, ok := d.GetOk("role_entity"); ok {
		role_entity = v.([]interface{})
	}

	if v, ok := d.GetOk("default_acl"); ok {
		default_acl = v.(string)
	}

	lockName, err := tpgresource.ReplaceVars(d, config, "storage/buckets/{{bucket}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	if len(predefined_acl) > 0 {
		res, err := config.NewStorageClient(userAgent).Buckets.Get(bucket).Do()

		if err != nil {
			return fmt.Errorf("Error reading bucket %s: %v", bucket, err)
		}

		_, err = config.NewStorageClient(userAgent).Buckets.Update(bucket,
			res).PredefinedAcl(predefined_acl).Do()

		if err != nil {
			return fmt.Errorf("Error updating bucket %s: %v", bucket, err)
		}

	}

	if len(role_entity) > 0 {
		current, err := config.NewStorageClient(userAgent).BucketAccessControls.List(bucket).Do()
		if err != nil {
			return fmt.Errorf("Error retrieving current ACLs: %s", err)
		}
		for _, v := range role_entity {
			pair, err := GetRoleEntityPair(v.(string))
			if err != nil {
				return err
			}
			var alreadyInserted bool
			for _, cur := range current.Items {
				if cur.Entity == pair.Entity && cur.Role == pair.Role {
					alreadyInserted = true
					break
				}
			}
			if alreadyInserted {
				log.Printf("[DEBUG]: pair %s-%s already exists, not trying to insert again\n", pair.Role, pair.Entity)
				continue
			}
			bucketAccessControl := &storage.BucketAccessControl{
				Role:   pair.Role,
				Entity: pair.Entity,
			}

			log.Printf("[DEBUG]: storing re %s-%s", pair.Role, pair.Entity)

			_, err = config.NewStorageClient(userAgent).BucketAccessControls.Insert(bucket, bucketAccessControl).Do()

			if err != nil {
				return fmt.Errorf("Error updating ACL for bucket %s: %v", bucket, err)
			}
		}

	}

	if len(default_acl) > 0 {
		res, err := config.NewStorageClient(userAgent).Buckets.Get(bucket).Do()

		if err != nil {
			return fmt.Errorf("Error reading bucket %s: %v", bucket, err)
		}

		_, err = config.NewStorageClient(userAgent).Buckets.Update(bucket,
			res).PredefinedDefaultObjectAcl(default_acl).Do()

		if err != nil {
			return fmt.Errorf("Error updating bucket %s: %v", bucket, err)
		}

	}

	d.SetId(getBucketAclId(bucket))
	return resourceStorageBucketAclRead(d, meta)
}

func resourceStorageBucketAclRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)

	// The API offers no way to retrieve predefined ACLs,
	// and we can't tell which access controls were created
	// by the predefined roles, so...
	//
	// This is, needless to say, a bad state of affairs and
	// should be fixed.
	if _, ok := d.GetOk("role_entity"); ok {
		res, err := config.NewStorageClient(userAgent).BucketAccessControls.List(bucket).Do()

		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Storage Bucket ACL for bucket %q", d.Get("bucket").(string)))
		}
		entities := make([]string, 0, len(res.Items))
		for _, item := range res.Items {
			entities = append(entities, item.Role+":"+item.Entity)
		}

		if err := d.Set("role_entity", entities); err != nil {
			return fmt.Errorf("Error setting role_entity: %s", err)
		}
	} else {
		// if we don't set `role_entity` to nil (effectively setting it
		// to empty in Terraform state), because it's computed now,
		// Terraform will think it's missing from state, is supposed
		// to be there, and throw up a diff for role_entity.#. So it
		// must always be set in state.
		if err := d.Set("role_entity", nil); err != nil {
			return fmt.Errorf("Error setting role_entity: %s", err)
		}
	}

	return nil
}

func resourceStorageBucketAclUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)

	lockName, err := tpgresource.ReplaceVars(d, config, "storage/buckets/{{bucket}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	if d.HasChange("role_entity") {
		bkt, err := config.NewStorageClient(userAgent).Buckets.Get(bucket).Do()
		if err != nil {
			return fmt.Errorf("Error reading bucket %q: %v", bucket, err)
		}

		project := strconv.FormatUint(bkt.ProjectNumber, 10)
		o, n := d.GetChange("role_entity")
		old_re, new_re := o.([]interface{}), n.([]interface{})

		old_re_map := make(map[string]string)
		for _, v := range old_re {
			res, err := GetRoleEntityPair(v.(string))

			if err != nil {
				return fmt.Errorf(
					"Old state has malformed Role/Entity pair: %v", err)
			}

			old_re_map[res.Entity] = res.Role
		}

		for _, v := range new_re {
			pair, err := GetRoleEntityPair(v.(string))

			bucketAccessControl := &storage.BucketAccessControl{
				Role:   pair.Role,
				Entity: pair.Entity,
			}

			// If the old state entity's role doesn't match the new one, it needs to be inserted
			if old_re_map[pair.Entity] != bucketAccessControl.Role {
				_, err = config.NewStorageClient(userAgent).BucketAccessControls.Insert(
					bucket, bucketAccessControl).Do()
			}

			// Now we only store the keys that have to be removed
			delete(old_re_map, pair.Entity)

			if err != nil {
				return fmt.Errorf("Error updating ACL for bucket %s: %v", bucket, err)
			}
		}

		for entity, role := range old_re_map {
			if entity == fmt.Sprintf("project-owners-%s", project) && role == "OWNER" {
				log.Printf("[WARN]: Skipping %s-%s; not deleting owner ACL.", role, entity)
				continue
			}
			log.Printf("[DEBUG]: removing entity %s", entity)
			err := config.NewStorageClient(userAgent).BucketAccessControls.Delete(bucket, entity).Do()

			if err != nil {
				return fmt.Errorf("Error updating ACL for bucket %s: %v", bucket, err)
			}
		}

		return resourceStorageBucketAclRead(d, meta)
	}

	if d.HasChange("default_acl") {
		default_acl := d.Get("default_acl").(string)

		res, err := config.NewStorageClient(userAgent).Buckets.Get(bucket).Do()

		if err != nil {
			return fmt.Errorf("Error reading bucket %s: %v", bucket, err)
		}

		_, err = config.NewStorageClient(userAgent).Buckets.Update(bucket,
			res).PredefinedDefaultObjectAcl(default_acl).Do()

		if err != nil {
			return fmt.Errorf("Error updating bucket %s: %v", bucket, err)
		}

		return resourceStorageBucketAclRead(d, meta)
	}

	return nil
}

func resourceStorageBucketAclDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)

	lockName, err := tpgresource.ReplaceVars(d, config, "storage/buckets/{{bucket}}")
	if err != nil {
		return err
	}
	transport_tpg.MutexStore.Lock(lockName)
	defer transport_tpg.MutexStore.Unlock(lockName)

	bkt, err := config.NewStorageClient(userAgent).Buckets.Get(bucket).Do()
	if err != nil {
		return fmt.Errorf("Error retrieving bucket %q: %v", bucket, err)
	}
	project := strconv.FormatUint(bkt.ProjectNumber, 10)

	re_local := d.Get("role_entity").([]interface{})
	for _, v := range re_local {
		res, err := GetRoleEntityPair(v.(string))
		if err != nil {
			return err
		}

		if res.Entity == fmt.Sprintf("project-owners-%s", project) && res.Role == "OWNER" {
			log.Printf("[WARN]: Skipping %s-%s; not deleting owner ACL.", res.Role, res.Entity)
			continue
		}

		log.Printf("[DEBUG]: removing entity %s", res.Entity)

		err = config.NewStorageClient(userAgent).BucketAccessControls.Delete(bucket, res.Entity).Do()

		if err != nil {
			return fmt.Errorf("Error deleting entity %s ACL: %s", res.Entity, err)
		}
	}

	return nil
}
