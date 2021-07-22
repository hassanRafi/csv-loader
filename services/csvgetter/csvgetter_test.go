package csvgetter

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/csv-loader/services"
	"github.com/csv-loader/types"
	"github.com/golang/mock/gomock"
)

func Test_GetByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCSVGetter := services.NewMockCSVGetter(ctrl)
	mockCSVGetterSvc := New(mockCSVGetter)

	tcs := []struct {
		key          string
		output       *types.KeyValPair
		err          error
		expectations []*gomock.Call
	}{
		// Success case
		{
			key: "two",
			output: &types.KeyValPair{
				Key:   "two",
				Value: 2,
			},
			expectations: []*gomock.Call{
				mockCSVGetter.EXPECT().GetByKey("two").Return(
					&types.KeyValPair{
						Key:   "two",
						Value: 2,
					}, nil,
				),
			},
		},
		// Case when key is empty
		{
			key: "",
			err: &types.CustomError{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
			},
		},
		// Case when store returns error
		{
			key: "three",
			err: errors.New("failed to get value for given key"),
			expectations: []*gomock.Call{
				mockCSVGetter.EXPECT().GetByKey("three").Return(
					nil, errors.New("failed to get value for given key"),
				),
			},
		},
	}

	for i := range tcs {
		output, err := mockCSVGetterSvc.GetByKey(tcs[i].key)
		if !reflect.DeepEqual(err, tcs[i].err) {
			t.Errorf("Test case [%d] failed. Expected %v, Got %v",
				i+1, tcs[i].err, err)
		}

		if err != nil && !reflect.DeepEqual(output, tcs[i].output) {
			t.Errorf("Test case [%d] failed. Expected %v, Got %v",
				i+1, tcs[i].output, output)
		}
	}
}
