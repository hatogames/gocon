package wireframe

import (
	"encoding/json"
	"fmt"
	connection "gocon/db"
	"gocon/db/mini"
	"net/http"

	"gorm.io/datatypes"
)

func Make(w http.ResponseWriter, r *http.Request) {
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

	newWireframe := connection.Wireframe{
		SchoolID: uint(session.Id),
		Name:     req.Wireframe + fmt.Sprint(session.Id),
		Data:     datatypes.JSON([]byte(`{"data": []}`)), // korrekt geschlossen
		Keys:     datatypes.JSON([]byte(`[]`)),           // leeres Array
		Activ:    false,
	}

	result := connection.DB.Create(&newWireframe)
	if result.Error != nil {
		http.Error(w, "Datenbank-Fehler: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Es wurde eine neues Formular erstellt")
}
