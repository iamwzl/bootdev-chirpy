package main

import (
    "net/http"
    "fmt"
)

func (m *metricFuncs)middlewareMetricsInc(next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        m.fileserverHits.Add(1)
        next.ServeHTTP(w,r)
    })
}

func (m *metricFuncs)metricsHandler(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
        w.Write([]byte(fmt.Sprintf(`<html>
    <body>
        <h1>Welcome, Chirpy Admin</h1>
        <p>Chirpy has been visited %d times!</p>
    </body>
</html>`, m.fileserverHits.Load())))
}

func (m *metricFuncs)metricsReset(){
    m.fileserverHits.Store(0)
}