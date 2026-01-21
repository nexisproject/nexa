// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-20, by liasica

package base

import "errors"

var (
	ErrNameMustStartWithUpper = errors.New("名称必须以大写字母开头")
	ErrUnknownPackageType     = errors.New("未知的包类型")
	ErrFileAlreadyExists      = errors.New("文件已存在, 使用 [-f | --force] 强制覆盖")
)
