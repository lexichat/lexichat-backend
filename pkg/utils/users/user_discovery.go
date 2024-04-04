package utils

import (
	"database/sql"
	"fmt"
	"lexichat-backend/pkg/models"
	"lexichat-backend/pkg/utils/logging"
)


func DiscoverUsersByUserId(db *sql.DB, partialUserID string) ([]models.User, error) {
    query := `SELECT user_id, user_name, phone_number FROM users WHERE user_id LIKE '%' || $1 || '%'`
    rows, err := db.Query(query, partialUserID)
    if err != nil {
		logging.ErrorLogger.Fatalln("Error in discovering users. ", err)
        return nil, fmt.Errorf("error querying database: %v", err)

    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.UserID, &user.Username, &user.PhoneNumber); err != nil {
            return nil, fmt.Errorf("error scanning row: %v", err)
        }
        users = append(users, user)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating over rows: %v", err)
    }

    return users, nil
}

func DiscoverUsersByPhoneContactList() {

}