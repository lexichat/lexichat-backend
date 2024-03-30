package main

import (
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

	// router.HandleFunc("/api/test", handlers.TestHandler)
	router.HandleFunc("/api/v1/user/create", handlers.CreateUser(dbClient))

	port := config.GOLANG_SERVER_PORT
	log.Printf("Server started on port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}