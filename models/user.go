package models

// User represents a user in the system.
type User struct {
	ID    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}
