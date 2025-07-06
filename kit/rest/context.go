// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-04, by liasica

package rest

import (
	"bytes"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"

	"nexis.run/nexa/kit"
)

type ContextWrapper interface {
	BindValidate(ptr any)
}

// Context Rest服务上下文
type Context struct {
	App string

	echo.Context
}

// NewContext 创建上下文
func NewContext(app string, c echo.Context) *Context {
	return &Context{
		App:     app,
		Context: c,
	}
}

// GetContext 获取上下文
func GetContext(c echo.Context) *Context {
	switch v := c.(type) {
	case *Context:
		return v
	default:
		return NewContext("UNKNOWN", c)
	}
}

// GetContextWrapper 获取上下文
func GetContextWrapper[C ContextWrapper](c echo.Context) C {
	ctx, ok := c.(C)
	if !ok {
		panic(NewError(http.StatusInternalServerError, kit.ErrInvalidContext.Error()))
	}
	return ctx
}

// BindValidate 绑定并校验
func (c *Context) BindValidate(ptr any) {
	err := c.Bind(ptr)
	if err != nil {
		panic(NewError(http.StatusBadRequest, err.Error()))
	}
	err = c.Validate(ptr)
	if err != nil {
		panic(NewError(http.StatusBadRequest, err.Error()))
	}
}

// BaseContextBinding 获取上下文并绑定
func BaseContextBinding[T any](c echo.Context) (ctx *Context, req *T) {
	ctx = GetContext(c)
	req = new(T)
	ctx.BindValidate(req)
	return
}

func ContextBinding[C ContextWrapper, T any](c echo.Context) (ctx C, req *T) {
	ctx = GetContextWrapper[C](c)
	req = new(T)
	ctx.BindValidate(req)
	return
}

// SendResponse 发送响应
func (c *Context) SendResponse(params ...any) error {
	buffer := &bytes.Buffer{}
	encoder := jsoniter.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(NewResponse().SetParams(params...))

	return c.JSONBlob(http.StatusOK, buffer.Bytes())
}
