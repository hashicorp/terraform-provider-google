package google

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/oauth2"
	iamcredentials "google.golang.org/api/iamcredentials/v1"
)

func dataSourceGoogleImpersonatedCredential() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceImpersonatedCredentialRead,
		Schema: map[string]*schema.Schema{
			"source_access_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_service_account": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRegexp("(" + strings.Join(PossibleServiceAccountNames, "|") + ")"),
			},
			"access_token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					StateFunc: func(v interface{}) string {
						return canonicalizeServiceScope(v.(string))
					},
				},
				// ValidateFunc is not yet supported on lists or sets.
			},
			"delegates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateRegexp(ServiceAccountLinkRegex),
				},
			},
			"lifetime": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDuration(), // duration <=3600s; TODO: support validteDuration(min,max)
				Default:      "3600s",
			},
		},
	}
}

func dataSourceImpersonatedCredentialRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	log.Printf("[INFO] Acquire Impersonated credentials for %s", d.Get("target_service_account").(string))

	d.SetId(time.Now().UTC().String())
	var client *http.Client
	if d.Get("source_access_token") != "" {
		rootTokenSource := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: d.Get("source_access_token").(string),
		})
		client = oauth2.NewClient(context.TODO(), rootTokenSource)
	} else {
		client = config.client
	}

	service, err := iamcredentials.New(client)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("projects/-/serviceAccounts/%s", d.Get("target_service_account").(string))
	tokenRequest := &iamcredentials.GenerateAccessTokenRequest{
		Lifetime:  d.Get("lifetime").(string),
		Delegates: convertStringSet(d.Get("delegates").(*schema.Set)),
		Scope:     canonicalizeServiceScopes(convertStringSet(d.Get("scopes").(*schema.Set))),
	}
	at, err := service.Projects.ServiceAccounts.GenerateAccessToken(name, tokenRequest).Do()
	if err != nil {
		return err
	}

	d.Set("access_token", at.AccessToken)

	return nil
}
