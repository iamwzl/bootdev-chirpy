package main

import (
    "net/http"
    "github.com/StupidWeasel/bootdev-chirpy/internal/database"
    "github.com/StupidWeasel/bootdev-chirpy/internal/auth"
    "errors"
    "github.com/google/uuid"
    "database/sql"
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

func (u *userFuncs)createRefreshToken(user uuid.UUID, r *http.Request)(string, error){
    for attempts := 0; attempts <= 3; attempts++ {
        token, err := auth.MakeRefreshToken()
        if err != nil {
            return "", err
        }

        params := database.AddRefreshTokenParams {
            Token: token,
            UserID: user,
        }

        err = u.cfg.database.AddRefreshToken(r.Context(), params)
        if err == nil {
            return token, nil
        }
        
        if !isDuplicateKeyError(err) {
            return "", err
        }
    }
    return "", errors.New("failed to generate unique token")
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

    token, err := auth.MakeJWT(result.ID, u.cfg.secret)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }

    refreshToken, err := u.createRefreshToken(result.ID, r)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }

    respondWithJSON(w, http.StatusOK, chirpUserLogin{
        ID: result.ID,
        CreatedAt: result.CreatedAt,
        UpdatedAt: result.UpdatedAt,
        Email: result.Email,
        Token: token,
        RefreshToken: refreshToken,
    })
}

func (u *userFuncs)RefreshAuth(w http.ResponseWriter, r *http.Request){
    refreshToken, err := auth.GetBearerToken(r.Header)
    if err != nil{
        respondWithError(w, http.StatusUnauthorized, "No refresh token", err)
        return
    }
    user_id, err := u.cfg.database.GetRefreshToken(r.Context(), refreshToken)
    if err != nil{
        if errors.Is(err, sql.ErrNoRows) {
            respondWithError(w, http.StatusUnauthorized, "No refresh token", err)
            return
        }
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }

    authToken, err := auth.MakeJWT(user_id, u.cfg.secret)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }

     respondWithJSON(w, http.StatusOK, chirpRefreshAuth{
        Token: authToken,
     })

}

func (u *userFuncs)RevokeRefreshToken(w http.ResponseWriter, r *http.Request){
    refreshToken, err := auth.GetBearerToken(r.Header)
    if err != nil{
        respondWithError(w, http.StatusUnauthorized, "No refresh token", err)
        return
    }

    err = u.cfg.database.RevokeRefreshToken(r.Context(), refreshToken)
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
        return
    }
    respondWithStatus(w, http.StatusNoContent)
}