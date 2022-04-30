package wrapper

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	DataField    = "data"
	TraceIDField = "trace_id"
	SuccessField = "success"
	// CodeField     = "code"
	MessageField  = "message"
	TotalField    = "total"
	MetadataField = "metadata"
)

// Response body
type Response struct {
	Error        error
	Data         interface{}
	Status       int
	Total        int64
	IncludeTotal bool
}

type EchoHandlerFunc func(c echo.Context) Response

// Wrap return new gin.HandlerFunc by GinHandlerFn
func Wrap(fn EchoHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// handle request
		res := fn(c)
		return Translate(c, res)
	}
}

// Translate response
func Translate(c echo.Context, res Response) error {
	result := map[string]interface{}{
		SuccessField: true,
		//TraceIDField: c.Get(constants.RequestIDContextKey),
	}

	status := http.StatusOK
	if res.Error != nil {
		result[MessageField] = res.Error.Error()
		result[SuccessField] = false
	}

	// get data
	if res.Data != nil {
		result[DataField] = res.Data
	}

	includeMetadata := res.IncludeTotal
	if includeMetadata {
		meta := map[string]interface{}{}
		if res.IncludeTotal {
			meta[TotalField] = res.Total
		}

		result[MetadataField] = meta
	}

	if res.Status > 0 {
		status = res.Status
	}

	return c.JSON(status, result)
}
