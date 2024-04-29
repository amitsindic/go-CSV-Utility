package handlers

import (
	"mime/multipart"
	"net/http"
	"strings"

	"ctlSolution.com/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	csvService services.CSVService
}

func NewHandler(csvService services.CSVService) *Handler {
	return &Handler{csvService: csvService}
}

func (h *Handler) UploadCSV(c *gin.Context) {
	// Parse multipart form
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// Get the CSV file from the request
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Please upload a CSV file")
		return
	}
	defer file.Close()

	if !IsCSVFile(fileHeader) {
		c.String(http.StatusBadRequest, "File type not supported. Required .csv type")
		return
	}

	// Process CSV data using the CSV service
	err = h.csvService.ProcessCSV(file)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Return success message
	c.String(http.StatusOK, "CSV file uploaded and data inserted into the database")
}

func IsCSVFile(fileHeader *multipart.FileHeader) bool {
	filename := fileHeader.Filename
	return strings.HasSuffix(strings.ToLower(filename), ".csv")
}
