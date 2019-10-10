package c2s_server

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kasworld/genprotocol/example/c2s_error"
	"github.com/kasworld/genprotocol/example/c2s_json"
	"github.com/kasworld/genprotocol/example/c2s_packet"
	"github.com/kasworld/genprotocol/example/c2s_wsgorilla"
)

// service const
const (
	SendBufferSize = 10

	PacketReadTimeoutSec  = 6 * time.Second
	PacketWriteTimeoutSec = 3 * time.Second
)

func (c2sc *ServeClientConn) String() string {
	return fmt.Sprintf(
		"ServeClientConn[%v SendCh:%v]",
		c2sc.RemoteAddr,
		len(c2sc.sendCh),
	)
}

type ServeClientConn struct {
	RemoteAddr   string
	sendCh       chan c2s_packet.Packet
	sendRecvStop func()
}

func NewServeClientConn(remoteAddr string) *ServeClientConn {
	c2sc := &ServeClientConn{
		RemoteAddr: remoteAddr,
		sendCh:     make(chan c2s_packet.Packet, SendBufferSize),
	}

	c2sc.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call %v\n", c2sc)
	}
	return c2sc
}

func (c2sc *ServeClientConn) StartServeClientConn(mainctx context.Context, wsConn *websocket.Conn) {

	sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
	c2sc.sendRecvStop = sendRecvCancel

	go func() {
		err := c2s_wsgorilla.RecvLoop(sendRecvCtx, c2sc.sendRecvStop, wsConn,
			PacketReadTimeoutSec, c2sc.HandleRecvPacket)
		if err != nil {
			fmt.Printf("end RecvLoop %v\n", err)
		}
	}()
	go func() {
		err := c2s_wsgorilla.SendLoop(sendRecvCtx, c2sc.sendRecvStop, wsConn,
			PacketWriteTimeoutSec, c2sc.sendCh,
			c2s_json.MarshalBodyFn, c2sc.handleSentPacket)
		if err != nil {
			fmt.Printf("end SendLoop %v\n", err)
		}
	}()

loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			break loop

		}
	}
}

func (c2sc *ServeClientConn) handleSentPacket(header c2s_packet.Header) error {
	return nil
}

func (c2sc *ServeClientConn) HandleRecvPacket(header c2s_packet.Header, rbody []byte) error {
	robj, err := c2s_json.UnmarshalPacket(header, rbody)
	if err != nil {
		return err
	}
	if header.FlowType != c2s_packet.Request {
		return fmt.Errorf("Unexpected header packet type: %v", header)
	}
	if int(header.Cmd) >= len(DemuxReq2APIFnMap) {
		return fmt.Errorf("Invalid header command %v", header)
	}
	response, errcode, apierr := DemuxReq2APIFnMap[header.Cmd](c2sc, header, robj)
	if errcode != c2s_error.Disconnect && apierr == nil {
		rhd := header
		rhd.FlowType = c2s_packet.Response
		rpk := c2s_packet.Packet{
			Header: rhd,
			Body:   response,
		}
		c2sc.enqueueSendPacket(rpk)
	}
	return apierr
}

func (c2sc *ServeClientConn) enqueueSendPacket(pk c2s_packet.Packet) error {
	trycount := 10
	for trycount > 0 {
		select {
		case c2sc.sendCh <- pk:
			return nil
		default:
			trycount--
		}
		fmt.Printf("Send delayed, %s send channel busy %v, retry %v\n",
			c2sc, len(c2sc.sendCh), 10-trycount)
		time.Sleep(1 * time.Millisecond)
	}

	return fmt.Errorf("Send channel full %v", c2sc)
}
