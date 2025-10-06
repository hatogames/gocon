package form

import (
	"encoding/json"
	"errors"
	"fmt"
	connection "gocon/db"
	"gocon/logger"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func Create(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		vars := mux.Vars(r)

		schoolstr := vars["school"]
		schoolID, err := strconv.Atoi(schoolstr)
		if err != nil {
			http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
			return
		}

		var wireframe connection.Wireframe
		result := connection.DB.First(&wireframe, connection.Wireframe{SchoolID: uint(schoolID), Name: "create"})
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Es konnte kein Eintrag gefunden werden", http.StatusNotFound)
			} else {
				http.Error(w, "Datenbank-Fehler", http.StatusInternalServerError)
			}
			return
		}
		if !wireframe.Activ {
			http.Error(w, "Der Eintrag ist momentan nicht aktiviert", http.StatusForbidden)
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

		logger.SendDiscord("Es wurde ein Wireframe ausgeliefert")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(reponse)

	case "POST":
		fmt.Fprint(w, "POST Request empfangen")
	}

}
