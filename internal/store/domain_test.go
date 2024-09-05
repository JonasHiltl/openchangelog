package store

import (
	"testing"

	"github.com/guregu/null/v5"
)

func TestParseDomain(t *testing.T) {
	tables := []struct {
		domain    string
		expected  null.String
		expectErr bool
	}{
		{
			domain:   "openchangelog",
			expected: null.NewString("openchangelog", true),
		},
		{
			domain:   "openchangelog.com",
			expected: null.NewString("openchangelog.com", true),
		},
		{
			domain:   "changelog.openchangelog.com",
			expected: null.NewString("changelog.openchangelog.com", true),
		},
		{
			domain:   "changelog.openchangelog.com:3000",
			expected: null.NewString("changelog.openchangelog.com:3000", true),
		},
		{
			domain:   "https://changelog.openchangelog.com",
			expected: null.NewString("changelog.openchangelog.com", true),
		},
		{
			domain:   "http://changelog.openchangelog.com",
			expected: null.NewString("changelog.openchangelog.com", true),
		},
	}

	for _, table := range tables {
		d, err := ParseDomain(null.NewString(table.domain, table.domain != ""))
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
