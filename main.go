package main
import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/StupidWeasel/bootdev-chirpy/internal/database"
	"log"
	"net/http"
	"os"
	_ "github.com/lib/pq"
)
func main(){

	// Env
	godotenv.Load()
	envDBURL := os.Getenv("DB_URL")
	if envDBURL == ""{
		log.Fatal("DB_URL is not set")
	}
	envPLATFORM := os.Getenv("PLATFORM")
	if envPLATFORM == ""{
		log.Fatal("DB_URL is not set")
	}

	db, err := sql.Open("postgres", envDBURL)
	if err != nil{
		log.Fatalf("Unable to connect to DB: %s", err)
	}
	fmt.Println("Connected to database!")
	dbQueries := database.New(db)
	
	// API Config
	ApiCFG := apiConfig{
		database: dbQueries,
		platform: envPLATFORM,
	}
	ApiCFG.users.cfg = &ApiCFG
	ApiCFG.admin.cfg = &ApiCFG
	ApiCFG.messages.cfg = &ApiCFG
	ApiCFG.metrics.cfg = &ApiCFG

	// Server Mux
	mux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	fs := http.FileServer(http.Dir("./webroot/"))
	
	// App
	mux.Handle("/app/", http.StripPrefix("/app/",ApiCFG.metrics.middlewareMetricsInc(fs)))
	
	// API
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("POST /api/users", ApiCFG.users.CreateUser)
	mux.HandleFunc("POST /api/chirps", ApiCFG.messages.CreateMessage)
	mux.HandleFunc("GET /api/chirps", ApiCFG.messages.GetMessages)
	mux.HandleFunc("GET /api/chirps/{id}", ApiCFG.messages.GetMessage)

	// Admin
	mux.HandleFunc("GET /admin/metrics", ApiCFG.metrics.metricsHandler)
	mux.HandleFunc("POST /admin/reset", ApiCFG.admin.ResetHandler)

	fmt.Println("Starting HTTPD on :8080")
	server.ListenAndServe()

}