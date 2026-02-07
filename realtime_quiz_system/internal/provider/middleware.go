package provider

import (
	"net/http"

	"github.com/gogf/gf/v2/net/ghttp"
)

// CustomResponseHandler is a custom middleware for handling API responses
// It wraps successful responses with code: 200 instead of code: 0
func CustomResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()

	// Skip if response already written
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		msg  string
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = http.StatusOK
	)

	if err != nil {
		code = http.StatusInternalServerError
		msg = err.Error()
	} else {
		msg = "OK"
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    code,
		Message: msg,
		Data:    res,
	})
}
