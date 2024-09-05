package store

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/guregu/null/v5"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"golang.org/x/exp/rand"
)

type Domain null.String

func (d Domain) ToNullString() null.String {
	return null.NewString(d.String, d.Valid)
}

func DomainFromSQL(d sql.NullString) Domain {
	return Domain(null.NewString(d.String, d.Valid))
}

// strips everything from d except the host
func ParseDomain(d null.String) (Domain, error) {
	if d.Valid && d.String != "" {
		if !strings.Contains(d.String, "://") {
			d.String = "http://" + d.String // Add a default scheme, else host is empty
		}

		parsedUrl, err := url.Parse(d.String)
		if err != nil {
			return Domain{}, errs.NewBadRequest(errors.New("domain not valid"))
		}

		d.String = parsedUrl.Host
	}
	return Domain(d), nil
}

type Subdomain string

func (s Subdomain) String() string {
	return string(s)
}

func (s Subdomain) NullString() null.String {
	return null.NewString(s.String(), s.String() != "")
}

func NewSubdomain(workspaceName string) Subdomain {
	wsName := strings.ReplaceAll(strings.ToLower(workspaceName), " ", "-")
	rnd := rand.Intn(100000)

	return Subdomain(fmt.Sprintf("%s-%d", wsName, rnd))
}

var subdomainRegex = regexp.MustCompile("^[a-z0-9-]*$")

func ParseSubdomain(subdomain string) Subdomain {
	return Subdomain(subdomain)
}

// Returns the subdomain from the host.
// Returns an error if the host doesn't have a subdomain
func SubdomainFromHost(host string) (Subdomain, error) {
	// add scheme, else parsed url won't include host
	if !strings.Contains(host, "://") {
		host = "https://" + host
	}

	parsedURL, err := url.Parse(host)
	if err != nil {
		return "", errs.NewBadRequest(errors.New("invalid URL"))
	}

	// Extract the host from the parsed URL
	host = parsedURL.Host
	parts := strings.Split(host, ".")
	if parts[0] == "www" {
		parts = parts[1:]
	}

	// subdomain exists, e.g. tenant.openchangelog.com
	if len(parts) > 2 {
		if !subdomainRegex.MatchString(parts[0]) {
			return "", errs.NewBadRequest(errors.New("subdomain not valid"))
		}
		return Subdomain(parts[0]), nil
	}
	return "", errs.NewBadRequest(errors.New("host has no subdomain"))
}
