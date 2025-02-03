package main

import (
	"sync/atomic"
	"strings"
	"github.com/StupidWeasel/bootdev-chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database *database.Queries
}

type chirpMessage struct{
	Body string `json:"body"`
}

type apiResponse struct{
	CleanedBody string `json:"cleaned_body"`
}

type apiErrorResponse struct{
	ErrorMsg string `json:"error"`
}

func (c chirpMessage) CleanedBody() string{
	badWords := map[string]struct{}{"kerfuffle":{},"sharbert":{},"fornax":{}}
	output := strings.Split(c.Body, " ")
	for i,word := range output{
		if _, exists := badWords[word]; exists {
			output[i] = "****"
			continue
		}
	}
	return strings.Join(output, " ")
}