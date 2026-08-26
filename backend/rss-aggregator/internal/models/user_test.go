package models

import (
	"database/sql"
	"testing"

	"rss-aggregator/internal/database"

	"github.com/google/uuid"
)

func TestDatabaseUserToUserNullFields(t *testing.T) {
	user := DatabaseUserToUser(database.User{
		ID:       uuid.New(),
		Fullname: "Matteo Campana",
	})

	if user.Fullname != "Matteo Campana" {
		t.Errorf("fullname = %q, want %q", user.Fullname, "Matteo Campana")
	}

	if user.Firstname != "" || user.Lastname != "" || user.Email != "" {
		t.Errorf("optional fields = %q/%q/%q, want empty strings", user.Firstname, user.Lastname, user.Email)
	}
}

func TestDatabaseUsersToUsers(t *testing.T) {
	users := DatabaseUsersToUsers([]database.User{
		{ID: uuid.New(), Fullname: "a", Email: sql.NullString{String: "a@example.com", Valid: true}},
		{ID: uuid.New(), Fullname: "b"},
	})

	if got, want := len(users), 2; got != want {
		t.Fatalf("len = %d, want %d", got, want)
	}

	if got, want := users[0].Email, "a@example.com"; got != want {
		t.Errorf("email = %q, want %q", got, want)
	}
}
