package google

import (
	"fmt"
	"log"
	"strings"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceComputeInstanceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInstanceGroupCreate,
		Read:   resourceComputeInstanceGroupRead,
		Update: resourceComputeInstanceGroupUpdate,
		Delete: resourceComputeInstanceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeInstanceGroupImportState,
		},

		SchemaVersion: 2,
		MigrateState:  resourceComputeInstanceGroupMigrateState,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"instances": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      selfLinkRelativePathHash,
			},

			"named_port": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"network": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				ForceNew:         true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func getInstanceReferences(instanceUrls []string) (refs []*compute.InstanceReference) {
	for _, v := range instanceUrls {
		refs = append(refs, &compute.InstanceReference{
			Instance: v,
		})
	}
	return refs
}

func validInstanceURLs(instanceUrls []string) bool {
	for _, v := range instanceUrls {
		if !strings.HasPrefix(v, "https://") {
			return false
		}
	}
	return true
}

func resourceComputeInstanceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	// Build the parameter
	instanceGroup := &compute.InstanceGroup{
		Name: name,
	}

	// Set optional fields
	if v, ok := d.GetOk("description"); ok {
		instanceGroup.Description = v.(string)
	}

	if v, ok := d.GetOk("named_port"); ok {
		instanceGroup.NamedPorts = getNamedPorts(v.([]interface{}))
	}

	if v, ok := d.GetOk("network"); ok {
		instanceGroup.Network = v.(string)
	}

	log.Printf("[DEBUG] InstanceGroup insert request: %#v", instanceGroup)
	op, err := config.clientCompute.InstanceGroups.Insert(
		project, zone, instanceGroup).Do()
	if err != nil {
		return fmt.Errorf("Error creating InstanceGroup: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(fmt.Sprintf("%s/%s", zone, name))

	// Wait for the operation to complete
	err = computeOperationWait(config.clientCompute, op, project, "Creating InstanceGroup")
	if err != nil {
		d.SetId("")
		return err
	}

	if v, ok := d.GetOk("instances"); ok {
		instanceUrls := convertStringArr(v.(*schema.Set).List())
		if !validInstanceURLs(instanceUrls) {
			return fmt.Errorf("Error invalid instance URLs: %v", instanceUrls)
		}

		addInstanceReq := &compute.InstanceGroupsAddInstancesRequest{
			Instances: getInstanceReferences(instanceUrls),
		}

		log.Printf("[DEBUG] InstanceGroup add instances request: %#v", addInstanceReq)
		op, err := config.clientCompute.InstanceGroups.AddInstances(
			project, zone, name, addInstanceReq).Do()
		if err != nil {
			return fmt.Errorf("Error adding instances to InstanceGroup: %s", err)
		}

		// Wait for the operation to complete
		err = computeOperationWait(config.clientCompute, op, project, "Adding instances to InstanceGroup")
		if err != nil {
			return err
		}
	}

	return resourceComputeInstanceGroupRead(d, meta)
}

func resourceComputeInstanceGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	// retrieve instance group
	instanceGroup, err := config.clientCompute.InstanceGroups.Get(
		project, zone, name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Instance Group %q", name))
	}

	// retrieve instance group members
	var memberUrls []string
	members, err := config.clientCompute.InstanceGroups.ListInstances(
		project, zone, name, &compute.InstanceGroupsListInstancesRequest{
			InstanceState: "ALL",
		}).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			// The resource doesn't have any instances
			d.Set("instances", nil)
		} else {
			// any other errors return them
			return fmt.Errorf("Error reading InstanceGroup Members: %s", err)
		}
	} else {
		for _, member := range members.Items {
			memberUrls = append(memberUrls, member.Instance)
		}
		log.Printf("[DEBUG] InstanceGroup members: %v", memberUrls)
		d.Set("instances", memberUrls)
	}

	d.Set("named_port", flattenNamedPorts(instanceGroup.NamedPorts))
	d.Set("description", instanceGroup.Description)

	// Set computed fields
	d.Set("network", instanceGroup.Network)
	d.Set("size", instanceGroup.Size)
	d.Set("project", project)
	d.Set("zone", zone)
	d.Set("self_link", instanceGroup.SelfLink)

	return nil
}
func resourceComputeInstanceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	d.Partial(true)

	if d.HasChange("instances") {
		// to-do check for no instances
		from_, to_ := d.GetChange("instances")

		from := convertStringArr(from_.(*schema.Set).List())
		to := convertStringArr(to_.(*schema.Set).List())

		if !validInstanceURLs(from) {
			return fmt.Errorf("Error invalid instance URLs: %v", from)
		}
		if !validInstanceURLs(to) {
			return fmt.Errorf("Error invalid instance URLs: %v", to)
		}

		add, remove := calcAddRemove(from, to)

		if len(remove) > 0 {
			removeReq := &compute.InstanceGroupsRemoveInstancesRequest{
				Instances: getInstanceReferences(remove),
			}

			log.Printf("[DEBUG] InstanceGroup remove instances request: %#v", removeReq)
			removeOp, err := config.clientCompute.InstanceGroups.RemoveInstances(
				project, zone, name, removeReq).Do()
			if err != nil {
				if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
					log.Printf("[WARN] Instances already removed from InstanceGroup: %s", remove)
				} else {
					return fmt.Errorf("Error removing instances from InstanceGroup: %s", err)
				}
			} else {
				// Wait for the operation to complete
				err = computeOperationWait(config.clientCompute, removeOp, project, "Updating InstanceGroup")
				if err != nil {
					return err
				}
			}
		}

		if len(add) > 0 {

			addReq := &compute.InstanceGroupsAddInstancesRequest{
				Instances: getInstanceReferences(add),
			}

			log.Printf("[DEBUG] InstanceGroup adding instances request: %#v", addReq)
			addOp, err := config.clientCompute.InstanceGroups.AddInstances(
				project, zone, name, addReq).Do()
			if err != nil {
				return fmt.Errorf("Error adding instances from InstanceGroup: %s", err)
			}

			// Wait for the operation to complete
			err = computeOperationWait(config.clientCompute, addOp, project, "Updating InstanceGroup")
			if err != nil {
				return err
			}
		}

		d.SetPartial("instances")
	}

	if d.HasChange("named_port") {
		namedPorts := getNamedPorts(d.Get("named_port").([]interface{}))

		namedPortsReq := &compute.InstanceGroupsSetNamedPortsRequest{
			NamedPorts: namedPorts,
		}

		log.Printf("[DEBUG] InstanceGroup updating named ports request: %#v", namedPortsReq)
		op, err := config.clientCompute.InstanceGroups.SetNamedPorts(
			project, zone, name, namedPortsReq).Do()
		if err != nil {
			return fmt.Errorf("Error updating named ports for InstanceGroup: %s", err)
		}

		err = computeOperationWait(config.clientCompute, op, project, "Updating InstanceGroup")
		if err != nil {
			return err
		}
		d.SetPartial("named_port")
	}

	d.Partial(false)

	return resourceComputeInstanceGroupRead(d, meta)
}

func resourceComputeInstanceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	zone, err := getZone(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	op, err := config.clientCompute.InstanceGroups.Delete(project, zone, name).Do()
	if err != nil {
		return fmt.Errorf("Error deleting InstanceGroup: %s", err)
	}

	err = computeOperationWait(config.clientCompute, op, project, "Deleting InstanceGroup")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceComputeInstanceGroupImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) == 2 {
		d.Set("zone", parts[0])
		d.Set("name", parts[1])
	} else if len(parts) == 3 {
		d.Set("project", parts[0])
		d.Set("zone", parts[1])
		d.Set("name", parts[2])
		d.SetId(parts[1] + "/" + parts[2])
	} else {
		return nil, fmt.Errorf("Invalid compute instance group specifier. Expecting {zone}/{name} or {project}/{zone}/{name}")
	}

	return []*schema.ResourceData{d}, nil
}
