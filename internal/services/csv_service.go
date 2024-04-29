package services

import (
	"encoding/csv"
	"io"

	//"log"

	"ctlSolution.com/internal/models"
)

type CSVService interface {
	ProcessCSV(file io.Reader) error
}

type csvService struct {
	dbService    DBService
	redisService RedisService
}

func NewCSVService(db DBService, redis RedisService) CSVService {
	return &csvService{dbService: db, redisService: redis}
}

func (s *csvService) ProcessCSV(file io.Reader) error {
	//var records []models.Record

	// Initialize CSV reader
	reader := csv.NewReader(file)

	// Skip header row
	_, err := reader.Read()
	if err != nil {
		return err
	}

	// Iterate over CSV records
	for {
		// Read a single record
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// Convert record to struct
		// Assuming column order matches struct fields order
		// Update the index if the order is different
		r := models.Record{
			FirstName:   record[0],
			LastName:    record[1],
			CompanyName: record[2],
			Address:     record[3],
			City:        record[4],
			County:      record[5],
			Postal:      record[6],
			Phone:       record[7],
			Email:       record[8],
			Web:         record[9],
		}

		// Create a Record object
		//records = append(records, r)

		//caching
		err = s.redisService.CacheRecord(&r)
		if err != nil {
			return err
		}

		//insert in db
		_ = s.dbService.Insert(&r)
	}

	return nil
}
