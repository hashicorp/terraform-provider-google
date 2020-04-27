package google

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
)

func resourceGoogleFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceGoogleFolderCreate,
		Read:   resourceGoogleFolderRead,
		Update: resourceGoogleFolderUpdate,
		Delete: resourceGoogleFolderDelete,

		Importer: &schema.ResourceImporter{
			State: resourceGoogleFolderImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Read:   schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Format is either folders/{folder_id} or organizations/{org_id}.
			"parent": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Must be unique amongst its siblings.
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Format is 'folders/{folder_id}.
			// The terraform id holds the same value.
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifecycle_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGoogleFolderCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	displayName := d.Get("display_name").(string)
	parent := d.Get("parent").(string)

	var op *resourceManagerV2Beta1.Operation
	err := retryTimeDuration(func() error {
		var reqErr error
		op, reqErr = config.clientResourceManagerV2Beta1.Folders.Create(&resourceManagerV2Beta1.Folder{
			DisplayName: displayName,
		}).Parent(parent).Do()
		return reqErr
	}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating folder '%s' in '%s': %s", displayName, parent, err)
	}

	opAsMap, err := ConvertToMap(op)
	if err != nil {
		return err
	}

	err = resourceManagerOperationWaitTime(config, opAsMap, "creating folder", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating folder '%s' in '%s': %s", displayName, parent, err)
	}

	// Since we waited above, the operation is guaranteed to have been successful by this point.
	waitOp, err := config.clientResourceManager.Operations.Get(op.Name).Do()
	if err != nil {
		return fmt.Errorf("The folder '%s' has been created but we could not retrieve its id. Delete the folder manually and retry or use 'terraform import': %s", displayName, err)
	}

	// Requires 3 successive checks for safety. Nested IFs are used to avoid 3 error statement with the same message.
	var responseMap map[string]interface{}
	if err := json.Unmarshal(waitOp.Response, &responseMap); err == nil {
		if val, ok := responseMap["name"]; ok {
			if name, ok := val.(string); ok {
				d.SetId(name)
				return resourceGoogleFolderRead(d, meta)
			}
		}
	}
	return fmt.Errorf("The folder '%s' has been created but we could not retrieve its id. Delete the folder manually and retry or use 'terraform import'", displayName)
}

func resourceGoogleFolderRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	folder, err := getGoogleFolder(d.Id(), d, config)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Folder Not Found : %s", d.Id()))
	}

	d.Set("name", folder.Name)
	d.Set("parent", folder.Parent)
	d.Set("display_name", folder.DisplayName)
	d.Set("lifecycle_state", folder.LifecycleState)
	d.Set("create_time", folder.CreateTime)

	return nil
}

func resourceGoogleFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	displayName := d.Get("display_name").(string)

	d.Partial(true)
	if d.HasChange("display_name") {
		err := retry(func() error {
			_, reqErr := config.clientResourceManagerV2Beta1.Folders.Patch(d.Id(), &resourceManagerV2Beta1.Folder{
				DisplayName: displayName,
			}).Do()
			return reqErr
		})
		if err != nil {
			return fmt.Errorf("Error updating display_name to '%s': %s", displayName, err)
		}

		d.SetPartial("display_name")
	}

	if d.HasChange("parent") {
		newParent := d.Get("parent").(string)

		var op *resourceManagerV2Beta1.Operation
		err := retry(func() error {
			var reqErr error
			op, reqErr = config.clientResourceManagerV2Beta1.Folders.Move(d.Id(), &resourceManagerV2Beta1.MoveFolderRequest{
				DestinationParent: newParent,
			}).Do()
			return reqErr
		})
		if err != nil {
			return fmt.Errorf("Error moving folder '%s' to '%s': %s", displayName, newParent, err)
		}

		opAsMap, err := ConvertToMap(op)
		if err != nil {
			return err
		}

		err = resourceManagerOperationWaitTime(config, opAsMap, "move folder", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error moving folder '%s' to '%s': %s", displayName, newParent, err)
		}

		d.SetPartial("parent")
	}

	d.Partial(false)

	return nil
}

func resourceGoogleFolderDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	displayName := d.Get("display_name").(string)

	err := retryTimeDuration(func() error {
		_, reqErr := config.clientResourceManagerV2Beta1.Folders.Delete(d.Id()).Do()
		return reqErr
	}, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("Error deleting folder '%s': %s", displayName, err)
	}
	return nil
}

func resourceGoogleFolderImportState(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	if !strings.HasPrefix(d.Id(), "folders/") {
		id = fmt.Sprintf("folders/%s", id)
	}

	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

// Util to get a Folder resource from API. Note that folder described by name is not necessarily the
// ResourceData resource.
func getGoogleFolder(folderName string, d *schema.ResourceData, config *Config) (*resourceManagerV2Beta1.Folder, error) {
	var folder *resourceManagerV2Beta1.Folder
	err := retryTimeDuration(func() error {
		var reqErr error
		folder, reqErr = config.clientResourceManagerV2Beta1.Folders.Get(folderName).Do()
		return reqErr
	}, d.Timeout(schema.TimeoutRead))
	if err != nil {
		return nil, err
	}
	return folder, nil
}
