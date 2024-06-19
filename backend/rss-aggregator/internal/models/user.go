package models

import (
	"rss-aggregator/internal/database"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Fullname  string    `json:"fullname"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
}

func DatabaseUserToUser(dbUser database.User) User {

	firstname := ""
	if dbUser.Firstname.Valid {
		firstname = dbUser.Firstname.String
	}

	lastname := ""
	if dbUser.Lastname.Valid {
		lastname = dbUser.Lastname.String
	}

	email := ""
	if dbUser.Email.Valid {
		email = dbUser.Email.String
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Fullname:  dbUser.Fullname,
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
	}

	return user
}

func DatabaseUsersToUsers(dbUsers []database.User) []User {
	users := []User{}

	for _, dbUser := range dbUsers {
		user := DatabaseUserToUser(dbUser)
		users = append(users, user)
	}

	return users
}
