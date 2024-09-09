package store

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/jonashiltl/openchangelog/apitypes"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"golang.org/x/exp/rand"
)

type Domain apitypes.NullString

func (d Domain) String() string {
	return d.NullString().String()
}

func (d Domain) NullString() apitypes.NullString {
	return apitypes.NullString(d)
}

var errInvalidDomain = errs.NewBadRequest(errors.New("domain is not valid"))

// strips everything from domain except the host
func ParseDomain(domain string) (Domain, error) {
	if !strings.Contains(domain, ".") {
		return Domain{}, errInvalidDomain
	}
	if !strings.Contains(domain, "://") {
		domain = "http://" + domain // Add a default scheme, else host is empty
	}

	parsedUrl, err := url.Parse(domain)
	if err != nil {
		return Domain{}, errInvalidDomain
	}

	domain = parsedUrl.Host
	return Domain(apitypes.NewString(domain)), nil
}

// if ns is valid, it parses the domain by stripping everything except the host from the string.
func ParseDomainNullString(ns apitypes.NullString) (Domain, error) {
	if !ns.IsValid() {
		return Domain(ns), nil
	}
	return ParseDomain(ns.String())
}

type Subdomain string

func (s Subdomain) String() string {
	return string(s)
}

func (s Subdomain) NullString() apitypes.NullString {
	return apitypes.NewString(s.String())
}

func NewSubdomain(workspaceName string) Subdomain {
	wsName := strings.ReplaceAll(strings.ToLower(workspaceName), " ", "-")
	rnd := rand.Intn(100000)

	return Subdomain(fmt.Sprintf("%s-%d", wsName, rnd))
}

var subdomainRegex = regexp.MustCompile("^[a-z0-9-]*$")

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
