// Copyright (C) micros. 2025-present.
//
// Created at 2025-02-11, by liasica

package micro

import "context"

func NewContext(app string) *Context {
	return &Context{
		Context: context.WithValue(context.Background(), "app", app),
		App:     app,
	}
}
