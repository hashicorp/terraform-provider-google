package google

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
)

func resourceComputeTargetSslProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeTargetSslProxyCreate,
		Read:   resourceComputeTargetSslProxyRead,
		Delete: resourceComputeTargetSslProxyDelete,
		Update: resourceComputeTargetSslProxyUpdate,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"backend_service": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"ssl_certificates": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: compareSelfLinkOrResourceName,
				},
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"proxy_header": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NONE",
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"proxy_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeTargetSslProxyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	sslCertificates, err := expandSslCertificates(d.Get("ssl_certificates").([]interface{}), d, config)
	if err != nil {
		return err
	}

	proxy := &compute.TargetSslProxy{
		Name:            d.Get("name").(string),
		Service:         d.Get("backend_service").(string),
		ProxyHeader:     d.Get("proxy_header").(string),
		Description:     d.Get("description").(string),
		SslCertificates: sslCertificates,
	}

	log.Printf("[DEBUG] TargetSslProxy insert request: %#v", proxy)
	op, err := config.clientCompute.TargetSslProxies.Insert(
		project, proxy).Do()
	if err != nil {
		return fmt.Errorf("Error creating TargetSslProxy: %s", err)
	}

	err = computeOperationWait(config, op, project, "Creating Target Ssl Proxy")
	if err != nil {
		return err
	}

	d.SetId(proxy.Name)

	return resourceComputeTargetSslProxyRead(d, meta)
}

func resourceComputeTargetSslProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	if d.HasChange("proxy_header") {
		proxy_header := d.Get("proxy_header").(string)
		proxy_header_payload := &compute.TargetSslProxiesSetProxyHeaderRequest{
			ProxyHeader: proxy_header,
		}
		op, err := config.clientCompute.TargetSslProxies.SetProxyHeader(
			project, d.Id(), proxy_header_payload).Do()
		if err != nil {
			return fmt.Errorf("Error updating proxy_header: %s", err)
		}

		err = computeOperationWait(config, op, project, "Updating Target SSL Proxy")
		if err != nil {
			return err
		}

		d.SetPartial("proxy_header")
	}

	if d.HasChange("backend_service") {
		op, err := config.clientCompute.TargetSslProxies.SetBackendService(project, d.Id(), &compute.TargetSslProxiesSetBackendServiceRequest{
			Service: d.Get("backend_service").(string),
		}).Do()

		if err != nil {
			return fmt.Errorf("Error updating backend_service: %s", err)
		}

		err = computeOperationWait(config, op, project, "Updating Target SSL Proxy")
		if err != nil {
			return err
		}

		d.SetPartial("backend_service")
	}

	if d.HasChange("ssl_certificates") {
		sslCertificates, err := expandSslCertificates(d.Get("ssl_certificates").([]interface{}), d, config)
		if err != nil {
			return err
		}

		op, err := config.clientCompute.TargetSslProxies.SetSslCertificates(project, d.Id(), &compute.TargetSslProxiesSetSslCertificatesRequest{
			SslCertificates: sslCertificates,
		}).Do()

		if err != nil {
			return fmt.Errorf("Error updating backend_service: %s", err)
		}

		err = computeOperationWait(config, op, project, "Updating Target SSL Proxy")
		if err != nil {
			return err
		}

		d.SetPartial("ssl_certificates")
	}

	d.Partial(false)

	return resourceComputeTargetSslProxyRead(d, meta)
}

func resourceComputeTargetSslProxyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	proxy, err := config.clientCompute.TargetSslProxies.Get(
		project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Target SSL Proxy %q", d.Get("name").(string)))
	}

	d.Set("name", proxy.Name)
	d.Set("description", proxy.Description)
	d.Set("proxy_header", proxy.ProxyHeader)
	d.Set("backend_service", proxy.Service)
	d.Set("ssl_certificates", proxy.SslCertificates)
	d.Set("self_link", proxy.SelfLink)
	d.Set("proxy_id", strconv.FormatUint(proxy.Id, 10))

	return nil
}

func resourceComputeTargetSslProxyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	op, err := config.clientCompute.TargetSslProxies.Delete(
		project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting TargetSslProxy: %s", err)
	}

	err = computeOperationWait(config, op, project, "Deleting Target SSL Proxy")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func expandSslCertificates(configured []interface{}, d *schema.ResourceData, config *Config) ([]string, error) {
	certs := make([]string, 0, len(configured))

	for _, sslCertificate := range configured {
		sslCertificateFieldValue, err := ParseSslCertificateFieldValue(sslCertificate.(string), d, config)
		if err != nil {
			return nil, fmt.Errorf("Invalid ssl certificate: %s", err)
		}

		certs = append(certs, sslCertificateFieldValue.RelativeLink())
	}

	return certs, nil
}
