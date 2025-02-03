package main

import (
	"net/http"
	"fmt"
	"os"
	"github.com/StupidWeasel/bootdev-chirpy/internal/database"
	"errors"
	"database/sql"
	"github.com/google/uuid"
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

func (cfg *apiConfig)metricsReset(){
	cfg.fileserverHits.Store(0)
}

func readinessHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func adminResetHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if os.Getenv("PLATFORM") != "dev"{
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden"))
		return
	} 			

	ApiCFG.metricsReset()
	err := ApiCFG.database.UsersClear(r.Context())
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unable to clear users database: %w", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("- Reset metrics to 0\n - Cleared users database"))
}

func APICreateUser(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	var user createChirpUser
	err := UnmarshalJSON(r.Body, &user)
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
	if len(user.Email)==0{
		w.WriteHeader(http.StatusBadRequest)
		response, err := MarshalJSONToString(apiErrorResponse{ErrorMsg: "Email is empty"})
		if err != nil{
			fmt.Println(err)
			return
		}
		w.Write([]byte(response))
		return
	}

	result, err := ApiCFG.database.CreateUser(r.Context(), user.Email)
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

	createdUser := chirpUser{
		ID: result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Email: result.Email,
	}

	response, err := MarshalJSONToString(createdUser)
	if err != nil{
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func APICreateMessage(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	var message createChirpMessage
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

	params := database.CreateMessageParams{
		Body:   message.Body,
		UserID: message.UserID,
	}

	result, err := ApiCFG.database.CreateMessage(r.Context(),params)
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

	createdMessage := chirpMessage{
		ID: result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Body: result.Body,
		UserID: result.UserID,
	}

	response, err := MarshalJSONToString(createdMessage)
	if err != nil{
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func APIGetMessages(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	results, err := ApiCFG.database.GetMessages_CreatedAtASC(r.Context())
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

	messages := make([]chirpMessage, len(results),len(results))
	for i, result := range results{
		messages[i] = chirpMessage{
			ID: result.ID,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
			Body: result.Body,
			UserID: result.UserID,
		}
	}

	response, err := MarshalJSONToString(messages)
	if err != nil{
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func APIGetMessage(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	messageID, err := uuid.Parse(r.PathValue("id"))
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No message id provided"))
		return
	}

	result, err := ApiCFG.database.GetMessage(r.Context(), messageID)
	if err != nil{
    if errors.Is(err, sql.ErrNoRows) {
    	w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not found"))
			return
  	}
		w.WriteHeader(http.StatusInternalServerError)
		response, err := MarshalJSONToString(apiErrorResponse{ErrorMsg: "Something went wrong"})
		if err != nil{
			fmt.Println(err)
			return
		}
		w.Write([]byte(response))
		return
	}

	message := chirpMessage{
		ID: result.ID,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		Body: result.Body,
		UserID: result.UserID,
	}

	response, err := MarshalJSONToString(message)
	if err != nil{
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}