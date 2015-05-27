package main

import (
	"github.com/go-errors/errors"
	"github.com/labstack/echo"
	"net/http"
)

type GenericError struct {
	echo.HTTPError
	StackTrace string `json:"StackTrace,omitempty"`
}

func AzureErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c *echo.Context) {
		ge := new(GenericError)
		ge.Code = http.StatusInternalServerError //default status code is 500
		ge.Message = http.StatusText(ge.Code)    // default message is 'Internal Server Error'
		switch error := err.(type) {
		case *errors.Error:
			if he, ok := error.Err.(*echo.HTTPError); ok {
				ge.Code = he.Code
				ge.Message = he.Message
			}
			if e.Debug() && ge.Code >= 500 {
				ge.StackTrace = error.ErrorStack()
			}
		case *echo.HTTPError:
			ge.Code = error.Code
			ge.Message = error.Message
		case error:
			if e.Debug() && ge.Code >= 500 {
				ge.Message = err.Error() //show original error message in case of debug mode https://github.com/labstack/echo/blob/1e117621e9006481bfc0fd8e6bafab48c1848639/echo.go#L161
			}
		}

		c.JSON(ge.Code, ge)
	}
}