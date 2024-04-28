package utils

import "github.com/rs/xid"

func NewSessionID() string {
	return xid.New().String()
}

func NewWorkspaceID() string {
	return "ws_" + xid.New().String()
}
