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
var ApiCFG apiConfig

func main(){

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == ""{
		log.Fatal("DB_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil{
		log.Fatalf("Unable to connect to DB: %s", err)
	}
	fmt.Println("Connected to database!")
	dbQueries := database.New(db)
	ApiCFG.database = dbQueries

	mux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	fs := http.FileServer(http.Dir("./webroot/"))
	// App
	mux.Handle("/app/", http.StripPrefix("/app/",ApiCFG.middlewareMetricsInc(fs)))
	
	// API
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("POST /api/users", APICreateUser)
	mux.HandleFunc("POST /api/chirps", APICreateMessage)
	mux.HandleFunc("GET /api/chirps", APIGetMessages)
	mux.HandleFunc("GET /api/chirps/{id}", APIGetMessage)

	// Admin
	mux.HandleFunc("GET /admin/metrics", ApiCFG.metricsHandler)
	mux.HandleFunc("POST /admin/reset", adminResetHandler)

	fmt.Println("Starting HTTPD on :8080")
	server.ListenAndServe()

}