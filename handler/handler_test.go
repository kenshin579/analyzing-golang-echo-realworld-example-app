package handler

import (
	"log"
	"os"
	"testing"

	"encoding/json"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kenshin579/analyzing-golang-echo-realworld-example-app/article"
	"github.com/kenshin579/analyzing-golang-echo-realworld-example-app/db"
	"github.com/kenshin579/analyzing-golang-echo-realworld-example-app/model"
	"github.com/kenshin579/analyzing-golang-echo-realworld-example-app/router"
	"github.com/kenshin579/analyzing-golang-echo-realworld-example-app/store"
	"github.com/kenshin579/analyzing-golang-echo-realworld-example-app/user"
	"github.com/labstack/echo/v4"
)

var (
	d  *gorm.DB
	us user.Store
	as article.Store
	h  *Handler
	e  *echo.Echo
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func authHeader(token string) string {
	return "Token " + token
}

func setup() {
	d = db.TestDB()
	db.AutoMigrate(d)
	us = store.NewUserStore(d)
	as = store.NewArticleStore(d)
	h = NewHandler(us, as)
	e = router.New()
	loadFixtures()
}

func tearDown() {
	_ = d.Close()
	if err := db.DropTestDB(); err != nil {
		log.Fatal(err)
	}
}

func responseMap(b []byte, key string) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m[key].(map[string]interface{})
}

func loadFixtures() error {
	user1BioText := "user1 bio"
	user1Image := "http://realworld.io/user1.jpg"
	user1 := model.User{
		Username: "user1",
		Email:    "user1@realworld.io",
		Bio:      &user1BioText,
		Image:    &user1Image,
	}
	user1.Password, _ = user1.HashPassword("secret")
	if err := us.Create(&user1); err != nil {
		return err
	}

	user2BioText := "user2 bio"
	user2Image := "http://realworld.io/user2.jpg"
	user2 := model.User{
		Username: "user2",
		Email:    "user2@realworld.io",
		Bio:      &user2BioText,
		Image:    &user2Image,
	}
	user2.Password, _ = user2.HashPassword("secret")
	if err := us.Create(&user2); err != nil {
		return err
	}
	us.AddFollower(&user2, user1.ID)

	article := model.Article{
		Slug:        "article1-slug",
		Title:       "article1 title",
		Description: "article1 description",
		Body:        "article1 body",
		AuthorID:    1,
		Tags: []model.Tag{
			{
				Tag: "tag1",
			},
			{
				Tag: "tag2",
			},
		},
	}
	as.CreateArticle(&article)
	as.AddComment(&article, &model.Comment{
		Body:      "article1 comment1",
		ArticleID: 1,
		UserID:    1,
	})

	article2 := model.Article{
		Slug:        "article2-slug",
		Title:       "article2 title",
		Description: "article2 description",
		Body:        "article2 body",
		AuthorID:    2,
		Favorites: []model.User{
			user1,
		},
		Tags: []model.Tag{
			{
				Tag: "tag1",
			},
		},
	}
	as.CreateArticle(&article2)
	as.AddComment(&article2, &model.Comment{
		Body:      "article2 comment1 by user1",
		ArticleID: 2,
		UserID:    1,
	})
	as.AddFavorite(&article2, 1)

	return nil
}
