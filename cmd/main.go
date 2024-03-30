package main

import (
	"fmt"
	"log"
	"net/http"

	"lexichat-backend/pkg/config"
	"lexichat-backend/pkg/handlers"
	"lexichat-backend/pkg/utils/db"
	logging "lexichat-backend/pkg/utils/logging"

	"github.com/gorilla/mux"
)

func main() {
	// setup logging
	logging.Initialize_logging()
	
	router := mux.NewRouter()

	dbClient, _ := db.Initialize_postgres()
	defer dbClient.Close()

    router.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "pong")
    })
	router.HandleFunc("/api/v1/user/create", handlers.CreateUser(dbClient))

	port := config.GOLANG_SERVER_PORT
	log.Printf("Server started on port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}