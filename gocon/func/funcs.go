package funcs

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func RandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil // hex = doppelt so lang
}

func CompareStringArray(arr1 []string, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	// Zähle alle Elemente von arr1
	count := make(map[string]int)
	for _, v := range arr1 {
		count[v]++
	}

	// Prüfe, ob arr2 dieselben Elemente enthält
	for _, v := range arr2 {
		if count[v] == 0 {
			return false
		}
		count[v]--
	}

	return true
}

func HashPassword(password string, cost int) (string, error) {
	if cost == 0 {
		cost = 14
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPassword vergleicht ein Klartext-Passwort mit einem gespeicherten Hash.
// Gibt true zurück, wenn das Passwort korrekt ist.
func CheckPassword(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
