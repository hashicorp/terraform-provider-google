package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

var GlobalForwardingRuleBaseApiVersion = v1
var GlobalForwardingRuleVersionedFeatures = []Feature{
	{Version: v0beta, Item: "labels"},
}

func resourceComputeGlobalForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeGlobalForwardingRuleCreate,
		Read:   resourceComputeGlobalForwardingRuleRead,
		Update: resourceComputeGlobalForwardingRuleUpdate,
		Delete: resourceComputeGlobalForwardingRuleDelete,

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
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"label_fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"port_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
	computeApiVersion := getComputeApiVersion(d, GlobalForwardingRuleBaseApiVersion, GlobalForwardingRuleVersionedFeatures)
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

	var op interface{}
	switch computeApiVersion {
	case v1:
		v1Frule := &compute.ForwardingRule{}
		err = Convert(frule, v1Frule)
		if err != nil {
			return err
		}

		op, err = config.clientCompute.GlobalForwardingRules.Insert(project, v1Frule).Do()
		if err != nil {
			return fmt.Errorf("Error creating Global Forwarding Rule: %s", err)
		}
	case v0beta:
		v0BetaFrule := &computeBeta.ForwardingRule{}
		err = Convert(frule, v0BetaFrule)
		if err != nil {
			return err
		}

		op, err = config.clientComputeBeta.GlobalForwardingRules.Insert(project, v0BetaFrule).Do()
		if err != nil {
			return fmt.Errorf("Error creating Global Forwarding Rule: %s", err)
		}
	}

	// It probably maybe worked, so store the ID now
	d.SetId(frule.Name)

	err = computeSharedOperationWait(config, op, project, "Creating Global Fowarding Rule")
	if err != nil {
		return err
	}

	// If we have labels to set, try to set those too
	if _, ok := d.GetOk("labels"); ok {
		labels := expandLabels(d)
		// Do a read to get the fingerprint value so we can update
		fingerprint, err := resourceComputeGlobalForwardingRuleReadLabelFingerprint(config, computeApiVersion, project, frule.Name)
		if err != nil {
			return err
		}

		err = resourceComputeGlobalForwardingRuleSetLabels(config, computeApiVersion, project, frule.Name, labels, fingerprint)
		if err != nil {
			return err
		}
	}

	return resourceComputeGlobalForwardingRuleRead(d, meta)
}

func resourceComputeGlobalForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersionUpdate(d, GlobalForwardingRuleBaseApiVersion, GlobalForwardingRuleVersionedFeatures, []Feature{})
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	if d.HasChange("target") {
		target := d.Get("target").(string)
		targetRef := &computeBeta.TargetReference{Target: target}

		var op interface{}
		switch computeApiVersion {
		case v1:
			v1TargetRef := &compute.TargetReference{}
			err = Convert(targetRef, v1TargetRef)
			if err != nil {
				return err
			}

			op, err = config.clientCompute.GlobalForwardingRules.SetTarget(
				project, d.Id(), v1TargetRef).Do()
			if err != nil {
				return fmt.Errorf("Error updating target: %s", err)
			}
		case v0beta:
			v0BetaTargetRef := &compute.TargetReference{}
			err = Convert(targetRef, v0BetaTargetRef)
			if err != nil {
				return err
			}

			op, err = config.clientCompute.GlobalForwardingRules.SetTarget(
				project, d.Id(), v0BetaTargetRef).Do()
			if err != nil {
				return fmt.Errorf("Error updating target: %s", err)
			}
		}

		err = computeSharedOperationWait(config, op, project, "Updating Global Forwarding Rule")
		if err != nil {
			return err
		}

		d.SetPartial("target")
	}
	if d.HasChange("labels") {
		labels := expandLabels(d)
		fingerprint := d.Get("label_fingerprint").(string)

		err = resourceComputeGlobalForwardingRuleSetLabels(config, computeApiVersion, project, d.Get("name").(string), labels, fingerprint)
		if err != nil {
			return err
		}

		d.SetPartial("labels")
	}

	d.Partial(false)

	return resourceComputeGlobalForwardingRuleRead(d, meta)
}

func resourceComputeGlobalForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, GlobalForwardingRuleBaseApiVersion, GlobalForwardingRuleVersionedFeatures)
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	frule := &computeBeta.ForwardingRule{}
	switch computeApiVersion {
	case v1:
		v1Frule, err := config.clientCompute.GlobalForwardingRules.Get(project, d.Id()).Do()
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Global Forwarding Rule %q", d.Get("name").(string)))
		}

		err = Convert(v1Frule, frule)
		if err != nil {
			return err
		}
	case v0beta:
		v0BetaFrule, err := config.clientComputeBeta.GlobalForwardingRules.Get(project, d.Id()).Do()
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("Global Forwarding Rule %q", d.Get("name").(string)))
		}

		err = Convert(v0BetaFrule, frule)
		if err != nil {
			return err
		}
	}

	d.Set("ip_address", frule.IPAddress)
	d.Set("ip_protocol", frule.IPProtocol)
	d.Set("ip_version", frule.IpVersion)
	d.Set("self_link", ConvertSelfLinkToV1(frule.SelfLink))
	d.Set("labels", frule.Labels)
	d.Set("label_fingerprint", frule.LabelFingerprint)

	return nil
}

func resourceComputeGlobalForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	computeApiVersion := getComputeApiVersion(d, GlobalForwardingRuleBaseApiVersion, GlobalForwardingRuleVersionedFeatures)
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the GlobalForwardingRule
	log.Printf("[DEBUG] GlobalForwardingRule delete request")
	var op interface{}
	switch computeApiVersion {
	case v1:
		op, err = config.clientCompute.GlobalForwardingRules.Delete(project, d.Id()).Do()
		if err != nil {
			return fmt.Errorf("Error deleting GlobalForwardingRule: %s", err)
		}
	case v0beta:
		op, err = config.clientComputeBeta.GlobalForwardingRules.Delete(project, d.Id()).Do()
		if err != nil {
			return fmt.Errorf("Error deleting GlobalForwardingRule: %s", err)
		}
	}

	err = computeSharedOperationWait(config, op, project, "Deleting GlobalForwarding Rule")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

// resourceComputeGlobalForwardingRuleReadLabelFingerprint performs a read on the remote resource and returns only the
// fingerprint. Used on create when setting labels as we don't know the label fingerprint initially.
func resourceComputeGlobalForwardingRuleReadLabelFingerprint(config *Config, computeApiVersion ComputeApiVersion,
	project, name string) (string, error) {
	switch computeApiVersion {
	case v0beta:
		frule, err := config.clientComputeBeta.GlobalForwardingRules.Get(project, name).Do()
		if err != nil {
			return "", fmt.Errorf("Unable to read global forwarding rule to update labels: %s", err)
		}

		return frule.LabelFingerprint, nil
	default:
		return "", fmt.Errorf(
			"Unable to read label fingerprint due to an internal error: can only handle v0beta but compute api logic indicates %d",
			computeApiVersion)
	}
}

// resourceComputeGlobalForwardingRuleSetLabels sets the Labels attribute on a forwarding rule.
func resourceComputeGlobalForwardingRuleSetLabels(config *Config, computeApiVersion ComputeApiVersion, project,
	name string, labels map[string]string, fingerprint string) error {
	var op interface{}
	var err error

	switch computeApiVersion {
	case v0beta:
		setLabels := computeBeta.GlobalSetLabelsRequest{
			Labels:           labels,
			LabelFingerprint: fingerprint,
		}
		op, err = config.clientComputeBeta.GlobalForwardingRules.SetLabels(project, name, &setLabels).Do()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf(
			"Unable to set labels due to an internal error: can only handle v0beta but compute api logic indicates %d",
			computeApiVersion)
	}

	err = computeSharedOperationWait(config, op, project, "Setting labels on Global Forwarding Rule")
	if err != nil {
		return err
	}

	return nil
}
