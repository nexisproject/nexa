// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package dlt6452007

import "fmt"

// AddressFromBytes 从字节转换为地址字符串
func AddressFromBytes(raw []byte) string {
	return fmt.Sprintf("%02x%02x%02x%02x%02x%02x", raw[5], raw[4], raw[3], raw[2], raw[1], raw[0])
}

// AddressToBytes 从地址字符串转换为字节
func AddressToBytes(str string) []byte {
	var data [6]byte
	_, _ = fmt.Sscanf(str, "%02x%02x%02x%02x%02x%02x", &data[5], &data[4], &data[3], &data[2], &data[1], &data[0])
	return data[:]
}
