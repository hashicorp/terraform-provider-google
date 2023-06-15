// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/services/billing"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudbilling/v1"
)

func ResourceBillingSubaccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceBillingSubaccountCreate,
		Read:   resourceBillingSubaccountRead,
		Delete: resourceBillingSubaccountDelete,
		Update: resourceBillingSubaccountUpdate,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"master_billing_account": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
			},
			"deletion_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: validation.StringInSlice([]string{"RENAME_ON_DESTROY", ""}, false),
			},
			"billing_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"open": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceBillingSubaccountCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	masterBillingAccount := d.Get("master_billing_account").(string)

	billingAccount := &cloudbilling.BillingAccount{
		DisplayName:          displayName,
		MasterBillingAccount: billing.CanonicalBillingAccountName(masterBillingAccount),
	}

	res, err := config.NewBillingClient(userAgent).BillingAccounts.Create(billingAccount).Do()
	if err != nil {
		return fmt.Errorf("Error creating billing subaccount '%s' in master account '%s': %s", displayName, masterBillingAccount, err)
	}

	d.SetId(res.Name)

	return resourceBillingSubaccountRead(d, meta)
}

func resourceBillingSubaccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	id := d.Id()

	billingAccount, err := config.NewBillingClient(userAgent).BillingAccounts.Get(d.Id()).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Billing Subaccount Not Found : %s", id))
	}

	if err := d.Set("name", billingAccount.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("display_name", billingAccount.DisplayName); err != nil {
		return fmt.Errorf("Error setting display_na,e: %s", err)
	}
	if err := d.Set("open", billingAccount.Open); err != nil {
		return fmt.Errorf("Error setting open: %s", err)
	}
	if err := d.Set("master_billing_account", billingAccount.MasterBillingAccount); err != nil {
		return fmt.Errorf("Error setting master_billing_account: %s", err)
	}
	if err := d.Set("billing_account_id", strings.TrimPrefix(d.Get("name").(string), "billingAccounts/")); err != nil {
		return fmt.Errorf("Error setting billing_account_id: %s", err)
	}

	return nil
}

func resourceBillingSubaccountUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	if ok := d.HasChange("display_name"); ok {
		billingAccount := &cloudbilling.BillingAccount{
			DisplayName: d.Get("display_name").(string),
		}
		_, err := config.NewBillingClient(userAgent).BillingAccounts.Patch(d.Id(), billingAccount).UpdateMask("display_name").Do()
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Error updating billing account : %s", d.Id()))
		}
	}
	return resourceBillingSubaccountRead(d, meta)
}

func resourceBillingSubaccountDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	deletionPolicy := d.Get("deletion_policy").(string)

	if deletionPolicy == "RENAME_ON_DESTROY" {
		t := time.Now()
		billingAccount := &cloudbilling.BillingAccount{
			DisplayName: "Terraform Destroyed " + t.Format("20060102150405"),
		}
		_, err := config.NewBillingClient(userAgent).BillingAccounts.Patch(d.Id(), billingAccount).UpdateMask("display_name").Do()
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Error updating billing account : %s", d.Id()))
		}
	}

	d.SetId("")

	return nil
}
