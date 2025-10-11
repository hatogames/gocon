package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	connection "gocon/db"
	funcs "gocon/func"
	mailer "gocon/mailer"
	"log"
	"net/http"
	"net/mail"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InitMailRequest struct {
	Email string `json:"email"`
}

func Initmail(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	/*case "GET":
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
	*/

	case "POST":

		var req InitMailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Ungültiger Request-Body", http.StatusBadRequest)
			return
		}

		_, err = mail.ParseAddress(req.Email)
		if err != nil {
			http.Error(w, "Ungültige E-Mail-Adresse", http.StatusBadRequest)
			return
		}

		code, err := funcs.RandomOTP(6)
		if err != nil {
			http.Error(w, "Interner Fehler", http.StatusInternalServerError)
			return
		}

		var user connection.User
		result := connection.DB.First(&user, "email = ?", req.Email)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Email existiert noch nicht → OK, weiter
		} else if result.Error != nil {
			http.Error(w, "Interner Fehler", http.StatusInternalServerError)
			return
		} else {
			// Email existiert bereits
			http.Error(w, "Diese Email wird bereits verwendet", http.StatusConflict)
			return
		}

		token := connection.EmailToken{
			Email:     req.Email,
			Code:      code,
			ExpiresAt: time.Now().Add(15 * time.Minute),
		}
		if err := connection.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},                         // hier prüfen wir Email
			DoUpdates: clause.AssignmentColumns([]string{"code", "expires_at"}), // Felder die aktualisiert werden
		}).Create(&token).Error; err != nil {
			log.Println("Fehler beim Speichern/Updaten des Tokens:", err)
			http.Error(w, "Interner Fehler", http.StatusInternalServerError)
			return
		}

		body := fmt.Sprintf(
			"<h2>Bitte verifizieren Sie Ihren Account</h2>"+
				"<p>Geben sie dieses OTP an, um Ihren Account zu verifizieren:</p>"+
				"%s", code,
		)

		err = mailer.SendMail(
			req.Email,
			"Verifizierungslink",
			body,
		)
		if err != nil {
			http.Error(w, "Interner Fehler", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "Ihnen wurde eine Email mit Verifizierungscode gesendet")
	}
}
