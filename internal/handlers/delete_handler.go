// internal/handlers/delete_handler.go

package handlers

import (
	"net/http"

	"ctlSolution.com/internal/services"

	"github.com/gin-gonic/gin"
)

type DeleteHandler struct {
	dbService    services.DBService
	redisService services.RedisService
}

func NewDeleteHandler(dbService services.DBService, redisService services.RedisService) *DeleteHandler {
	return &DeleteHandler{dbService: dbService, redisService: redisService}
}

func (h *DeleteHandler) DeleteData(c *gin.Context) {
	// Get the first_name parameter from the request
	firstName := c.Query("first_name")
	if firstName == "" {
		c.String(http.StatusBadRequest, "Missing first_name parameter")
		return
	}

	// Delete the record from the database
	if err := h.dbService.DeleteRecordByFirstName(firstName); err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete record from the database")
		return
	}

	// Delete the record from Redis (if present)
	if err := h.redisService.DeleteRecordByFirstName(firstName); err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete record from Redis")
		return
	}

	c.String(http.StatusOK, "Record deleted successfully")
}
