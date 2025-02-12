// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package dlt6452007

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddress(t *testing.T) {
	data := []byte{0x77, 0x32, 0x09, 0x19, 0x04, 0x22}
	address := "220419093277"
	require.Equal(t, address, AddressFromBytes(data))
	require.Equal(t, data, AddressToBytes(address))
	x := AddressToBytes("22041909327")
	t.Logf("%#v", x)
	x = AddressToBytes("220419093207")
	t.Logf("%#v", x)
}
