// Copyright (C) micros. 2025-present.
//
// Created at 2025-02-11, by liasica

package micro

import "context"

type Context struct {
	context.Context

	Name string
}

func NewContext(name string) {

}
