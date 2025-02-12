// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-10, by liasica

package ammeter

import (
	"errors"
	"fmt"
	"os"

	"github.com/panjf2000/gnet/v2"
	"go.uber.org/zap"
)

type Delegate interface {
	OnOpen(c *Client)
	OnClose(c *Client, err error)
	OnMessage(c *Client, raw []byte)
}

type MessageHandler func(c *Client, raw []byte)

type Handler struct {
	gnet.BuiltinEventEngine

	codec    Codec
	delegate Delegate
}

func NewHandler(codec Codec, delegate Delegate) *Handler {
	return &Handler{
		codec:    codec,
		delegate: delegate,
	}
}

func (h *Handler) OnBoot(_ gnet.Engine) (action gnet.Action) {
	return gnet.None
}

func (h *Handler) OnOpen(conn gnet.Conn) (out []byte, action gnet.Action) {
	c := NewClient(conn, h)

	zap.L().Info("新增客户端连接", zap.String("addr", conn.RemoteAddr().String()))
	go h.delegate.OnOpen(c)

	// 设置连接上下文信息
	conn.SetContext(c)
	return
}

func (h *Handler) OnClose(conn gnet.Conn, err error) (action gnet.Action) {
	// TODO: 关闭连接后续处理
	zap.L().Info("客户端断开连接", zap.Error(err))
	c, ok := conn.Context().(*Client)
	if ok {
		c.Close()
		go h.delegate.OnClose(c, err)
	}
	return
}

func (h *Handler) OnTraffic(conn gnet.Conn) (action gnet.Action) {
	// 获取客户端
	c, ok := conn.Context().(*Client)
	if !ok {
		return gnet.Shutdown
	}

	var (
		b   []byte
		err error
	)

	for {
		b, err = h.codec.Decode(conn)

		if len(b) > 0 {
			if os.Getenv("LOCAL_DEV") == "true" {
				fmt.Printf("%d\t%s\n", len(b), b)
			}
		}

		if errors.Is(err, ErrorIncompletePacket) {
			break
		}

		if err != nil {
			zap.L().Error("消息读取失败", zap.Error(err))
			return
		}

		// 处理消息
		if len(b) > 0 {
			go h.delegate.OnMessage(c, b)
		}
	}

	return
}
