package google

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecretManagerSecretVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretManagerSecretVersionRead,
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
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destroy_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
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

func dataSourceSecretManagerSecretVersionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	fv, err := parseProjectFieldValue("secrets", d.Get("secret").(string), "project", d, config, false)
	if err != nil {
		return err
	}
	if d.Get("project").(string) != "" && d.Get("project").(string) != fv.Project {
		return fmt.Errorf("The project set on this secret version (%s) is not equal to the project where this secret exists (%s).", d.Get("project").(string), fv.Project)
	}
	project := fv.Project
	d.Set("project", project)
	d.Set("secret", fv.Name)

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

	var version map[string]interface{}
	version, err = sendRequest(config, "GET", project, url, nil)
	if err != nil {
		return fmt.Errorf("Error retrieving available secret manager secret versions: %s", err.Error())
	}

	secretVersionRegex := regexp.MustCompile("projects/(.+)/secrets/(.+)/versions/(.+)$")

	parts := secretVersionRegex.FindStringSubmatch(version["name"].(string))
	// should return [full string, project number, secret name, version number]
	if len(parts) != 4 {
		panic(fmt.Sprintf("secret name, %s, does not match format, projects/{{project}}/secrets/{{secret}}/versions/{{version}}", version["name"].(string)))
	}

	log.Printf("[DEBUG] Received Google SecretManager Version: %q", version)

	d.Set("version", parts[3])

	url = fmt.Sprintf("%s:access", url)
	resp, err := sendRequest(config, "GET", project, url, nil)
	if err != nil {
		return fmt.Errorf("Error retrieving available secret manager secret version access: %s", err.Error())
	}

	d.Set("create_time", version["createTime"].(string))
	if version["destroyTime"] != nil {
		d.Set("destroy_time", version["destroyTime"].(string))
	}
	d.Set("name", version["name"].(string))
	d.Set("enabled", true)

	data := resp["payload"].(map[string]interface{})
	secretData, err := base64.StdEncoding.DecodeString(data["data"].(string))
	if err != nil {
		return fmt.Errorf("Error decoding secret manager secret version data: %s", err.Error())
	}
	d.Set("secret_data", string(secretData))

	d.SetId(version["name"].(string))
	return nil
}
