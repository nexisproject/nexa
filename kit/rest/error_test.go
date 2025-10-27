// Copyright (C) nexa. 2025-present.
//
// Created at 2025-10-27, by liasica

package rest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"nexis.run/nexa/kit/authz"
)

func TestWrapError(t *testing.T) {
	err := WrapError(http.StatusUnauthorized, authz.ErrUnauthorized)
	require.ErrorIs(t, err, authz.ErrUnauthorized)
}
