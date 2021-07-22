package redisloader

import (
	"errors"
	"reflect"
	"testing"

	"github.com/csv-loader/types"
	redismock "github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
)

func Test_Store(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client, mock := redismock.NewClientMock()

	mockRedisLoaderSvc := New(client)

	tcs := []struct {
		pairs []*types.KeyValPair
		err   error
	}{
		// Successful case of setting keys
		{
			pairs: []*types.KeyValPair{
				{
					Key:   "one",
					Value: 1,
				},
				{
					Key:   "two",
					Value: 2,
				},
			},
			err: nil,
		},
	}

	for i := range tcs {
		mock.ExpectMSet("one", 1, "two", 2).SetVal("ok")
		mock.ExpectMSet("one", 1, "two", 2).SetErr(nil)

		if err := mockRedisLoaderSvc.Store(tcs[i].pairs); err != nil {
			if !reflect.DeepEqual(err, tcs[i].err) {
				t.Errorf("Test case [%d] failed. Expected %v, Got %v",
					i+1, tcs[i].err, err)
			}
		}
	}
}

func Test_Store_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client, mock := redismock.NewClientMock()

	mockRedisLoaderSvc := New(client)

	tcs := []struct {
		pairs []*types.KeyValPair
		err   error
	}{
		// Case when some error occur's while setting the keys
		{
			pairs: []*types.KeyValPair{
				{
					Key:   "five",
					Value: 5,
				},
				{
					Key:   "six",
					Value: 6,
				},
			},
			err: errors.New("some error occured"),
		},
	}

	for i := range tcs {
		mock.ExpectMSet("five", 5, "six", 6).SetErr(errors.New("some error occured"))

		if err := mockRedisLoaderSvc.Store(tcs[i].pairs); err != nil {
			if !reflect.DeepEqual(err, tcs[i].err) {
				t.Errorf("Test case [%d] failed. Expected %v, Got %v",
					i+1, tcs[i].err, err)
			}
		}
	}
}
