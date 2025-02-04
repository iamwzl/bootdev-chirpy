package main

import (
    "net/http"
)

func (u *userFuncs)CreateUser(w http.ResponseWriter, r *http.Request){

    var user createChirpUser
    err := UnmarshalJSON(r.Body, &user)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }
    if len(user.Email)==0{
        respondWithError(w, http.StatusBadRequest, "Email is empty", nil)
        return
    }
    result, err := u.cfg.database.CreateUser(r.Context(), user.Email)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }
    respondWithJSON(w, http.StatusCreated, chirpUser{
        ID: result.ID,
        CreatedAt: result.CreatedAt,
        UpdatedAt: result.UpdatedAt,
        Email: result.Email,
    })
}