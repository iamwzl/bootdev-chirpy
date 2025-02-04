package main

import (
	"net/http"
	"github.com/StupidWeasel/bootdev-chirpy/internal/database"
	"errors"
	"database/sql"
	"github.com/google/uuid"
)

func (m *msgFuncs)CreateMessage(w http.ResponseWriter, r *http.Request){

	var message createChirpMessage
	err := UnmarshalJSON(r.Body, &message)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if len(message.Body)>140{
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	params := database.CreateMessageParams{
		Body:   message.Body,
		UserID: message.UserID,
	}

	result, err := m.cfg.database.CreateMessage(r.Context(),params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirpMessage{
		ID: result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Body: result.Body,
		UserID: result.UserID,
	})
}

func (m *msgFuncs)GetMessages(w http.ResponseWriter, r *http.Request){

	results, err := m.cfg.database.GetMessages_CreatedAtASC(r.Context())
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	messages := make([]chirpMessage, len(results),len(results))
	for i, result := range results{
		messages[i] = chirpMessage{
			ID: result.ID,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
			Body: result.Body,
			UserID: result.UserID,
		}
	}

	respondWithJSON(w, http.StatusCreated, messages)
}

func (m *msgFuncs)GetMessage(w http.ResponseWriter, r *http.Request){

	messageID, err := uuid.Parse(r.PathValue("id"))
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "No message id provided", err)
		return
	}

	result, err := m.cfg.database.GetMessage(r.Context(), messageID)
	if err != nil{
	    if errors.Is(err, sql.ErrNoRows) {
	    	respondWithError(w, http.StatusNotFound, "Message not found", err)
			return
	  	}
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	respondWithJSON(w, http.StatusOK, chirpMessage{
		ID: result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Body: result.Body,
		UserID: result.UserID,
	})
}