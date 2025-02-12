// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package convert

func Reverse[S ~[]E, E any](s S) (result S) {
	result = make(S, len(s))
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = s[j], s[i]
	}
	return result
}
