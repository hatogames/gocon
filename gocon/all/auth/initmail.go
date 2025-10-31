package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	connection "gocon/db"
	funcs "gocon/func"
	mailer "gocon/mailer"
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

	var req InitMailRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Ungültiger Request-Body", http.StatusBadRequest)
		return
		//HACKER-err
	}

	_, err = mail.ParseAddress(req.Email)
	if err != nil {
		http.Error(w, "Ungültige E-Mail-Adresse", http.StatusBadRequest)
		return
		//HACKER-err
	}

	code, err := funcs.RandomOTP(6)
	if err != nil {
		http.Error(w, "Interner Fehler", http.StatusInternalServerError)
		return
		//FATAL-err
	}

	var user connection.Registration
	result := connection.DB.First(&user, "email = ?", req.Email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) { //continue
	} else if result.Error != nil {
		http.Error(w, "Interner Fehler", http.StatusInternalServerError)
		//FATAL-err
		return
	} else {
		http.Error(w, "Diese Email wird bereits verwendet", http.StatusConflict)
		//NORMAL-err
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
		http.Error(w, "Interner Fehler", http.StatusInternalServerError)
		//FATAL-err
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
		//FATAL-err
		return
	}

	fmt.Fprint(w, "Ihnen wurde eine Email mit Verifizierungscode gesendet")

}
