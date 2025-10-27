// Copyright (C) nexa. 2025-present.
//
// Created at 2025-10-27, by liasica

package micro

type Option struct {
	middlewares []Handler
}

type OptionFunc func(*Option)
