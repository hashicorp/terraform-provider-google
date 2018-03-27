package google

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
)

const (
	canonicalSslCertificateTemplate = "https://www.googleapis.com/compute/v1/projects/%s/global/sslCertificates/%s"
)

func resourceComputeTargetHttpsProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeTargetHttpsProxyCreate,
		Read:   resourceComputeTargetHttpsProxyRead,
		Delete: resourceComputeTargetHttpsProxyDelete,
		Update: resourceComputeTargetHttpsProxyUpdate,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ssl_certificates": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: compareSelfLinkOrResourceName,
				},
			},

			"url_map": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkRelativePaths,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"proxy_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeTargetHttpsProxyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	sslCertificates, err := expandSslCertificates(d, config)
	if err != nil {
		return err
	}

	proxy := &compute.TargetHttpsProxy{
		Name:            d.Get("name").(string),
		UrlMap:          d.Get("url_map").(string),
		SslCertificates: sslCertificates,
	}

	if v, ok := d.GetOk("description"); ok {
		proxy.Description = v.(string)
	}

	log.Printf("[DEBUG] TargetHttpsProxy insert request: %#v", proxy)
	op, err := config.clientCompute.TargetHttpsProxies.Insert(
		project, proxy).Do()
	if err != nil {
		return fmt.Errorf("Error creating TargetHttpsProxy: %s", err)
	}

	err = computeOperationWait(config.clientCompute, op, project, "Creating Target Https Proxy")
	if err != nil {
		return err
	}

	d.SetId(proxy.Name)

	return resourceComputeTargetHttpsProxyRead(d, meta)
}

func resourceComputeTargetHttpsProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	if d.HasChange("url_map") {
		url_map := d.Get("url_map").(string)
		url_map_ref := &compute.UrlMapReference{UrlMap: url_map}
		op, err := config.clientCompute.TargetHttpsProxies.SetUrlMap(
			project, d.Id(), url_map_ref).Do()
		if err != nil {
			return fmt.Errorf("Error updating Target HTTPS proxy URL map: %s", err)
		}

		err = computeOperationWait(config.clientCompute, op, project, "Updating Target Https Proxy URL Map")
		if err != nil {
			return err
		}

		d.SetPartial("url_map")
	}

	if d.HasChange("ssl_certificates") {
		certs, err := expandSslCertificates(d, config)
		if err != nil {
			return err
		}
		cert_ref := &compute.TargetHttpsProxiesSetSslCertificatesRequest{
			SslCertificates: certs,
		}
		op, err := config.clientCompute.TargetHttpsProxies.SetSslCertificates(
			project, d.Id(), cert_ref).Do()
		if err != nil {
			return fmt.Errorf("Error updating Target Https Proxy SSL Certificates: %s", err)
		}

		err = computeOperationWait(config.clientCompute, op, project, "Updating Target Https Proxy SSL certificates")
		if err != nil {
			return err
		}

		d.SetPartial("ssl_certificate")
	}

	d.Partial(false)

	return resourceComputeTargetHttpsProxyRead(d, meta)
}

func resourceComputeTargetHttpsProxyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	proxy, err := config.clientCompute.TargetHttpsProxies.Get(
		project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Target HTTPS proxy %q", d.Get("name").(string)))
	}

	d.Set("ssl_certificates", proxy.SslCertificates)
	d.Set("proxy_id", strconv.FormatUint(proxy.Id, 10))
	d.Set("self_link", proxy.SelfLink)
	d.Set("description", proxy.Description)
	d.Set("url_map", proxy.UrlMap)
	d.Set("name", proxy.Name)
	d.Set("project", project)

	return nil
}

func resourceComputeTargetHttpsProxyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the TargetHttpsProxy
	log.Printf("[DEBUG] TargetHttpsProxy delete request")
	op, err := config.clientCompute.TargetHttpsProxies.Delete(
		project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting TargetHttpsProxy: %s", err)
	}

	err = computeOperationWait(config.clientCompute, op, project, "Deleting Target Https Proxy")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
