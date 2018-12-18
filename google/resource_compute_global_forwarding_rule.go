package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	computeBeta "google.golang.org/api/compute/v0.beta"
	compute "google.golang.org/api/compute/v1"
)

func resourceComputeGlobalForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeGlobalForwardingRuleCreate,
		Read:   resourceComputeGlobalForwardingRuleRead,
		Update: resourceComputeGlobalForwardingRuleUpdate,
		Delete: resourceComputeGlobalForwardingRuleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"target": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkRelativePaths,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"ip_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"labels": {
				Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"label_fingerprint": {
				Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
				Type:     schema.TypeString,
				Computed: true,
			},

			"port_range": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: portRangeDiffSuppress,
			},

			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV6"}, false),
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Removed:  "Please remove this attribute (it was never used)",
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeGlobalForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	frule := &computeBeta.ForwardingRule{
		IPAddress:   d.Get("ip_address").(string),
		IPProtocol:  d.Get("ip_protocol").(string),
		IpVersion:   d.Get("ip_version").(string),
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		PortRange:   d.Get("port_range").(string),
		Target:      d.Get("target").(string),
	}

	op, err := config.clientComputeBeta.GlobalForwardingRules.Insert(project, frule).Do()
	if err != nil {
		return fmt.Errorf("Error creating Global Forwarding Rule: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(frule.Name)

	err = computeSharedOperationWait(config.clientCompute, op, project, "Creating Global Fowarding Rule")
	if err != nil {
		return err
	}

	return resourceComputeGlobalForwardingRuleRead(d, meta)
}

func resourceComputeGlobalForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	if d.HasChange("target") {
		target := d.Get("target").(string)
		targetRef := &compute.TargetReference{Target: target}

		op, err := config.clientCompute.GlobalForwardingRules.SetTarget(
			project, d.Id(), targetRef).Do()
		if err != nil {
			return fmt.Errorf("Error updating target: %s", err)
		}

		err = computeSharedOperationWait(config.clientCompute, op, project, "Updating Global Forwarding Rule")
		if err != nil {
			return err
		}

		d.SetPartial("target")
	}

	d.Partial(false)

	return resourceComputeGlobalForwardingRuleRead(d, meta)
}

func resourceComputeGlobalForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	frule, err := config.clientComputeBeta.GlobalForwardingRules.Get(project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Global Forwarding Rule %q", d.Get("name").(string)))
	}

	d.Set("name", frule.Name)
	d.Set("description", frule.Description)
	d.Set("target", frule.Target)
	d.Set("port_range", frule.PortRange)
	d.Set("ip_address", frule.IPAddress)
	d.Set("ip_protocol", frule.IPProtocol)
	d.Set("ip_version", frule.IpVersion)
	d.Set("self_link", ConvertSelfLinkToV1(frule.SelfLink))
	// removed lists need something set
	d.Set("labels", nil)
	d.Set("project", project)

	return nil
}

func resourceComputeGlobalForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the GlobalForwardingRule
	log.Printf("[DEBUG] GlobalForwardingRule delete request")
	op, err := config.clientCompute.GlobalForwardingRules.Delete(project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting GlobalForwardingRule: %s", err)
	}
	err = computeSharedOperationWait(config.clientCompute, op, project, "Deleting GlobalForwarding Rule")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
