package main

import (
	"net/http"
	"fmt"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig)metricsHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig)metricsResetHandler(w http.ResponseWriter, r *http.Request){
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
    w.Write([]byte("His reset to 0"))
}

func readinessHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func APIvalidateChirp(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	
	var message chirpMessage
	err := UnmarshalJSON(r.Body, &message)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		response, err := MarshalJSONToString(apiErrorResponse{ErrorMsg: "Something went wrong"})
		if err != nil{
			fmt.Println(err)
			return
		}
		w.Write([]byte(response))
		return
	}

	if len(message.Body)>140{
		w.WriteHeader(http.StatusBadRequest)
		response, err := MarshalJSONToString(apiErrorResponse{ErrorMsg: "Chirp is too long"})
		if err != nil{
			fmt.Println(err)
			return
		}
		w.Write([]byte(response))
		return
	}
	response, err := MarshalJSONToString(apiResponse{CleanedBody: message.CleanedBody()})
	if err != nil{
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}