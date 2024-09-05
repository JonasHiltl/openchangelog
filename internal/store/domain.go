package store

import (
	"database/sql"
	"errors"
	"net/url"
	"strings"

	"github.com/guregu/null/v5"
	"github.com/jonashiltl/openchangelog/internal/errs"
)

type Domain null.String

func (d Domain) ToNullString() null.String {
	return null.NewString(d.String, d.Valid)
}

func DomainFromSQL(d sql.NullString) Domain {
	return Domain(null.NewString(d.String, d.Valid))
}

func ParseDomain(d null.String) (Domain, error) {
	if d.Valid && d.String != "" {
		domain := d.String
		if !strings.Contains(domain, "://") {
			domain = "http://" + domain // Add a default scheme, else host is empty
		}

		parsedUrl, err := url.Parse(domain)
		if err != nil {
			return Domain{}, errs.NewBadRequest(errors.New("domain not valid"))
		}

		d.String = parsedUrl.Host
	}
	return Domain(d), nil
}
