package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func parseQueryInt(value, fieldName string, w http.ResponseWriter) (int, bool) {
	i, err := strconv.Atoi(value)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid %s type provided: %s", fieldName, value))
		return 0, false
	}
	return i, true
}
