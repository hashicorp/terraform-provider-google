package google

import (
	"context"
	"os"

	"cloud.google.com/go/bigtable"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

type BigtableClientFactory struct {
	UserAgent           string
	TokenSource         oauth2.TokenSource
	BillingProject      string
	UserProjectOverride bool
}

func (s BigtableClientFactory) NewInstanceAdminClient(project string) (*bigtable.InstanceAdminClient, error) {
	var opts []option.ClientOption
	if requestReason := os.Getenv("CLOUDSDK_CORE_REQUEST_REASON"); requestReason != "" {
		opts = append(opts, option.WithRequestReason(requestReason))
	}

	if s.UserProjectOverride && s.BillingProject != "" {
		opts = append(opts, option.WithQuotaProject(s.BillingProject))
	}

	opts = append(opts, option.WithTokenSource(s.TokenSource), option.WithUserAgent(s.UserAgent))
	return bigtable.NewInstanceAdminClient(context.Background(), project, opts...)
}

func (s BigtableClientFactory) NewAdminClient(project, instance string) (*bigtable.AdminClient, error) {
	var opts []option.ClientOption
	if requestReason := os.Getenv("CLOUDSDK_CORE_REQUEST_REASON"); requestReason != "" {
		opts = append(opts, option.WithRequestReason(requestReason))
	}

	if s.UserProjectOverride && s.BillingProject != "" {
		opts = append(opts, option.WithQuotaProject(s.BillingProject))
	}

	opts = append(opts, option.WithTokenSource(s.TokenSource), option.WithUserAgent(s.UserAgent))
	return bigtable.NewAdminClient(context.Background(), project, instance, opts...)
}
