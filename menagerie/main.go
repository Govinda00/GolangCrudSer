package main

import (
	"log"
	"menagerie/db"
	"menagerie/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	r := mux.NewRouter()

	r.HandleFunc("/pets", handlers.GetAllPets).Methods("GET")
	r.HandleFunc("/pets/{id}", handlers.GetPetByID).Methods("GET")
	r.HandleFunc("/pets", handlers.AddPet).Methods("POST")
	r.HandleFunc("/pets/{id}", handlers.UpdatePet).Methods("PUT")
	r.HandleFunc("/pets/{id}", handlers.DeletePet).Methods("DELETE")
	r.HandleFunc("/pets/{id}/events", handlers.AddEvent).Methods("POST")

	log.Println("Server started at :8002")
	if err := http.ListenAndServe(":8002", r); err != nil {
		log.Fatal(err)
	}
}
