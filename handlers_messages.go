package main

import (
    "net/http"
    "github.com/StupidWeasel/bootdev-chirpy/internal/database"
    "github.com/StupidWeasel/bootdev-chirpy/internal/auth"
    "errors"
    "database/sql"
    "github.com/google/uuid"
    "strings"
)

func (m *msgFuncs)CreateMessage(w http.ResponseWriter, r *http.Request){
    token, err := auth.GetBearerToken(r.Header)
    if err != nil{
        respondWithError(w, http.StatusUnauthorized, "No auth token", err)
        return
    }

    var message createChirpMessage
    err = UnmarshalJSON(r.Body, &message)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }

    userid, err := auth.ValidateJWT(token, m.cfg.secret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Invalid auth token", err)
        return
    }

    if len(message.Body)>140{
        respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
        return
    }

    params := database.CreateMessageParams{
        Body:   message.Body,
        UserID: userid,
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

    authorString := r.URL.Query().Get("author_id")
    sortDesc := strings.ToLower(r.URL.Query().Get("sort")) == "desc"

    var results []database.Message
    if authorString == ""{
        theseResults, err := m.cfg.database.GetMessages_CreatedAtASC(r.Context())
        if err != nil{
            respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
            return
        }
        results = theseResults
    }else{
        authorID, err := uuid.Parse(authorString)
        if err != nil{
            respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
            return
        }
        theseResults, err := m.cfg.database.GetMessages_ByAuthor_CreatedAtASC(r.Context(),authorID)
        if err != nil{
            respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
            return
        }
        results = theseResults
    }

    if sortDesc{
        for i := 0; i < len(results)/2; i++ {
            j := len(results) - i - 1
            results[i], results[j] = results[j], results[i]
        }
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

    respondWithJSON(w, http.StatusOK, messages)
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

func (m *msgFuncs)DeleteMessage(w http.ResponseWriter, r *http.Request){
    authToken, err := auth.GetBearerToken(r.Header)
    if err != nil{
        respondWithError(w, http.StatusUnauthorized, "No auth token", err)
        return
    }

    userid, err := auth.ValidateJWT(authToken, m.cfg.secret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Invalid auth token", err)
        return
    }

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

    if result.UserID != userid{
        respondWithError(w, http.StatusForbidden, "Not your message", err)
        return
    }

    params := database.DeleteMessageParams{
        ID: messageID,
        UserID: userid,
    }

    numRows, err := m.cfg.database.DeleteMessage(r.Context(),params)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }

    if numRows==0{
        respondWithError(w, http.StatusNotFound, "Message not found", err)
        return
    }

    respondWithStatus(w, http.StatusNoContent)
}