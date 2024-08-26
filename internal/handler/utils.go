package handler

import "strings"

const (
	WS_ID_QUERY = "wid"
	CL_ID_QUERY = "cid"
)

func ParseSubdomain(host string) string {
	// Remove port if present
	host = strings.Split(host, ":")[0]
	parts := strings.Split(host, ".")
	if parts[0] == "www" {
		parts = parts[1:]
	}

	// subdomain exists, e.g. tenant.openchangelog.com
	if len(parts) > 2 {
		return parts[0]
	}
	return ""
}
