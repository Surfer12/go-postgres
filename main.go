package main

import (
	"fmt"
	"log"
	"my-go-postgres-project/db"
	"my-go-postgres-project/models"
)

func main() {
	// Initialize the database connection
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Example usage: Create a user
	newUser := models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	createdUser, err := db.CreateUser(dbConn, newUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created user: %+v\n", createdUser)

	// Example usage: Get a user by ID
	retrievedUser, err := db.GetUserByID(dbConn, createdUser.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Retrieved user: %+v\n", retrievedUser)

	// Example usage: Update a user
	updatedUser := models.User{
		ID:       retrievedUser.ID,
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	rowsAffected, err := db.UpdateUser(dbConn, updatedUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Rows affected by update: %d\n", rowsAffected)

	// Example usage: Get all users.
	allUsers, err := db.GetAllUsers(dbConn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("All users: %+v\n", allUsers)

	// Example usage: Delete a user
	err = db.DeleteUser(dbConn, createdUser.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User deleted successfully")
}
