package port

import (
	"context"
)

type RepositoryAtomicCallback func(r Repository) error

type Repository interface {
	Atomic(context.Context, RepositoryAtomicCallback) error
	User() UserRepository
	Article() ArticleRepository
}
