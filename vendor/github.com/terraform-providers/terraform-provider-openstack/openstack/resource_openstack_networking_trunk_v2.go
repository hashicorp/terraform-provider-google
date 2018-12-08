package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunks"
)

func resourceNetworkingTrunkV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingTrunkV2Create,
		Read:   resourceNetworkingTrunkV2Read,
		Update: resourceNetworkingTrunkV2Update,
		Delete: resourceNetworkingTrunkV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"port_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"admin_state_up": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"sub_port": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"segmentation_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"segmentation_id": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"tags": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNetworkingTrunkV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := trunks.CreateOpts{
		Name:         d.Get("name").(string),
		AdminStateUp: resourceNetworkingTrunkV2AdminStateUp(d),
		PortID:       d.Get("port_id").(string),
		TenantID:     d.Get("tenant_id").(string),
		Subports:     resourceNetworkingTrunkV2Subports(d.Get("sub_port").(*schema.Set)),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	trunk, err := trunks.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenStack Neutron trunk: %s", err)
	}

	log.Printf("[DEBUG] Waiting for OpenStack Neutron trunk (%s) to become available.", trunk.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForNetworkTrunkActive(client, trunk.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for trunk to become active: %s", err)
	}

	tags := networkV2AttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(client, "trunks", trunk.ID, tagOpts).Extract()
		if err != nil {
			return fmt.Errorf("Unable to set trunk tags: %v", err)
		}
		log.Printf("[DEBUG] Set Tags = %+v on trunk %s", tags, trunk.ID)
	}

	d.SetId(trunk.ID)
	return resourceNetworkingTrunkV2Read(d, meta)
}

func resourceNetworkingTrunkV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	trunk, err := trunks.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "trunk")
	}
	log.Printf("[DEBUG] Retrieved trunk %s: %+v", d.Id(), trunk)

	d.Set("region", GetRegion(d, config))
	d.Set("name", trunk.Name)
	d.Set("port_id", trunk.PortID)
	d.Set("admin_state_up", trunk.AdminStateUp)
	d.Set("tenant_id", trunk.TenantID)
	d.Set("tags", trunk.Tags)

	subports := make([]map[string]interface{}, len(trunk.Subports))
	for i, trunkSubport := range trunk.Subports {
		subports[i] = make(map[string]interface{})
		subports[i]["port_id"] = trunkSubport.PortID
		subports[i]["segmentation_type"] = trunkSubport.SegmentationType
		subports[i]["segmentation_id"] = trunkSubport.SegmentationID
	}
	if err = d.Set("sub_port", subports); err != nil {
		return fmt.Errorf("Unable to set sub_port for trunk %s: %s", d.Id(), err)
	}

	return nil
}

func resourceNetworkingTrunkV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var updateTrunk bool
	var updateOpts trunks.UpdateOpts

	if d.HasChange("name") {
		updateTrunk = true
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("admin_state_up") {
		updateTrunk = true
		updateOpts.AdminStateUp = resourceNetworkingTrunkV2AdminStateUp(d)
	}

	if updateTrunk {
		log.Printf("[DEBUG] Updating trunk %s with options: %+v", d.Id(), updateOpts)
		_, err = trunks.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenStack Neutron trunk: %s", err)
		}
	}

	if d.HasChange("sub_port") {
		old, new := d.GetChange("sub_port")

		toRemove := old.(*schema.Set).Difference(new.(*schema.Set))
		toAdd := new.(*schema.Set).Difference(old.(*schema.Set))

		if len(toRemove.List()) > 0 {
			var subportsToRemove trunks.RemoveSubportsOpts

			subports := resourceNetworkingTrunkV2Subports(toRemove)
			log.Printf("[DEBUG] Removing subports from trunk %s: %#v", d.Id(), subports)
			for _, subport := range subports {
				v := trunks.RemoveSubport{PortID: subport.PortID}
				subportsToRemove.Subports = append(subportsToRemove.Subports, v)
			}
			_, err := trunks.RemoveSubports(client, d.Id(), subportsToRemove).Extract()
			if err != nil {
				return fmt.Errorf("Error removing subports when updating OpenStack Neutron trunk: %s", err)
			}
		}

		if len(toAdd.List()) > 0 {
			var subportsToAdd trunks.AddSubportsOpts
			subportsToAdd.Subports = resourceNetworkingTrunkV2Subports(toAdd)

			log.Printf("[DEBUG] Adding subports to trunk %s: %#v", d.Id(), subportsToAdd.Subports)

			_, err := trunks.AddSubports(client, d.Id(), subportsToAdd).Extract()
			if err != nil {
				return fmt.Errorf("Error adding subports when updating OpenStack Neutron trunk: %s", err)
			}
		}
	}

	if d.HasChange("tags") {
		tags := networkV2AttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(client, "trunks", d.Id(), tagOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating tags on trunk: %s", err)
		}
		log.Printf("[DEBUG] Updated tags = %+v on trunk %s", tags, d.Id())
	}

	return resourceNetworkingTrunkV2Read(d, meta)
}

func resourceNetworkingTrunkV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForNetworkTrunkDelete(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting OpenStack Neutron trunk: %s", err)
	}

	return nil
}

func resourceNetworkingTrunkV2Subports(d *schema.Set) (subports []trunks.Subport) {
	rawSubports := d.List()
	subports = make([]trunks.Subport, len(rawSubports))

	for i, raw := range rawSubports {
		rawMap := raw.(map[string]interface{})
		subports[i] = trunks.Subport{
			PortID:           rawMap["port_id"].(string),
			SegmentationType: rawMap["segmentation_type"].(string),
			SegmentationID:   rawMap["segmentation_id"].(int),
		}
	}
	return
}

func resourceNetworkingTrunkV2AdminStateUp(d *schema.ResourceData) *bool {
	value := false

	if v, ok := d.GetOkExists("admin_state_up"); ok {
		value = v.(bool)
	}

	return &value
}

func waitForNetworkTrunkActive(client *gophercloud.ServiceClient, trunkID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		trunk, err := trunks.Get(client, trunkID).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenStack Neutron trunk: %+v", trunk)
		if trunk.Status == "DOWN" || trunk.Status == "ACTIVE" {
			return trunk, "ACTIVE", nil
		}

		return trunk, trunk.Status, nil
	}
}

func waitForNetworkTrunkDelete(client *gophercloud.ServiceClient, trunkID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenStack Neutron trunk %s", trunkID)

		trunk, err := trunks.Get(client, trunkID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenStack trunk %s", trunkID)
				return trunk, "DELETED", nil
			}
			return trunk, "ACTIVE", err
		}

		err = trunks.Delete(client, trunkID).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenStack trunk %s", trunkID)
				return trunk, "DELETED", nil
			}
			return trunk, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenStack trunk %s still active.\n", trunkID)
		return trunk, "ACTIVE", nil
	}
}
