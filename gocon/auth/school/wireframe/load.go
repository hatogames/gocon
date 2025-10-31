package wireframe

import (
	"encoding/json"
	"errors"
	"fmt"
	connection "gocon/db"

	"gocon/db/mini"
	"net/http"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type requestLoad struct {
	Wireframe string `json:"wireframe"`
}

func Load(w http.ResponseWriter, r *http.Request) {
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

	var req requestLoad
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Ungültiger Request-Body", http.StatusBadRequest)
		return
	}

	if req.Wireframe == "" {
		http.Error(w, "Ungültiger Request-Body", http.StatusBadRequest)
		return
	}

	var wireframe connection.Wireframe
	result := connection.DB.Where("school_id = ? AND name = ?", session.Id, req.Wireframe+fmt.Sprint(session.Id)).First(&wireframe)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Es konnte kein Eintrag gefunden werden", http.StatusNotFound)
		} else {
			http.Error(w, "Datenbank-Fehler", http.StatusInternalServerError)
		}
		return
	}

	type WireframeResponse struct {
		Keys datatypes.JSON `json:"keys"`
		Data datatypes.JSON `json:"data"`
	}

	reponse := WireframeResponse{
		Keys: wireframe.Keys,
		Data: wireframe.Data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reponse)
}
