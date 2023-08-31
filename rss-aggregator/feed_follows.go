package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dnieln7/go-examples/rss-aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func tbFeedFollowtoFeedFollow(tbFeedFollow database.TbFeedFollow) FeedFollow {
	return FeedFollow{
		ID:        tbFeedFollow.ID,
		UserID:    tbFeedFollow.UserID,
		FeedID:    tbFeedFollow.FeedID,
		CreatedAt: tbFeedFollow.CreatedAt,
		UpdatedAt: tbFeedFollow.UpdatedAt,
	}
}

func tbFeedFollowstoFeedFollows(tbFeedFollow []database.TbFeedFollow) []FeedFollow {
	feedFollows := []FeedFollow{}

	for _, feedFollow := range tbFeedFollow {
		feedFollows = append(feedFollows, tbFeedFollowtoFeedFollow(feedFollow))
	}

	return feedFollows
}

type PostFeedFollowBody struct {
	FeedID uuid.UUID `json:"feed_id"`
}

func (apiConfig *ApiConfig) postFeedFollow(writer http.ResponseWriter, request *http.Request, tbUser database.TbUser) {
	decoder := json.NewDecoder(request.Body)

	body := PostFeedFollowBody{}

	err := decoder.Decode(&body)

	if err != nil {
		message := fmt.Sprintf("Could not parse JSON: %v", err)
		responseError(writer, 400, message)
		return
	}

	tbFeedFollow, err := apiConfig.DB.CreateFeedFollow(request.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    tbUser.ID,
		FeedID:    body.FeedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		message := fmt.Sprintf("Could not create feed follow: %v", err)
		responseError(writer, 400, message)
		return
	}

	responseJson(writer, 201, tbFeedFollowtoFeedFollow(tbFeedFollow))
}

func (apiConfig *ApiConfig) getFeedFollows(writer http.ResponseWriter, request *http.Request, tbUser database.TbUser) {
	tbFeedFollows, err := apiConfig.DB.GetFeedFollows(request.Context(), tbUser.ID)

	if err != nil {
		message := fmt.Sprintf("Could not get feed follows: %v", err)
		responseError(writer, 400, message)
		return
	}

	responseJson(writer, 200, tbFeedFollowstoFeedFollows(tbFeedFollows))
}

func (apiConfig *ApiConfig) deleteFeedFollow(writer http.ResponseWriter, request *http.Request, tbUser database.TbUser) {
	path := request.URL.Path
	query := request.URL.Query()

	log.Println("path: ", path)
	log.Println("query: ", query)

	feedFollowIDString := chi.URLParam(request, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDString)

	if err != nil {
		message := fmt.Sprintf("Could not parse feedFollowID: %v", err)
		responseError(writer, 400, message)
		return
	}

	err = apiConfig.DB.DeleteFeedFollow(request.Context(), database.DeleteFeedFollowParams{
		ID: feedFollowID,
		UserID: tbUser.ID,
	})

	if err != nil {
		message := fmt.Sprintf("Could not delete feed follow: %v", err)
		responseError(writer, 400, message)
		return
	}

	responseJson(writer, 200, struct{} {})
}
