package main

import (
	//"log"
	"net/http"
	"fmt"
)

func main(){

	mux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	fs := http.FileServer(http.Dir("./webroot/"))
	apiCFG := apiConfig{}
	mux.Handle("/app/", http.StripPrefix("/app/",apiCFG.middlewareMetricsInc(fs)))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCFG.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCFG.metricsResetHandler)
	mux.HandleFunc("POST /api/validate_chirp", APIvalidateChirp)
	fmt.Println("Starting server on :8080")
	server.ListenAndServe()

}