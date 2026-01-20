// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-19, by liasica

package base

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetModule(t *testing.T) {
	// 测试获取当前项目的模块信息
	module, err := GetModule("../../../..")

	require.NoError(t, err)

	t.Logf("Module: %s", module)

	expected := "nexis.run/nexa"
	require.Equal(t, expected, module)
}

func TestGetModuleNotFound(t *testing.T) {
	// 测试不存在的目录
	_, err := GetModule("/nonexistent")
	require.Error(t, err)
	t.Logf("Expected error: %v", err)
}
