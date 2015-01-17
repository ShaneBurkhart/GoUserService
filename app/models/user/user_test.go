package user

// TODO Validate email has an @
// TODO Validate email is unique

import (
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

func createDefaultUser() int {
	return Create(email, password)
}

var countTests = []struct {
	in  int
	out int
}{
	{1, 1},
	{3, 3},
}

func TestCount(t *testing.T) {
	initDB()
	defer closeDB()

	for _, test := range countTests {
		clearDB()
		for i := 0; i < test.in; i++ {
			createDefaultUser()
		}
		if c := Count(); c != test.out {
			t.Error("Incorrect count. Expected:", test.out, "Got:", c)
		}
	}
}

func TestFind(t *testing.T) {
	initDB()
	defer closeDB()
	i := createDefaultUser()

	u := Find(i)

	if u == nil {
		t.Error("Didn't find user with Id:", i)
	} else {
		if u.Email != email {
			t.Error("Emails don't match. Expected:", email, "Got:", u.Email)
		}
		if u.Id != i {
			t.Error("Ids don't match. Expected:", i, "Got:", u.Id)
		}
		if err := comparePassword(u, password); err != nil {
			t.Error("Passwords don't match.")
		}
	}

	// Returns nil when doesn't exist
	u = Find(i + 1)
	if u != nil {
		t.Error("User should be nil. Expected:", nil, "Got:", u)
	}
}

var createTests = []struct {
	email    string
	password string
}{
	{"email@domain.com", "password"},
	{"user@anotherdomain.com", "H3llo"},
}

func TestCreate(t *testing.T) {
	initDB()
	defer closeDB()

	for _, test := range createTests {
		clearDB()
		initialCount := Count()

		i := Create(test.email, test.password)

		u := Find(i)
		if u == nil {
			t.Error("User didn't get created.")
		} else {
			if u.Email != test.email {
				t.Error("Emails don't match. Expected:", test.email, "Got:", u.Email)
			}
			if err := comparePassword(u, test.password); err != nil {
				t.Error("Password wasn't properly hashed.")
			}
		}

		if c := Count(); c != initialCount+1 {
			t.Error("Initial user count was", initialCount, "and final user count was", c)
		}
	}
}
