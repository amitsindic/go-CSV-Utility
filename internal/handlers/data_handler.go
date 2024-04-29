package handlers

import (
	"net/http"

	"ctlSolution.com/internal/services"
	"github.com/gin-gonic/gin"
)

type DataHandler struct {
	dbService    services.DBService
	redisService services.RedisService
}

func NewDataHandler(dbService services.DBService, redisService services.RedisService) *DataHandler {
	return &DataHandler{dbService: dbService, redisService: redisService}
}

func (h *DataHandler) GetData(c *gin.Context) {

	// If url has path param. Return only one record for incoming param
	key := c.Query("first_name")
	if key != "" {
		data, err := h.redisService.GetDataByKey(key)
		if err == nil && len(data) != 0 {
			c.JSON(http.StatusOK, data)
			return
		}

		// If data not found in Redis, fetch from the database
		records, err := h.dbService.GetRecordByFirstName(key)
		if err != nil || records == nil {
			c.JSON(http.StatusInternalServerError, "Failed to retrieve data from the database")
			return
		}

		// Update cache
		err = h.redisService.CacheRecord(records)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Failed to cache data in Redis")
			return
		}
	}

	// Url param is empty, return list of all records
	// Try to fetch data from Redis
	data, err := h.redisService.GetData()
	if err == nil && len(data) != 0 {
		c.JSON(http.StatusOK, data)
		return
	}

	// If data not found in Redis, fetch from the database
	records, err := h.dbService.GetAllRecords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to retrieve data from the database")
		return
	}

	// If no records found in db : ask user to upload csv first
	if len(records) == 0 {
		c.JSON(http.StatusNoContent, "No records. Please Upload CSV file")
		return
	}

	// Cache each record individually in Redis with an appropriate TTL
	for _, record := range records {
		err = h.redisService.CacheRecord(&record)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Failed to cache data in Redis")
			return
		}
	}

	// Return the data to the client
	c.JSON(http.StatusOK, records)
}
