// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package dlt6452007

import (
	"fmt"
	"strconv"
)

// Control 控制码
// D7 D6 D5 D4 D3 D2 D1 D0
// D7 传送方向标识
// D6 从站应答标识
// D5 后续帧标识
// D4 D3 D2 D1 D0 功能码
type Control struct {
	Function  Function  // 功能码
	NextFrame NextFrame // 后续帧标识
	Response  Response  // 从站应答标识
	Direction Direction // 传送方向标识
}

// Function 功能码 D4 D3 D2 D1 D0
type Function byte

const (
	CFunctionKeep           Function = 0b00000 // 保留
	CFunctionTiming         Function = 0b01000 // 校时
	CFunctionRead           Function = 0b10001 // 读数据
	CFunctionReadContinue   Function = 0b10010 // 读后续数据
	CFunctionReadAddress    Function = 0b10011 // 读通信地址
	CFunctionWrite          Function = 0b10100 // 写数据
	CFunctionWriteAddress   Function = 0b10101 // 写通信地址
	CFunctionFreeze         Function = 0b10110 // 冻结命令
	CFunctionRate           Function = 0b10111 // 修改通信速率
	CFunctionPassword       Function = 0b11000 // 修改密码
	CFunctionClearMaxDemand Function = 0b11001 // 最大需量清零
	CFunctionClearMeter     Function = 0b11010 // 电表清零
	CFunctionClearEvent     Function = 0b11011 // 事件清零
)

// NextFrame D5 后续帧标识
type NextFrame byte

const (
	CNextEnd      NextFrame = iota // 无后续帧
	CNextContiune                  // 有后续帧
)

// Response D6 从站应答标识
type Response byte

const (
	CResponseSuccess Response = iota // 正确应答
	CResponseError                   // 异常应答
)

// Direction D7 传送方向标识
type Direction byte

const (
	CDirectionMaster Direction = iota // 主站发出的命令帧
	CDirectionSlave                   // 从站发出的应答帧
)

func ControlFromByte(raw byte) Control {
	return Control{
		Function:  Function(raw & 0x1F),
		NextFrame: NextFrame(raw >> 5 & 0x01),
		Response:  Response(raw >> 6 & 0x01),
		Direction: Direction(raw >> 7 & 0x01),
	}
}

func (c Control) Byte() byte {
	return byte(c.Direction)<<7 | byte(c.Response)<<6 | byte(c.NextFrame)<<5 | byte(c.Function)
}

func (c Control) String() string {
	return fmt.Sprintf("%d%d%d%s", c.Direction, c.Response, c.NextFrame, strconv.FormatUint(uint64(c.Function), 2))
}
