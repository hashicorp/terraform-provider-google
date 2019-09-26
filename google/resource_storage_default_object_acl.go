package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/storage/v1"
)

func resourceStorageDefaultObjectAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageDefaultObjectAclCreateUpdate,
		Read:   resourceStorageDefaultObjectAclRead,
		Update: resourceStorageDefaultObjectAclCreateUpdate,
		Delete: resourceStorageDefaultObjectAclDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"role_entity": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateRoleEntityPair,
				},
			},
		},
	}
}

func resourceStorageDefaultObjectAclCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)
	defaultObjectAcl := []*storage.ObjectAccessControl{}
	for _, v := range d.Get("role_entity").(*schema.Set).List() {
		pair := getValidatedRoleEntityPair(v.(string))
		defaultObjectAcl = append(defaultObjectAcl, &storage.ObjectAccessControl{
			Role:   pair.Role,
			Entity: pair.Entity,
		})
	}

	res, err := config.clientStorage.Buckets.Get(bucket).Do()
	if err != nil {
		return fmt.Errorf("Error reading bucket %s: %v", bucket, err)
	}

	// Even with ForceSendFields the empty array wasn't working. Luckily, this is the same thing
	if len(defaultObjectAcl) == 0 {
		_, err = config.clientStorage.Buckets.Update(bucket, res).IfMetagenerationMatch(res.Metageneration).PredefinedDefaultObjectAcl("private").Do()
		if err != nil {
			return fmt.Errorf("Error updating default object acl to empty for bucket %s: %v", bucket, err)
		}
	} else {
		res.DefaultObjectAcl = defaultObjectAcl
		_, err = config.clientStorage.Buckets.Update(bucket, res).IfMetagenerationMatch(res.Metageneration).Do()
		if err != nil {
			return fmt.Errorf("Error updating default object acl for bucket %s: %v", bucket, err)
		}
	}

	return resourceStorageDefaultObjectAclRead(d, meta)
}

func resourceStorageDefaultObjectAclRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)
	res, err := config.clientStorage.Buckets.Get(bucket).Projection("full").Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Default Storage Object ACL for Bucket %q", d.Get("bucket").(string)))
	}

	var roleEntities []string
	for _, roleEntity := range res.DefaultObjectAcl {
		role := roleEntity.Role
		entity := roleEntity.Entity
		roleEntities = append(roleEntities, fmt.Sprintf("%s:%s", role, entity))
	}

	err = d.Set("role_entity", roleEntities)
	if err != nil {
		return err
	}

	d.SetId(bucket)
	return nil
}

func resourceStorageDefaultObjectAclDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)
	res, err := config.clientStorage.Buckets.Get(bucket).Do()
	if err != nil {
		return fmt.Errorf("Error reading bucket %s: %v", bucket, err)
	}

	_, err = config.clientStorage.Buckets.Update(bucket, res).IfMetagenerationMatch(res.Metageneration).PredefinedDefaultObjectAcl("private").Do()
	if err != nil {
		return fmt.Errorf("Error deleting (updating to private) default object acl for bucket %s: %v", bucket, err)
	}

	return nil
}
