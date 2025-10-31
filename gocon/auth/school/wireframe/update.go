package wireframe

import (
	"encoding/json"
	"fmt"
	connection "gocon/db"
	"gocon/db/mini"
	"net/http"
)

type request struct {
	Wireframe string                 `json:"wireframe"`
	Data      map[string]interface{} `json:"data"`
	Keys      []string               `json:"keys"`
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

	// Marshal data and keys to JSON for datatypes.JSON columns
	jsonData, err := json.Marshal(req.Data)
	if err != nil {
		http.Error(w, "Fehler beim Verarbeiten der Daten", http.StatusInternalServerError)
		return
	}
	jsonKeys, err := json.Marshal(req.Keys)
	if err != nil {
		http.Error(w, "Fehler beim Verarbeiten der Keys", http.StatusInternalServerError)
		return
	}

	updates := map[string]interface{}{
		"data": jsonData,
		"keys": jsonKeys,
	}

	result := connection.DB.
		Model(&connection.Wireframe{}).
		Where("school_id = ? AND name = ?", session.Id, req.Wireframe+fmt.Sprint(session.Id)).
		Updates(updates)
	if result.Error != nil {
		http.Error(w, "Fehler beim Speichern", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Erfolgreich gespeichert"))
}
