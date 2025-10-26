package form

import (
	"encoding/json"
	"errors"
	"fmt"
	connection "gocon/db"
	funcs "gocon/func"
	"gocon/logger"
	mailer "gocon/mailer"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Request struct {
	Data map[string]string `json:"data"`
	User User              `json:"user"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		vars := mux.Vars(r)

		//school to int
		schoolstr := vars["school"]
		schoolID, err := strconv.Atoi(schoolstr)
		if err != nil {
			http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
			return
		}

		//get "create"-wireframe
		var wireframe connection.Wireframe
		result := connection.DB.Where("school_id = ? AND name = ?", schoolID, "create").First(&wireframe)
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

		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Ungültiger Request-Body", http.StatusBadRequest)
			fmt.Print(req)
			return
		}
		_, err = mail.ParseAddress(req.User.Email)
		if err != nil {
			http.Error(w, "Ungültige E-Mail-Adresse", http.StatusBadRequest)
			fmt.Print(req)
			return
		}

		if req.User.Code == "" {
			http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
			fmt.Print(req)
			return

		}
		if req.User.Email == "" {
			http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
			fmt.Print(req)
			return

		}
		if req.User.Password == "" {
			http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
			fmt.Print(req)
			return
		}

		vars := mux.Vars(r)
		schoolstr := vars["school"]
		schoolID, err := strconv.Atoi(schoolstr)
		if err != nil {
			http.Error(w, "Ungültige Anfrage", http.StatusBadRequest)
			return
		}

		//Check code
		var token connection.EmailToken
		result := connection.DB.First(&token, "email = ?", req.User.Email)
		if result.Error != nil {
			http.Error(w, "Email nicht vorhanden oder bereits registriert", http.StatusNotFound)
			return
		}

		if token.Code != req.User.Code {
			http.Error(w, "Email ist nicht korekt verifiziert", http.StatusNotFound)
			return
		}

		//Compare keys
		var wireframe connection.Wireframe
		result = connection.DB.Where("school_id = ? AND name = ?", schoolID, "create").First(&wireframe)
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

		var keysREQ []string
		for k := range req.Data {
			keysREQ = append(keysREQ, k)
		}

		var keysWIREFRAME []string
		err = json.Unmarshal(wireframe.Keys, &keysWIREFRAME)
		if err != nil {
			http.Error(w, "Fehler beim Parsen der Keys", http.StatusInternalServerError)
			return
		}

		if !funcs.CompareStringArray(keysREQ, keysWIREFRAME) {
			http.Error(w, "Fehlerhafte http-Anfrage", http.StatusBadRequest)
			return
		}

		//Hash Password
		hash, err := funcs.HashPassword(req.User.Password, 14)
		if err != nil {
			http.Error(w, "Fehler beim Hashen des Passworts", http.StatusInternalServerError)
			return
		}
		//add missingKeys
		var school connection.School
		result = connection.DB.First(&school, "id = ?", schoolID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Interner Fehler", http.StatusInternalServerError)
			return
		} else if result.Error != nil {
			http.Error(w, "Interner Fehler", http.StatusInternalServerError)
			return
		}

		var schoolKeys []string
		err = json.Unmarshal(school.Keys, &schoolKeys)
		if err != nil {
			http.Error(w, "Fehler beim Lesen der School-Keys", http.StatusInternalServerError)
			return
		}

		for _, key := range schoolKeys {
			if _, exists := req.Data[key]; !exists {
				req.Data[key] = "" // fehlenden Key hinzufügen
			}
		}

		//make DB
		dataJSON, err := json.Marshal(req.Data)
		if err != nil {
			http.Error(w, "Fehler beim Parsen der Daten", http.StatusInternalServerError)
			return
		}

		newUser := connection.Registration{
			Email:    req.User.Email,
			Phash:    hash,
			Data:     datatypes.JSON(dataJSON),
			SchoolID: uint(schoolID),
		}

		if err := connection.DB.Create(&newUser).Error; err != nil {
			http.Error(w, "Fehler beim Speichern des Benutzers", http.StatusInternalServerError)
			return
		}

		result = connection.DB.Where("email = ?", req.User.Email).Delete(&connection.EmailToken{})
		if result.Error != nil {
			http.Error(w, "Fehler beim Löschen des Tokens", http.StatusInternalServerError)
			return
		}

		mailer.SendMail(req.User.Email, "Schulregistrierung", string(newUser.Data))

		fmt.Fprint(w, "Ihre Registrierung wurde erfolgreich durch geführt. Eine Email mit allen Angaben wurde versendet")

	}

}
