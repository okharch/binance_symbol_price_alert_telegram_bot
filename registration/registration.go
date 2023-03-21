package registration

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleUserRegistration(w http.ResponseWriter, r *http.Request) {
	// Parse user information from request body
	var user struct {
		TelegramID int    `json:"telegram_id"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Username   string `json:"username"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert new user into database
	_, err = db.Exec(`
        INSERT INTO users (telegram_id, first_name, last_name, username)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (telegram_id) DO NOTHING
    `, user.TelegramID, user.FirstName, user.LastName, user.Username)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Send response back to client
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User registered successfully")
}
