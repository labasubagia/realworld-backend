package domain

import (
	"github.com/oklog/ulid/v2"
)

type ID string

func (id ID) String() string {
	return string(id)
}

func NewID() ID {
	return ID(ulid.Make().String())
}

func RandomID() ID {
	return NewID()
}

func ParseID(value string) (ID, error) {
	id, err := ulid.Parse(value)
	if err != nil {
		return ID(ulid.ULID{}.String()), err
	}
	return ID(id.String()), nil
}
