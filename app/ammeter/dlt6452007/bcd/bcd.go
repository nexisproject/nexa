// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-15, by liasica

package bcd

import (
	"strconv"
	"unsafe"
)

// ToString BCD码转换为字符串
// 字符串长度为BCD码长度的2倍
// BCD码的高位在前，低位在后
// 例如：[]byte{0x78, 0x56, 0x34, 0x12} -> "12345678"
func ToString(b []byte) (result string) {
	if len(b) == 0 {
		return
	}

	j := 0
	out := make([]byte, len(b)*2)
	for i := len(b) - 1; i >= 0; i-- {
		c0 := b[i] & 0xF
		c1 := (b[i] >> 4) & 0xF

		if c1 <= 9 {
			out[j] = c1 + '0'
		} else {
			out[j] = c1 + 'A'
		}
		j++

		if c0 <= 9 {
			out[j] = c0 + '0'
		} else {
			out[j] = c0 + 'A'
		}
		j++
	}
	return *(*string)(unsafe.Pointer(&out))
}

// FromString 字符串转换为BCD码
// BCD码的高位在前，低位在后
// 字符串中的字符必须是0-9、A-F、a-f
// 例如："12345678" -> []byte{0x78, 0x56, 0x34, 0x12}
// 字符串必须是偶数个字符，否则最后一个字符会被忽略
// 例如："1234567" -> []byte{0x56, 0x34, 0x12}
func FromString(str string) (result []byte) {
	var tmpValue byte
	var i, j, m int
	var sLen int

	n := len(str) / 2
	result = make([]byte, n)

	sLen = len(str)
	for i = 0; i < sLen; i++ {
		if (str[i] < '0') ||
			((str[i] > '9') && (str[i] < 'A')) ||
			((str[i] > 'F') && (str[i] < 'a')) ||
			(str[i] > 'f') {
			sLen = i
			break
		}
	}

	if sLen > n*2 {
		sLen = n * 2
	}

	for i, j, m = sLen-1, 0, 0; (i >= 0) && (m < n); i, j = i-1, j+1 {
		switch {
		case str[i] >= '0' && str[i] <= '9':
			tmpValue = str[i] - '0'
		case str[i] >= 'A' && str[i] <= 'F':
			tmpValue = str[i] - 'A' + 0x0A
		case str[i] >= 'a' && str[i] <= 'f':
			tmpValue = str[i] - 'a' + 0x0A
		default:
			tmpValue = 0
		}

		if j%2 == 0 {
			result[m] = tmpValue
		} else {
			result[m] |= tmpValue << 4
			m++
		}

		if tmpValue == 0 && str[i] != '0' {
			break
		}
	}

	return
}

func pow100(power byte) uint64 {
	res := uint64(1)
	for i := byte(0); i < power; i++ {
		res *= 100
	}
	return res
}

// ToUint64X BCD码转换为十进制数
// BCD码的高位在前，低位在后
// 例如：[]byte{0x56, 0x34, 0x12} -> 123456
// 最大支持位数 20, 最大数为 18446744073709551615
func ToUint64X(value []byte) (result uint64) {
	vlen := len(value)
	for i, b := range value {
		hi, lo := b>>4, b&0x0f
		if hi > 9 || lo > 9 {
			return 0
		}
		result += uint64(hi*10+lo) * pow100(byte(vlen-i)-1)
	}
	return
}

// ToUint64 BCD码转换为十进制数
// BCD码的高位在前，低位在后
// 例如：[]byte{0x56, 0x34, 0x12} -> 123456
// 最大支持位数 20, 最大数为 18446744073709551615
func ToUint64(b []byte) (result uint64) {
	for i := len(b) - 1; i >= 0; i-- {
		result += uint64(b[i]>>4)*10 + uint64(b[i]&0x0f)
		if i != 0 {
			result *= 100
		}
	}
	return
}

// FromUint64 十进制数转换为BCD码
// BCD码的高位在前，低位在后
// 例如：123456 -> []byte{0x56, 0x34, 0x12}
// 最大支持位数 20, 最大数为 18446744073709551615
// 数字长度必须是偶数位，否则最后一个数字会被忽略
// 例如：1234567 -> []byte{0x56, 0x34, 0x12}
func FromUint64(value uint64) []byte {
	str := strconv.FormatUint(value, 10)
	size := len(str) / 2

	buf := make([]byte, size)
	if value > 0 {
		remainder := value
		for pos := 0; pos < size && remainder > 0; pos++ {
			tail := byte(remainder % 100)
			hi, lo := tail/10, tail%10
			buf[pos] = hi<<4 + lo
			remainder = remainder / 100
		}
	}

	return buf
}
