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

// UserWithApiKey exposes the API key alongside the user. It is only served to
// the owner of the key: registration and GET /users/me.
type UserWithApiKey struct {
	User
	ApiKey string `json:"api_key"`
}

func DatabaseUserToUserWithApiKey(dbUser database.User) UserWithApiKey {
	return UserWithApiKey{
		User:   DatabaseUserToUser(dbUser),
		ApiKey: dbUser.ApiKey,
	}
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

// DatabaseListUsersRowToUser maps a row of the paginated listing. It goes
// through DatabaseUserToUser, so the API key stays out of the response.
func DatabaseListUsersRowToUser(row database.ListUsersRow) User {
	return DatabaseUserToUser(database.User{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		Fullname:  row.Fullname,
		Firstname: row.Firstname,
		Lastname:  row.Lastname,
		Email:     row.Email,
		ApiKey:    row.ApiKey,
	})
}
