package handler

import (
	"time"

	"github.com/kenshin579/analyzing-golang-echo-realworld-example-app/model"
	"github.com/kenshin579/analyzing-golang-echo-realworld-example-app/user"
	"github.com/kenshin579/analyzing-golang-echo-realworld-example-app/utils"
	"github.com/labstack/echo/v4"
)

type userResponse struct {
	User struct {
		Username string  `json:"username"`
		Email    string  `json:"email"`
		Bio      *string `json:"bio"`
		Image    *string `json:"image"`
		Token    string  `json:"token"`
	} `json:"user"`
}

func newUserResponse(user *model.User) *userResponse {
	userResponse := new(userResponse)
	userResponse.User.Username = user.Username
	userResponse.User.Email = user.Email
	userResponse.User.Bio = user.Bio
	userResponse.User.Image = user.Image
	userResponse.User.Token = utils.GenerateJWT(user.ID)
	return userResponse
}

type profileResponse struct {
	Profile struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"profile"`
}

func newProfileResponse(us user.Store, userID uint, user *model.User) *profileResponse {
	profileResponse := new(profileResponse)
	profileResponse.Profile.Username = user.Username
	profileResponse.Profile.Bio = user.Bio
	profileResponse.Profile.Image = user.Image
	profileResponse.Profile.Following, _ = us.IsFollower(user.ID, userID)
	return profileResponse
}

type articleResponse struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"author"`
}

type singleArticleResponse struct {
	Article *articleResponse `json:"article"`
}

type articleListResponse struct {
	Articles      []*articleResponse `json:"articles"`
	ArticlesCount int                `json:"articlesCount"`
}

func newArticleResponse(c echo.Context, article *model.Article) *singleArticleResponse {
	articleResponse := new(articleResponse)
	articleResponse.TagList = make([]string, 0)
	articleResponse.Slug = article.Slug
	articleResponse.Title = article.Title
	articleResponse.Description = article.Description
	articleResponse.Body = article.Body
	articleResponse.CreatedAt = article.CreatedAt
	articleResponse.UpdatedAt = article.UpdatedAt
	for _, t := range article.Tags {
		articleResponse.TagList = append(articleResponse.TagList, t.Tag)
	}
	for _, u := range article.Favorites {
		if u.ID == userIDFromToken(c) {
			articleResponse.Favorited = true
		}
	}
	articleResponse.FavoritesCount = len(article.Favorites)
	articleResponse.Author.Username = article.Author.Username
	articleResponse.Author.Image = article.Author.Image
	articleResponse.Author.Bio = article.Author.Bio
	articleResponse.Author.Following = article.Author.FollowedBy(userIDFromToken(c))
	return &singleArticleResponse{articleResponse}
}

func newArticleListResponse(us user.Store, userID uint, articles []model.Article, count int) *articleListResponse {
	articleListResponse := new(articleListResponse)
	articleListResponse.Articles = make([]*articleResponse, 0)
	for _, article := range articles {
		articleResponse := new(articleResponse)
		articleResponse.TagList = make([]string, 0)
		articleResponse.Slug = article.Slug
		articleResponse.Title = article.Title
		articleResponse.Description = article.Description
		articleResponse.Body = article.Body
		articleResponse.CreatedAt = article.CreatedAt
		articleResponse.UpdatedAt = article.UpdatedAt
		for _, tag := range article.Tags {
			articleResponse.TagList = append(articleResponse.TagList, tag.Tag)
		}
		for _, u := range article.Favorites {
			if u.ID == userID {
				articleResponse.Favorited = true
			}
		}
		articleResponse.FavoritesCount = len(article.Favorites)
		articleResponse.Author.Username = article.Author.Username
		articleResponse.Author.Image = article.Author.Image
		articleResponse.Author.Bio = article.Author.Bio
		articleResponse.Author.Following, _ = us.IsFollower(article.AuthorID, userID)
		articleListResponse.Articles = append(articleListResponse.Articles, articleResponse)
	}
	articleListResponse.ArticlesCount = count
	return articleListResponse
}

type commentResponse struct {
	ID        uint      `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Author    struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"author"`
}

type singleCommentResponse struct {
	Comment *commentResponse `json:"comment"`
}

type commentListResponse struct {
	Comments []commentResponse `json:"comments"`
}

func newCommentResponse(c echo.Context, cm *model.Comment) *singleCommentResponse {
	comment := new(commentResponse)
	comment.ID = cm.ID
	comment.Body = cm.Body
	comment.CreatedAt = cm.CreatedAt
	comment.UpdatedAt = cm.UpdatedAt
	comment.Author.Username = cm.User.Username
	comment.Author.Image = cm.User.Image
	comment.Author.Bio = cm.User.Bio
	comment.Author.Following = cm.User.FollowedBy(userIDFromToken(c))
	return &singleCommentResponse{comment}
}

func newCommentListResponse(c echo.Context, comments []model.Comment) *commentListResponse {
	commentListResponse := new(commentListResponse)
	cr := commentResponse{}
	commentListResponse.Comments = make([]commentResponse, 0)
	for _, comment := range comments {
		cr.ID = comment.ID
		cr.Body = comment.Body
		cr.CreatedAt = comment.CreatedAt
		cr.UpdatedAt = comment.UpdatedAt
		cr.Author.Username = comment.User.Username
		cr.Author.Image = comment.User.Image
		cr.Author.Bio = comment.User.Bio
		cr.Author.Following = comment.User.FollowedBy(userIDFromToken(c))

		commentListResponse.Comments = append(commentListResponse.Comments, cr)
	}
	return commentListResponse
}

type tagListResponse struct {
	Tags []string `json:"tags"`
}

func newTagListResponse(tags []model.Tag) *tagListResponse {
	tagListResponse := new(tagListResponse)
	for _, tag := range tags {
		tagListResponse.Tags = append(tagListResponse.Tags, tag.Tag)
	}
	return tagListResponse
}
