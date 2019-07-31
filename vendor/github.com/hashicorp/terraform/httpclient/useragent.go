package httpclient

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/version"
)

const userAgentFormat = "Terraform/%s"
const uaEnvVar = "TF_APPEND_USER_AGENT"

// Deprecated: Use UserAgent(version) instead
func UserAgentString() string {
	ua := fmt.Sprintf(userAgentFormat, version.Version)

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}

type userAgentRoundTripper struct {
	inner     http.RoundTripper
	userAgent string
}

func (rt *userAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", rt.userAgent)
	}
	return rt.inner.RoundTrip(req)
}

func UserAgent(version string) *userAgent {
	return newUserAgent([]*UserAgentProduct{
		{"HashiCorp", "1.0", ""},
		{"Terraform", version, "+https://www.terraform.io"},
	})
}

type UserAgentProduct struct {
	Name    string
	Version string
	Comment string
}

func (uap *UserAgentProduct) String() string {
	var b strings.Builder
	b.WriteString(uap.Name)
	if uap.Version != "" {
		b.WriteString(fmt.Sprintf("/%s", uap.Version))
	}
	if uap.Comment != "" {
		b.WriteString(fmt.Sprintf(" (%s)", uap.Comment))
	}
	return b.String()
}

func (uap *UserAgentProduct) Equal(p *UserAgentProduct) bool {
	if uap.Name == p.Name &&
		uap.Version == p.Version &&
		uap.Comment == p.Comment {
		return true
	}
	return false
}

type userAgent struct {
	products []*UserAgentProduct
}

func (ua *userAgent) Products() []*UserAgentProduct {
	return ua.products
}

func (ua *userAgent) String() string {
	var b strings.Builder
	for i, p := range ua.products {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(p.String())
	}

	return b.String()
}

func (ua *userAgent) Append(uap ...*UserAgentProduct) *userAgent {
	ua.products = append(ua.products, uap...)
	return ua
}

func (ua *userAgent) AppendString(uap ...string) *userAgent {
	for _, uaString := range uap {
		products, err := ParseUserAgentString(uaString)
		if err != nil {
			log.Printf("[WARN] Unable to append User-Agent string %q: %s",
				uaString, err)
			continue
		}

		ua.products = append(ua.products, products...)
	}
	return ua
}

func (ua *userAgent) Equal(userAgent *userAgent) bool {
	if len(ua.products) != len(userAgent.products) {
		return false
	}

	for i, p := range ua.products {
		if !p.Equal(userAgent.products[i]) {
			return false
		}
	}

	return true
}

func newUserAgent(products []*UserAgentProduct) *userAgent {
	ua := &userAgent{products}

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			parsedUAs, err := ParseUserAgentString(add)
			if err != nil {
				log.Printf("[WARN] Unable to parse User-Agent string %q: %s",
					add, err)
				return ua
			}
			ua = ua.Append(parsedUAs...)
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}

func ParseUserAgentString(uaString string) ([]*UserAgentProduct, error) {
	products := make([]*UserAgentProduct, 0)

	// Parse "Product/version (comment)"
	re := regexp.MustCompile(`([^/]+)/([^\s]+)(\s\([^\)]+\))?`)
	matches := re.FindAllStringSubmatch(uaString, -1)

	if len(matches) == 0 {
		return nil, fmt.Errorf("Invalid User-Agent format: %q", uaString)
	}

	for _, match := range matches {
		products = append(products, &UserAgentProduct{
			Name:    strings.TrimSpace(match[1]),
			Version: match[2],
			Comment: strings.Trim(strings.TrimSpace(match[3]), "()"),
		})
	}

	return products, nil
}
