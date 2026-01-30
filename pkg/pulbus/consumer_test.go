// Copyright (C) nexa. 2026-present.
//
// Created at 2026-01-30, by liasica

package pulbus

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConsumer(t *testing.T) {
	var m sync.Map

	key1 := ConsumerKey{
		Topic:        "test-Topic-1",
		Subscription: "test-Subscription-1",
	}
	key2 := ConsumerKey{
		Topic:        "test-Topic-2",
		Subscription: "test-Subscription-2",
	}

	consumer1 := "test-consumer"
	consumer2 := "test-consumer-2"

	m.Store(key1, consumer1)
	m.Store(key2, consumer2)

	value, ok := m.Load(key1)
	require.True(t, ok)
	require.Equal(t, value, consumer1)

	value, ok = m.Load(key2)
	require.True(t, ok)
	require.Equal(t, value, consumer2)
}

func TestConsumerOption(t *testing.T) {
	options := &ConsumerOptions{}

	WithConsumerChannelSize(100)(options)

	require.Equal(t, options.channelSize, 100)
}
