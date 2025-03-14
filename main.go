package main

import (
	"database/sql"
	"fmt"
	"github.com/StupidWeasel/bootdev-chirpy/internal/auth"
	"github.com/StupidWeasel/bootdev-chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {

	// Env
	godotenv.Load()
	envDBURL := os.Getenv("DB_URL")
	if envDBURL == "" {
		log.Fatal("DB_URL is not set")
	}
	envPLATFORM := os.Getenv("PLATFORM")
	if envPLATFORM == "" {
		log.Fatal("DB_URL is not set")
	}
	envSECRET := os.Getenv("SECRET")
	if envSECRET == "" {
		log.Fatal("SECRET is not set")
	}
	if envSECRET == "biglongsecrethere" {
		log.Fatal("Set a proper .env SECRET")
	}
	envPOLKAKEY := os.Getenv("POLKAKEY")
	if envPOLKAKEY == "" {
		log.Fatal("POLKAKEY is not set")
	}
	if envPOLKAKEY == "keyhere" {
		log.Fatal("Set a proper .env POLKAKEY")
	}

	db, err := sql.Open("postgres", envDBURL)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %s", err)
	}
	fmt.Println("Connected to database!")
	dbQueries := database.New(db)

	// API Config
	ApiCFG := apiConfig{
		database: dbQueries,
		platform: envPLATFORM,
		secret:   envSECRET,
		polkakey: envPOLKAKEY,
	}
	ApiCFG.users.cfg = &ApiCFG
	dummyHash, err := auth.HashPassword("I love to sing-a, About a moon-a and a June-a and a spring-a")
	if err != nil {
		log.Fatalf("Unable to generate dummyhash: %s", err)
	}
	ApiCFG.users.dummyHash = dummyHash
	ApiCFG.admin.cfg = &ApiCFG
	ApiCFG.messages.cfg = &ApiCFG
	ApiCFG.metrics.cfg = &ApiCFG
	ApiCFG.polka.cfg = &ApiCFG

	// Server Mux
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fs := http.FileServer(http.Dir("./webroot/"))

	// App
	mux.Handle("/app/", http.StripPrefix("/app/", ApiCFG.metrics.middlewareMetricsInc(fs)))

	// API
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	// API - Users
	mux.HandleFunc("POST /api/users", ApiCFG.users.CreateUser)
	mux.HandleFunc("PUT /api/users", ApiCFG.users.UserUpdateSelf)
	mux.HandleFunc("POST /api/login", ApiCFG.users.LoginUser)
	mux.HandleFunc("POST /api/refresh", ApiCFG.users.RefreshAuth)
	mux.HandleFunc("POST /api/revoke", ApiCFG.users.RevokeRefreshToken)
	// API - Messages
	mux.HandleFunc("POST /api/chirps", ApiCFG.messages.CreateMessage)
	mux.HandleFunc("GET /api/chirps", ApiCFG.messages.GetMessages)
	mux.HandleFunc("GET /api/chirps/{id}", ApiCFG.messages.GetMessage)
	mux.HandleFunc("DELETE /api/chirps/{id}", ApiCFG.messages.DeleteMessage)
	// API - Polka
	mux.HandleFunc("POST /api/polka/webhooks", ApiCFG.polka.Webhook)

	// Admin
	mux.HandleFunc("GET /admin/metrics", ApiCFG.metrics.metricsHandler)
	mux.HandleFunc("POST /admin/reset", ApiCFG.admin.ResetHandler)

	fmt.Println("Starting HTTPD on :8080")
	server.ListenAndServe()

}
