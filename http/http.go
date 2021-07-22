package http

import (
	"net/http"

	"github.com/csv-loader/services"
	"github.com/csv-loader/types"
	"github.com/labstack/echo/v4"
)

type Server struct {
	service services.CSVGetter
}

func New(service services.CSVGetter) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) Read(c echo.Context) error {
	key := c.Param("key")

	keyVal, err := s.service.GetByKey(key)
	if err != nil {
		if e, ok := err.(*types.CustomError); ok {
			return c.JSON(
				e.Code,
				e,
			)
		}

		return c.JSON(
			http.StatusInternalServerError,
			&types.CustomError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			},
		)
	}

	return c.JSON(http.StatusOK, keyVal)
}
