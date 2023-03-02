package google

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecretManagerSecretVersionAccess() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretManagerSecretVersionAccessRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secret": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_data": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceSecretManagerSecretVersionAccessRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	fv, err := parseProjectFieldValue("secrets", d.Get("secret").(string), "project", d, config, false)
	if err != nil {
		return err
	}
	if d.Get("project").(string) != "" && d.Get("project").(string) != fv.Project {
		return fmt.Errorf("The project set on this secret version (%s) is not equal to the project where this secret exists (%s).", d.Get("project").(string), fv.Project)
	}
	project := fv.Project
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("secret", fv.Name); err != nil {
		return fmt.Errorf("Error setting secret: %s", err)
	}

	var url string
	versionNum := d.Get("version")

	if versionNum != "" {
		url, err = replaceVars(d, config, "{{SecretManagerBasePath}}projects/{{project}}/secrets/{{secret}}/versions/{{version}}")
		if err != nil {
			return err
		}
	} else {
		url, err = replaceVars(d, config, "{{SecretManagerBasePath}}projects/{{project}}/secrets/{{secret}}/versions/latest")
		if err != nil {
			return err
		}
	}

	url = fmt.Sprintf("%s:access", url)
	resp, err := SendRequest(config, "GET", project, url, userAgent, nil)
	if err != nil {
		return fmt.Errorf("Error retrieving available secret manager secret version access: %s", err.Error())
	}

	if err := d.Set("name", resp["name"].(string)); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	secretVersionRegex := regexp.MustCompile("projects/(.+)/secrets/(.+)/versions/(.+)$")

	parts := secretVersionRegex.FindStringSubmatch(resp["name"].(string))
	// should return [full string, project number, secret name, version number]
	if len(parts) != 4 {
		panic(fmt.Sprintf("secret name, %s, does not match format, projects/{{project}}/secrets/{{secret}}/versions/{{version}}", resp["name"].(string)))
	}

	log.Printf("[DEBUG] Received Google SecretManager Version: %q", parts[3])

	if err := d.Set("version", parts[3]); err != nil {
		return fmt.Errorf("Error setting version: %s", err)
	}

	data := resp["payload"].(map[string]interface{})
	secretData, err := base64.StdEncoding.DecodeString(data["data"].(string))
	if err != nil {
		return fmt.Errorf("Error decoding secret manager secret version data: %s", err.Error())
	}
	if err := d.Set("secret_data", string(secretData)); err != nil {
		return fmt.Errorf("Error setting secret_data: %s", err)
	}

	d.SetId(resp["name"].(string))
	return nil
}
