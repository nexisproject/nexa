// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-20, by liasica

package command

import "strings"

func examples(ex ...string) string {
	for i := range ex {
		ex[i] = "  " + ex[i] // indent each row with 2 spaces.
	}
	return strings.Join(ex, "\n")
}
