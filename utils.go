package main

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"strings"
)

func UnmarshalJSON[T any](r io.Reader, v *T) error {
	decoder := json.NewDecoder(r)
	//decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("unmarshal JSON: %w", err)
	}
	return nil
}

// These are stolen from the lesson solution & are much cleaner than what I was doing before!
// https://www.boot.dev/courses/learn-http-servers-golang
func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	respondWithJSON(w, code, apiErrorResponse{
		ErrorMsg: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithStatus(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func isDuplicateKeyError(err error) bool {
	pgErr, ok := err.(*pq.Error)
	return ok && pgErr.Code == "23505"
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("No authorized header")
	}
	splitHeader := strings.Split(authHeader, " ")
	if len(splitHeader) != 2 || splitHeader[0] != "ApiKey" {
		return "", fmt.Errorf("Malformed header")
	}
	return splitHeader[1], nil
}
