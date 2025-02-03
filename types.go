package main

import (
	"sync/atomic"
	"strings"
	"github.com/StupidWeasel/bootdev-chirpy/internal/database"
	"github.com/google/uuid"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database *database.Queries
}

type createChirpMessage struct{
 	Body string `json:"body"`
 	UserID uuid.UUID `json:"user_id"`
}

type chirpMessage struct{
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
  	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type getchirpMessage struct{
	ID uuid.UUID `json:"id"`
}


type chirpUser struct{
  ID uuid.UUID `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  Email string `json:"email"`
}

type createChirpUser struct{
  Email string `json:"email"`
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