// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	iamcredentials "google.golang.org/api/iamcredentials/v1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleServiceAccountJwt() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleServiceAccountJwtRead,
		Schema: map[string]*schema.Schema{
			"payload": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `A JSON-encoded JWT claims set that will be included in the signed JWT.`,
			},
			"expires_in": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of seconds until the JWT expires. If set and non-zero an `exp` claim will be added to the payload derived from the current timestamp plus expires_in seconds.",
			},
			"target_service_account": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidateRegexp("(" + strings.Join(verify.PossibleServiceAccountNames, "|") + ")"),
			},
			"delegates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidateRegexp(verify.ServiceAccountLinkRegex),
				},
			},
			"jwt": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

var (
	DataSourceGoogleServiceAccountJwtNow = time.Now
)

func dataSourceGoogleServiceAccountJwtRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)

	if err != nil {
		return err
	}

	payload := d.Get("payload").(string)

	if expiresIn := d.Get("expires_in").(int); expiresIn != 0 {
		var decoded map[string]interface{}

		if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
			return fmt.Errorf("error decoding `payload` while adding `exp` field: %w", err)
		}

		decoded["exp"] = DataSourceGoogleServiceAccountJwtNow().Add(time.Duration(expiresIn) * time.Second).Unix()

		payloadBytesWithExp, err := json.Marshal(decoded)

		if err != nil {
			return fmt.Errorf("error re-encoding `payload` while adding `exp` field: %w", err)
		}

		payload = string(payloadBytesWithExp)
	}

	name := fmt.Sprintf("projects/-/serviceAccounts/%s", d.Get("target_service_account").(string))

	jwtRequest := &iamcredentials.SignJwtRequest{
		Payload:   payload,
		Delegates: tpgresource.ConvertStringSet(d.Get("delegates").(*schema.Set)),
	}

	service := config.NewIamCredentialsClient(userAgent)

	jwtResponse, err := service.Projects.ServiceAccounts.SignJwt(name, jwtRequest).Do()

	if err != nil {
		return fmt.Errorf("error calling iamcredentials.SignJwt: %w", err)
	}

	d.SetId(name)

	if err := d.Set("jwt", jwtResponse.SignedJwt); err != nil {
		return fmt.Errorf("error setting jwt attribute: %w", err)
	}

	return nil
}
