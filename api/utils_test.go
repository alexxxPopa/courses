package api

import (
	"io"
	"net/http/httptest"
	"github.com/labstack/echo"
)

func (api *API) NewRequest(method string, url string, body io.Reader) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, url, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	api.echo.ServeHTTP(rec, req)

	return rec
}