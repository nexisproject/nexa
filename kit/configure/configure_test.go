// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-25, by liasica

package configure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	type test struct {
		Configure
		Version string
	}

	c, err := Load[test]("../../tests/config.yaml")
	require.NoError(t, err)
	require.NotNil(t, c)

	require.Equal(t, "v1.0.0", c.Version)
	require.Equal(t, "test-app", c.AppName)
}
