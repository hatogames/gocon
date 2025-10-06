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

	/*
		owner := connection.Owner{
			Email: "mail",
			Phash: "pass",
		}

		connection.DB.Create(&owner)

		school := connection.School{
			OwnerID: owner.ID,
		}

		school2 := connection.School{
			OwnerID: owner.ID,
		}

		connection.DB.Create(&school)
		connection.DB.Create(&school2)

		us := connection.User{
			Email: "mail",
			Phash: "pass",
		}

		connection.DB.Create(&us)

		reg := connection.Registration{
			SchoolID: 1,
			UserID:   us.ID,
		}

		connection.DB.Create(&reg)
	*/

	/*
		owner := connection.Owner{
			Email: "nl@hg-ka.de",
			Phash: "pass",
		}
		if err := connection.DB.Create(&owner).Error; err != nil {
			fmt.Println("Fehler beim Erstellen des Owners:", err)
			return
		}

		fmt.Println("Owner erstellt:", owner)

		// 2️⃣ School erstellen
		school := connection.School{
			OwnerID: owner.ID,
		}
		if err := connection.DB.Create(&school).Error; err != nil {
			fmt.Println("Fehler beim Erstellen der School:", err)
			return
		}

		fmt.Println("School erstellt:", school)

		// 3️⃣ Wireframe erstellen
		wireframe := connection.Wireframe{
			Name:     "create",
			SchoolID: school.ID,
			Keys:     datatypes.JSON([]byte(`["create"]`)),
			Data:     datatypes.JSON([]byte(`{"example":"value"}`)),
		}
		if err := connection.DB.Create(&wireframe).Error; err != nil {
			fmt.Println("Fehler beim Erstellen des Wireframes:", err)
			return
		}
	*/

	// Routen registrieren
	r.HandleFunc("/login", auth.Login).Methods("POST")
	r.HandleFunc("/initmail", auth.Initmail).Methods("POST")
	r.HandleFunc("/form/{school}", form.Create).Methods("POST", "GET")

	log.Println("Server läuft auf http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server Fehler:", err)
	}

}
