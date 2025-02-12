// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package dlt6452007

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestControl(t *testing.T) {
	// b := byte(0x91)
	b := byte(0x11)
	control := ControlFromByte(b)
	t.Logf("control: %v", control)
	x := control.Byte()
	require.Equal(t, b, x)
}
