package csvgetter

import (
	"log"
	"net/http"

	"github.com/csv-loader/stores"
	"github.com/csv-loader/types"
)

type Service struct {
	csvGetter stores.CSVGetter
}

func New(csvGetter stores.CSVGetter) *Service {
	return &Service{
		csvGetter: csvGetter,
	}
}

func (s *Service) GetByKey(key string) (*types.KeyValPair, error) {
	if key == "" {
		return nil, &types.CustomError{
			Code:    http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
		}
	}

	keyVal, err := s.csvGetter.GetByKey(key)
	if err != nil {
		log.Printf("Error while getting key from store, %v", err)
		return nil, err
	}

	return keyVal, nil
}
