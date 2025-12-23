package models

import "github.com/traceway/go-lightning/lmy"

/*
	If you don't have this table already you should create it:
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL
	);
*/

type User struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
}

func init() {
	lmy.Register[User]()
}
