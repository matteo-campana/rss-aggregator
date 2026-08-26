package server

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"rss-aggregator/internal/database"
	"rss-aggregator/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	defaultPageSize = 50
	maxPageSize     = 200

	sortRecent  = "recent"
	sortSeeders = "seeders"
	sortOldest  = "oldest"
)

// listItemsQuery is the parsed and validated form of the query string.
type listItemsQuery struct {
	params  database.ListItemsParams
	page    int
	perPage int
}

// parseListItemsQuery validates the query string. It is a pure function so the
// validation can be tested without a database.
func parseListItemsQuery(query map[string][]string) (listItemsQuery, error) {
	get := func(key string) string {
		if values, ok := query[key]; ok && len(values) > 0 {
			return strings.TrimSpace(values[0])
		}

		return ""
	}

	parsed := listItemsQuery{page: 1, perPage: defaultPageSize}

	if raw := get("page"); raw != "" {
		page, err := strconv.Atoi(raw)
		if err != nil || page < 1 {
			return listItemsQuery{}, &badRequestError{"page must be a positive integer"}
		}

		parsed.page = page
	}

	if raw := get("per_page"); raw != "" {
		perPage, err := strconv.Atoi(raw)
		if err != nil || perPage < 1 {
			return listItemsQuery{}, &badRequestError{"per_page must be a positive integer"}
		}

		if perPage > maxPageSize {
			perPage = maxPageSize
		}

		parsed.perPage = perPage
	}

	sort := get("sort")
	switch sort {
	case "":
		sort = sortRecent
	case sortRecent, sortSeeders, sortOldest:
	default:
		return listItemsQuery{}, &badRequestError{"sort must be one of: recent, seeders, oldest"}
	}

	params := database.ListItemsParams{
		Sort:       sort,
		PageSize:   int32(parsed.perPage),
		PageOffset: int32((parsed.page - 1) * parsed.perPage),
	}

	if search := get("search"); search != "" {
		params.Search = sql.NullString{String: search, Valid: true}
	}

	if category := get("category"); category != "" {
		params.Category = sql.NullString{String: category, Valid: true}
	}

	if raw := get("min_seeders"); raw != "" {
		minSeeders, err := strconv.Atoi(raw)
		if err != nil || minSeeders < 0 {
			return listItemsQuery{}, &badRequestError{"min_seeders must be a non-negative integer"}
		}

		params.MinSeeders = sql.NullInt32{Int32: int32(minSeeders), Valid: true}
	}

	if raw := get("channel_id"); raw != "" {
		channelID, err := uuid.Parse(raw)
		if err != nil {
			return listItemsQuery{}, &badRequestError{"channel_id is not a valid UUID"}
		}

		params.ChannelID = uuid.NullUUID{UUID: channelID, Valid: true}
	}

	parsed.params = params

	return parsed, nil
}

type badRequestError struct {
	message string
}

func (e *badRequestError) Error() string {
	return e.message
}

// itemResponse adds the source channel to the stored item.
type itemResponse struct {
	models.Item
	ChannelTitle string `json:"channel_title"`
}

// ListItemsHandler serves the items the scraper has already stored, filtered,
// sorted and paginated.
func (apiCfg *ApiConfig) ListItemsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		query, err := parseListItemsQuery(c.Request.URL.Query())

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		rows, err := apiCfg.queries.ListItems(c, query.params)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// The channels are few, so one lookup is cheaper than joining on every
		// row and mapping a wider generated type.
		channelTitles, err := apiCfg.channelTitles(c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		total := int64(0)
		items := make([]itemResponse, 0, len(rows))

		for _, row := range rows {
			total = row.TotalCount

			items = append(items, itemResponse{
				Item:         models.DatabaseListItemsRowToItem(row),
				ChannelTitle: channelTitles[row.ChannelID],
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"items":    items,
			"page":     query.page,
			"per_page": query.perPage,
			"total":    total,
		})
	}
}

// ListItemCategoriesHandler serves the distinct categories, for the filters.
func (apiCfg *ApiConfig) ListItemCategoriesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		rows, err := apiCfg.queries.ListItemCategories(c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		categories := make([]string, 0, len(rows))

		for _, row := range rows {
			if row.Valid {
				categories = append(categories, row.String)
			}
		}

		c.JSON(http.StatusOK, gin.H{"categories": categories})
	}
}

// GetChannelsHandler serves the stored channels.
func (apiCfg *ApiConfig) GetChannelsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		channels, err := apiCfg.queries.GetChannels(c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models.DatabaseChannelsToChannels(channels))
	}
}

func (apiCfg *ApiConfig) channelTitles(c *gin.Context) (map[uuid.UUID]string, error) {
	channels, err := apiCfg.queries.GetChannels(c)
	if err != nil {
		return nil, err
	}

	titles := make(map[uuid.UUID]string, len(channels))

	for _, channel := range channels {
		titles[channel.ID] = channel.Title
	}

	return titles, nil
}
