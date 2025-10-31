package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	connection "gocon/db"
	funcs "gocon/func"
	"net/http"
	"net/mail"
	"time"

	mini "gocon/db/mini"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string        `json:"email"`
	Password string        `json:"password"`
	Type     mini.UserType `json:"type"`
}

func Login(w http.ResponseWriter, r *http.Request) {

	var req LoginRequest
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

	if req.Password == "" {
		http.Error(w, "Ungültiges Passwort", http.StatusBadRequest)
		fmt.Print(req)
		return
		//HACKER-err
	}

	if req.Type == "" {
		http.Error(w, "Ungültige Benutzerrolle", http.StatusBadRequest)
		fmt.Print(req)
		return
		//HACKER-err
	}

	if req.Type == mini.UserType("school") {
		var owner connection.Owner
		result := connection.DB.First(&owner, "email = ?", req.Email)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			bcrypt.CompareHashAndPassword([]byte("$2a$14$rTeS9s7JOAlCjUuy/YZ1meIdLTgy30sR4gJ.RDQaiGoTUiIN6jYri"), []byte("req.Password"))
			http.Error(w, "Ungültige Anmeldedaten", http.StatusUnauthorized)
			return
			//LOGGER-err
		} else if result.Error != nil {
			http.Error(w, "Interner Fehler", http.StatusInternalServerError)
			return
			//FATAL-err
		}

		err := bcrypt.CompareHashAndPassword([]byte(owner.Phash), []byte(req.Password))
		if err != nil {
			http.Error(w, "Ungültige Anmeldedaten", http.StatusUnauthorized)
			return
			//LOGGER-err
		}

		SessionStr, err := funcs.RandomString(20)
		if err != nil {
			http.Error(w, "Interner Fehler", http.StatusInternalServerError)
			return
			//FATAL-err
		}

		new := mini.Session{
			Session: SessionStr,
			Typ:     mini.UserType("school"),
			Id:      int(owner.ID),
			Expires: time.Now().Add(10 * time.Hour),
		}

		mini.Sessions = append(mini.Sessions, new)

		cookie := &http.Cookie{
			Name:     "token",
			Value:    SessionStr,
			Path:     "/",
			HttpOnly: true, // schützt vor Zugriff durch JavaScript
			Secure:   true, // nur über HTTPS senden
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(1 * time.Hour), // 1 Tag gültig
		}

		// Cookie im Response-Header setzen
		http.SetCookie(w, cookie)

		fmt.Fprint(w, "Sie werden nun angemeldet")
	}

}
