package main

import (
	"fmt"
	"log"
	"net/http"

	"lexichat-backend/pkg/config"
	"lexichat-backend/pkg/handlers"
	"lexichat-backend/pkg/utils/db"
	logging "lexichat-backend/pkg/utils/logging"

	fcm "lexichat-backend/pkg/utils/fcm"

	"github.com/gorilla/mux"
)

func main() {
	config.LoadEnvVariables(".env.dev")

	// setup logging
	logging.Initialize_logging()

    // setup FCM
    fcmClient, _ := fcm.SetupFCM()

	router := mux.NewRouter()

	dbClient, _ := db.Initialize_postgres()
	defer dbClient.Close()

    router.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "pong")
    })

	router.HandleFunc("/api/v1/users/create", handlers.CreateUser(dbClient))

    router.HandleFunc("/api/v1/channels/active-clients", handlers.GetActiveClientsHandler)

    router.HandleFunc("/api/v1/users/discover", func(w http.ResponseWriter, r *http.Request) {
        handlers.DiscoverUsersByUserIdHandler(w, r, dbClient)
    }).Methods("GET")

    router.HandleFunc("/api/v1/channels/create", func(w http.ResponseWriter, r *http.Request) {
        handlers.CreateChannelHandler(w,r, dbClient, fcmClient)
    })

    router.HandleFunc("/api/v1/users/{userid}", func(w http.ResponseWriter, r *http.Request) {
        userid := mux.Vars(r)["userid"]
        handlers.FetchUserDetailsHandler(w,r, dbClient, userid)
    })


	log.Printf("App Server started on port: %s", config.APP_SERVER_PORT)
    go func() {
        log.Fatal(http.ListenAndServe(":"+config.APP_SERVER_PORT, router))
    }()


    ///// WS ROUTER ////
	
    log.Printf("WebSocket handler started on port: %s", config.WS_SERVER_PORT)
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        handlers.HandleConnections(w, r, dbClient, fcmClient)
    })
        http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {fmt.Fprintln(w,"test")})

    go func() {
        err := http.ListenAndServe(":"+ config.WS_SERVER_PORT, nil)
        if err != nil {
            log.Fatal("ListenAndServe: ", err)
        }
    }()



    // Wait forever
    select {}
}