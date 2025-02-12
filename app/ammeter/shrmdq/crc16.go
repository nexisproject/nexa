// Copyright (C) micros. 2024-present.
//
// Created at 2024-12-11, by liasica

package shrmdq

type CRC16 struct{}

func (*CRC16) CheckSum(raw []byte) uint16 {
	crc := uint16(0xffff)
	offset := 0
	size := len(raw)

	for i := offset; i < size; i++ {
		crc ^= uint16(raw[i])

		for j := 0; j < 8; j++ {
			if crc&0x0001 > 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}

var crc16 = &CRC16{}

func checkSum(data []byte) uint16 {
	return crc16.CheckSum(data)
}

func sumBytes(data []byte) [2]byte {
	sum := checkSum(data)
	return [2]byte{byte(sum), byte(sum >> 8)}
}
