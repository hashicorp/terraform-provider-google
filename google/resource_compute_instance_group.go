package google

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"

	"github.com/hashicorp/terraform/helper/schema"
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
				Set:      schema.HashString,
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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

		CustomizeDiff: customDiffInstanceGroupInstancesField,
	}
}

func customDiffInstanceGroupInstancesField(diff *schema.ResourceDiff, meta interface{}) error {
	// This deals with an interesting problem that deserves some attention.
	// When an instance is destroyed and recreated, its membership in
	// instance groups disappears.  However, its recreated `self_link` will
	// be the same as the `self_link` from before the destroy/recreate.
	// Therefore, if some instances which are set in the `instances` field
	// are destroyed and recreated, although in *reality* there is a diff
	// between the GCP state and the desired state, Terraform cannot *see*
	// the diff without running a full Read cycle.  There's no Read implicit
	// in the `Apply` stage, so we need to trick Terraform into calling
	// Update when things like this happen.

	// This function will be called in 3 different states.
	// 1) it will be called on a new resource which hasn't been created yet.
	//		We shouldn't do anything interesting in that case.
	// 2) it will be called on an updated resource during "plan" time - that is,
	//		before anything has actually been done.  In that case, we need to show
	//		a diff on the resource if there's a chance that any of the instances
	//		will be destroyed and recreated.  Fortunately, in this case, the
	//		upstream logic will show that there is a diff.  This will be a
	//		"Computed" diff - as if we had called diff.SetComputed("instances"),
	//		and that's a good response.
	// 3) it will be called on an updated resource at "apply" time - that is,
	//		right when we're about to do something with this resource.  That
	//		is designed to check whether there really is a diff on the Computed
	//		field "instances".  Here, we have to get tricky.  We need to show
	//		a diff, and it can't be a Computed diff (`apply` skips `Update` if
	//		the only diffs are Computed).  It also can't be a ForceNew diff,
	//		because Terraform crashes if there's a ForceNew diff at apply time
	//		after not seeing one at plan time.  We're in a pickle - the Terraform
	//		state matches our desired state, but is *wrong*.  We add a fake item
	//		to the "instances" set, so that Terraform sees a diff between the
	//		state and the desired state.

	oldI, newI := diff.GetChange("instances")
	oldInstanceSet := oldI.(*schema.Set)
	newInstanceSet := newI.(*schema.Set)
	oldInstances := convertStringArr(oldInstanceSet.List())
	newInstances := convertStringArr(newInstanceSet.List())

	log.Printf("[DEBUG] InstanceGroup CustomizeDiff old: %#v, new: %#v", oldInstances, newInstances)
	var memberUrls []string
	config := meta.(*Config)
	// We can't use getProject() or getZone(), because we only have a schema.ResourceDiff,
	// not a schema.ResourceData.  We'll have to emulate them like this.
	project := diff.Get("project").(string)
	if project == "" {
		project = config.Project
	}
	zone := diff.Get("zone").(string)
	if zone == "" {
		project = config.Zone
	}

	// We need to see what instances are present in the instance group.  There are a few
	// possible results.
	// 1) The instance group doesn't exist.  We don't change the diff in this case -
	//		if the instance is being created, that's the right thing to do.
	// 2) The instance group exists, and the GCP state matches the terraform state.  In this
	//		case, we should do nothing.
	// 3) The instance group exists, and the GCP state does not match the terraform state.
	//		In this case, we add the string "FORCE_UPDATE" to the list of instances, to convince
	//		Terraform to execute an update even though there's no diff between the terraform
	//		state and the desired state.
	members, err := config.clientCompute.InstanceGroups.ListInstances(
		project, zone, diff.Get("name").(string), &compute.InstanceGroupsListInstancesRequest{
			InstanceState: "ALL",
		}).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			// This is where we'll end up at either plan time or apply time on first creation.
			return nil
		} else {
			// Any other errors return them
			return fmt.Errorf("Error reading InstanceGroup Members: %s", err)
		}
	}
	for _, member := range members.Items {
		memberUrls = append(memberUrls, member.Instance)
	}
	sort.Strings(memberUrls)
	sort.Strings(oldInstances)
	sort.Strings(newInstances)
	log.Printf("[DEBUG] InstanceGroup members: %#v.  OldInstances: %#v", memberUrls, oldInstances)
	if !reflect.DeepEqual(memberUrls, oldInstances) && reflect.DeepEqual(newInstances, oldInstances) {
		// This is where we'll end up at apply-time only if an instance is
		// somehow removed from the set of instances between refresh and update.
		newInstancesList := append(newInstances, "FORCE_UPDATE")
		diff.SetNew("instances", newInstancesList)
	}
	// This is where we'll end up if the GCP state matches the Terraform state.
	return nil
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
		if !strings.HasPrefix(v, "https://www.googleapis.com/compute/v1/") {
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
	err := resourceComputeInstanceGroupRead(d, meta)
	if err != nil {
		return err
	}
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
		_, to_ := d.GetChange("instances")
		// We need to get the current state from d directly because
		// it is likely to have been changed by the Read() above.
		from_ := d.Get("instances")
		to_.(*schema.Set).Remove("FORCE_UPDATE")

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
		// Important to fetch via GetChange, because the above Read() will
		// have reset the value retrieved via Get() to its current value.
		_, namedPorts_ := d.GetChange("named_port")
		namedPorts := getNamedPorts(namedPorts_.([]interface{}))

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
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid compute instance group specifier. Expecting {zone}/{name}")
	}

	d.Set("zone", parts[0])
	d.Set("name", parts[1])

	return []*schema.ResourceData{d}, nil
}
