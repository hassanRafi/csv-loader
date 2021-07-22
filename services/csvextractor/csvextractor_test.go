package csvextractor

import (
	"errors"
	"reflect"
	"testing"

	"github.com/csv-loader/services"
	"github.com/golang/mock/gomock"
)

func Test_GetRow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCSVReader := services.NewMockCSVReader(ctrl)
	mockCSVExtractor := New(mockCSVReader)

	tcs := []struct {
		output       []string
		err          error
		expectations []*gomock.Call
	}{
		// Successful read from csv file
		{
			output: []string{"one", "1"},
			err:    nil,
			expectations: []*gomock.Call{
				mockCSVReader.EXPECT().Read().Return(
					[]string{"one", "1"}, nil,
				),
			},
		},
		// Case when read from csv returns an error
		{
			output: nil,
			err:    errors.New("some error occured while reading"),
			expectations: []*gomock.Call{
				mockCSVReader.EXPECT().Read().Return(
					nil, errors.New("some error occured while reading"),
				),
			},
		},
	}

	for i := range tcs {
		row, err := mockCSVExtractor.GetRow()
		if !reflect.DeepEqual(err, tcs[i].err) {
			t.Errorf("Test case [%d] failed. Expected %v, Got %v",
				i+1, tcs[i].err, err)
		}

		if err != nil && !reflect.DeepEqual(row, tcs[i].output) {
			t.Errorf("Test case [%d] failed. Expected %v, Got %v",
				i+1, tcs[i].output, row)
		}
	}
}
