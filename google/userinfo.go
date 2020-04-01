package google

import (
	"fmt"
)

func GetCurrentUserEmail(config *Config) (string, error) {
	// See https://github.com/golang/oauth2/issues/306 for a recommendation to do this from a Go maintainer
	// URL retrieved from https://accounts.google.com/.well-known/openid-configuration
	res, err := sendRequest(config, "GET", "", "https://openidconnect.googleapis.com/v1/userinfo", nil)
	if err != nil {
		return "", fmt.Errorf("error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope? error: %s", err)
	}
	return res["email"].(string), nil
}
