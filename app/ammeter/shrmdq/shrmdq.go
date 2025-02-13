// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-16, by liasica

package shrmdq

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"nexis.run/nexa/app/ammeter"
	"nexis.run/nexa/app/ammeter/dlt6452007"
	"nexis.run/nexa/pkg/channel"
)

type Shrmdq struct {
	stop chan struct{}

	// 透传消息等待结果
	// 结构为: [4]byte -> chan *PassthroughMessage
	passthrough sync.Map
}

func New() *Shrmdq {
	return &Shrmdq{
		stop:        make(chan struct{}),
		passthrough: sync.Map{},
	}
}

func (s *Shrmdq) OnOpen(_ *ammeter.Client) {
}

func (s *Shrmdq) OnClose(_ *ammeter.Client, _ error) {
	s.stop <- struct{}{}
}

func (s *Shrmdq) OnMessage(c *ammeter.Client, raw []byte) {
	p, err := NewPacket().Decode(raw)
	if err != nil {
		zap.L().Error("解码数据包失败", zap.Error(err))
		return
	}

	zap.L().Info("收到数据包", zap.Reflect("packet", p))

	if len(p.Data) > 0 {
		var msg Message
		msg, err = p.GetMessage()
		if err != nil {
			zap.L().Error("解析消息失败", zap.Error(err))
			return
		}
		zap.L().Info("消息解包", zap.Reflect("message", msg))

		switch p.Code {
		case CodeRegisterU:
			// TODO: 注册成功
			// 例如透传消息
			// go func() {
			// 	time.Sleep(10 * time.Second)
			// 	// _ = CommandPassthroughAsync(c, dlt6452007.New(msg.(*RegisterMessage).Address, dlt6452007.IForwardPowerTotal, nil).Bytes())
			// 	// err = CommandPassthroughAsync(c, dlt6452007.New(msg.(*RegisterMessage).Address, dlt6452007.IForwardMaxDemand, nil).Bytes())
			// 	// err = CommandPassthroughAsync(c, dlt6452007.New(msg.(*RegisterMessage).Address, dlt6452007.IForwardPowerTotal, nil).Bytes())
			// 	// _ = CommandPassthroughAsync(c, dlt6452007.New(msg.(*RegisterMessage).Address, dlt6452007.IVoltage, nil).Bytes())
			// 	_ = CommandPassthroughAsync(c, dlt6452007.New(msg.(*RegisterMessage).Address, dlt6452007.IHavePowerRate, nil).Bytes())
			// 	// fmt.Println(err)
			// }()

			// 定时抄表 (1h)
			go s.meterReading(c, msg.(*RegisterMessage).Address)

		case CodePassthroughU:
			// 透传消息返回
			pt, ok := s.passthrough.Load(p.Sn)
			if ok {
				s.passthrough.Delete(p.Sn)
				channel.SafeSend(pt.(chan *PassthroughMessage), msg.(*PassthroughMessage))
			}
		}
	}

	b, ok := p.Answer()
	if ok {
		err = c.Send(b)
		if err != nil {
			zap.L().Error("应答失败", zap.Error(err), zap.Uint8("code", uint8(p.Code)))
			return
		}
	}
}

func (s *Shrmdq) meterReading(c *ammeter.Client, address string) {
	// 延时五秒钟
	time.Sleep(5 * time.Second)

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	go func() {
		for ; true; <-ticker.C {
			_, err := s.CommandPassthroughSync(
				c,
				dlt6452007.New(address, dlt6452007.IForwardPowerTotal, nil).Bytes(),
				// dlt6452007.New(address, dlt6452007.IVoltage, nil).Bytes(),
				// dlt6452007.New(address, dlt6452007.IElectricCurrent, nil).Bytes(),
			)
			if err != nil {
				zap.L().Error("抄表失败", zap.Error(err))
			}

			_, err = s.CommandPassthroughSync(
				c,
				dlt6452007.New(address, dlt6452007.IVoltage, nil).Bytes(),
			)
			if err != nil {
				zap.L().Error("电压读取失败", zap.Error(err))
			}

			_, err = s.CommandPassthroughSync(
				c,
				dlt6452007.New(address, dlt6452007.IElectricCurrent, nil).Bytes(),
			)
			if err != nil {
				zap.L().Error("电流读取失败", zap.Error(err))
			}
		}
	}()

	<-s.stop
}
