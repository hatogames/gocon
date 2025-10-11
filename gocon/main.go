package main

import (
	form "gocon/all"
	auth "gocon/all/auth"
	connection "gocon/db"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()
	connection.Setup()

	// Routen registrieren
	r.HandleFunc("/login", auth.Login).Methods("POST")
	r.HandleFunc("/initmail", auth.Initmail).Methods("POST")
	r.HandleFunc("/form/{school}", form.Create).Methods("POST", "GET")

	// CORS Middleware konfigurieren
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",          // Vue Dev Server
			"https://deine-frontend-url.com", // Produktion
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(r)

	// PORT aus Environment Variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server läuft auf Port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
