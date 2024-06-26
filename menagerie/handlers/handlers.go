package handlers

import (
	"database/sql"
	"encoding/json"
	"menagerie/db"
	"menagerie/models"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func GetAllPets(w http.ResponseWriter, r *http.Request) {
	species := r.URL.Query().Get("species")

	var query string
	var rows *sql.Rows
	var err error

	if species != "" {
		query = "SELECT id, name, owner, species, birth, death FROM pets WHERE species = ?"
		rows, err = db.DB.Query(query, species)
	} else {
		query = "SELECT id, name, owner, species, birth, death FROM pets"
		rows, err = db.DB.Query(query)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pets []models.Pet
	for rows.Next() {
		var pet models.Pet
		var birth, death sql.NullString

		if err := rows.Scan(&pet.ID, &pet.Name, &pet.Owner, &pet.Species, &birth, &death); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Convert sql.NullString to *time.Time
		if birth.Valid {
			parsedBirth, err := time.Parse("2006-01-02", birth.String)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			pet.Birth = &parsedBirth
		}
		if death.Valid {
			parsedDeath, err := time.Parse("2006-01-02", death.String)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			pet.Death = &parsedDeath
		}
		pets = append(pets, pet)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pets)
}

func GetPetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	var pet models.Pet
	var birth, death sql.NullString

	err = db.DB.QueryRow("SELECT id, name, owner, species, birth, death FROM pets WHERE id = ?", id).Scan(
		&pet.ID, &pet.Name, &pet.Owner, &pet.Species, &birth, &death)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Pet not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if birth.Valid {
		parsedBirth, err := time.Parse("2006-01-02", birth.String)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pet.Birth = &parsedBirth
	}
	if death.Valid {
		parsedDeath, err := time.Parse("2006-01-02", death.String)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pet.Death = &parsedDeath
	}

	rows, err := db.DB.Query("SELECT id, pet_id, date, type, remark FROM events WHERE pet_id = ? ORDER BY date DESC", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.PetID, &event.Date, &event.Type, &event.Remark); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		events = append(events, event)
	}
	pet.Events = events

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pet)
}

func AddPet(w http.ResponseWriter, r *http.Request) {
	var pet models.Pet
	if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(pet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("INSERT INTO pets (name, owner, species, birth, death) VALUES (?, ?, ?, ?, ?)",
		pet.Name, pet.Owner, pet.Species, pet.Birth, pet.Death)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pet.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pet)
}

func UpdatePet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	var pet models.Pet
	if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(pet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("UPDATE pets SET name = ?, owner = ?, species = ?, birth = ?, death = ? WHERE id = ?",
		pet.Name, pet.Owner, pet.Species, pet.Birth, pet.Death, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pet.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pet)
}

func DeletePet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("DELETE FROM pets WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pet successfully deleted"))
}

func AddEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event.PetID = id
	if err := validate.Struct(event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("SELECT id FROM pets WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Pet not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := db.DB.Exec("INSERT INTO events (pet_id, date, type, remark) VALUES (?, ?, ?, ?)",
		event.PetID, event.Date, event.Type, event.Remark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	event.ID = int(eventID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}
