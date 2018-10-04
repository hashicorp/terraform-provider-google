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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"target": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkRelativePaths,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"ip_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"labels": &schema.Schema{
				Deprecated: "This field is in beta and will be removed from this provider. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
				Type:       schema.TypeMap,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Set:        schema.HashString,
			},

			"label_fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"port_range": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: portRangeDiffSuppress,
			},

			"ip_version": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPV4", "IPV6"}, false),
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Removed:  "Please remove this attribute (it was never used)",
			},

			"self_link": &schema.Schema{
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

	// If we have labels to set, try to set those too
	if _, ok := d.GetOk("labels"); ok {
		labels := expandLabels(d)
		// Do a read to get the fingerprint value so we can update
		fingerprint, err := resourceComputeGlobalForwardingRuleReadLabelFingerprint(config, project, frule.Name)
		if err != nil {
			return err
		}

		err = resourceComputeGlobalForwardingRuleSetLabels(config, project, frule.Name, labels, fingerprint)
		if err != nil {
			return err
		}
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
	if d.HasChange("labels") {
		labels := expandLabels(d)
		fingerprint := d.Get("label_fingerprint").(string)

		err = resourceComputeGlobalForwardingRuleSetLabels(config, project, d.Get("name").(string), labels, fingerprint)
		if err != nil {
			return err
		}

		d.SetPartial("labels")
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
	d.Set("labels", frule.Labels)
	d.Set("label_fingerprint", frule.LabelFingerprint)
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

// resourceComputeGlobalForwardingRuleReadLabelFingerprint performs a read on the remote resource and returns only the
// fingerprint. Used on create when setting labels as we don't know the label fingerprint initially.
func resourceComputeGlobalForwardingRuleReadLabelFingerprint(config *Config, project, name string) (string, error) {
	frule, err := config.clientComputeBeta.GlobalForwardingRules.Get(project, name).Do()
	if err != nil {
		return "", fmt.Errorf("Unable to read global forwarding rule to update labels: %s", err)
	}

	return frule.LabelFingerprint, nil
}

// resourceComputeGlobalForwardingRuleSetLabels sets the Labels attribute on a forwarding rule.
func resourceComputeGlobalForwardingRuleSetLabels(config *Config, project, name string, labels map[string]string, fingerprint string) error {
	setLabels := computeBeta.GlobalSetLabelsRequest{
		Labels:           labels,
		LabelFingerprint: fingerprint,
	}
	op, err := config.clientComputeBeta.GlobalForwardingRules.SetLabels(project, name, &setLabels).Do()
	if err != nil {
		return err
	}

	err = computeSharedOperationWait(config.clientCompute, op, project, "Setting labels on Global Forwarding Rule")
	if err != nil {
		return err
	}

	return nil
}
