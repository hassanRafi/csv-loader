package csvbuilder

import (
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/csv-loader/services"
	"github.com/csv-loader/types"
)

type Service struct {
	csvExtractor services.CSVExtractor
	csvLoader    map[string]services.CSVLoader
	workers      int
	chunkSize    int
}

func New(csvExtractor services.CSVExtractor, csvLoader map[string]services.CSVLoader, workers int, chunkSize int) *Service {
	return &Service{
		csvExtractor: csvExtractor,
		csvLoader:    csvLoader,
		workers:      workers,
		chunkSize:    chunkSize,
	}
}

type dataChan struct {
	err  error
	data *types.KeyValPair
}

func (s *Service) BuildCSV() {
	ch := s.getCSVChan()

	log.Printf("Starting %d workers", s.workers)

	wg := &sync.WaitGroup{}
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go s.populateStore(ch, wg)
	}

	wg.Wait()

	log.Printf("Successfully loaded the csv file")
}

func (s *Service) getCSVChan() <-chan dataChan {
	ch := make(chan dataChan)

	go func() {
		defer close(ch)

		// Skipping the first row
		_, _ = s.csvExtractor.GetRow()

		for {
			row, err := s.csvExtractor.GetRow()
			if err != nil {
				if err == io.EOF {
					break
				}

				log.Printf("Got an error while reading, error: %v", err)
			}

			value, err := strconv.Atoi(row[1])
			if err == nil {
				ch <- dataChan{err: nil, data: &types.KeyValPair{Key: row[0], Value: value}}
			} else {
				log.Printf("Can't convert the value to int, error: %s", err)
			}
		}
	}()

	return ch
}

func (s *Service) populateStore(ch <-chan dataChan, wg *sync.WaitGroup) {
	var data []*types.KeyValPair
	defer wg.Done()

	for d := range ch {
		data = append(data, d.data)

		// Store in chunks of 50
		if len(data) == s.chunkSize {
			for dataStore := range s.csvLoader {
				if err := s.csvLoader[dataStore].Store(data); err != nil {
					log.Printf("Failed to store the entries: %v", data)
				}
			}

			data = []*types.KeyValPair{}
		}
	}

	// Populate leftover data
	if len(data) != 0 {
		for dataStore := range s.csvLoader {
			if err := s.csvLoader[dataStore].Store(data); err != nil {
				log.Printf("Failed to store the entries: %v", data)
			}
		}
	}
}
