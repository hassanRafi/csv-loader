package csvextractor

import (
	"github.com/csv-loader/services"
)

type Service struct {
	csvReader services.CSVReader
}

func New(csvReader services.CSVReader) services.CSVExtractor {
	return &Service{
		csvReader: csvReader,
	}
}

func (s *Service) GetRow() ([]string, error) {
	record, err := s.csvReader.Read()
	if err != nil {
		return nil, err
	}

	return record, nil
}
