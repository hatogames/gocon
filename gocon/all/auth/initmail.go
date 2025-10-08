package auth

import (
	"encoding/json"
	"fmt"
	connection "gocon/db"
	funcs "gocon/func"
	mailer "gocon/mailer"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"
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

		//TODO email darf nicht in user sein

		token := connection.EmailToken{
			Email:     req.Email,
			Code:      code,
			ExpiresAt: time.Now().Add(15 * time.Minute),
		}

		if err := connection.DB.Save(&token).Error; err != nil {
			if strings.Contains(err.Error(), "23505") {
				http.Error(w, "Diese Email wird bereits verwendet", http.StatusFound)
				return
			}

			// Alles andere
			log.Println("Fehler beim Speichern des Tokens:", err)
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
	}
}
