package google

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
)

func resourceComputeTargetTcpProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeTargetTcpProxyCreate,
		Read:   resourceComputeTargetTcpProxyRead,
		Delete: resourceComputeTargetTcpProxyDelete,
		Update: resourceComputeTargetTcpProxyUpdate,

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

			"proxy_header": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NONE",
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"proxy_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeTargetTcpProxyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	proxy := &compute.TargetTcpProxy{
		Name:        d.Get("name").(string),
		Service:     d.Get("backend_service").(string),
		ProxyHeader: d.Get("proxy_header").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[DEBUG] TargetTcpProxy insert request: %#v", proxy)
	op, err := config.clientCompute.TargetTcpProxies.Insert(
		project, proxy).Do()
	if err != nil {
		return fmt.Errorf("Error creating TargetTcpProxy: %s", err)
	}

	err = computeOperationWait(config, op, project, "Creating Target Tcp Proxy")
	if err != nil {
		return err
	}

	d.SetId(proxy.Name)

	return resourceComputeTargetTcpProxyRead(d, meta)
}

func resourceComputeTargetTcpProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	if d.HasChange("proxy_header") {
		proxy_header := d.Get("proxy_header").(string)
		proxy_header_payload := &compute.TargetTcpProxiesSetProxyHeaderRequest{
			ProxyHeader: proxy_header,
		}
		op, err := config.clientCompute.TargetTcpProxies.SetProxyHeader(
			project, d.Id(), proxy_header_payload).Do()
		if err != nil {
			return fmt.Errorf("Error updating target: %s", err)
		}

		err = computeOperationWait(config, op, project, "Updating Target Tcp Proxy")
		if err != nil {
			return err
		}

		d.SetPartial("proxy_header")
	}

	d.Partial(false)

	return resourceComputeTargetTcpProxyRead(d, meta)
}

func resourceComputeTargetTcpProxyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	proxy, err := config.clientCompute.TargetTcpProxies.Get(
		project, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Target TCP Proxy %q", d.Get("name").(string)))
	}

	d.Set("name", proxy.Name)
	d.Set("backend_service", proxy.Service)
	d.Set("proxy_header", proxy.ProxyHeader)
	d.Set("description", proxy.Description)
	d.Set("self_link", proxy.SelfLink)
	d.Set("proxy_id", strconv.FormatUint(proxy.Id, 10))

	return nil
}

func resourceComputeTargetTcpProxyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the TargetTcpProxy
	log.Printf("[DEBUG] TargetTcpProxy delete request")
	op, err := config.clientCompute.TargetTcpProxies.Delete(
		project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting TargetTcpProxy: %s", err)
	}

	err = computeOperationWait(config, op, project, "Deleting Target Tcp Proxy")
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
