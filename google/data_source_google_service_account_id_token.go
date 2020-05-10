package google

import (
	"log"
	"time"

	"fmt"
	"reflect"
	"strings"

	iamcredentials "google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/idtoken"

	"github.com/hashicorp/terraform-plugin-sdk/helper/pathorcontents"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	metadataIdentityDocURL = "http://metadata/computeMetadata/v1/instance/service-accounts/default/identity"
	userInfoScope          = "https://www.googleapis.com/auth/userinfo.email"
	tokenEndpoint          = "https://www.googleapis.com/oauth2/v4/token"
)

func dataSourceGoogleServiceAccountIdToken() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceGoogleServiceAccountIdTokenRead,
		Schema: map[string]*schema.Schema{
			"target_audience": {
				Type:     schema.TypeString,
				Required: true,
				//ValidateFunc: validateRegexp("(" + strings.Join(PossibleServiceAccountNames, "|") + ")"),
			},
			"access_token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
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
			"id_token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func getCredentials(c *Config, clientScopes []string) (google.Credentials, error) {
	if c.AccessToken != "" {
		contents, _, err := pathorcontents.Read(c.AccessToken)
		if err != nil {
			return google.Credentials{}, fmt.Errorf("Error loading access token: %s", err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'access_token'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)

		token := &oauth2.Token{AccessToken: contents}

		return google.Credentials{
			TokenSource: oauth2.StaticTokenSource(token),
		}, nil
	}

	if c.Credentials != "" {
		contents, _, err := pathorcontents.Read(c.Credentials)
		if err != nil {
			return google.Credentials{}, fmt.Errorf("Error loading credentials: %s", err)
		}

		creds, err := google.CredentialsFromJSON(c.context, []byte(contents), clientScopes...)
		if err != nil {
			return google.Credentials{}, fmt.Errorf("Unable to parse credentials from '%s': %s", contents, err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'credentials'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return *creds, nil
	}

	log.Printf("[INFO] Authenticating using DefaultClient...")
	log.Printf("[INFO]   -- Scopes: %s", clientScopes)

	creds, err := google.FindDefaultCredentials(c.context, clientScopes...)
	if err != nil {
		return google.Credentials{}, fmt.Errorf("Unable FindDefaultCredentials '%s'", err)
	}
	return *creds, nil

}

func dataSourceGoogleServiceAccountIdTokenRead(d *schema.ResourceData, meta interface{}) error {
	var idToken string

	config := meta.(*Config)
	targetAudience := d.Get("target_audience").(string)
	log.Printf("[INFO] Acquire  IdToken for Audience %s", targetAudience)

	var ts oauth2.TokenSource
	var creds google.Credentials
	var err error

	ts, err = config.getTokenSource([]string{userInfoScope})
	if err != nil {
		return fmt.Errorf("Unable to acquire TokenSource from credentials: %v", err)
	}

	// ok, if a token was provided either directly as access_token parameter
	// or inderectly as an impersonated token provider
	// either way, generate a credential object for this
	accessToken := d.Get("access_token").(string)
	if accessToken != "" {
		token := &oauth2.Token{AccessToken: accessToken}
		ts = oauth2.StaticTokenSource(token)
		creds = google.Credentials{
			TokenSource: ts,
		}
	} else if reflect.TypeOf(ts) != reflect.TypeOf(oauth2.StaticTokenSource) {
		creds = google.Credentials{
			TokenSource: ts,
		}
	}

	tok, err := ts.Token()
	if err != nil {
		return fmt.Errorf("Unable to get Token() from tokenSource: %v", err)
	}

	// TODO: improve all this
	// ok, if we've got a statictoken, all we can do is use the iamcredentials api to get an idtoken
	if fmt.Sprintf("%s", reflect.TypeOf(ts)) == "oauth2.staticTokenSource" {
		// all we have here is an access_token so we have to use
		// https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/generateIdToken
		service := config.clientIamCredentials
		if err != nil {
			return fmt.Errorf("data_source_google_service_account_id_token: Error creating IAMCredentials: %v", err)
		}
		name := fmt.Sprintf("projects/-/serviceAccounts/%s", d.Get("target_service_account").(string))
		tokenRequest := &iamcredentials.GenerateIdTokenRequest{
			Audience:     targetAudience,
			IncludeEmail: d.Get("include_email").(bool),
			Delegates:    convertStringSet(d.Get("delegates").(*schema.Set)),
		}
		at, err := service.Projects.ServiceAccounts.GenerateIdToken(name, tokenRequest).Do()
		if err != nil {
			return fmt.Errorf("data_source_google_service_account_id_token:: Error calling iamcredentials.GenerateIdToken: %v", err)
		}
		idToken = at.Token

		d.SetId(time.Now().UTC().String())
		d.Set("id_token", idToken)

		return nil
	}

	// otherwise, it could be either a service account JSON or compute engine metadata
	creds, err = getCredentials(config, []string{userInfoScope})
	if err != nil {
		return fmt.Errorf("data_source_google_service_account_id_token: Error calling getCredentials(): %v", err)
	}

	if creds.JSON != nil {
		ctx := context.Background()
		ts, err := idtoken.NewTokenSource(ctx, targetAudience, idtoken.WithCredentialsJSON(creds.JSON))
		if err != nil {
			return fmt.Errorf("data_source_google_service_account Unable to init idTokenSource%v", err)
		}
		tok, err := ts.Token()
		if err != nil {
			return fmt.Errorf("unable to retrieve Token: %v", err)
		}
		idToken = tok.AccessToken

	} else if tok.RefreshToken == "" {
		// if the token isn't a json cert, it should be a ReuseTokenSource from MetadataServer
		ctx := context.Background()
		ts, err := idtoken.NewTokenSource(ctx, targetAudience)
		if err != nil {
			return fmt.Errorf("data_source_google_service_account Unable to init idTokenSource%v", err)
		}
		tok, err := ts.Token()
		if err != nil {
			return fmt.Errorf("unable to retrieve Token: %v", err)
		}
		idToken = tok.AccessToken
	} else {
		// bail, this shoudn't happen
		return fmt.Errorf("data_source_google_service_account_id_token: Unsupported Credential Type supplied: got %v", reflect.TypeOf(creds.TokenSource))
	}

	d.SetId(time.Now().UTC().String())
	d.Set("id_token", idToken)

	return nil
}
