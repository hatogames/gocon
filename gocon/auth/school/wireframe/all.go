package wireframe

import (
	"encoding/json"
	"errors"
	connection "gocon/db"
	"gocon/db/mini"
	"net/http"

	"gorm.io/gorm"
)

func All(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Bitte loggen Sie sich ein", http.StatusUnauthorized)
		return
	}

	session, err := mini.GetSession(cookie.Value)
	if err != nil {
		http.Error(w, "Bitte loggen Sie sich ein", http.StatusUnauthorized)
		return
	}

	if session.Typ != mini.UserType("school") {
		http.Error(w, "Zugriff verweigert", http.StatusForbidden)
		return
	}

	var wireframes []connection.Wireframe
	result := connection.DB.Where("school_id = ?", session.Id).Find(&wireframes)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Es wurden noch keine Formulare gefunden", http.StatusNotFound)
		} else {
			http.Error(w, "Datenbank-Fehler", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wireframes)
}
