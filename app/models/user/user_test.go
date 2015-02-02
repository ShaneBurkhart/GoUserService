package user

// TODO Validate email has an @
// TODO Validate email is unique

import (
	"errors"
	"fmt"
	"github.com/ShaneBurkhart/GoUserService/config"
	"log"
	"testing"
)

const email = "user@example.com"
const password = "password"
const id = 1

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
	clearDB()
}

func clearDB() {
	if _, err := config.DB.Exec("DELETE FROM users"); err != nil {
		log.Fatal(err)
		return
	}
}

func closeDB() {
	config.CloseDB()
}

func createDefaultUser() (int, error) {
	return Create(email, password)
}

var countTests = []struct {
	email string
}{
	{"user@domain.com"},
	{"example@domain.com"},
	{"me@example.com"},
	{"you@example.com"},
}

func TestCount(t *testing.T) {
	initDB()
	defer closeDB()

	for i, test := range countTests {
		Create(test.email, password)
		if c := Count(); c != i+1 {
			t.Error("Incorrect count. Expected:", i+1, "Got:", c)
		}
	}
}

func TestFind(t *testing.T) {
	initDB()
	defer closeDB()
	i, _ := createDefaultUser()

	u := Find(i)

	err := verifyUser(u, i, email, password)
	if err != nil {
		t.Error(err)
	}

	// Returns nil when doesn't exist
	u = Find(i + 1)
	if u != nil {
		t.Error("User should be nil. Expected:", nil, "Got:", u)
	}
}

var createTests = []struct {
	email       string
	password    string
	shouldError bool
}{
	{"email@domain.com", "password", false},
	{"email@domain.com", "foo", true}, // Duplicate email
	{"user@anotherdomain.com", "H3llo", false},
	{"useranotherdomain.com", "H3lloWorld", true},  // Invalid email format
	{"@useranotherdomain.com", "H3lloWorld", true}, // Invalid email format
	{"useranotherdomain@", "H3lloWorld", true},     // Invalid email format
	{"", "H3lloWorld", true},
	{"user@anotherdomain.com", "", true},
}

func TestCreate(t *testing.T) {
	initDB()
	defer closeDB()

	for _, test := range createTests {
		initialCount := Count()

		i, err := Create(test.email, test.password)

		if err != nil && test.shouldError {
			continue
		}
		if err != nil && !test.shouldError {
			t.Error("There was an error when there shouldn't have been. Error:",
				err, "Test:", test)
			continue
		}
		if err == nil && test.shouldError {
			t.Error("There wasn't an error when there should have been. Error:",
				err, "Test:", test)
			continue
		}

		u := Find(i)
		err = verifyUser(u, i, test.email, test.password)
		if err != nil {
			t.Error(err)
		}

		if c := Count(); c != initialCount+1 {
			t.Error("Initial user count was", initialCount,
				"and final user count was", c)
		}
	}
}

func verifyUser(u *User, id int, email string, password string) error {
	if u == nil {
		return errors.New("User is nil.")
	}
	if u.Id != id {
		return errors.New(fmt.Sprintf("Ids don't match. Expected: %d, Got: %d", id, u.Id))
	}
	if u.Email != email {
		return errors.New(fmt.Sprintf("Emails don't match. Expected: %s, Got: %s", email, u.Email))
	}
	if err := comparePassword(u, password); err != nil {
		return errors.New("Passwords don't match.")
	}
	return nil
}
