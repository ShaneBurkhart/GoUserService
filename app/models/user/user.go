package user

import (
	"database/sql"
	"github.com/ShaneBurkhart/GoUserService/config"
	"golang.org/x/crypto/bcrypt"
	"log"
)

const countSQL = `
	SELECT COUNT(*)
	FROM users`

const findByIdSQL = `
	SELECT id, email
	FROM users
	WHERE id = $1`

const insertUserSQL = `
	INSERT INTO users(
		email,
		password_digest
	) VALUES (
		$1, $2
	) RETURNING id`

type User struct {
	Id    int
	Email string
}

func Count() int {
	var c int
	err := config.DB.QueryRow(countSQL).Scan(&c)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	return c
}

func Find(id int) *User {
	u := User{}
	err := config.DB.QueryRow(findByIdSQL, id).Scan(&u.Id, &u.Email)
	if err != nil && err != sql.ErrNoRows {
		//TODO What to do if no rows
		//TODO Logging
		log.Fatal(err)
	}
	return &u

}

func Create(email string, password string) int {
	// TODO bcrypt password hashing
	var id int
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	err = config.DB.QueryRow(insertUserSQL, email, passwordHash).Scan(&id)
	if err != nil {
		// TODO Logging
		log.Fatal(err)
	}
	return id
}
