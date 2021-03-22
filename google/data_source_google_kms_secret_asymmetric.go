package google

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"hash/crc32"
	"regexp"
	"strconv"
)

var (
	cryptoKeyVersionRegexp = regexp.MustCompile(`^(//[^/]*/[^/]*/)?(projects/[^/]+/locations/[^/]+/keyRings/[^/]+/cryptoKeys/[^/]+/cryptoKeyVersions/[^/]+)$`)
)

func dataSourceGoogleKmsSecretAsymmetric() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGoogleKmsSecretAsymmetricReadContext,
		Schema: map[string]*schema.Schema{
			"crypto_key_version": {
				Type:         schema.TypeString,
				Description:  "The fully qualified KMS crypto key version name",
				ValidateFunc: validateRegexp(cryptoKeyVersionRegexp.String()),
				Required:     true,
			},
			"ciphertext": {
				Type:         schema.TypeString,
				Description:  "The public key encrypted ciphertext in base64 encoding",
				ValidateFunc: validateBase64WithWhitespaces,
				Required:     true,
			},
			"crc32": {
				Type:         schema.TypeString,
				Description:  "The crc32 checksum of the ciphertext, hexadecimal encoding",
				ValidateFunc: validateHexadecimalUint32,
				Optional:     true,
			},
			"plaintext": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceGoogleKmsSecretAsymmetricReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	err := dataSourceGoogleKmsSecretAsymmetricRead(ctx, d, meta)
	if err != nil {
		diags = diag.FromErr(err)
	}
	return diags
}

func dataSourceGoogleKmsSecretAsymmetricRead(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	// `google_kms_crypto_key_version` returns an id with the prefix
	// //cloudkms.googleapis.com/v1, which is an invalid name. To allow for the most elegant
	// configuration, we will allow it as an input.
	keyVersion := cryptoKeyVersionRegexp.FindStringSubmatch(d.Get("crypto_key_version").(string))
	cryptoKeyVersion := keyVersion[len(keyVersion)-1]

	base64CipherText := removeWhiteSpaceFromString(d.Get("ciphertext").(string))
	ciphertext, err := base64.StdEncoding.DecodeString(base64CipherText)
	if err != nil {
		return err
	}

	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}

	ciphertextCRC32C := crc32c(ciphertext)
	if s, ok := d.Get("crc32").(string); ok && s != "" {
		u, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return fmt.Errorf("failed to convert crc32 into uint32, %s", err)
		}
		ciphertextCRC32C = uint32(u)
	} else {
		if err := d.Set("crc32", fmt.Sprintf("%x", ciphertextCRC32C)); err != nil {
			return fmt.Errorf("failed to set crc32, %s", err)
		}
	}

	req := &kmspb.AsymmetricDecryptRequest{
		Name:             cryptoKeyVersion,
		Ciphertext:       ciphertext,
		CiphertextCrc32C: wrapperspb.Int64(int64(ciphertextCRC32C)),
	}

	client := config.NewKeyManagementClient(ctx, userAgent)
	result, err := client.AsymmetricDecrypt(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to decrypt ciphertext: %v", err)
	}

	if !result.VerifiedCiphertextCrc32C || int64(crc32c(result.Plaintext)) != result.PlaintextCrc32C.Value {
		return fmt.Errorf("asymmetricDecrypt request corrupted in-transit")
	}

	if err := d.Set("plaintext", string(result.Plaintext)); err != nil {
		return fmt.Errorf("error setting plaintext: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%x:%s", cryptoKeyVersion, ciphertextCRC32C, base64CipherText))
	return nil
}

func removeWhiteSpaceFromString(s string) string {
	whitespaceRegexp := regexp.MustCompile(`(?m)[\s]+`)
	return whitespaceRegexp.ReplaceAllString(s, "")
}

func validateBase64WithWhitespaces(i interface{}, val string) ([]string, []error) {
	_, err := base64.StdEncoding.DecodeString(removeWhiteSpaceFromString(i.(string)))
	if err != nil {
		return nil, []error{fmt.Errorf("could not decode %q as a valid base64 value. Please use the terraform base64 functions such as base64encode() or filebase64() to supply a valid base64 string", val)}
	}
	return nil, nil
}

func validateHexadecimalUint32(i interface{}, val string) ([]string, []error) {
	_, err := strconv.ParseUint(i.(string), 16, 32)
	if err != nil {
		return nil, []error{fmt.Errorf("could not decode %q as a unsigned 32 bit hexadecimal integer", val)}
	}
	return nil, nil
}
