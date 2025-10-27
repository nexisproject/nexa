// Copyright (C) nexa. 2025-present.
//
// Created at 2025-09-07, by liasica

package rest

// MiddlewareKey 定义中间件使用的key
type MiddlewareKey = string

const (
	MiddlewareKeyDumpSkip MiddlewareKey = "DUMP_SKIP" // 跳过dump, value bool
)
