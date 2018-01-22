package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/storage/v1"
)

func resourceStorageDefaultObjectAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageDefaultObjectAclCreate,
		Read:   resourceStorageDefaultObjectAclRead,
		Update: resourceStorageDefaultObjectAclUpdate,
		Delete: resourceStorageDefaultObjectAclDelete,

		Schema: map[string]*schema.Schema{
			"bucket": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"role_entity": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func getDefaultObjectAclId(bucket string) string {
	return bucket + "-default-object-acl"
}

func resourceStorageDefaultObjectAclCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)
	role_entity := make([]interface{}, 0)

	if v, ok := d.GetOk("role_entity"); ok {
		role_entity = v.([]interface{})
	}

	if len(role_entity) > 0 {
		for _, v := range role_entity {
			pair, err := getRoleEntityPair(v.(string))

			ObjectAccessControl := &storage.ObjectAccessControl{
				Role:   pair.Role,
				Entity: pair.Entity,
			}

			log.Printf("[DEBUG]: setting role = %s, entity = %s on bucket %s", pair.Role, pair.Entity, bucket)

			_, err = config.clientStorage.DefaultObjectAccessControls.Insert(bucket, ObjectAccessControl).Do()

			if err != nil {
				return fmt.Errorf("Error setting Default Object ACL for %s on bucket %s: %v", pair.Entity, bucket, err)
			}
		}

		return resourceStorageDefaultObjectAclRead(d, meta)
	}
	return nil
}

func resourceStorageDefaultObjectAclRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)

	if _, ok := d.GetOk("role_entity"); ok {
		role_entity := make([]interface{}, 0)
		re_local := d.Get("role_entity").([]interface{})
		re_local_map := make(map[string]string)
		for _, v := range re_local {
			res, err := getRoleEntityPair(v.(string))

			if err != nil {
				return fmt.Errorf(
					"Old state has malformed Role/Entity pair: %v", err)
			}

			re_local_map[res.Entity] = res.Role
		}

		res, err := config.clientStorage.DefaultObjectAccessControls.List(bucket).Do()

		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Storage Default Object ACL for bucket %q", d.Get("bucket").(string)))
		}

		for _, v := range res.Items {
			role := v.Role
			entity := v.Entity
			// We only store updates to the locally defined access controls
			if _, in := re_local_map[entity]; in {
				role_entity = append(role_entity, fmt.Sprintf("%s:%s", role, entity))
				log.Printf("[DEBUG]: saving re %s-%s", v.Role, v.Entity)
			}
		}

		d.Set("role_entity", role_entity)
	}

	d.SetId(getDefaultObjectAclId(bucket))
	return nil
}

func resourceStorageDefaultObjectAclUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)

	if d.HasChange("role_entity") {
		o, n := d.GetChange("role_entity")
		old_re := o.([]interface{})
		new_re := n.([]interface{})

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

			ObjectAccessControl := &storage.ObjectAccessControl{
				Role:   pair.Role,
				Entity: pair.Entity,
			}

			// If the old state is missing for this entity, it needs to
			// be created. Otherwise it is updated
			if _, ok := old_re_map[pair.Entity]; ok {
				_, err = config.clientStorage.DefaultObjectAccessControls.Update(
					bucket, pair.Entity, ObjectAccessControl).Do()
			} else {
				_, err = config.clientStorage.DefaultObjectAccessControls.Insert(
					bucket, ObjectAccessControl).Do()
			}

			// Now we only store the keys that have to be removed
			delete(old_re_map, pair.Entity)

			if err != nil {
				return fmt.Errorf("Error updating Storage Default Object ACL for bucket %s: %v", bucket, err)
			}
		}

		for entity, _ := range old_re_map {
			log.Printf("[DEBUG]: removing entity %s", entity)
			err := config.clientStorage.DefaultObjectAccessControls.Delete(bucket, entity).Do()

			if err != nil {
				return fmt.Errorf("Error updating Storage Default Object ACL for bucket %s: %v", bucket, err)
			}
		}

		return resourceStorageDefaultObjectAclRead(d, meta)
	}

	return nil
}

func resourceStorageDefaultObjectAclDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)

	re_local := d.Get("role_entity").([]interface{})
	for _, v := range re_local {
		res, err := getRoleEntityPair(v.(string))
		if err != nil {
			return err
		}

		log.Printf("[DEBUG]: removing entity %s", res.Entity)

		err = config.clientStorage.DefaultObjectAccessControls.Delete(bucket, res.Entity).Do()

		if err != nil {
			return fmt.Errorf("Error deleting entity %s ACL: %s", res.Entity, err)
		}
	}

	return nil
}
