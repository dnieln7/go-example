package main

import (
	"fmt"
	"net/http"

	"github.com/dnieln7/go-examples/rss-aggregator/internal/database"
)

type Post struct {
	Title           string  `json:"title"`
	Description     *string `json:"description"` // If pointer is nil Marshall will return null
	Url             string  `json:"url"`
	PublicationDate string  `json:"publication_date"`
}

func dbPostToPost(dbPost database.TbPost) Post {
	return Post{
		Title:           dbPost.Title,
		Description:     &dbPost.Description.String,
		Url:             dbPost.Url,
		PublicationDate: dbPost.PublishedAt.String(),
	}
}

func dbPostsToPosts(dbPosts []database.TbPost) []Post {
	posts := []Post{}

	for _, dbPost := range dbPosts {
		posts = append(posts, dbPostToPost(dbPost))
	}

	return posts
}

func (apiConfig *ApiConfig) getPostsByUserID(writer http.ResponseWriter, request *http.Request, tbUser database.TbUser) {
	dbPosts, err := apiConfig.DB.GetPostForUser(request.Context(), database.GetPostForUserParams{
		UserID: tbUser.ID,
		Limit:  10,
	})

	if err != nil {
		message := fmt.Sprintf("Could not get posts: %v", err)
		responseError(writer, 403, message)
		return
	}

	responseJson(writer, 200, dbPostsToPosts(dbPosts))
}
