package main

import (
    "net/http"
    "github.com/StupidWeasel/bootdev-chirpy/internal/database"
    "github.com/StupidWeasel/bootdev-chirpy/internal/auth"
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

    passwordHash, err := auth.HashPassword(user.Password)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }

    params := database.CreateUserParams{
        Email:   user.Email,
        HashedPassword: passwordHash,
    }

    result, err := u.cfg.database.CreateUser(r.Context(), params)
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

func (u *userFuncs)LoginUser(w http.ResponseWriter, r *http.Request){

    var user loginChirpUser
    err := UnmarshalJSON(r.Body, &user)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }

    result, err := u.cfg.database.LoginUser(r.Context(), user.Email)
    if err != nil{
        _ = auth.CheckPasswordHash(u.dummyHash, "I am a password")
        respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
        return
    }
    if auth.CheckPasswordHash(result.HashedPassword, user.Password) != nil{
        respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
        return
    }

    respondWithJSON(w, http.StatusOK, chirpUser{
        ID: result.ID,
        CreatedAt: result.CreatedAt,
        UpdatedAt: result.UpdatedAt,
        Email: result.Email,
    })
}