package google

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CA related utilities.

func enableCA(config *Config, d *schema.ResourceData, project string, billingProject string, userAgent string) error {
	enableUrl, err := replaceVars(d, config, "{{PrivatecaBasePath}}projects/{{project}}/locations/{{location}}/caPools/{{pool}}/certificateAuthorities/{{certificate_authority_id}}:enable")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Enabling CertificateAuthority")

	res, err := SendRequest(config, "POST", billingProject, enableUrl, userAgent, nil)
	if err != nil {
		return fmt.Errorf("Error enabling CertificateAuthority: %s", err)
	}

	var opRes map[string]interface{}
	err = PrivatecaOperationWaitTimeWithResponse(
		config, res, &opRes, project, "Enabling CertificateAuthority", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error waiting to enable CertificateAuthority: %s", err)
	}
	return nil
}

func disableCA(config *Config, d *schema.ResourceData, project string, billingProject string, userAgent string) error {
	disableUrl, err := replaceVars(d, config, "{{PrivatecaBasePath}}projects/{{project}}/locations/{{location}}/caPools/{{pool}}/certificateAuthorities/{{certificate_authority_id}}:disable")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Disabling CA")

	dRes, err := SendRequest(config, "POST", billingProject, disableUrl, userAgent, nil)
	if err != nil {
		return fmt.Errorf("Error disabling CA: %s", err)
	}

	var opRes map[string]interface{}
	err = PrivatecaOperationWaitTimeWithResponse(
		config, dRes, &opRes, project, "Disabling CA", userAgent,
		d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("Error waiting to disable CA: %s", err)
	}
	return nil
}

func activateSubCAWithThirdPartyIssuer(config *Config, d *schema.ResourceData, project string, billingProject string, userAgent string) error {
	// 1. prepare parameters
	signedCACert := d.Get("pem_ca_certificate").(string)

	sc, ok := d.GetOk("subordinate_config")
	if !ok {
		return fmt.Errorf("subordinate_config is required to activate subordinate CA")
	}
	c := sc.([]interface{})
	if len(c) == 0 || c[0] == nil {
		return fmt.Errorf("subordinate_config is required to activate subordinate CA")
	}
	chain, ok := c[0].(map[string]interface{})["pem_issuer_chain"]
	if !ok {
		return fmt.Errorf("subordinate_config.pem_issuer_chain is required to activate subordinate CA with third party issuer")
	}
	issuerChain := chain.([]interface{})
	if len(issuerChain) == 0 || issuerChain[0] == nil {
		return fmt.Errorf("subordinate_config.pem_issuer_chain is required to activate subordinate CA with third party issuer")
	}
	pc := issuerChain[0].(map[string]interface{})["pem_certificates"].([]interface{})
	pemIssuerChain := make([]string, 0, len(pc))
	for _, pem := range pc {
		pemIssuerChain = append(pemIssuerChain, pem.(string))
	}

	// 2. activate CA
	activateObj := make(map[string]interface{})
	activateObj["pemCaCertificate"] = signedCACert
	activateObj["subordinateConfig"] = make(map[string]interface{})
	activateObj["subordinateConfig"].(map[string]interface{})["pemIssuerChain"] = make(map[string]interface{})
	activateObj["subordinateConfig"].(map[string]interface{})["pemIssuerChain"].(map[string]interface{})["pemCertificates"] = pemIssuerChain

	activateUrl, err := replaceVars(d, config, "{{PrivatecaBasePath}}projects/{{project}}/locations/{{location}}/caPools/{{pool}}/certificateAuthorities/{{certificate_authority_id}}:activate")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Activating CertificateAuthority: %#v", activateObj)
	res, err := SendRequest(config, "POST", billingProject, activateUrl, userAgent, activateObj)
	if err != nil {
		return fmt.Errorf("Error enabling CertificateAuthority: %s", err)
	}

	var opRes map[string]interface{}
	err = PrivatecaOperationWaitTimeWithResponse(
		config, res, &opRes, project, "Activating CertificateAuthority", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error waiting to actiavte CertificateAuthority: %s", err)
	}
	return nil
}

func activateSubCAWithFirstPartyIssuer(config *Config, d *schema.ResourceData, project string, billingProject string, userAgent string) error {
	// 1. get issuer
	sc, ok := d.GetOk("subordinate_config")
	if !ok {
		return fmt.Errorf("subordinate_config is required to activate subordinate CA")
	}
	c := sc.([]interface{})
	if len(c) == 0 || c[0] == nil {
		return fmt.Errorf("subordinate_config is required to activate subordinate CA")
	}
	ca, ok := c[0].(map[string]interface{})["certificate_authority"]
	if !ok {
		return fmt.Errorf("subordinate_config.certificate_authority is required to activate subordinate CA with first party issuer")
	}
	issuer := ca.(string)

	// 2. fetch CSR
	fetchCSRUrl, err := replaceVars(d, config, "{{PrivatecaBasePath}}projects/{{project}}/locations/{{location}}/caPools/{{pool}}/certificateAuthorities/{{certificate_authority_id}}:fetch")
	if err != nil {
		return err
	}
	res, err := SendRequest(config, "GET", billingProject, fetchCSRUrl, userAgent, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch CSR: %v", err)
	}
	csr := res["pemCsr"]

	// 3. sign the CSR with first party issuer
	genCertId := func() string {
		currentTime := time.Now()
		dateStr := currentTime.Format("20060102")

		rand.Seed(time.Now().UnixNano())
		const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		rand1 := make([]byte, 3)
		for i := range rand1 {
			rand1[i] = letters[rand.Intn(len(letters))]
		}
		rand2 := make([]byte, 3)
		for i := range rand2 {
			rand2[i] = letters[rand.Intn(len(letters))]
		}
		return fmt.Sprintf("subordinate-%v-%v-%v", dateStr, string(rand1), string(rand2))
	}

	// parseCAName parses a CA name and return the CaPool name and CaId.
	parseCAName := func(n string) (string, string, error) {
		parts := regexp.MustCompile(`(projects/[a-z0-9-]+/locations/[a-z0-9-]+/caPools/[a-zA-Z0-9-]+)/certificateAuthorities/([a-zA-Z0-9-]+)`).FindStringSubmatch(n)
		if len(parts) != 3 {
			return "", "", fmt.Errorf("failed to parse CA name: %v, parts: %v", n, parts)
		}
		return parts[1], parts[2], err
	}

	obj := make(map[string]interface{})
	obj["pemCsr"] = csr
	obj["lifetime"] = d.Get("lifetime")

	certId := genCertId()
	poolName, issuerId, err := parseCAName(issuer)
	if err != nil {
		return err
	}

	PrivatecaBasePath, err := replaceVars(d, config, "{{PrivatecaBasePath}}")
	if err != nil {
		return err
	}
	signUrl := fmt.Sprintf("%v%v/certificates?certificateId=%v", PrivatecaBasePath, poolName, certId)
	signUrl, err = addQueryParams(signUrl, map[string]string{"issuingCertificateAuthorityId": issuerId})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Signing CA Certificate: %#v", obj)
	res, err = SendRequestWithTimeout(config, "POST", billingProject, signUrl, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Certificate: %s", err)
	}
	signedCACert := res["pemCertificate"]

	// 4. activate sub CA with the signed CA cert.
	activateObj := make(map[string]interface{})
	activateObj["pemCaCertificate"] = signedCACert
	activateObj["subordinateConfig"] = make(map[string]interface{})
	activateObj["subordinateConfig"].(map[string]interface{})["certificateAuthority"] = issuer

	activateUrl, err := replaceVars(d, config, "{{PrivatecaBasePath}}projects/{{project}}/locations/{{location}}/caPools/{{pool}}/certificateAuthorities/{{certificate_authority_id}}:activate")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Activating CertificateAuthority: %#v", activateObj)
	res, err = SendRequest(config, "POST", billingProject, activateUrl, userAgent, activateObj)
	if err != nil {
		return fmt.Errorf("Error enabling CertificateAuthority: %s", err)
	}

	var opRes map[string]interface{}
	err = PrivatecaOperationWaitTimeWithResponse(
		config, res, &opRes, project, "Enabling CertificateAuthority", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error waiting to actiavte CertificateAuthority: %s", err)
	}
	return nil
}
