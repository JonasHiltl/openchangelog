package store

import (
	"testing"

	"github.com/guregu/null/v5"
)

func TestParseDomain(t *testing.T) {
	tables := []struct {
		host      string
		expected  null.String
		expectErr bool
	}{
		{
			host:     "openchangelog",
			expected: null.NewString("openchangelog", true),
		},
		{
			host:     "openchangelog.com",
			expected: null.NewString("openchangelog.com", true),
		},
		{
			host:     "changelog.openchangelog.com",
			expected: null.NewString("changelog.openchangelog.com", true),
		},
		{
			host:     "changelog.openchangelog.com:3000",
			expected: null.NewString("changelog.openchangelog.com:3000", true),
		},
		{
			host:     "https://changelog.openchangelog.com",
			expected: null.NewString("changelog.openchangelog.com", true),
		},
		{
			host:     "http://changelog.openchangelog.com",
			expected: null.NewString("changelog.openchangelog.com", true),
		},
	}

	for _, table := range tables {
		d, err := ParseDomain(null.NewString(table.host, table.host != ""))
		if table.expectErr && err == nil {
			t.Error("expected to error but no error returned")
		}
		if d.String != table.expected.String {
			t.Errorf("expected %s to equal %s", d.String, table.expected.String)
		}
		if d.Valid != table.expected.Valid {
			t.Errorf("expected valid %t to equal%t", d.Valid, table.expected.Valid)
		}
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
		s, _ := SubdomainFromHost(table.host)
		if table.subdomain != s.String() {
			t.Fatalf("expected %s to equal %s", s, table.subdomain)
		}
	}
}
