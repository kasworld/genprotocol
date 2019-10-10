// +build ignore

// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/kasworld/genprotocol/example/c2s_idcmd"
	"github.com/kasworld/genprotocol/example/c2s_json"
	"github.com/kasworld/genprotocol/example/c2s_obj"
	"github.com/kasworld/genprotocol/example/c2s_packet"
	"github.com/kasworld/genprotocol/example/c2s_tcploop"
)

// service const
const (
	SendBufferSize = 10

	// for client
	PacketReadTimeoutSec  = 6 * time.Second
	PacketWriteTimeoutSec = 3 * time.Second
)

func main() {
	addr := flag.String("addr", "localhost:8080", "server addr")
	flag.Parse()
	fmt.Printf("addr %v \n", *addr)

	ctx := context.Background()
	c2sc := NewConnection(*addr)

	c2sc.enqueueSendPacket(c2sc.makePacket())
	c2sc.enqueueSendPacket(c2sc.makePacket())
	c2sc.Connect(ctx)
}

///////////////////

func (c2sc *Connection) String() string {
	return fmt.Sprintf(
		"Connection[%v SendCh:%v]",
		c2sc.RemoteAddr,
		len(c2sc.sendCh),
	)
}

type Connection struct {
	RemoteAddr   string
	sendCh       chan c2s_packet.Packet
	sendRecvStop func()
	pid          uint32
}

func (c2sc *Connection) makePacket() c2s_packet.Packet {
	body := c2s_obj.ReqHeartbeat_data{}
	hd := c2s_packet.Header{
		Cmd:      uint16(c2s_idcmd.Heartbeat),
		ID:       c2sc.pid,
		FlowType: c2s_packet.Request,
	}
	c2sc.pid++

	return c2s_packet.Packet{
		Header: hd,
		Body:   body,
	}
}

func NewConnection(remoteAddr string) *Connection {
	c2sc := &Connection{
		RemoteAddr: remoteAddr,
		sendCh:     make(chan c2s_packet.Packet, SendBufferSize),
	}

	c2sc.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call %v\n", c2sc)
	}
	return c2sc
}

func (c2sc *Connection) Connect(mainctx context.Context) {

	tcpaddr, err := net.ResolveTCPAddr("tcp", c2sc.RemoteAddr)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
	c2sc.sendRecvStop = sendRecvCancel

	go func() {
		err := c2s_tcploop.RecvLoop(sendRecvCtx, c2sc.sendRecvStop, conn,
			PacketReadTimeoutSec, c2sc.HandleRecvPacket)
		if err != nil {
			fmt.Printf("end RecvLoop %v\n", err)
		}
	}()
	go func() {
		err := c2s_tcploop.SendLoop(sendRecvCtx, c2sc.sendRecvStop, conn,
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

func (c2sc *Connection) handleSentPacket(header c2s_packet.Header) error {
	return nil
}

func (c2sc *Connection) HandleRecvPacket(header c2s_packet.Header, body []byte) error {
	robj, err := c2s_json.UnmarshalPacket(header, body)
	_ = robj
	// fmt.Println(header, robj, err)
	if err == nil {
		c2sc.enqueueSendPacket(c2sc.makePacket())
	}
	return err
}

func (c2sc *Connection) enqueueSendPacket(pk c2s_packet.Packet) error {
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
