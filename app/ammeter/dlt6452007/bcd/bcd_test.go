// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-15, by liasica

package bcd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertString(t *testing.T) {
	b := []byte{0x77, 0x32, 0x09, 0x19, 0x04, 0x22}
	str := "220419093277"
	require.Equal(t, str, ToString(b))
	require.Equal(t, b, FromString(str))
}

func TestConvertInt64(t *testing.T) {
	b := []byte{0x78, 0x56, 0x34, 0x12}
	var d uint64 = 12345678
	t.Logf("b = %#v (%08d)", b, b)
	t.Logf("d = %d", d)
	require.Equal(t, d, ToUint64(b))
	require.Equal(t, b, FromUint64(d))
}
