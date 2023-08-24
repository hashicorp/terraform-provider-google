// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
)

const (
	jwtTestSubject          = "custom-subject"
	jwtTestFoo              = "bar"
	jwtTestComplexFooNested = "baz"
	jwtTestExpiresIn        = 60
)

type jwtTestPayload struct {
	Subject string `json:"sub"`

	Foo string `json:"foo"`

	ComplexFoo struct {
		Nested string `json:"nested"`
	} `json:"complexFoo"`

	Expiration int64 `json:"exp"`
}

func testAccCheckServiceAccountJwtValue(name, audience string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ms := s.RootModule()

		rs, ok := ms.Resources[name]

		if !ok {
			return fmt.Errorf("can't find %s in state", name)
		}

		jwtString, ok := rs.Primary.Attributes["jwt"]

		if !ok {
			return fmt.Errorf("jwt not found")
		}

		jwtParts := strings.Split(jwtString, ".")

		if len(jwtParts) != 3 {
			return errors.New("jwt does not appear well-formed")
		}

		decoded, err := base64.RawURLEncoding.DecodeString(jwtParts[1])

		if err != nil {
			return fmt.Errorf("could not base64 decode jwt body: %w", err)
		}

		var payload jwtTestPayload

		err = json.NewDecoder(bytes.NewBuffer(decoded)).Decode(&payload)

		if err != nil {
			return fmt.Errorf("could not decode jwt payload: %w", err)
		}

		if payload.Subject != jwtTestSubject {
			return fmt.Errorf("invalid 'sub', expected '%s', got '%s'", jwtTestSubject, payload.Subject)
		}

		if payload.Foo != jwtTestFoo {
			return fmt.Errorf("invalid 'foo', expected '%s', got '%s'", jwtTestFoo, payload.Foo)
		}

		if payload.ComplexFoo.Nested != jwtTestComplexFooNested {
			return fmt.Errorf("invalid 'foo', expected '%s', got '%s'", jwtTestComplexFooNested, payload.ComplexFoo.Nested)
		}

		expectedExpiration := resourcemanager.DataSourceGoogleServiceAccountJwtNow().Add(jwtTestExpiresIn * time.Second).Unix()

		if payload.Expiration != expectedExpiration {
			return fmt.Errorf("invalid 'exp', expected '%d', got '%d'", expectedExpiration, payload.Expiration)
		}

		return nil
	}
}

func TestAccDataSourceGoogleServiceAccountJwt(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_service_account_jwt.default"
	serviceAccount := envvar.GetTestServiceAccountFromEnv(t)
	targetServiceAccountEmail := acctest.BootstrapServiceAccount(t, envvar.GetTestProjectFromEnv(), serviceAccount)

	staticTime := time.Now()

	// Override the current time with one that is set to a static value, to compare against later.
	resourcemanager.DataSourceGoogleServiceAccountJwtNow = func() time.Time {
		return staticTime
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleServiceAccountJwt(targetServiceAccountEmail),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceAccountJwtValue(resourceName, targetAudience),
				),
			},
			{
				PreConfig: func() {
					// Bump the hardcoded time to ensure terraform responds well to the JWT expiration changing.
					staticTime = time.Now().Add(10 * time.Second)
				},
				Config: testAccCheckGoogleServiceAccountJwt(targetServiceAccountEmail),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceAccountJwtValue(resourceName, targetAudience),
				),
			},
		},
	})
}

func testAccCheckGoogleServiceAccountJwt(targetServiceAccount string) string {
	return fmt.Sprintf(`
data "google_service_account_jwt" "default" {
	target_service_account = "%s"

    payload = jsonencode({
      sub: "%s",
      foo: "%s",
      complexFoo: {
        nested: "%s"
      }
    })

    expires_in = %d
}
`, targetServiceAccount, jwtTestSubject, jwtTestFoo, jwtTestComplexFooNested, jwtTestExpiresIn)
}
