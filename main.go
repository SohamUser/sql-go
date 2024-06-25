package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// HTTP route handling
	http.HandleFunc("/form", handleForm)
	http.HandleFunc("/api/addmovie", handleAddMovieAPI)

	// CORS middleware
	corsHandler := corsMiddleware(http.DefaultServeMux)

	// Start the server with CORS middleware
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

// handleForm handles requests to display the form
func handleForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/form.html")
}

// Movie struct to represent JSON data
type Movie struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

// handleAddMovieAPI handles requests to add a movie to the database via API
func handleAddMovieAPI(w http.ResponseWriter, r *http.Request) {
	// Replace 'username', 'password', 'dbname' with your MySQL credentials and database name
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/db_name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Allow CORS (Cross-Origin Resource Sharing)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Decode JSON request body into Movie struct
	var movie Movie
	err = json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Example insert query
	insertQuery := "INSERT INTO movies (title, director, year) VALUES (?, ?, ?)"
	result, err := db.Exec(insertQuery, movie.Title, movie.Director, movie.Year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the ID of the inserted row
	insertID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with JSON indicating success
	response := map[string]interface{}{
		"message":    "Movie added successfully",
		"insertedID": insertID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// corsMiddleware adds CORS headers to HTTP requests
func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Allow options method to handle preflight requests
		if r.Method == "OPTIONS" {
			return
		}

		handler.ServeHTTP(w, r)
	})
}
