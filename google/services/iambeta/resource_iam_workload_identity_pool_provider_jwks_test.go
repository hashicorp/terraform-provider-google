// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package iambeta

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestNormalizeJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "compact JSON",
			input:    `{"keys":[{"kty":"RSA","alg":"RS256"}]}`,
			expected: `{"keys":[{"alg":"RS256","kty":"RSA"}]}`,
			wantErr:  false,
		},
		{
			name: "formatted JSON with whitespace",
			input: `{
  "keys": [
    {
      "kty": "RSA",
      "alg": "RS256"
    }
  ]
}`,
			expected: `{"keys":[{"alg":"RS256","kty":"RSA"}]}`,
			wantErr:  false,
		},
		{
			name: "JWKS with multiple keys",
			input: `{
  "keys": [
    {
      "kty": "RSA",
      "alg": "RS256",
      "kid": "key1",
      "use": "sig"
    },
    {
      "kty": "EC",
      "alg": "ES256",
      "kid": "key2",
      "use": "sig"
    }
  ]
}`,
			expected: `{"keys":[{"alg":"RS256","kid":"key1","kty":"RSA","use":"sig"},{"alg":"ES256","kid":"key2","kty":"EC","use":"sig"}]}`,
			wantErr:  false,
		},
		{
			name:     "invalid JSON",
			input:    `{"keys": [{"kty": "RSA", "alg": "RS256"}`,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeJSON(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("normalizeJSON() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("normalizeJSON() unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("normalizeJSON() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestJwksJsonDiffSuppressFunc(t *testing.T) {
	tests := []struct {
		name     string
		old      string
		new      string
		expected bool
	}{
		{
			name:     "identical compact JSON",
			old:      `{"keys":[{"kty":"RSA","alg":"RS256"}]}`,
			new:      `{"keys":[{"kty":"RSA","alg":"RS256"}]}`,
			expected: true,
		},
		{
			name: "equivalent JSON with different whitespace",
			old:  `{"keys":[{"kty":"RSA","alg":"RS256"}]}`,
			new: `{
  "keys": [
    {
      "kty": "RSA",
      "alg": "RS256"
    }
  ]
}`,
			expected: true,
		},
		{
			name:     "equivalent JSON with different key order",
			old:      `{"keys":[{"kty":"RSA","alg":"RS256"}]}`,
			new:      `{"keys":[{"alg":"RS256","kty":"RSA"}]}`,
			expected: true,
		},
		{
			name:     "different JSON content",
			old:      `{"keys":[{"kty":"RSA","alg":"RS256"}]}`,
			new:      `{"keys":[{"kty":"EC","alg":"ES256"}]}`,
			expected: false,
		},
		{
			name: "JWKS with GCP-style formatting vs user input",
			old:  `{"keys":[{"alg":"RS256","e":"AQAB","kid":"sif0AR","kty":"RSA","n":"ylH1Chl1tpfti3lh51E1g5dPogzXDaQseqjsefGLknaNl5W6Wd4frBhHyE2t41Q5zgz_Ll0-NvWm0FlaG6brhrN9QZu6sJP1bM8WPfJVPgXOanxi7d7TXCkeNubGeiLTf5R3UXtS9Lm_guemU7MxDjDTelxnlgGCihOVTcL526suNJUdfXtpwUsvdU6_ZnAp9IpsuYjCtwPm9hPumlcZGMbxstdh07O4y4O90cVQClJOKSGQjAUCKJWXIQ0cqffGS_HuS_725CPzQ85SzYZzaNpgfhAER7kx_9P16ARM3BJz0PI5fe2hECE61J4GYU_BY43sxDfs7HyJpEXKLU9eWw","use":"sig"}]}`,
			new: `{
  "keys": [
    {
      "kty": "RSA",
      "alg": "RS256",
      "kid": "sif0AR",
      "use": "sig",
      "e": "AQAB",
      "n": "ylH1Chl1tpfti3lh51E1g5dPogzXDaQseqjsefGLknaNl5W6Wd4frBhHyE2t41Q5zgz_Ll0-NvWm0FlaG6brhrN9QZu6sJP1bM8WPfJVPgXOanxi7d7TXCkeNubGeiLTf5R3UXtS9Lm_guemU7MxDjDTelxnlgGCihOVTcL526suNJUdfXtpwUsvdU6_ZnAp9IpsuYjCtwPm9hPumlcZGMbxstdh07O4y4O90cVQClJOKSGQjAUCKJWXIQ0cqffGS_HuS_725CPzQ85SzYZzaNpgfhAER7kx_9P16ARM3BJz0PI5fe2hECE61J4GYU_BY43sxDfs7HyJpEXKLU9eWw"
    }
  ]
}`,
			expected: true,
		},
		{
			name:     "empty strings",
			old:      "",
			new:      "",
			expected: true,
		},
		{
			name:     "old empty, new has content",
			old:      "",
			new:      `{"keys":[]}`,
			expected: false,
		},
		{
			name:     "invalid JSON in old",
			old:      `{"keys": [{"kty": "RSA", "alg": "RS256"}`,
			new:      `{"keys":[{"kty":"RSA","alg":"RS256"}]}`,
			expected: false,
		},
		{
			name:     "invalid JSON in new",
			old:      `{"keys":[{"kty":"RSA","alg":"RS256"}]}`,
			new:      `{"keys": [{"kty": "RSA", "alg": "RS256"}`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := jwksJsonDiffSuppressFunc("oidc.0.jwks_json", tt.old, tt.new, &schema.ResourceData{})
			if result != tt.expected {
				t.Errorf("jwksJsonDiffSuppressFunc() = %v, want %v", result, tt.expected)
			}
		})
	}
}
