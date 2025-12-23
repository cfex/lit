package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"userscrudmysql/models"

	"github.com/tracewayapp/go-lightning/lmy"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/crud")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v\n", err)
	}
	defer tx.Rollback()

	insertUser(tx)
	readUser(tx)
	readMultipleUsers(tx)
	updateUser(tx)
	deleteUser(tx)

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit: %v\n", err)
	}

	fmt.Println("CRUD example completed successfully!")
}

// this will insert into the users table
// the insert is based on the struct name
// you can override this by changing the lmy.NamingStrategy
func insertUser(tx *sql.Tx) {
	id, err := lmy.InsertGeneric(tx, &models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     fmt.Sprintf("john.doe.%d@example.com", os.Getpid()),
	})
	if err != nil {
		log.Fatalf("Insert failed: %v\n", err)
	}
	fmt.Printf("Inserted user with ID: %d\n", id)
}

func readUser(tx *sql.Tx) {
	foundUser, err := lmy.SelectGenericSingle[models.User](tx, "SELECT * FROM users WHERE id = ?", 1)
	if err != nil {
		log.Fatalf("Select single failed: %v\n", err)
	}
	if foundUser != nil {
		fmt.Printf("Found user: %s %s (%s)\n", foundUser.FirstName, foundUser.LastName, foundUser.Email)
	}
}
func readMultipleUsers(tx *sql.Tx) {
	users, err := lmy.SelectGeneric[models.User](tx, "SELECT * FROM users LIMIT 5")
	if err != nil {
		log.Fatalf("Select multiple failed: %v\n", err)
	}
	fmt.Printf("Found %d users in total\n", len(users))
	for _, u := range users {
		fmt.Printf(" - %d: %s %s\n", u.Id, u.FirstName, u.LastName)
	}
}
func updateUser(tx *sql.Tx) {
	foundUser, err := lmy.SelectGenericSingle[models.User](tx, "SELECT * FROM users WHERE id = ?", 1)
	if err != nil {
		log.Fatalf("Select single failed: %v\n", err)
	}
	// if no user is found this will error
	foundUser.FirstName = "Jane"
	err = lmy.UpdateGeneric(tx, foundUser, "id = ?", 1)
	if err != nil {
		log.Fatalf("Update failed: %v\n", err)
	}
	fmt.Println("Updated user name to Jane")
}
func deleteUser(tx *sql.Tx) {
	err := lmy.Delete(tx, "DELETE FROM users WHERE id = ?", 1)
	if err != nil {
		log.Fatalf("Delete failed: %v\n", err)
	}
	fmt.Printf("Deleted user %d\n", 1)
}
