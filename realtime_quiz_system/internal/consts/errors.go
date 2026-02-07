package consts

import (
	"fmt"
	"net/http"
)

var (
	CodeInvalidToken = customCode{code: 201, message: "Invalid token", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeTokenExpired = customCode{code: 202, message: "Token expired", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeUnauthorized = customCode{code: 203, message: "Unauthorized", detail: nil, httpStatus: http.StatusUnauthorized}
)

type customCode struct {
	code       int
	message    string
	detail     interface{}
	httpStatus int
}

func (c customCode) Code() int {
	return c.code
}

func (c customCode) Message() string {
	return c.message
}

func (c customCode) Detail() interface{} {
	return c.detail
}

func (c customCode) String() string {
	if c.detail != nil {
		return fmt.Sprintf(`%d:%s %v`, c.code, c.message, c.detail)
	}
	if c.message != "" {
		return fmt.Sprintf(`%d:%s`, c.code, c.message)
	}
	return fmt.Sprintf(`%d`, c.code)
}

func (c customCode) HttpStatus() int {
	return c.httpStatus
}
