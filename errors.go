package echoex

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

var reAccepts = regexp.MustCompile(`\s*(?P<base>[^;+]+)(\+(?P<suffix>[^;]+))?(\s*;\s*(?P<opts>.*))?\s*`)

// CustomHTTPErrorHandler implements an error handler to format errors according to client's accept, preferring JSON by default.
// JSON response contains "error" key with error description. XML response returns <error> tag with message attribute.
// Additional keys and attributes can be returned in the future.
func CustomHTTPErrorHandler(err error, c echo.Context) {
	var internalErr error
	code := http.StatusInternalServerError
	message := err.Error()

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if msg, ok := he.Message.(string); ok {
			message = msg
		}
		internalErr = he.Internal
	} else if ve, ok := err.(validator.ValidationErrors); ok && len(ve) > 0 {
		message = fmt.Sprintf("Field '%s' failed validation '%s'.", ve[0].Field(), ve[0].ActualTag())
	}

	accepts := strings.Split(strings.ToLower(c.Request().Header.Get("Accept")), ",")
	acceptsAny := false
	acceptsJSON := false
	acceptsXML := false
	for _, a := range accepts {
		m := reAccepts.FindStringSubmatch(a)
		if len(m) < 6 {
			continue
		}

		base := m[1]
		suffix := m[3]
		//opts := m[5]

		switch base {
		case "*/*":
			acceptsAny = true
		case "application/json":
			acceptsJSON = true
		case "application/xml", "text/xml":
			acceptsXML = true
		}

		switch suffix {
		case "json":
			acceptsJSON = true
		case "xml":
			acceptsXML = true
		}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			_ = c.NoContent(code)
		} else {
			type Error struct {
				XMLName struct{} `json:"-" xml:"error"`
				Error   string   `json:"error" xml:"message,attr"`
			}

			errObj := Error{Error: message}

			if acceptsJSON || acceptsAny {
				_ = c.JSON(code, &errObj)
			} else if acceptsXML {
				_ = c.XML(code, &errObj)
			} else {
				_ = c.String(code, fmt.Sprintf("Error: %s\n", message))
			}
		}
	}

	// Log if internal error set
	if internalErr != nil {
		c.Logger().Error(internalErr)
	}
}

func ServerErr(message string, internal ...error) error {
	err := echo.NewHTTPError(http.StatusInternalServerError, message)
	if len(internal) > 0 {
		err = err.SetInternal(internal[0])
	}
	return err
}
