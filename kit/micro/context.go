// Copyright (C) nexa. 2025-present.
//
// Created at 2025-07-06, by liasica

package micro

import "context"

type Context struct {
	context.Context

	App string
}

func NewContext(app string) *Context {
	return &Context{
		Context: context.WithValue(context.Background(), "app", app),
		App:     app,
	}
}
