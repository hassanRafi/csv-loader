package services

import "github.com/csv-loader/types"

type CSVLoader interface {
	Store(row []*types.KeyValPair) error
}

type CSVExtractor interface {
	GetRow() ([]string, error)
}

type CSVReader interface {
	Read() ([]string, error)
}

type CSVGetter interface {
	GetByKey(key string) (*types.KeyValPair, error)
}
