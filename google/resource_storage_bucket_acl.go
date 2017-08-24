package google

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"google.golang.org/api/storage/v1"
)

func resourceStorageBucketAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageBucketAclCreate,
		Read:   resourceStorageBucketAclRead,
		Update: resourceStorageBucketAclUpdate,
		Delete: resourceStorageBucketAclDelete,

		Schema: map[string]*schema.Schema{
			"bucket": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"default_acl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"predefined_acl": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"role_entity"},
			},

			"role_entity": &schema.Schema{
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"predefined_acl"},
			},
		},
	}
}

type RoleEntity struct {
	Role   string
	Entity string
}

func getBucketAclId(bucket string) string {
	return bucket + "-acl"
}

func getRoleEntityPair(role_entity string) (*RoleEntity, error) {
	split := strings.Split(role_entity, ":")
	if len(split) != 2 {
		return nil, fmt.Errorf("Error, each role entity pair must be " +
			"formatted as ROLE:entity")
	}

	return &RoleEntity{Role: split[0], Entity: split[1]}, nil
}

func resourceStorageBucketAclCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

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

	if len(predefined_acl) > 0 {
		res, err := config.clientStorage.Buckets.Get(bucket).Do()

		if err != nil {
			return fmt.Errorf("Error reading bucket %s: %v", bucket, err)
		}

		res, err = config.clientStorage.Buckets.Update(bucket,
			res).PredefinedAcl(predefined_acl).Do()

		if err != nil {
			return fmt.Errorf("Error updating bucket %s: %v", bucket, err)
		}

	}
	if len(role_entity) > 0 {
		for _, v := range role_entity {
			pair, err := getRoleEntityPair(v.(string))

			bucketAccessControl := &storage.BucketAccessControl{
				Role:   pair.Role,
				Entity: pair.Entity,
			}

			log.Printf("[DEBUG]: storing re %s-%s", pair.Role, pair.Entity)

			_, err = config.clientStorage.BucketAccessControls.Insert(bucket, bucketAccessControl).Do()

			if err != nil {
				return fmt.Errorf("Error updating ACL for bucket %s: %v", bucket, err)
			}
		}

	}

	if len(default_acl) > 0 {
		res, err := config.clientStorage.Buckets.Get(bucket).Do()

		if err != nil {
			return fmt.Errorf("Error reading bucket %s: %v", bucket, err)
		}

		res, err = config.clientStorage.Buckets.Update(bucket,
			res).PredefinedDefaultObjectAcl(default_acl).Do()

		if err != nil {
			return fmt.Errorf("Error updating bucket %s: %v", bucket, err)
		}

	}

	d.SetId(getBucketAclId(bucket))
	return resourceStorageBucketAclRead(d, meta)
}

func resourceStorageBucketAclRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)

	// The API offers no way to retrieve predefined ACLs,
	// and we can't tell which access controls were created
	// by the predefined roles, so...
	//
	// This is, needless to say, a bad state of affairs and
	// should be fixed.
	if _, ok := d.GetOk("role_entity"); ok {
		res, err := config.clientStorage.BucketAccessControls.List(bucket).Do()

		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Storage Bucket ACL for bucket %q", d.Get("bucket").(string)))
		}
		entities := make([]string, 0, len(res.Items))
		for _, item := range res.Items {
			entities = append(entities, item.Role+":"+item.Entity)
		}

		d.Set("role_entity", entities)
	}

	return nil
}

func resourceStorageBucketAclUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)

	if d.HasChange("role_entity") {
		o, n := d.GetChange("role_entity")
		old_re, new_re := o.([]interface{}), n.([]interface{})

		old_re_map := make(map[string]string)
		for _, v := range old_re {
			res, err := getRoleEntityPair(v.(string))

			if err != nil {
				return fmt.Errorf(
					"Old state has malformed Role/Entity pair: %v", err)
			}

			old_re_map[res.Entity] = res.Role
		}

		for _, v := range new_re {
			pair, err := getRoleEntityPair(v.(string))

			bucketAccessControl := &storage.BucketAccessControl{
				Role:   pair.Role,
				Entity: pair.Entity,
			}

			// If the old state is missing this entity, it needs to
			// be created. Otherwise it is updated
			if _, ok := old_re_map[pair.Entity]; ok {
				_, err = config.clientStorage.BucketAccessControls.Update(
					bucket, pair.Entity, bucketAccessControl).Do()
			} else {
				_, err = config.clientStorage.BucketAccessControls.Insert(
					bucket, bucketAccessControl).Do()
			}

			// Now we only store the keys that have to be removed
			delete(old_re_map, pair.Entity)

			if err != nil {
				return fmt.Errorf("Error updating ACL for bucket %s: %v", bucket, err)
			}
		}

		for entity, _ := range old_re_map {
			log.Printf("[DEBUG]: removing entity %s", entity)
			err := config.clientStorage.BucketAccessControls.Delete(bucket, entity).Do()

			if err != nil {
				return fmt.Errorf("Error updating ACL for bucket %s: %v", bucket, err)
			}
		}

		return resourceStorageBucketAclRead(d, meta)
	}

	if d.HasChange("default_acl") {
		default_acl := d.Get("default_acl").(string)

		res, err := config.clientStorage.Buckets.Get(bucket).Do()

		if err != nil {
			return fmt.Errorf("Error reading bucket %s: %v", bucket, err)
		}

		res, err = config.clientStorage.Buckets.Update(bucket,
			res).PredefinedDefaultObjectAcl(default_acl).Do()

		if err != nil {
			return fmt.Errorf("Error updating bucket %s: %v", bucket, err)
		}

		return resourceStorageBucketAclRead(d, meta)
	}

	return nil
}

func resourceStorageBucketAclDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)

	re_local := d.Get("role_entity").([]interface{})
	for _, v := range re_local {
		res, err := getRoleEntityPair(v.(string))
		if err != nil {
			return err
		}

		log.Printf("[DEBUG]: removing entity %s", res.Entity)

		err = config.clientStorage.BucketAccessControls.Delete(bucket, res.Entity).Do()

		if err != nil {
			return fmt.Errorf("Error deleting entity %s ACL: %s", res.Entity, err)
		}
	}

	return nil
}
