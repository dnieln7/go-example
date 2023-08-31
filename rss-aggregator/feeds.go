package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dnieln7/go-examples/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type Feed struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	CreatorID uuid.UUID `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func tbFeedtoFeed(tbFeed database.TbFeed) Feed {
	return Feed{
		ID:        tbFeed.ID,
		Name:      tbFeed.Name.String,
		Url:       tbFeed.Url,
		CreatorID: tbFeed.UserID,
		CreatedAt: tbFeed.CreatedAt,
		UpdatedAt: tbFeed.UpdatedAt,
	}
}

func tbFeedstoFeeds(tbFeeds []database.TbFeed) []Feed {
	feeds := []Feed{}

	for _, feed := range tbFeeds {
		feeds = append(feeds, tbFeedtoFeed(feed))
	}

	return feeds
}

type PostFeedBody struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (apiConfig *ApiConfig) postFeed(writer http.ResponseWriter, request *http.Request, tbUser database.TbUser) {
	decoder := json.NewDecoder(request.Body)

	body := PostFeedBody{}

	err := decoder.Decode(&body)

	if err != nil {
		message := fmt.Sprintf("Could not parse JSON: %v", err)
		responseError(writer, 400, message)
		return
	}

	tbFeed, err := apiConfig.DB.CreateFeed(request.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		Name: sql.NullString{
			String: body.Name,
			Valid:  true,
		},
		Url:       body.Url,
		UserID:    tbUser.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		message := fmt.Sprintf("Could not create feed: %v", err)
		responseError(writer, 400, message)
		return
	}

	responseJson(writer, 201, tbFeedtoFeed(tbFeed))
}

func (apiConfig *ApiConfig) getFeeds(writer http.ResponseWriter, request *http.Request) {
	tbFeeds, err := apiConfig.DB.GetFeeds(request.Context())

	if err != nil {
		message := fmt.Sprintf("Could not get feeds: %v", err)
		responseError(writer, 400, message)
		return
	}

	responseJson(writer, 200, tbFeedstoFeeds(tbFeeds))
}
