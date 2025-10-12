package funcs

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func RandomOTP(length int) (string, error) {
	digits := "0123456789"
	otp := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[n.Int64()]
	}
	return string(otp), nil
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

func RandomString(length int) (string, error) {
	digits := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	otp := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[n.Int64()]
	}
	return string(otp), nil
}
