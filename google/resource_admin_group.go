package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	adminDirectory "google.golang.org/api/admin/directory/v1"
)

func resourceAdminGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAdminGroupCreate,
		Read:   resourceAdminGroupRead,
		Update: resourceAdminGroupUpdate,
		Delete: resourceAdminGroupDelete,

		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"direct_members_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"admin_created": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},

			"aliases": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"non_editable_aliases": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAdminGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	group := &adminDirectory.Group{
		Email: d.Get("email").(string),
	}

	if v, ok := d.GetOk("name"); ok {
		log.Printf("[DEBUG] Setting group name: %s", v.(string))
		group.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		log.Printf("[DEBUG] Setting group description: %s", v.(string))
		group.Description = v.(string)
	}

	createdGroup, err := config.clientAdminDirectory.Groups.Insert(group).Do()
	if err != nil {
		return fmt.Errorf("Error creating group: %s", err)
	}

	d.SetId(createdGroup.Id)
	log.Printf("[INFO] Created group: %s", createdGroup.Email)
	return resourceAdminGroupRead(d, meta)
}

func resourceAdminGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	group := &adminDirectory.Group{}
	nullFields := []string{}

	if d.HasChange("email") {
		log.Printf("[DEBUG] Updating group email: %s", d.Get("email").(string))
		group.Email = d.Get("email").(string)
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			log.Printf("[DEBUG] Updating group name: %s", v.(string))
			group.Name = v.(string)
		} else {
			log.Printf("[DEBUG] Removing group name")
			group.Name = ""
			nullFields = append(nullFields, "name")
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			log.Printf("[DEBUG] Updating group description: %s", v.(string))
			group.Description = v.(string)
		} else {
			log.Printf("[DEBUG] Removing group description")
			group.Description = ""
			nullFields = append(nullFields, "description")
		}
	}

	if len(nullFields) > 0 {
		group.NullFields = nullFields
	}

	updatedGroup, err := config.clientAdminDirectory.Groups.Patch(d.Id(), group).Do()
	if err != nil {
		return fmt.Errorf("Error updating group: %s", err)
	}

	log.Printf("[INFO] Updated group: %s", updatedGroup.Email)
	return resourceAdminGroupRead(d, meta)
}

func resourceAdminGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	group, err := config.clientAdminDirectory.Groups.Get(d.Id()).Do()
	if err != nil {
		return err
	}

	d.Set("id", group.Id)
	d.Set("direct_members_count", group.DirectMembersCount)
	d.Set("admin_created", group.AdminCreated)
	d.Set("aliases", group.Aliases)
	d.Set("non_editable_aliases", group.NonEditableAliases)

	return nil
}

func resourceAdminGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	err := config.clientAdminDirectory.Groups.Delete(d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting group: %s", err)
	}

	d.SetId("")
	return nil
}
