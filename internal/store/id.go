package store

import (
	"errors"
	"strings"

	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/rs/xid"
)

const (
	wid_prefix   = "ws"
	cid_prefix   = "cl"
	ghid_prefix  = "gh"
	id_separator = "_"
)

type WorkspaceID string

func NewWID() WorkspaceID {
	return WorkspaceID(wid_prefix + id_separator + xid.New().String())
}

var errWSFormat = errs.NewError(errs.ErrBadRequest, errors.New("wrong workspace id format"))

func ParseWID(id string) (WorkspaceID, error) {
	if id == WS_DEFAULT_ID.String() {
		return WS_DEFAULT_ID, nil
	}

	parts := strings.Split(id, id_separator)
	if len(parts) != 2 {
		return "", errWSFormat
	}
	if parts[0] != wid_prefix {
		return "", errs.NewError(errs.ErrBadRequest, errors.New("invalid workspace id prefix"))
	}
	_, err := xid.FromString(parts[1])
	if err != nil {
		return "", errWSFormat
	}
	return WorkspaceID(id), nil
}

func (i WorkspaceID) String() string {
	return string(i)
}

type ChangelogID string

func NewCID() ChangelogID {
	return ChangelogID(cid_prefix + id_separator + xid.New().String())
}

var errCLFormat = errs.NewError(errs.ErrBadRequest, errors.New("wrong changelog id format"))

func ParseCID(id string) (ChangelogID, error) {
	if id == CL_DEFAULT_ID.String() {
		return CL_DEFAULT_ID, nil
	}
	parts := strings.Split(id, id_separator)
	if len(parts) != 2 {
		return "", errCLFormat
	}
	if parts[0] != cid_prefix {
		return "", errs.NewError(errs.ErrBadRequest, errors.New("invalid changelog id prefix"))
	}
	_, err := xid.FromString(parts[1])
	if err != nil {
		return "", errCLFormat
	}
	return ChangelogID(id), nil
}

func (i ChangelogID) String() string {
	return string(i)
}

type GHSourceID string

func NewGHID() GHSourceID {
	return GHSourceID(ghid_prefix + id_separator + xid.New().String())
}

var errGHFormat = errs.NewError(errs.ErrBadRequest, errors.New("wrong github source id format"))

func ParseGHID(id string) (GHSourceID, error) {
	if id == GH_DEFAULT_ID.String() {
		return GH_DEFAULT_ID, nil
	}
	parts := strings.Split(id, id_separator)
	if len(parts) != 2 {
		return "", errGHFormat
	}
	if parts[0] != ghid_prefix {
		return "", errs.NewError(errs.ErrBadRequest, errors.New("invalid gh source id prefix"))
	}
	_, err := xid.FromString(parts[1])
	if err != nil {
		return "", errGHFormat
	}
	return GHSourceID(id), nil
}

func IsGHID(id string) bool {
	return strings.HasPrefix(id, ghid_prefix+id_separator)
}

func (i GHSourceID) String() string {
	return string(i)
}
