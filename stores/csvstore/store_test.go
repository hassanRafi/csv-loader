package csvstore

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/csv-loader/types"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
)

func Test_GetByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client, mock := redismock.NewClientMock()

	mockStore := New(client)

	tcs := []struct {
		key    string
		err    error
		output *types.KeyValPair
	}{
		// Success case of getting a key
		{
			key:    "one",
			err:    nil,
			output: &types.KeyValPair{Key: "one", Value: 1},
		},
	}

	for i := range tcs {
		mock.ExpectGet(tcs[i].key).SetVal("1")

		output, err := mockStore.GetByKey(tcs[i].key)
		if !reflect.DeepEqual(err, tcs[i].err) {
			t.Errorf("Test case [%d] failed. Expected %v, Got %v",
				i+1, tcs[i].err, err)
		}

		if !reflect.DeepEqual(tcs[i].output, output) {
			t.Errorf("Test case [%d] failed. Expected %v, Got %v",
				i+1, tcs[i].output, output)
		}
	}
}

func Test_GetByKey_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client, mock := redismock.NewClientMock()

	mockStore := New(client)

	tcs := []struct {
		key    string
		err    error
		output *types.KeyValPair
	}{
		// Case when there is some error while getting a key
		{
			key: "eight",
			err: errors.New("error while getting a key"),
		},
	}

	for i := range tcs {
		mock.ExpectGet(tcs[i].key).SetErr(errors.New("error while getting a key"))

		output, err := mockStore.GetByKey(tcs[i].key)
		if !reflect.DeepEqual(err, tcs[i].err) {
			t.Errorf("Test case [%d] failed. Expected %v, Got %v",
				i+1, tcs[i].err, err)
		}

		if !reflect.DeepEqual(tcs[i].output, output) {
			t.Errorf("Test case [%d] failed. Expected %v, Got %v",
				i+1, tcs[i].output, output)
		}
	}
}

func Test_GetByKey_KeyNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client, mock := redismock.NewClientMock()

	mockStore := New(client)

	tcs := []struct {
		key    string
		err    error
		output *types.KeyValPair
	}{
		// Case when the key doesn't exist
		{
			key: "eleven",
			err: &types.CustomError{
				Code:    http.StatusNotFound,
				Message: http.StatusText(http.StatusNotFound),
			},
		},
	}

	for i := range tcs {
		mock.ExpectGet(tcs[i].key).SetErr(redis.Nil)

		output, err := mockStore.GetByKey(tcs[i].key)
		if !reflect.DeepEqual(err, tcs[i].err) {
			t.Errorf("Test case [%d] failed. Expected %v, Got %v",
				i+1, tcs[i].err, err)
		}

		if !reflect.DeepEqual(tcs[i].output, output) {
			t.Errorf("Test case [%d] failed. Expected %v, Got %v",
				i+1, tcs[i].output, output)
		}
	}
}
