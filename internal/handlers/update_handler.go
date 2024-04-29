// internal/handlers/update_handler.go

package handlers

import (
	"net/http"

	"ctlSolution.com/internal/models"
	"ctlSolution.com/internal/services"

	"github.com/gin-gonic/gin"
)

type UpdateHandler struct {
	dbService    services.DBService
	redisService services.RedisService
}

func NewUpdateHandler(dbService services.DBService, redisService services.RedisService) *UpdateHandler {
	return &UpdateHandler{dbService: dbService, redisService: redisService}
}

func (h *UpdateHandler) UpdateData(c *gin.Context) {
	// Get the first_name parameter from the request
	firstName := c.Query("first_name")
	if firstName == "" {
		c.String(http.StatusBadRequest, "Missing first_name parameter")
		return
	}

	// Check if the record exists in the database
	existingRecord, err := h.dbService.GetRecordByFirstName(firstName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve record from the database")
		return
	}
	if existingRecord == nil {
		c.String(http.StatusNotFound, "Record not found in the database")
		return
	}

	var updatedRecord models.Record
	if err := c.BindJSON(&updatedRecord); err != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		return
	}

	// Map updated fields from request to existing record
	getMappedRecord(existingRecord, &updatedRecord)

	// Update the record in the database
	if err := h.dbService.UpdateRecord(existingRecord); err != nil {
		c.String(http.StatusInternalServerError, "Failed to update record in the database")
		return
	}

	// Update the record in Redis (if present)
	if err := h.redisService.CacheRecord(&updatedRecord); err != nil {
		c.String(http.StatusInternalServerError, "Failed to update record in Redis")
		return
	}

	c.String(http.StatusOK, "Record updated successfully")
}



func getMappedRecord(existingRecord *models.Record, updatedRecord *models.Record) {
	// Update the record with the data from the request body
	existingRecord.FirstName = updatedRecord.FirstName
	existingRecord.LastName = updatedRecord.LastName
	existingRecord.CompanyName = updatedRecord.CompanyName
	existingRecord.Address = updatedRecord.Address
	existingRecord.City = updatedRecord.City
	existingRecord.County = updatedRecord.County
	existingRecord.Postal = updatedRecord.Postal
	existingRecord.Phone = updatedRecord.Phone
	existingRecord.Email = updatedRecord.Email
	existingRecord.Web = updatedRecord.Web
}
