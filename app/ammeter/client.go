// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-16, by liasica

package ammeter

import (
	"log"

	"github.com/panjf2000/gnet/v2"

	"nexis.run/nexa/pkg/dump"
)

type Client struct {
	gnet.Conn

	Handler *Handler
}

func NewClient(conn gnet.Conn, handler *Handler) (c *Client) {
	c = &Client{
		Conn:    conn,
		Handler: handler,
	}
	return
}

// Send 发送数据
// 使用编码器对数据进行编码
func (c *Client) Send(data []byte) (err error) {
	var n int
	n, err = c.Conn.Write(c.Handler.codec.Encode(data))
	log.Printf("<<<< N: %d, DATA: %s\n", n, dump.Bytes(data))
	return
}

func (c *Client) Close() {

}
