package models

import "github.com/traceway/go-lightning/lpg"

/*
	If you don't have this table already you should create it:
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	);
*/

type User struct {
	Id        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

func init() {
	lpg.Register[User]()
}
