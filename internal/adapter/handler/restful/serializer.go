package restful

import "github.com/labasubagia/realworld-backend/internal/core/domain"

type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

type ProfileResponse struct {
	Profile Profile `json:"profile"`
}

func serializeProfile(arg domain.User) Profile {
	return Profile{
		Username:  arg.Username,
		Image:     arg.Image,
		Bio:       arg.Bio,
		Following: arg.IsFollowed,
	}
}

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

type UserResponse struct {
	User User `json:"user"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}

func serializeUser(arg domain.User) User {
	return User{
		Email:    arg.Email,
		Username: arg.Username,
		Bio:      arg.Bio,
		Image:    arg.Email,
		Token:    arg.Token,
	}
}

type Article struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Body           string   `json:"body"`
	TagList        []string `json:"tagList"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount int      `json:"favoritesCount"`
	Author         Profile  `json:"author"`
}

type ArticleResponse struct {
	Article Article `json:"article"`
}

type ArticlesResponse struct {
	Articles []Article `json:"articles"`
	Count    int       `json:"articlesCount"`
}

func serializeArticle(arg domain.Article) Article {
	tags := []string{}
	if len(arg.TagNames) > 0 {
		tags = arg.TagNames
	}
	return Article{
		Slug:           arg.Slug,
		Title:          arg.Title,
		Description:    arg.Description,
		Body:           arg.Body,
		TagList:        tags,
		Favorited:      arg.IsFavorite,
		FavoritesCount: arg.FavoriteCount,
		Author:         serializeProfile(arg.Author),
		CreatedAt:      timeString(arg.CreatedAt),
		UpdatedAt:      timeString(arg.UpdatedAt),
	}
}

type Comment struct {
	ID        int64   `json:"id"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Body      string  `json:"body"`
	Author    Profile `json:"author"`
}

type CommentResponse struct {
	Comment Comment `json:"comment"`
}

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}

func serializeComment(arg domain.Comment) Comment {
	return Comment{
		ID:        int64(arg.ID),
		Body:      arg.Body,
		Author:    serializeProfile(arg.Author),
		CreatedAt: timeString(arg.CreatedAt),
		UpdatedAt: timeString(arg.UpdatedAt),
	}
}
