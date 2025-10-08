package main

import (
	form "gocon/all"
	auth "gocon/all/auth"
	connection "gocon/db"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	connection.Setup()

	// Routen registrieren
	r.HandleFunc("/login", auth.Login).Methods("POST")
	r.HandleFunc("/initmail", auth.Initmail).Methods("POST")           // email : "email"  -> Emailversand: Verifizierungslink
	r.HandleFunc("/initmail/{code}", auth.Initmail).Methods("GET")     // link + code -> verified Email
	r.HandleFunc("/form/{school}", form.Create).Methods("POST", "GET") // school -> Wireframe

	log.Println("Server läuft auf http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server Fehler:", err)
	}

}
