package main

import (
	"log"
	"net/http"

	"github.com/yourusername/yourprojectname/api" // Import the api package
	"github.com/yourusername/yourprojectname/db"  // Replace with your actual module path
	// "github.com/yourusername/yourprojectname/models" // No longer directly used here
)

func main() {
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbConn.Close()

	// Set up HTTP handlers
	http.HandleFunc("/users", api.CreateUserHandler)        // POST /users
	http.HandleFunc("/users/", api.GetUserHandler)          // GET /users/{id}
	http.HandleFunc("/users/update", api.UpdateUserHandler) // PUT /users/update?id={id}
	http.HandleFunc("/users/delete", api.DeleteUserHandler) // DELETE /users/delete?id={id}

	// Use ListenAndServeTLS for HTTPS
	log.Println("Server listening on port 8080 with HTTPS")
	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil))
}
