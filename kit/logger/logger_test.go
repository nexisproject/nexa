// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-06, by liasica

package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	l, err := zap.NewProduction()
	require.NoError(t, err)
	l.Info("test")
	l.Named("xtest").Info("test")

	var ld *zap.Logger
	ld, err = zap.NewDevelopment()
	require.NoError(t, err)
	ld.Info("test")
	ld.Named("xtest").Info("test")
}
