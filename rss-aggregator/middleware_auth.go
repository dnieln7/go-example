package main

import (
	"fmt"
	"net/http"

	"github.com/dnieln7/go-examples/rss-aggregator/internal/auth"
	"github.com/dnieln7/go-examples/rss-aggregator/internal/database"
)

type authorized func(writer http.ResponseWriter, request *http.Request, tbUser database.TbUser)

func (apiConfig *ApiConfig) middlewareAuth(block authorized) http.HandlerFunc  {
	return func(writer http.ResponseWriter, request *http.Request) {
		apiKey, err := auth.GetAPIKey(request.Header)

		if err != nil {
			message := fmt.Sprintf("AuthError: %v", err)
			responseError(writer, 403, message)
			return
		}
	
		tbUser, err := apiConfig.DB.GetUserByAPIKey(request.Context(), apiKey)
	
		if err != nil {
			message := fmt.Sprintf("Could not get user: %v", err)
			responseError(writer, 400, message)
			return
		}
	
		block(writer, request, tbUser)
	}
}
