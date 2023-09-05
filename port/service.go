package port

type Service interface {
	User() UserService
	Article() ArticleService
}
