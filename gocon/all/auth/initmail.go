package auth

import (
	"encoding/json"
	"fmt"
	connection "gocon/db"
	funcs "gocon/func"
	mail "gocon/mailer"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type InitMailRequest struct {
	Email string `json:"email"`
}

func Initmail(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		code := mux.Vars(r)["code"]

		var token connection.EmailToken
		result := connection.DB.First(&token, "code = ?", code)
		if result.Error != nil {
			http.Error(w, "Ungültiger oder abgelaufener Code", http.StatusNotFound)
			return
		}

		if time.Now().After(token.ExpiresAt) {
			http.Error(w, "Dieser Verifizierungslink ist abgelaufen", http.StatusGone)
			return
		}

		token.Verified = true
		if err := connection.DB.Save(&token).Error; err != nil {
			log.Println("Fehler beim Aktualisieren des Tokens:", err)
			http.Error(w, "Interner Fehler", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "E-Mail %s wurde erfolgreich verifiziert ✅", token.Email)

	case "POST":

		var req InitMailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Ungültiger Request-Body", http.StatusBadRequest)
			return
		}

		code, err := funcs.RandomString(10)
		if err != nil {
			panic(err)
		}

		//TODO email darf nicht in user sein

		var verified bool
		var existing connection.EmailToken
		if err := connection.DB.First(&existing, "email = ?", req.Email).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				verified = false
			}
		} else {
			verified = existing.Verified
		}

		token := connection.EmailToken{
			Email:     req.Email,
			Code:      code,
			ExpiresAt: time.Now().Add(15 * time.Minute),
			Verified:  verified,
		}

		if err := connection.DB.Save(&token).Error; err != nil {
			log.Println("Fehler beim Speichern des Tokens:", err)
			http.Error(w, "Interner Fehler", http.StatusInternalServerError)
			return
		}
		link := fmt.Sprintf("http://localhost:8080/initmail/%s", code)
		body := fmt.Sprintf(
			"<h2>Bitte verifizieren Sie Ihren Account</h2>"+
				"<p>Klicken Sie auf den Link, um Ihren Account zu verifizieren:</p>"+
				"<a href=\"%s\">%s</a>", link, link,
		)

		err = mail.SendMail(
			"hatogames@yahoo.com",
			"Verifizierungslink",
			body,
		)
		if err != nil {
			log.Fatal(err)
		}

	}
}
