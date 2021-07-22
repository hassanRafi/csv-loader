package types

import "fmt"

type KeyValPair struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func (k *KeyValPair) String() string {
	return fmt.Sprintf("Key = [%s], Value = [%d]", k.Key, k.Value)
}

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("Code = [%d], Message = [%s]", ce.Code, ce.Message)
}
