package main

import (
	form "gocon/all"
	auth "gocon/all/auth"
	auth_all "gocon/auth"
	school "gocon/auth/school"
	"gocon/auth/school/wireframe"
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

	r.HandleFunc("/registrations", school.Registrations).Methods("GET")
	r.HandleFunc("/wireframe/update", wireframe.Update).Methods("POST")
	r.HandleFunc("/wireframe/load", wireframe.Load).Methods("POST")

	r.HandleFunc("/whati", auth_all.Whati).Methods("GET")

	// CORS Middleware konfigurieren
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://localhost:8080",
			"https://deine-frontend-url.com",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // ← wichtig!
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
