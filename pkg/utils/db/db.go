package db

import (
	"database/sql"
	"fmt"
	"lexichat-backend/pkg/config"

	log "lexichat-backend/pkg/utils/logging"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func Initialize_postgres() (*sql.DB, error) {

    connStr := "user=" + config.DBUser + " password=" + config.DBPassword +" dbname=" + config.DBName + " sslmode=disable" + " host=postgres" + " port=5432"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.ErrorLogger.Printf("Error in postgres init: %v", err)
        return nil, err
    }

    err = db.Ping()
    if err != nil {
    log.ErrorLogger.Printf("Error in test ping: %v", err)
        return nil, err
    }

    fmt.Println("PostgreSQL initialized and connected.")
    log.InfoLogger.Println("PostgreSQL initiialized and connected.")
    return db, nil
}
