package http

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/csv-loader/services"
	"github.com/csv-loader/types"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func Test_Read(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCSVGetter := services.NewMockCSVGetter(ctrl)
	mockSrvr := New(mockCSVGetter)

	tcs := []struct {
		key          string
		output       string
		statusCode   int
		expectations []*gomock.Call
	}{
		// Success response when key exist's
		{
			key:        "one",
			output:     `{"key":"one","value":1}`,
			statusCode: http.StatusOK,
			expectations: []*gomock.Call{
				mockCSVGetter.EXPECT().GetByKey(gomock.Any()).Return(
					&types.KeyValPair{Key: "one", Value: 1}, nil,
				),
			},
		},
		// Case when key doesn't exist
		{
			key:        "two",
			output:     `{"code":404,"message":"Not Found"}`,
			statusCode: http.StatusNotFound,
			expectations: []*gomock.Call{
				mockCSVGetter.EXPECT().GetByKey(gomock.Any()).Return(
					nil,
					&types.CustomError{
						Code:    http.StatusNotFound,
						Message: http.StatusText(http.StatusNotFound),
					},
				),
			},
		},
		// Case when csvGetter service returns error
		{
			key:        "three",
			output:     `{"code":500,"message":"Internal Server Error"}`,
			statusCode: http.StatusInternalServerError,
			expectations: []*gomock.Call{
				mockCSVGetter.EXPECT().GetByKey(gomock.Any()).Return(
					nil,
					errors.New("some error occured while fetching the key"),
				),
			},
		},
	}

	for i := range tcs {
		e := echo.New()
		req := httptest.NewRequest("GET", "/csv/"+tcs[i].key, nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		ctx.SetParamNames([]string{"key"}...)
		ctx.SetParamValues([]string{tcs[i].key}...)

		_ = mockSrvr.Read(ctx)

		resp := rec.Result()
		statusCode := resp.StatusCode
		body, _ := ioutil.ReadAll(resp.Body)

		if strings.TrimRight(string(body), "\n") != tcs[i].output || statusCode != tcs[i].statusCode {
			t.Errorf("Test case [%d] failed. Expected [Body = %v, Status Code = %d], Got [Body = %v, Status Code = %d]",
				i+1, tcs[i].output, tcs[i].statusCode, string(body), statusCode)
		}
	}
}
