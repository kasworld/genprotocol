// Code generated by "genprotocol -ver=1.0 -prefix=c2s -basedir example"

package c2s_const

import "time"

const (
	// MaxBodyLen set to max body len, affect send/recv buffer size
	MaxBodyLen = 0xfffff
	// PacketBufferPoolSize max size of pool packet buffer
	PacketBufferPoolSize = 10

	// ServerAPICallTimeOutDur api call watchdog timer
	ServerAPICallTimeOutDur = time.Second * 2
)
