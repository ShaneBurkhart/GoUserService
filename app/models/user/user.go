package user

import (
	"database/sql"
	"errors"
	"github.com/ShaneBurkhart/GoUserService/config"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
)

const countSQL = `
	SELECT COUNT(*)
	FROM users`

const findByIdSQL = `
	SELECT *
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
	Id             int
	Email          string
	PasswordDigest string
}

func Count() int {
	var c int
	err := config.DB.QueryRow(countSQL).Scan(&c)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func Find(id int) *User {
	u := User{}
	err := config.DB.QueryRow(findByIdSQL, id).Scan(&u.Id, &u.Email, &u.PasswordDigest)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		//TODO Logging
		log.Fatal(err)
	}
	return &u

}

func Create(email string, password string) (int, error) {
	if err := validate(email, password); err != nil {
		return 0, err
	}

	var id int
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	err = config.DB.QueryRow(insertUserSQL, email, passwordHash).Scan(&id)
	if err != nil {
		// Most likely this error is due to index constraints not being met.
		// In this case, it is most likely that the email is not unique.
		// TODO Logging
		return 0, err
	}
	return id, nil
}

func validate(email string, password string) error {
	if matched, _ := regexp.MatchString("^\\S+@\\S+$", email); !matched {
		return errors.New("Email is invalid.")
	}
	if len(password) == 0 {
		return errors.New("Password cannot be blank.")
	}
	return nil
}

func comparePassword(u *User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password))
}
