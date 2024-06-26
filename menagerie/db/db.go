package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "root:Govinda@123@tcp(localhost:3306)/servicaGolang")
	if err != nil {
		log.Fatal(err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database successfully")

	// Create Pets table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS pets (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		owner VARCHAR(255) NOT NULL,
		species VARCHAR(255) NOT NULL,
		birth DATE,
		death DATE
	)`
	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	// Create Events table if it doesn't exist
	query = `
	CREATE TABLE IF NOT EXISTS events (
		id INT AUTO_INCREMENT PRIMARY KEY,
		pet_id INT NOT NULL,
		date DATE NOT NULL,
		type VARCHAR(255) NOT NULL,
		remark TEXT,
		FOREIGN KEY (pet_id) REFERENCES pets(id)
	)`
	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
