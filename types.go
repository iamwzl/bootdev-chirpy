package main

import (
  "sync/atomic"
  "github.com/StupidWeasel/bootdev-chirpy/internal/database"
  "github.com/google/uuid"
  "time"
)

type apiConfig struct {
  database *database.Queries
  platform string
  users userFuncs
  admin adminFuncs
  messages msgFuncs
  metrics metricFuncs
}

type userFuncs struct{
  cfg *apiConfig
}
type adminFuncs struct{
  cfg *apiConfig
}
type msgFuncs struct{
  cfg *apiConfig
}
type metricFuncs struct{
  fileserverHits atomic.Int32
  cfg *apiConfig
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

