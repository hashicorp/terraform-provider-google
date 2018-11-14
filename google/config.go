package google

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/pathorcontents"
	"github.com/hashicorp/terraform/httpclient"
	"github.com/terraform-providers/terraform-provider-google/version"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	appengine "google.golang.org/api/appengine/v1"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/cloudiot/v1"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
	"google.golang.org/api/composer/v1"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	containerBeta "google.golang.org/api/container/v1beta1"
	"google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/dns/v1"
	dnsBeta "google.golang.org/api/dns/v1beta2"
	file "google.golang.org/api/file/v1beta1"
	"google.golang.org/api/iam/v1"
	cloudlogging "google.golang.org/api/logging/v2"
	"google.golang.org/api/pubsub/v1"
	"google.golang.org/api/redis/v1beta1"
	"google.golang.org/api/runtimeconfig/v1beta1"
	"google.golang.org/api/servicemanagement/v1"
	"google.golang.org/api/serviceusage/v1beta1"
	"google.golang.org/api/sourcerepo/v1"
	"google.golang.org/api/spanner/v1"
	"google.golang.org/api/sqladmin/v1beta4"
	"google.golang.org/api/storage/v1"
)

// Config is the configuration structure used to instantiate the Google
// provider.
type Config struct {
	Credentials string
	Project     string
	Region      string
	Zone        string

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
	clientRedis                  *redis.Service
	clientResourceManager        *cloudresourcemanager.Service
	clientResourceManagerV2Beta1 *resourceManagerV2Beta1.Service
	clientRuntimeconfig          *runtimeconfig.Service
	clientSpanner                *spanner.Service
	clientSourceRepo             *sourcerepo.Service
	clientStorage                *storage.Service
	clientSqlAdmin               *sqladmin.Service
	clientIAM                    *iam.Service
	clientServiceMan             *servicemanagement.APIService
	clientServiceUsage           *serviceusage.APIService
	clientBigQuery               *bigquery.Service
	clientCloudFunctions         *cloudfunctions.Service
	clientCloudIoT               *cloudiot.Service
	clientAppEngine              *appengine.APIService

	bigtableClientFactory *BigtableClientFactory
}

func (c *Config) loadAndValidate() error {
	var account accountFile
	clientScopes := []string{
		"https://www.googleapis.com/auth/compute",
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/ndev.clouddns.readwrite",
		"https://www.googleapis.com/auth/devstorage.full_control",
	}

	var client *http.Client
	var tokenSource oauth2.TokenSource

	if c.Credentials != "" {
		contents, _, err := pathorcontents.Read(c.Credentials)
		if err != nil {
			return fmt.Errorf("Error loading credentials: %s", err)
		}

		// Assume account_file is a JSON string
		if err := parseJSON(&account, contents); err != nil {
			return fmt.Errorf("Error parsing credentials '%s': %s", contents, err)
		}

		// Get the token for use in our requests
		log.Printf("[INFO] Requesting Google token...")
		log.Printf("[INFO]   -- Email: %s", account.ClientEmail)
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		log.Printf("[INFO]   -- Private Key Length: %d", len(account.PrivateKey))

		conf := jwt.Config{
			Email:      account.ClientEmail,
			PrivateKey: []byte(account.PrivateKey),
			Scopes:     clientScopes,
			TokenURL:   "https://accounts.google.com/o/oauth2/token",
		}

		// Initiate an http.Client. The following GET request will be
		// authorized and authenticated on the behalf of
		// your service account.
		client = conf.Client(context.Background())

		tokenSource = conf.TokenSource(context.Background())
	} else {
		log.Printf("[INFO] Authenticating using DefaultClient")
		err := error(nil)
		client, err = google.DefaultClient(context.Background(), clientScopes...)
		if err != nil {
			return err
		}

		tokenSource, err = google.DefaultTokenSource(context.Background(), clientScopes...)
		if err != nil {
			return err
		}
	}

	c.tokenSource = tokenSource

	client.Transport = logging.NewTransport("Google", client.Transport)

	terraformVersion := httpclient.UserAgentString()
	providerVersion := fmt.Sprintf("terraform-provider-google/%s", version.ProviderVersion)
	terraformWebsite := "(+https://www.terraform.io)"
	userAgent := fmt.Sprintf("%s %s %s", terraformVersion, terraformWebsite, providerVersion)

	c.client = client
	c.userAgent = userAgent

	var err error

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

	log.Printf("[INFO] Instantiating Google Cloud Redis Client...")
	c.clientRedis, err = redis.New(client)
	if err != nil {
		return err
	}
	c.clientRedis.UserAgent = userAgent

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

	return nil
}

// accountFile represents the structure of the account file JSON file.
type accountFile struct {
	PrivateKeyId string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientId     string `json:"client_id"`
}

func parseJSON(result interface{}, contents string) error {
	r := strings.NewReader(contents)
	dec := json.NewDecoder(r)

	return dec.Decode(result)
}
