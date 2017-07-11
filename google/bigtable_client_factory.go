package google

import (
	"context"

	"cloud.google.com/go/bigtable"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

type BigtableClientFactory struct {
	UserAgent   string
	TokenSource oauth2.TokenSource
}

func (s BigtableClientFactory) NewInstanceAdminClient(project string) (*bigtable.InstanceAdminClient, error) {
	return bigtable.NewInstanceAdminClient(context.Background(), project, option.WithTokenSource(s.TokenSource), option.WithUserAgent(s.UserAgent))
}

func (s BigtableClientFactory) NewAdminClient(project, instance string) (*bigtable.AdminClient, error) {
	return bigtable.NewAdminClient(context.Background(), project, instance, option.WithTokenSource(s.TokenSource), option.WithUserAgent(s.UserAgent))
}
