package stores

import "github.com/csv-loader/types"

type CSVGetter interface {
	GetByKey(key string) (*types.KeyValPair, error)
}
