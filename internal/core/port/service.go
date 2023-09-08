package port

import "github.com/labasubagia/realworld-backend/internal/core/util/token"

type Service interface {
	TokenMaker() token.Maker
	User() UserService
	Article() ArticleService
}
