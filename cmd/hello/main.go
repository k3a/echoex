package main

import (
	"errors"
	"net/http"
	"github.com/k3a/echoex"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/*
curl  -H 'Accept: application/xml' -d 'email=test@gmail.com' -d 'str=mystring' -d 'query=inform' -d 'form=inform' 'localhost:8888/inpath?query=inquery&form=inquery&int=123' | less

results in:
{
  "form": "inform",
  "ok": "good!",
  "params": {
    "Str": "mystring",
    "Int": 123,
    "email": "test@gmail.com",
    "Form": "inform",
    "Query": "inquery",
    "Path": "inpath"
  },
  "path": "inpath",
  "query": "inquery"
}
*/

func main() {
	e := echoex.New()
	e.Use(middleware.BodyLimit("32M") /*, middleware.Logger()*/, middleware.Recover())

	h := func(c echo.Context) (err error) {
		var params struct {
			Str   string
			Int   int
			Email string `json:"email" validate:"required,email"`
			Form  string // `form:"ren"` to rename
			Query string `form:"-"` // to require query
			Path  string //`param:"ren"` to rename
		}
		if err = c.Bind(&params); err != nil {
			return
		}
		if err = c.Validate(&params); err != nil {
			return
		}

		if params.Str == "err" {
			return errors.New("common go error returned")
		} else if params.Str == "interr" {
			return echoex.ServerErr("Some server error happened, sorry.", errors.New("internal technical error details"))
		}

		//c.String(http.StatusOK, "Good")
		return c.JSONPretty(http.StatusOK, echo.Map{
			"ok":     "good!",
			"form":   c.FormValue("form"),
			"query":  c.QueryParam("query"),
			"path":   c.Param("path"),
			"params": params,
		}, "  ")
	}

	e.POST("/", h)
	e.POST("/:path", h)

	e.Logger.Fatal(e.Start(":8888"))
}
