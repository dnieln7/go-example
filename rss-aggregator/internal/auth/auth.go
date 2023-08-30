package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Authorization: ApiKey {value}
func GetAPIKey(headers http.Header) (string, error) {
	value := headers.Get("Authorization")

	if value == "" {
		return "", errors.New("Missing Authorization header")
	}

	values := strings.Split(value, " ")

	if len(values) != 2 {
		return "", errors.New("Malformed Authorization header")
	}

	if values[0] != "ApiKey" {
		return "", errors.New("First occurrence MUST be ApiKey")
	}

	return values[1], nil
}
