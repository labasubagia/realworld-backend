package port

import (
	"context"
	"errors"
)

var (
	ErrIsolation  = errors.New("Error Isolation")
	ErrUniqueKey  = errors.New("Error Unique Key")
	ErrForeignKey = errors.New("Error Foreign Key")
)

type RepositoryAtomicCallback func(r Repository) error

type Repository interface {
	Atomic(context.Context, RepositoryAtomicCallback) error
	User() UserRepository
	Article() ArticleRepository
}
