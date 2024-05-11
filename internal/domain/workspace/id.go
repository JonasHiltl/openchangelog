package workspace

import (
	"errors"
	"strings"

	"github.com/jonashiltl/openchangelog/internal/domain"
	"github.com/rs/xid"
)

type ID string

func NewID() ID {
	return ID("ws_" + xid.New().String())
}

func ParseID(id string) (ID, error) {
	parts := strings.Split(id, "_")
	if len(parts) != 2 {
		return "", domain.NewError(domain.ErrBadRequest, errors.New("wrong workspace id format"))
	}
	if parts[0] != "ws" {
		return "", domain.NewError(domain.ErrBadRequest, errors.New("invalid workspace id prefix"))
	}
	_, err := xid.FromString(parts[1])
	if err != nil {
		return "", domain.NewError(domain.ErrBadRequest, err)
	}
	return ID(id), nil
}

func (i ID) String() string {
	return string(i)
}
