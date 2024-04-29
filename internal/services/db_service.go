package services

import (
	"ctlSolution.com/internal/models"
	"github.com/jinzhu/gorm"
)

type DBService interface {
	Insert(record *models.Record) error
	BulkInsert(record []models.Record) error
	GetAllRecords() ([]models.Record, error)
	DeleteRecordByFirstName(firstName string) error
	UpdateRecord(record *models.Record) error
	GetRecordByFirstName(firstName string) (*models.Record, error)
}

type dbService struct {
	db *gorm.DB
}

func NewDBService(db *gorm.DB) DBService {
	return &dbService{db: db}
}

func (s *dbService) Insert(record *models.Record) error {
	err := s.db.Create(record).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *dbService) BulkInsert(record []models.Record) error {
	err := s.db.Create(record).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *dbService) GetAllRecords() ([]models.Record, error) {
	var records []models.Record
	if err := s.db.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (s *dbService) GetRecordByFirstName(firstName string) (*models.Record, error) {
	var record models.Record
	if err := s.db.Where("first_name = ?", firstName).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *dbService) UpdateRecord(record *models.Record) error {
	query := "UPDATE records SET last_name=?, company_name=?, address=?, city=?, county=?, postal=?, phone=?, email=?, web=? WHERE first_name=?"
	args := []interface{}{
		record.LastName,
		record.CompanyName,
		record.Address,
		record.City,
		record.County,
		record.Postal,
		record.Phone,
		record.Email,
		record.Web,
		record.FirstName,
	}

	// Execute the SQL query
	if err := s.db.Exec(query, args...).Error; err != nil {
		return err
	}
	return nil
}

func (s *dbService) DeleteRecordByFirstName(firstName string) error {
	if err := s.db.Where("first_name = ?", firstName).Delete(&models.Record{}).Error; err != nil {
		return err
	}
	return nil
}
