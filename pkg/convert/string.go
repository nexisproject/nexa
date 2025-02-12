// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package convert

import "unsafe"

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func String2Bytes(s string) (b []byte) {
	return *(*[]byte)(unsafe.Pointer(&s))
}
