// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-04, by liasica

package rest

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func NewError(code int, message string) *Error {
	err := &Error{
		Code:    code,
		Message: message,
	}
	if err.Message == "" {
		err.Message = http.StatusText(code)
	}
	return err
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewResponse() *Response {
	return &Response{
		Code:    http.StatusOK,
		Message: "",
	}
}

// SetCode 设置code
func (r *Response) SetCode(code int) *Response {
	r.Code = code
	return r
}

// SetMessage 设置message
func (r *Response) SetMessage(message string) *Response {
	r.Message = message
	return r
}

// SetData 设置data
func (r *Response) SetData(data any) *Response {
	if !reflect.ValueOf(data).IsNil() {
		r.Data = data
	}
	return r
}

// SetParams 设置响应参数
func (r *Response) SetParams(params ...any) *Response {
	for i := 0; i < len(params); i++ {
		switch v := params[i].(type) {
		case int:
			r.SetCode(v)
		case string:
			r.SetMessage(v)
		case *Error:
			r.SetCode(v.Code).SetMessage(v.Message)
		case error:
			message := v.Error()
			var he *echo.HTTPError
			if errors.As(v, &he) {
				message = fmt.Sprintf("%v", he.Message)
			}
			r.SetMessage(message)
		default:
			if r.Data == nil {
				r.SetData(v)
			}
		}
	}

	if r.Message == "" {
		r.Message = http.StatusText(r.Code)
	}

	return r
}
