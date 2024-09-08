package store

import (
	"testing"

	"github.com/jonashiltl/openchangelog/apitypes"
)

func TestParseDomain(t *testing.T) {
	tables := []struct {
		host      string
		expected  apitypes.NullString
		expectErr bool
	}{
		{
			host:     "openchangelog",
			expected: apitypes.NewString("openchangelog"),
		},
		{
			host:     "openchangelog.com",
			expected: apitypes.NewString("openchangelog.com"),
		},
		{
			host:     "changelog.openchangelog.com",
			expected: apitypes.NewString("changelog.openchangelog.com"),
		},
		{
			host:     "changelog.openchangelog.com:3000",
			expected: apitypes.NewString("changelog.openchangelog.com:3000"),
		},
		{
			host:     "https://changelog.openchangelog.com",
			expected: apitypes.NewString("changelog.openchangelog.com"),
		},
		{
			host:     "http://changelog.openchangelog.com",
			expected: apitypes.NewString("changelog.openchangelog.com"),
		},
		{
			host:      "https://test com",
			expectErr: true,
		},
	}

	for _, table := range tables {
		t.Run(table.host, func(t *testing.T) {
			d, err := ParseDomain(table.host)
			if table.expectErr && err == nil {
				t.Error("expected to error but no error returned")
			}
			if d.String() != table.expected.String() {
				t.Errorf("expected %s to equal %s", d.String(), table.expected.String())
			}
		})
	}
}

func TestParseSubdomain(t *testing.T) {
	tables := []struct {
		host      string
		subdomain string
	}{
		{
			host:      "tenant.openchangelog.com",
			subdomain: "tenant",
		},
		{
			host:      "tenant-2.openchangelog.com",
			subdomain: "tenant-2",
		},
		{
			host:      "https://changelog.test.com",
			subdomain: "changelog",
		},
		{
			host:      "https://changelog.test.com:6001",
			subdomain: "changelog",
		},
		{
			host:      "openchangelog.com",
			subdomain: "",
		},
		{
			host:      "www.openchangelog.com",
			subdomain: "",
		},
		{
			host:      "",
			subdomain: "",
		},
		{
			host:      ".",
			subdomain: "",
		},
		{
			host:      ".com",
			subdomain: "",
		},
	}

	for _, table := range tables {
		t.Run(table.host, func(t *testing.T) {
			s, _ := SubdomainFromHost(table.host)
			if table.subdomain != s.String() {
				t.Fatalf("expected %s to equal %s", s, table.subdomain)
			}
		})
	}
}
