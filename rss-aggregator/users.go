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

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func tbUsertoUser(tbUser database.TbUser) User {
	return User{
		ID:        tbUser.ID,
		Name:      tbUser.Name.String,
		CreatedAt: tbUser.CreatedAt,
		UpdatedAt: tbUser.UpdatedAt,
	}
}

type PostUserBody struct {
	Name      string `json:"name"`
}

func (apiConfig *ApiConfig) postUser(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	createUserBody := PostUserBody{}

	err := decoder.Decode(&createUserBody)

	if err != nil {
		message := fmt.Sprintf("Could not parse JSON: %v", err)
		responseError(writer, 400, message)
		return
	}

	if createUserBody.Name == "" {
		responseError(writer, 400, "Name is empty")
		return
	}

	tbUser, err := apiConfig.DB.CreateUser(request.Context(), database.CreateUserParams{
		ID: uuid.New(),
		Name: sql.NullString{
			String: createUserBody.Name,
			Valid:  true,
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		message := fmt.Sprintf("Could not create user: %v", err)
		responseError(writer, 400, message)
		return
	}

	responseJson(writer, 201, tbUsertoUser(tbUser))
}
