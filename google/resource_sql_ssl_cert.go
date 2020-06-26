package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func resourceSqlSslCert() *schema.Resource {
	return &schema.Resource{
		Create: resourceSqlSslCertCreate,
		Read:   resourceSqlSslCertRead,
		Delete: resourceSqlSslCertDelete,

		SchemaVersion: 1,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"common_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The common name to be used in the certificate to identify the client. Constrained to [a-zA-Z.-_ ]+. Changing this forces a new resource to be created.`,
			},

			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the Cloud SQL instance. Changing this forces a new resource to be created.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"cert": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The actual certificate data for this client certificate.`,
			},

			"cert_serial_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The serial number extracted from the certificate data.`,
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the certificate was created in RFC 3339 format, for example 2012-11-15T16:19:00.094Z.`,
			},

			"expiration_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the certificate expires in RFC 3339 format, for example 2012-11-15T16:19:00.094Z.`,
			},

			"private_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `The private key associated with the client certificate.`,
			},

			"server_ca_cert": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The CA cert of the server this client cert was generated from.`,
			},

			"sha1_fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SHA1 Fingerprint of the certificate.`,
			},
		},
	}
}

func resourceSqlSslCertCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instance := d.Get("instance").(string)
	commonName := d.Get("common_name").(string)

	sslCertsInsertRequest := &sqladmin.SslCertsInsertRequest{
		CommonName: commonName,
	}

	mutexKV.Lock(instanceMutexKey(project, instance))
	defer mutexKV.Unlock(instanceMutexKey(project, instance))
	resp, err := config.clientSqlAdmin.SslCerts.Insert(project, instance, sslCertsInsertRequest).Do()
	if err != nil {
		return fmt.Errorf("Error, failed to insert "+
			"ssl cert %s into instance %s: %s", commonName, instance, err)
	}

	err = sqlAdminOperationWaitTime(config, resp.Operation, project, "Create Ssl Cert", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error, failure waiting for creation of %q "+
			"in %q: %s", commonName, instance, err)
	}

	fingerprint := resp.ClientCert.CertInfo.Sha1Fingerprint
	d.SetId(fmt.Sprintf("projects/%s/instances/%s/sslCerts/%s", project, instance, fingerprint))
	d.Set("sha1_fingerprint", fingerprint)

	// The private key is only returned on the initial insert so set it here.
	d.Set("private_key", resp.ClientCert.CertPrivateKey)
	d.Set("server_ca_cert", resp.ServerCaCert.Cert)

	return resourceSqlSslCertRead(d, meta)
}

func resourceSqlSslCertRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instance := d.Get("instance").(string)
	commonName := d.Get("common_name").(string)
	fingerprint := d.Get("sha1_fingerprint").(string)

	sslCerts, err := config.clientSqlAdmin.SslCerts.Get(project, instance, fingerprint).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SQL Ssl Cert %q in instance %q", commonName, instance))
	}

	if sslCerts == nil {
		log.Printf("[WARN] Removing SQL Ssl Cert %q because it's gone", commonName)
		d.SetId("")

		return nil
	}

	d.Set("instance", sslCerts.Instance)
	d.Set("project", project)
	d.Set("sha1_fingerprint", sslCerts.Sha1Fingerprint)
	d.Set("common_name", sslCerts.CommonName)
	d.Set("cert", sslCerts.Cert)
	d.Set("cert_serial_number", sslCerts.CertSerialNumber)
	d.Set("create_time", sslCerts.CreateTime)
	d.Set("expiration_time", sslCerts.ExpirationTime)

	d.SetId(fmt.Sprintf("projects/%s/instances/%s/sslCerts/%s", project, instance, fingerprint))
	return nil
}

func resourceSqlSslCertDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instance := d.Get("instance").(string)
	commonName := d.Get("common_name").(string)
	fingerprint := d.Get("sha1_fingerprint").(string)

	mutexKV.Lock(instanceMutexKey(project, instance))
	defer mutexKV.Unlock(instanceMutexKey(project, instance))
	op, err := config.clientSqlAdmin.SslCerts.Delete(project, instance, fingerprint).Do()

	if err != nil {
		return fmt.Errorf("Error, failed to delete "+
			"ssl cert %q in instance %q: %s", commonName,
			instance, err)
	}

	err = sqlAdminOperationWaitTime(config, op, project, "Delete Ssl Cert", d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return fmt.Errorf("Error, failure waiting for deletion of ssl cert %q "+
			"in %q: %s", commonName, instance, err)
	}

	return nil
}
