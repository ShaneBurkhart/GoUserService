package user

import (
	"github.com/ShaneBurkhart/GoUserService/config"
	"log"
	"testing"
)

const email = "user@example.com"
const password = "password"

func initDB() {
	err := config.SetupDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := config.VerifyDB(); err != nil {
		log.Fatal(err)
		return
	}
	if _, err := config.DB.Exec("DELETE FROM users"); err != nil {
		log.Fatal(err)
		return
	}
}

func closeDB() {
	config.CloseDB()
}

func createUser(email string, password string) int {
	return Create(email, password)
}

func TestCount(t *testing.T) {
	initDB()

	createUser(email, password)

	if c := Count(); c != 1 {
		t.Error("Expected 1 but got", c)
	}

	closeDB()
}

func TestCreate(t *testing.T) {
	initDB()

	i := Count()
	createUser(email, password)
	if c := Count(); c != i+1 {
		t.Error("User not created. Initial count was", i, "and final count was", c)
	}

	closeDB()
}

func TestPasswordHash(t *testing.T) {
	initDB()

	var digest string
	i := createUser(email, password)
	err := config.DB.QueryRow(`
		SELECT password_digest
		FROM users
		WHERE id = $1
	`, i).Scan(&digest)
	if err != nil {
		t.Error("There was a problem getting user with id", i, "from the database.")
	}

	if digest == password {
		t.Error("The password was not hashed when put into the table.")
	}

	closeDB()
}

func TestFind(t *testing.T) {
	initDB()

	i := createUser(email, password)
	u := Find(i)
	if u == nil || u.Email != email || u.Id != i {
		t.Error("Didn't get correct user. Expected: &{", i, email, "} Got:", u)
	}

	closeDB()
}
