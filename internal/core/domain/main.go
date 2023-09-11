package domain

import (
	"github.com/oklog/ulid/v2"
)

type ID string

func NewID() ID {
	return ID(ulid.Make().String())
}

func RandomID() ID {
	return NewID()
}

func ParseID(value string) (ID, error) {
	id, err := ulid.Parse(value)
	if err != nil {
		return ID(id.String()), nil
	}
	return ID(ulid.ULID{}.String()), err
}
