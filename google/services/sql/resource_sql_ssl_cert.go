// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func ResourceSqlSslCert() *schema.Resource {
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
		UseJSONNumber: true,
	}
}

func resourceSqlSslCertCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instance := d.Get("instance").(string)
	commonName := d.Get("common_name").(string)

	sslCertsInsertRequest := &sqladmin.SslCertsInsertRequest{
		CommonName: commonName,
	}

	transport_tpg.MutexStore.Lock(instanceMutexKey(project, instance))
	defer transport_tpg.MutexStore.Unlock(instanceMutexKey(project, instance))
	resp, err := config.NewSqlAdminClient(userAgent).SslCerts.Insert(project, instance, sslCertsInsertRequest).Do()
	if err != nil {
		return fmt.Errorf("Error, failed to insert "+
			"ssl cert %s into instance %s: %s", commonName, instance, err)
	}

	err = SqlAdminOperationWaitTime(config, resp.Operation, project, "Create Ssl Cert", userAgent, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error, failure waiting for creation of %q "+
			"in %q: %s", commonName, instance, err)
	}

	fingerprint := resp.ClientCert.CertInfo.Sha1Fingerprint
	d.SetId(fmt.Sprintf("projects/%s/instances/%s/sslCerts/%s", project, instance, fingerprint))
	if err := d.Set("sha1_fingerprint", fingerprint); err != nil {
		return fmt.Errorf("Error setting sha1_fingerprint: %s", err)
	}

	// The private key is only returned on the initial insert so set it here.
	if err := d.Set("private_key", resp.ClientCert.CertPrivateKey); err != nil {
		return fmt.Errorf("Error setting private_key: %s", err)
	}
	if err := d.Set("server_ca_cert", resp.ServerCaCert.Cert); err != nil {
		return fmt.Errorf("Error setting server_ca_cert: %s", err)
	}

	return resourceSqlSslCertRead(d, meta)
}

func resourceSqlSslCertRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instance := d.Get("instance").(string)
	commonName := d.Get("common_name").(string)
	fingerprint := d.Get("sha1_fingerprint").(string)

	sslCerts, err := config.NewSqlAdminClient(userAgent).SslCerts.Get(project, instance, fingerprint).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SQL Ssl Cert %q in instance %q", commonName, instance))
	}

	if sslCerts == nil {
		log.Printf("[WARN] Removing SQL Ssl Cert %q because it's gone", commonName)
		d.SetId("")

		return nil
	}

	if err := d.Set("instance", sslCerts.Instance); err != nil {
		return fmt.Errorf("Error setting instance: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("sha1_fingerprint", sslCerts.Sha1Fingerprint); err != nil {
		return fmt.Errorf("Error setting sha1_fingerprint: %s", err)
	}
	if err := d.Set("common_name", sslCerts.CommonName); err != nil {
		return fmt.Errorf("Error setting common_name: %s", err)
	}
	if err := d.Set("cert", sslCerts.Cert); err != nil {
		return fmt.Errorf("Error setting cert: %s", err)
	}
	if err := d.Set("cert_serial_number", sslCerts.CertSerialNumber); err != nil {
		return fmt.Errorf("Error setting cert_serial_number: %s", err)
	}
	if err := d.Set("create_time", sslCerts.CreateTime); err != nil {
		return fmt.Errorf("Error setting create_time: %s", err)
	}
	if err := d.Set("expiration_time", sslCerts.ExpirationTime); err != nil {
		return fmt.Errorf("Error setting expiration_time: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/instances/%s/sslCerts/%s", project, instance, fingerprint))
	return nil
}

func resourceSqlSslCertDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instance := d.Get("instance").(string)
	commonName := d.Get("common_name").(string)
	fingerprint := d.Get("sha1_fingerprint").(string)

	transport_tpg.MutexStore.Lock(instanceMutexKey(project, instance))
	defer transport_tpg.MutexStore.Unlock(instanceMutexKey(project, instance))
	op, err := config.NewSqlAdminClient(userAgent).SslCerts.Delete(project, instance, fingerprint).Do()

	if err != nil {
		return fmt.Errorf("Error, failed to delete "+
			"ssl cert %q in instance %q: %s", commonName,
			instance, err)
	}

	err = SqlAdminOperationWaitTime(config, op, project, "Delete Ssl Cert", userAgent, d.Timeout(schema.TimeoutDelete))

	if err != nil {
		return fmt.Errorf("Error, failure waiting for deletion of ssl cert %q "+
			"in %q: %s", commonName, instance, err)
	}

	return nil
}
