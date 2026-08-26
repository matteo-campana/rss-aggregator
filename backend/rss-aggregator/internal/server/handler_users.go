package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"rss-aggregator/internal/database"
	"time"

	"rss-aggregator/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) CreateUserHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		// parse and check the request parameters to create a new user

		type parameters struct {
			Fullname  string `json:"fullname"`
			Firstname string `json:"firstname"`
			Lastname  string `json:"lastname"`
			Email     string `json:"email"`
		}

		decoder := json.NewDecoder(c.Request.Body)

		params := parameters{}

		err := decoder.Decode(&params)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if params.Fullname == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fullname is required"})
			return
		}

		user, err := apiCfg.queries.CreateUser(c, database.CreateUserParams{
			Fullname:  params.Fullname,
			Firstname: sql.NullString{String: params.Firstname, Valid: true},
			Lastname:  sql.NullString{String: params.Lastname, Valid: true},
			Email:     sql.NullString{String: params.Email, Valid: true},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			ID:        uuid.New(),
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// The key is only returned here: it is the one chance the caller has
		// to store it.
		c.JSON(http.StatusCreated, models.DatabaseUserToUserWithApiKey(user))
	}

}

// GetCurrentUserHandler returns the user the request authenticated as.
func (apiCfg *ApiConfig) GetCurrentUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		user, ok := UserFromContext(c)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseUserToUserWithApiKey(user))
	}
}

func (apiCfg *ApiConfig) GetUsersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		page, err := parsePagination(c.Request.URL.Query())

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		rows, err := apiCfg.queries.ListUsers(c, database.ListUsersParams{
			PageSize:   page.limit(),
			PageOffset: page.offset(),
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		total := int64(0)
		users := make([]models.User, 0, len(rows))

		for _, row := range rows {
			total = row.TotalCount
			users = append(users, models.DatabaseListUsersRowToUser(row))
		}

		c.JSON(http.StatusOK, gin.H{
			"users":    users,
			"page":     page.page,
			"per_page": page.perPage,
			"total":    total,
		})
	}

}

func (apiCfg *ApiConfig) GetUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Params.ByName("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		id_uudi, err := uuid.Parse(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid uuid"})
			return
		}

		user, err := apiCfg.queries.GetUser(c, id_uudi)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseUserToUser(user))
	}
}

func (apiCfg *ApiConfig) UpdateUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Params.ByName("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		id_uudi, err := uuid.Parse(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid uuid"})
			return
		}

		type parameters struct {
			Fullname  string `json:"fullname"`
			Firstname string `json:"firstname"`
			Lastname  string `json:"lastname"`
			Email     string `json:"email"`
		}

		decoder := json.NewDecoder(c.Request.Body)

		params := parameters{}

		err = decoder.Decode(&params)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if params.Fullname == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fullname is required"})
			return
		}

		user, err := apiCfg.queries.UpdateUser(c, database.UpdateUserParams{
			Fullname:  params.Fullname,
			Firstname: sql.NullString{String: params.Firstname, Valid: true},
			Lastname:  sql.NullString{String: params.Lastname, Valid: true},
			Email:     sql.NullString{String: params.Email, Valid: true},
			UpdatedAt: time.Now().UTC(),
			ID:        id_uudi,
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseUserToUser(user))
	}
}

func (apiCfg *ApiConfig) DeleteUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Params.ByName("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		id_uudi, err := uuid.Parse(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid uuid"})
			return
		}

		user, err := apiCfg.queries.GetUser(c, id_uudi)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		err = apiCfg.queries.DeleteUser(c, id_uudi)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message": "user deleted",
			"user":    models.DatabaseUserToUser(user),
		})
	}
}
