package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/storage/v1"
)

func resourceStorageDefaultObjectAcl() *schema.Resource {
	return &schema.Resource{
		Create:        resourceStorageDefaultObjectAclCreate,
		Read:          resourceStorageDefaultObjectAclRead,
		Update:        resourceStorageDefaultObjectAclUpdate,
		Delete:        resourceStorageDefaultObjectAclDelete,
		CustomizeDiff: resourceStorageRoleEntityCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"bucket": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"role_entity": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
		},
	}
}

func resourceStorageDefaultObjectAclCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)
	roleEntity := d.Get("role_entity").([]interface{})

	for _, v := range roleEntity {
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
	d.SetId(bucket)
	return resourceStorageDefaultObjectAclRead(d, meta)
}

func resourceStorageDefaultObjectAclRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)

	roleEntities := make([]interface{}, 0)
	reLocal := d.Get("role_entity").([]interface{})
	reLocalMap := make(map[string]string)
	for _, v := range reLocal {
		res, err := getRoleEntityPair(v.(string))

		if err != nil {
			return fmt.Errorf(
				"Old state has malformed Role/Entity pair: %v", err)
		}

		reLocalMap[res.Entity] = res.Role
	}

	res, err := config.clientStorage.DefaultObjectAccessControls.List(bucket).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Storage Default Object ACL for bucket %q", d.Get("bucket").(string)))
	}

	for _, v := range res.Items {
		role := v.Role
		entity := v.Entity
		// We only store updates to the locally defined access controls
		if _, in := reLocalMap[entity]; in {
			roleEntities = append(roleEntities, fmt.Sprintf("%s:%s", role, entity))
			log.Printf("[DEBUG]: saving re %s-%s", v.Role, v.Entity)
		}
	}

	d.Set("role_entity", roleEntities)

	return nil
}

func resourceStorageDefaultObjectAclUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)

	if !d.HasChange("role_entity") {
		return nil
	}
	o, n := d.GetChange("role_entity")
	oldRe := o.([]interface{})
	newRe := n.([]interface{})

	oldReMap := make(map[string]string)
	for _, v := range oldRe {
		res, err := getRoleEntityPair(v.(string))

		if err != nil {
			return fmt.Errorf(
				"Old state has malformed Role/Entity pair: %v", err)
		}

		oldReMap[res.Entity] = res.Role
	}

	for _, v := range newRe {
		pair, err := getRoleEntityPair(v.(string))

		ObjectAccessControl := &storage.ObjectAccessControl{
			Role:   pair.Role,
			Entity: pair.Entity,
		}

		// If the old state is present for the  entity, it is updated
		// If the old state is missing, it is inserted
		if _, ok := oldReMap[pair.Entity]; ok {
			_, err = config.clientStorage.DefaultObjectAccessControls.Update(
				bucket, pair.Entity, ObjectAccessControl).Do()
		} else {
			_, err = config.clientStorage.DefaultObjectAccessControls.Insert(
				bucket, ObjectAccessControl).Do()
		}

		// Now we only store the keys that have to be removed
		delete(oldReMap, pair.Entity)

		if err != nil {
			return fmt.Errorf("Error updating Storage Default Object ACL for bucket %s: %v", bucket, err)
		}
	}

	for entity := range oldReMap {
		log.Printf("[DEBUG]: removing entity %s", entity)
		err := config.clientStorage.DefaultObjectAccessControls.Delete(bucket, entity).Do()

		if err != nil {
			return fmt.Errorf("Error updating Storage Default Object ACL for bucket %s: %v", bucket, err)
		}
	}

	return resourceStorageDefaultObjectAclRead(d, meta)
}

func resourceStorageDefaultObjectAclDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)

	reLocal := d.Get("role_entity").([]interface{})
	for _, v := range reLocal {
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
