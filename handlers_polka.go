package main

import (
    "net/http"
    "github.com/google/uuid"
)

func (p *polkaFuncs)Webhook(w http.ResponseWriter, r *http.Request){

    apiKey, err := GetAPIKey(r.Header)
    if err != nil || apiKey != p.cfg.polkakey{
        respondWithError(w, http.StatusUnauthorized, "Invalid key", err)
        return
    }
    
    var request PolkaRequest
    err = UnmarshalJSON(r.Body, &request)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }

    if request.Event != "user.upgraded"{
        respondWithStatus(w, http.StatusNoContent)
        return
    }

    userID, err := uuid.Parse(request.Data.UserID)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }

    numRows, err := p.cfg.database.UserUpgradeToRed(r.Context(), userID)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }

    if numRows==0{
        respondWithError(w, http.StatusNotFound, "User not found", err)
        return
    }
    respondWithStatus(w, http.StatusNoContent)
}