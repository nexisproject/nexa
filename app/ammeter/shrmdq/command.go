// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-12, by liasica

package shrmdq

import (
	"time"

	"go.uber.org/zap"

	"nexis.run/nexa/app/ammeter"
	"nexis.run/nexa/pkg/channel"
)

const (
	CommandDefaultKeep = 0x01
)

// CommandSignal 查询信号强度
func CommandSignal(c *ammeter.Client) error {
	return c.Send(NewPacket().SetCode(CodeSignalD).SetSn().SetKeep(CommandDefaultKeep).Build())
}

// CommandPassthroughAsync 异步透传指令
// items: 透传指令
// 返回值为透传指令的SN
func (*Shrmdq) CommandPassthroughAsync(c *ammeter.Client, items ...[]byte) (sn [4]byte, err error) {
	if len(items) == 0 {
		err = ammeter.ErrorInvalidDataLength
		return
	}

	// 拼接数据
	b := NewMontage(nil).
		AddDecomposed([]byte{byte(len(items) & 0xFF)}). // 透传指令数量
		AddDecomposed(items...).                        // 透传指令
		Compose().
		GetComposed()

	p := NewPacket().SetCode(CodePassthroughD).SetSn().SetKeep(CommandDefaultKeep).SetData(b)

	return p.Sn, c.Send(p.Build())
}

// CommandPassthroughSync 同步透传指令
// items: 透传指令
// 返回透传结果
func (s *Shrmdq) CommandPassthroughSync(c *ammeter.Client, items ...[]byte) (*PassthroughMessage, error) {
	sn, err := s.CommandPassthroughAsync(c, items...)
	if err != nil {
		return nil, nil
	}

	ch := make(chan *PassthroughMessage)
	s.passthrough.Store(sn, ch)

	time.AfterFunc(10*time.Second, func() {
		if _, ok := s.passthrough.Load(sn); ok {
			s.passthrough.Delete(sn)
			zap.L().Error("透传指令超时，关闭channel", zap.Binary("sn", sn[:]))
			channel.SafeClose(ch)
		}
	})

	return <-ch, nil
}
