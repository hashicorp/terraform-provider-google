package google

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/pathorcontents"
	"github.com/hashicorp/terraform/httpclient"
	"github.com/terraform-providers/terraform-provider-google/version"

	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"
	appengine "google.golang.org/api/appengine/v1"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/cloudiot/v1"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
	composer "google.golang.org/api/composer/v1beta1"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	containerBeta "google.golang.org/api/container/v1beta1"
	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/dns/v1"
	dnsBeta "google.golang.org/api/dns/v1beta2"
	file "google.golang.org/api/file/v1beta1"
	"google.golang.org/api/iam/v1"
	cloudlogging "google.golang.org/api/logging/v2"
	"google.golang.org/api/pubsub/v1"
	runtimeconfig "google.golang.org/api/runtimeconfig/v1beta1"
	"google.golang.org/api/servicemanagement/v1"
	"google.golang.org/api/serviceusage/v1"
	"google.golang.org/api/sourcerepo/v1"
	"google.golang.org/api/spanner/v1"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
	"google.golang.org/api/storage/v1"
	"google.golang.org/api/storagetransfer/v1"
)

// Config is the configuration structure used to instantiate the Google
// provider.
type Config struct {
	Credentials string
	AccessToken string
	Project     string
	Region      string
	Zone        string
	Scopes      []string

	client    *http.Client
	userAgent string

	tokenSource oauth2.TokenSource

	clientBilling                *cloudbilling.APIService
	clientBuild                  *cloudbuild.Service
	clientComposer               *composer.Service
	clientCompute                *compute.Service
	clientComputeBeta            *computeBeta.Service
	clientContainer              *container.Service
	clientContainerBeta          *containerBeta.Service
	clientDataproc               *dataproc.Service
	clientDataflow               *dataflow.Service
	clientDns                    *dns.Service
	clientDnsBeta                *dnsBeta.Service
	clientFilestore              *file.Service
	clientKms                    *cloudkms.Service
	clientLogging                *cloudlogging.Service
	clientPubsub                 *pubsub.Service
	clientResourceManager        *cloudresourcemanager.Service
	clientResourceManagerV2Beta1 *resourceManagerV2Beta1.Service
	clientRuntimeconfig          *runtimeconfig.Service
	clientSpanner                *spanner.Service
	clientSourceRepo             *sourcerepo.Service
	clientStorage                *storage.Service
	clientSqlAdmin               *sqladmin.Service
	clientIAM                    *iam.Service
	clientServiceMan             *servicemanagement.APIService
	clientServiceUsage           *serviceusage.Service
	clientBigQuery               *bigquery.Service
	clientCloudFunctions         *cloudfunctions.Service
	clientCloudIoT               *cloudiot.Service
	clientAppEngine              *appengine.APIService
	clientStorageTransfer        *storagetransfer.Service

	bigtableClientFactory *BigtableClientFactory
}

var defaultClientScopes = []string{
	"https://www.googleapis.com/auth/compute",
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/ndev.clouddns.readwrite",
	"https://www.googleapis.com/auth/devstorage.full_control",
}

func (c *Config) LoadAndValidate() error {
	if len(c.Scopes) == 0 {
		c.Scopes = defaultClientScopes
	}

	tokenSource, err := c.getTokenSource(c.Scopes)
	if err != nil {
		return err
	}
	c.tokenSource = tokenSource

	client := oauth2.NewClient(context.Background(), tokenSource)
	client.Transport = logging.NewTransport("Google", client.Transport)
	// Each individual request should return within 30s - timeouts will be retried.
	// This is a timeout for, e.g. a single GET request of an operation - not a
	// timeout for the maximum amount of time a logical request can take.
	client.Timeout, _ = time.ParseDuration("30s")

	terraformVersion := httpclient.UserAgentString()
	providerVersion := fmt.Sprintf("terraform-provider-google/%s", version.ProviderVersion)
	terraformWebsite := "(+https://www.terraform.io)"
	userAgent := fmt.Sprintf("%s %s %s", terraformVersion, terraformWebsite, providerVersion)

	c.client = client
	c.userAgent = userAgent

	log.Printf("[INFO] Instantiating GCE client...")
	c.clientCompute, err = compute.New(client)
	if err != nil {
		return err
	}
	c.clientCompute.UserAgent = userAgent

	log.Printf("[INFO] Instantiating GCE Beta client...")
	c.clientComputeBeta, err = computeBeta.New(client)
	if err != nil {
		return err
	}
	c.clientComputeBeta.UserAgent = userAgent

	log.Printf("[INFO] Instantiating GKE client...")
	c.clientContainer, err = container.New(client)
	if err != nil {
		return err
	}
	c.clientContainer.UserAgent = userAgent

	log.Printf("[INFO] Instantiating GKE Beta client...")
	c.clientContainerBeta, err = containerBeta.New(client)
	if err != nil {
		return err
	}
	c.clientContainerBeta.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud DNS client...")
	c.clientDns, err = dns.New(client)
	if err != nil {
		return err
	}
	c.clientDns.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud DNS Beta client...")
	c.clientDnsBeta, err = dnsBeta.New(client)
	if err != nil {
		return err
	}
	c.clientDnsBeta.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud KMS Client...")
	c.clientKms, err = cloudkms.New(client)
	if err != nil {
		return err
	}
	c.clientKms.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Stackdriver Logging client...")
	c.clientLogging, err = cloudlogging.New(client)
	if err != nil {
		return err
	}
	c.clientLogging.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Storage Client...")
	c.clientStorage, err = storage.New(client)
	if err != nil {
		return err
	}
	c.clientStorage.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google SqlAdmin Client...")
	c.clientSqlAdmin, err = sqladmin.New(client)
	if err != nil {
		return err
	}
	c.clientSqlAdmin.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Pubsub Client...")
	c.clientPubsub, err = pubsub.New(client)
	if err != nil {
		return err
	}
	c.clientPubsub.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Dataflow Client...")
	c.clientDataflow, err = dataflow.New(client)
	if err != nil {
		return err
	}
	c.clientDataflow.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud ResourceManager Client...")
	c.clientResourceManager, err = cloudresourcemanager.New(client)
	if err != nil {
		return err
	}
	c.clientResourceManager.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud ResourceManager V Client...")
	c.clientResourceManagerV2Beta1, err = resourceManagerV2Beta1.New(client)
	if err != nil {
		return err
	}
	c.clientResourceManagerV2Beta1.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud Runtimeconfig Client...")
	c.clientRuntimeconfig, err = runtimeconfig.New(client)
	if err != nil {
		return err
	}
	c.clientRuntimeconfig.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud IAM Client...")
	c.clientIAM, err = iam.New(client)
	if err != nil {
		return err
	}
	c.clientIAM.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud Service Management Client...")
	c.clientServiceMan, err = servicemanagement.New(client)
	if err != nil {
		return err
	}
	c.clientServiceMan.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud Service Usage Client...")
	c.clientServiceUsage, err = serviceusage.New(client)
	if err != nil {
		return err
	}
	c.clientServiceUsage.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud Billing Client...")
	c.clientBilling, err = cloudbilling.New(client)
	if err != nil {
		return err
	}
	c.clientBilling.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud Build Client...")
	c.clientBuild, err = cloudbuild.New(client)
	if err != nil {
		return err
	}
	c.clientBuild.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud BigQuery Client...")
	c.clientBigQuery, err = bigquery.New(client)
	if err != nil {
		return err
	}
	c.clientBigQuery.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud CloudFunctions Client...")
	c.clientCloudFunctions, err = cloudfunctions.New(client)
	if err != nil {
		return err
	}
	c.clientCloudFunctions.UserAgent = userAgent

	c.bigtableClientFactory = &BigtableClientFactory{
		UserAgent:   userAgent,
		TokenSource: tokenSource,
	}

	log.Printf("[INFO] Instantiating Google Cloud Source Repo Client...")
	c.clientSourceRepo, err = sourcerepo.New(client)
	if err != nil {
		return err
	}
	c.clientSourceRepo.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud Spanner Client...")
	c.clientSpanner, err = spanner.New(client)
	if err != nil {
		return err
	}
	c.clientSpanner.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud Dataproc Client...")
	c.clientDataproc, err = dataproc.New(client)
	if err != nil {
		return err
	}
	c.clientDataproc.UserAgent = userAgent

	c.clientFilestore, err = file.New(client)
	if err != nil {
		return err
	}
	c.clientFilestore.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud IoT Core Client...")
	c.clientCloudIoT, err = cloudiot.New(client)
	if err != nil {
		return err
	}
	c.clientCloudIoT.UserAgent = userAgent

	log.Printf("[INFO] Instantiating App Engine Client...")
	c.clientAppEngine, err = appengine.New(client)
	if err != nil {
		return err
	}
	c.clientAppEngine.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Cloud Composer Client...")
	c.clientComposer, err = composer.New(client)
	if err != nil {
		return err
	}
	c.clientComposer.UserAgent = userAgent

	log.Printf("[INFO] Instantiating Google Cloud Storage Transfer Client...")
	c.clientStorageTransfer, err = storagetransfer.New(client)
	if err != nil {
		return err
	}
	c.clientStorageTransfer.UserAgent = userAgent

	return nil
}

func (c *Config) getTokenSource(clientScopes []string) (oauth2.TokenSource, error) {
	if c.AccessToken != "" {
		contents, _, err := pathorcontents.Read(c.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("Error loading access token: %s", err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'access_token'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		token := &oauth2.Token{AccessToken: contents}
		return oauth2.StaticTokenSource(token), nil
	}

	if c.Credentials != "" {
		contents, _, err := pathorcontents.Read(c.Credentials)
		if err != nil {
			return nil, fmt.Errorf("Error loading credentials: %s", err)
		}

		creds, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(contents), clientScopes...)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse credentials from '%s': %s", contents, err)
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'credentials'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return creds.TokenSource, nil
	}

	log.Printf("[INFO] Authenticating using DefaultClient...")
	log.Printf("[INFO]   -- Scopes: %s", clientScopes)
	return googleoauth.DefaultTokenSource(context.Background(), clientScopes...)
}
