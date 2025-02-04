package main

import (
    "net/http"
)

func (a *adminFuncs)ResetHandler(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    if a.cfg.platform != "dev"{
        respondWithError(w, http.StatusForbidden, "Forbidden", nil)
        return
    }           

    a.cfg.metrics.metricsReset()
    err := a.cfg.database.UsersClear(r.Context())
    if err != nil{
        respondWithError(w, http.StatusInternalServerError, "Unable to clear users database", err)
        return
    }

    respondWithJSON(w, http.StatusOK, "Reset metrics & cleaered database")
}