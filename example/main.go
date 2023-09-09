package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/muhwyndhamhp/tmax"
)

func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.HTTPErrorHandler = httpErrorHandler
	tmax.NewEchoTemplateRenderer(e, "root", "body", "public/components", "public/views")

	e.GET("/", func(c echo.Context) error {
		if isHXRequest(c) {
			return c.Render(http.StatusOK, "index", nil)
		} else {
			return c.Render(http.StatusOK, "root#index", nil)
		}
	})

	e.GET("/content", func(c echo.Context) error {
		if isHXRequest(c) {
			return c.Render(http.StatusOK, "content", nil)
		} else {
			return c.Render(http.StatusOK, "root#content", nil)
		}
	})

	e.Logger.Fatal(e.Start(":4002"))
}

func isHXRequest(c echo.Context) bool {
	hx_request, err := strconv.ParseBool(c.Request().Header.Get("Hx-Request"))
	if err != nil {
		return false
	}

	return hx_request
}

func httpErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	if code != http.StatusInternalServerError {
		_ = c.JSON(code, err)
	} else {
		log.Error(err)
		_ = c.JSON(http.StatusInternalServerError, err)
	}
}
