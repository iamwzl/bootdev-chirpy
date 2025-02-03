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
	apiCFG := apiConfig{
		database:	dbQueries,
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	fs := http.FileServer(http.Dir("./webroot/"))
	mux.Handle("/app/", http.StripPrefix("/app/",apiCFG.middlewareMetricsInc(fs)))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCFG.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCFG.metricsResetHandler)
	mux.HandleFunc("POST /api/validate_chirp", APIvalidateChirp)
	fmt.Println("Starting HTTPD on :8080")
	server.ListenAndServe()

}