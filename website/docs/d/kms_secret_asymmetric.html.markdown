---
subcategory: "Cloud Key Management Service"
description: |-
  Provides access to secret data encrypted with Google Cloud KMS asymmetric key
---

# google_kms_secret_asymmetric

This data source allows you to use data encrypted with a Google Cloud KMS asymmetric key
within your resource definitions.

For more information see
[the official documentation](https://cloud.google.com/kms/docs/encrypt-decrypt-rsa).

~> **NOTE:** Using this data provider will allow you to conceal secret data within your
resource definitions, but it does not take care of protecting that data in the
logging output, plan output, or state output.  Please take care to secure your secret
data outside of resource definitions.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

## Example Usage

First, create a KMS KeyRing and CryptoKey using the resource definitions:

```hcl
resource "google_kms_key_ring" "my_key_ring" {
  project  = "my-project"
  name     = "my-key-ring"
  location = "us-central1"
}

resource "google_kms_crypto_key" "my_crypto_key" {
  name     = "my-crypto-key"
  key_ring = google_kms_key_ring.my_key_ring.id
  purpose  = "ASYMMETRIC_DECRYPT"
  version_template {
    algorithm = "RSA_DECRYPT_OAEP_4096_SHA256"
  }
}

data "google_kms_crypto_key_version" "my_crypto_key" {
  crypto_key = google_kms_crypto_key.my_crypto_key.id
}
```

Next, use the [Cloud SDK](https://cloud.google.com/kms/docs/encrypt-decrypt-rsa#kms-encrypt-asymmetric-cli) to encrypt 
some sensitive information:

```bash
## get the public key to encrypt the secret with
$ gcloud kms keys versions get-public-key 1 \
  --project my-project \
  --location us-central1 \
  --keyring my-key-ring \
  --key my-crypto-key \
  --output-file public-key.pem

## encrypt secret with the public key
$ echo -n my-secret-password | \
  openssl pkeyutl -in - \
    -encrypt \
    -pubin \
    -inkey public-key.pem \
    -pkeyopt rsa_padding_mode:oaep \
    -pkeyopt rsa_oaep_md:sha256 \
    -pkeyopt rsa_mgf1_md:sha256 > \
  my-secret-password.enc
  
## base64 encode the ciphertext  
$ openssl base64 -in my-secret-password.enc
M7nUoba9EGVTu2LjNjBKGdGVBYjyS/i/AY+4yQMQF0Qf/RfUfX31Jw6+VO9OuThq
ylu/7ihX9XD4bM7yYdXnMv9p1OHQUlorSBSbb/J6n1W9UJhcp6um8Tw8/Isx4f75
4PskYS6f8Y2ItliGt1/A9iR5BTgGtJBwOxMlgoX2Ggq+Nh4E5SbdoaE5o6CO1nBx
eIPsPEebQ6qC4JehQM3IGuV/lrm58+hZhaXAqNzX1cEYyAt5GYqJIVCiI585SUYs
wRToGyTgaN+zthF0HP9IWlR4Am4LmJ/1OcePTnYw11CkU8wNRbDzVAzogwNH+rXr
LTmf7hxVjBm6bBSVSNFcBKAXFlllubSfIeZ5hgzGqn54OmSf6odO12L5JxllddHc
yAd54vWKs2kJtnsKV2V4ZdkI0w6y1TeI67baFZDNGo6qsCpFMPnvv7d46Pg2VOp1
J6Ivner0NnNHE4MzNmpZRk8WXMwqq4P/gTiT7F/aCX6oFCUQ4AWPQhJYh2dkcOmL
IP+47Veb10aFn61F1CJwpmOOiGNXKdDT1vK8CMnnwhm825K0q/q9Zqpzc1+1ae1z
mSqol1zCoa88CuSN6nTLQlVnN/dzfrGbc0boJPaM0iGhHtSzHk4SWg84LhiJB1q9
A9XFJmOVdkvRY9nnz/iVLAdd0Q3vFtLqCdUYsNN2yh4=

## optionally calculate the CRC32 of the ciphertext
$ go get github.com/binxio/crc32 
$ $GOPATH/bin/crc32 -polynomial castagnoli < my-secret-password.enc
12c59e54
```

Finally, reference the encrypted ciphertext in your resource definitions:

```hcl
data "google_kms_secret_asymmetric" "sql_user_password" {
  crypto_key_version = data.google_kms_crypto_key_version.my_crypto_key.id
  crc32              = "12c59e54"
  ciphertext         = <<EOT
    M7nUoba9EGVTu2LjNjBKGdGVBYjyS/i/AY+4yQMQF0Qf/RfUfX31Jw6+VO9OuThq
    ylu/7ihX9XD4bM7yYdXnMv9p1OHQUlorSBSbb/J6n1W9UJhcp6um8Tw8/Isx4f75
    4PskYS6f8Y2ItliGt1/A9iR5BTgGtJBwOxMlgoX2Ggq+Nh4E5SbdoaE5o6CO1nBx
    eIPsPEebQ6qC4JehQM3IGuV/lrm58+hZhaXAqNzX1cEYyAt5GYqJIVCiI585SUYs
    wRToGyTgaN+zthF0HP9IWlR4Am4LmJ/1OcePTnYw11CkU8wNRbDzVAzogwNH+rXr
    LTmf7hxVjBm6bBSVSNFcBKAXFlllubSfIeZ5hgzGqn54OmSf6odO12L5JxllddHc
    yAd54vWKs2kJtnsKV2V4ZdkI0w6y1TeI67baFZDNGo6qsCpFMPnvv7d46Pg2VOp1
    J6Ivner0NnNHE4MzNmpZRk8WXMwqq4P/gTiT7F/aCX6oFCUQ4AWPQhJYh2dkcOmL
    IP+47Veb10aFn61F1CJwpmOOiGNXKdDT1vK8CMnnwhm825K0q/q9Zqpzc1+1ae1z
    mSqol1zCoa88CuSN6nTLQlVnN/dzfrGbc0boJPaM0iGhHtSzHk4SWg84LhiJB1q9
    A9XFJmOVdkvRY9nnz/iVLAdd0Q3vFtLqCdUYsNN2yh4=
  EOT
}

resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "main" {
  name             = "main-instance-${random_id.db_name_suffix.hex}"
  database_version = "MYSQL_5_7"
  
  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "users" {
  name     = "me"
  instance = google_sql_database_instance.main.name
  host     = "me.com"
  password = data.google_kms_secret.sql_user_password.plaintext
}
```

This will result in a Cloud SQL user being created with password `my-secret-password`.

## Argument Reference

The following arguments are supported:

* `ciphertext` (Required) - The ciphertext to be decrypted, encoded in base64
* `crypto_key_version` (Required) - The id of the CryptoKey version that will be used to
  decrypt the provided ciphertext. This is represented by the format
  `projects/{project}/locations/{location}/keyRings/{keyring}/cryptoKeys/{key}/cryptoKeyVersions/{version}`.
* `crc32` (Optional) - The crc32 checksum of the `ciphertext` in hexadecimal notation. If not specified, it will be computed.

## Attributes Reference

The following attribute is exported:

* `plaintext` - Contains the result of decrypting the provided ciphertext.
* `crc32` - Contains the crc32 checksum of the provided ciphertext.
