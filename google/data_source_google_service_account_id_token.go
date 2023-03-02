package google

import (
	"fmt"
	"strings"

	iamcredentials "google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/net/context"
)

const (
	userInfoScope = "https://www.googleapis.com/auth/userinfo.email"
)

func DataSourceGoogleServiceAccountIdToken() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceGoogleServiceAccountIdTokenRead,
		Schema: map[string]*schema.Schema{
			"target_audience": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_service_account": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp("(" + strings.Join(PossibleServiceAccountNames, "|") + ")"),
			},
			"delegates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateRegexp(ServiceAccountLinkRegex),
				},
			},
			"include_email": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// Not used currently
			// https://github.com/googleapis/google-api-go-client/issues/542
			// "format": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// 	ValidateFunc: validation.StringInSlice([]string{
			// 		"FULL", "STANDARD"}, true),
			// 	Default: "STANDARD",
			// },
			"id_token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func dataSourceGoogleServiceAccountIdTokenRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	targetAudience := d.Get("target_audience").(string)
	creds, err := config.GetCredentials([]string{userInfoScope}, false)
	if err != nil {
		return fmt.Errorf("error calling getCredentials(): %v", err)
	}

	// If the source credential is not a service account key, use the API to generate the idToken
	if creds.JSON == nil {
		// Use
		// https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/generateIdToken
		service := config.NewIamCredentialsClient(userAgent)
		name := fmt.Sprintf("projects/-/serviceAccounts/%s", d.Get("target_service_account").(string))
		tokenRequest := &iamcredentials.GenerateIdTokenRequest{
			Audience:     targetAudience,
			IncludeEmail: d.Get("include_email").(bool),
			Delegates:    convertStringSet(d.Get("delegates").(*schema.Set)),
		}
		at, err := service.Projects.ServiceAccounts.GenerateIdToken(name, tokenRequest).Do()
		if err != nil {
			return fmt.Errorf("error calling iamcredentials.GenerateIdToken: %v", err)
		}

		d.SetId(d.Get("target_service_account").(string))
		if err := d.Set("id_token", at.Token); err != nil {
			return fmt.Errorf("Error setting id_token: %s", err)
		}

		return nil
	}

	ctx := context.Background()
	co := []option.ClientOption{}
	if creds.JSON != nil {
		co = append(co, idtoken.WithCredentialsJSON(creds.JSON))
	}

	idTokenSource, err := idtoken.NewTokenSource(ctx, targetAudience, co...)
	if err != nil {
		return fmt.Errorf("unable to retrieve TokenSource: %v", err)
	}
	idToken, err := idTokenSource.Token()
	if err != nil {
		return fmt.Errorf("unable to retrieve Token: %v", err)
	}

	d.SetId(targetAudience)
	if err := d.Set("id_token", idToken.AccessToken); err != nil {
		return fmt.Errorf("Error setting id_token: %s", err)
	}

	return nil
}
