package school

import (
	"encoding/json"
	connection "gocon/db"
	"gocon/db/mini"
	"net/http"
)

func Registrations(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Bitte loggen sie sich ein", http.StatusUnauthorized)
		return
	}

	session, err := mini.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, "Bitte loggen sie sich ein", http.StatusUnauthorized)
		return
	}

	if session.Typ != mini.UserType("school") {
		http.Error(w, "Bitte loggen sie sich ein", http.StatusUnauthorized)
		return
	}

	//Verified

	var users []connection.Registration
	result := connection.DB.
		Select("data").
		Find(&users, "school_id = ?", session.Id)

	if result.Error != nil {
		http.Error(w, "Interner Serverfehler", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
