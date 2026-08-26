package server

import (
	"database/sql"
	"net/http"
	"net/url"
	"strconv"

	"rss-aggregator/internal/database"
	"rss-aggregator/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
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
func parseListItemsQuery(query url.Values) (listItemsQuery, error) {
	page, err := parsePagination(query)
	if err != nil {
		return listItemsQuery{}, err
	}

	parsed := listItemsQuery{page: page.page, perPage: page.perPage}

	get := func(key string) string {
		return queryValue(query, key)
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
		PageSize:   page.limit(),
		PageOffset: page.offset(),
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

		apiCfg.serveCached(c, "items", func() (any, error) {
			rows, err := apiCfg.queries.ListItems(c, query.params)

			if err != nil {
				return nil, err
			}

			// The channels are few, so one lookup is cheaper than joining on
			// every row and mapping a wider generated type.
			channelTitles, err := apiCfg.channelTitles(c)

			if err != nil {
				return nil, err
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

			return gin.H{
				"items":    items,
				"page":     query.page,
				"per_page": query.perPage,
				"total":    total,
			}, nil
		})
	}
}

// ListItemCategoriesHandler serves the distinct categories, for the filters.
func (apiCfg *ApiConfig) ListItemCategoriesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		apiCfg.serveCached(c, "categories", func() (any, error) {
			rows, err := apiCfg.queries.ListItemCategories(c)

			if err != nil {
				return nil, err
			}

			categories := make([]string, 0, len(rows))

			for _, row := range rows {
				if row.Valid {
					categories = append(categories, row.String)
				}
			}

			return gin.H{"categories": categories}, nil
		})
	}
}

// GetChannelsHandler serves the stored channels.
func (apiCfg *ApiConfig) GetChannelsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		apiCfg.serveCached(c, "channels", func() (any, error) {
			channels, err := apiCfg.queries.GetChannels(c)

			if err != nil {
				return nil, err
			}

			return models.DatabaseChannelsToChannels(channels), nil
		})
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
