package google

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func init() {
	resource.AddTestSweepers("CertificateAuthority", &resource.Sweeper{
		Name: "CertificateAuthority",
		F:    testSweepCertificateAuthority,
	})
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepCertificateAuthority(region string) error {
	resourceName := "CertificateAuthority"
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for %s", resourceName)

	config, err := acctest.SharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}

	// Setup variables to replace in list template
	d := &acctest.ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"project":  config.Project,
			"location": region,
		},
	}

	caPoolsUrl, err := ReplaceVars(d, config, "{{PrivatecaBasePath}}projects/{{project}}/locations/{{location}}/caPools")
	if err != nil {
		return err
	}

	res, err := transport_tpg.SendRequest(config, "GET", config.Project, caPoolsUrl, config.UserAgent, nil)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error in response from request %s: %s", caPoolsUrl, err)
		return nil
	}

	resourceList, ok := res["caPools"]
	if !ok {
		log.Printf("[INFO][SWEEPER_LOG] Nothing found in response.")
		return nil
	}

	rl := resourceList.([]interface{})

	log.Printf("[INFO][SWEEPER_LOG] Found %d items in %s list response.", len(rl), resourceName)
	// Count items that weren't sweeped.
	nonPrefixCount := 0
	for _, ri := range rl {
		obj := ri.(map[string]interface{})

		poolName := obj["name"].(string)

		caListUrl := config.PrivatecaBasePath + poolName + "/certificateAuthorities"

		res, err := transport_tpg.SendRequest(config, "GET", config.Project, caListUrl, config.UserAgent, nil)
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] Error in response from request %s: %s", caPoolsUrl, err)
			return nil
		}

		caResourceList, ok := res["certificateAuthorities"]
		if !ok {
			log.Printf("[INFO][SWEEPER_LOG] Nothing found in certificate authority list response.")
			continue
		}

		carl := caResourceList.([]interface{})
		for _, cai := range carl {
			obj := cai.(map[string]interface{})
			caName := obj["name"].(string)

			// Increment count and skip if resource is not sweepable.
			nameParts := strings.Split(caName, "/")
			id := nameParts[len(nameParts)-1]
			if !acctest.IsSweepableTestResource(id) {
				nonPrefixCount++
				continue
			}

			if obj["state"] == "DELETED" {
				continue
			}

			if obj["state"] == "ENABLED" {
				disableUrl := fmt.Sprintf("%s%s:disable", config.PrivatecaBasePath, caName)
				_, err = transport_tpg.SendRequest(config, "POST", config.Project, disableUrl, config.UserAgent, nil)
				if err != nil {
					log.Printf("[INFO][SWEEPER_LOG] Error disabling for url %s : %s", disableUrl, err)
				} else {
					log.Printf("[INFO][SWEEPER_LOG] Disabling %s resource: %s", resourceName, caName)
				}
			}

			deleteUrl := config.PrivatecaBasePath + caName
			_, err = transport_tpg.SendRequest(config, "DELETE", config.Project, deleteUrl, config.UserAgent, nil)
			if err != nil {
				log.Printf("[INFO][SWEEPER_LOG] Error deleting for url %s : %s", deleteUrl, err)
			} else {
				log.Printf("[INFO][SWEEPER_LOG] Deleted a %s resource: %s", resourceName, caName)
			}
		}
	}

	if nonPrefixCount > 0 {
		log.Printf("[INFO][SWEEPER_LOG] %d items without tf-test prefix remain.", nonPrefixCount)
	}

	return nil
}
