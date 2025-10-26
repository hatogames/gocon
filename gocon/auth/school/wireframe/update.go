package wireframe

import (
	"encoding/json"
	connection "gocon/db"
	"gocon/db/mini"
	"net/http"
)

type request struct {
	Wireframe string                 `json:"wireframe"`
	Data      map[string]interface{} `json:"data"`
}

// TODO: form validation
func Update(w http.ResponseWriter, r *http.Request) {
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

	var req request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Ungültiger Request-Body", http.StatusBadRequest)
		return
	}

	if req.Wireframe == "" {
		http.Error(w, "Ungültiger Request-Body", http.StatusBadRequest)
		return
	}

	result := connection.DB.
		Model(&connection.Wireframe{}).
		Where("school_id = ? AND name = ?", session.Id, req.Wireframe).
		Update("data", req.Data)
	if result.Error != nil {
		http.Error(w, "Fehler beim Speichern", http.StatusInternalServerError)
		return
	}

}
