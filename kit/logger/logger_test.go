// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-06, by liasica

package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	l, _ := zap.NewProduction()
	l.Info("test")
	l.Named("xtest").Info("test")

	ld, _ := zap.NewDevelopment()
	ld.Info("test")
	ld.Named("xtest").Info("test")
}
