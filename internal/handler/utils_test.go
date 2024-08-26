package handler

import "testing"

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
		s := ParseSubdomain(table.host)
		if table.subdomain != s {
			t.Fatalf("expected %s to equal %s", s, table.subdomain)
		}
	}
}
