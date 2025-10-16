package school

import (
	"encoding/json"
	"fmt"
	connection "gocon/db"
	"gocon/db/mini"
	"net/http"
)

type WhereType string

const (
	Student      WhereType = "student"
	Registration WhereType = "registration"
	Selected     WhereType = "selected"
)

type request struct {
	Users []int     `json:"users"`
	Where WhereType `json:"where"`
}

func Move(w http.ResponseWriter, r *http.Request) {
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

	switch req.Where {
	case Student, Registration, Selected:
		if Student == req.Where {
			http.Error(w, "Ungültiger Request-Body", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Ungültiger Request-Body", http.StatusBadRequest)
		return
	}

	if len(req.Users) < 1 {
		http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
		fmt.Print(req)
		return

	}

	result := connection.DB.
		Model(&connection.User{}).
		Where("school_id = ? AND id IN ?", session.Id, req.Users).
		Update("role", string(req.Where))
	if result.Error != nil {
		http.Error(w, "Fehler beim Aktualisieren der Rollen", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Verschiebung wurde durchgeführt")
}
