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
	"time"

	"github.com/kasworld/genprotocol/example/c2s_connwsgorilla"
	"github.com/kasworld/genprotocol/example/c2s_idcmd"
	"github.com/kasworld/genprotocol/example/c2s_json"
	"github.com/kasworld/genprotocol/example/c2s_obj"
	"github.com/kasworld/genprotocol/example/c2s_packet"
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

	app := new(App)
	app.c2sc = c2s_connwsgorilla.New(
		PacketReadTimeoutSec, PacketWriteTimeoutSec,
		c2s_json.MarshalBodyFn,
		app.HandleRecvPacket,
		app.handleSentPacket,
	)
	if err := app.c2sc.ConnectTo(*addr); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	app.c2sc.EnqueueSendPacket(app.makePacket())
	app.c2sc.EnqueueSendPacket(app.makePacket())
	ctx := context.Background()
	app.c2sc.Run(ctx)
}

type App struct {
	pid  uint32
	c2sc *c2s_connwsgorilla.Connection
}

func (app *App) handleSentPacket(header c2s_packet.Header) error {
	return nil
}

func (app *App) HandleRecvPacket(header c2s_packet.Header, body []byte) error {
	robj, err := c2s_json.UnmarshalPacket(header, body)
	_ = robj
	// fmt.Println(header, robj, err)
	if err == nil {
		app.c2sc.EnqueueSendPacket(app.makePacket())
	}
	return err
}

func (app *App) makePacket() c2s_packet.Packet {
	body := c2s_obj.ReqHeartbeat_data{}
	hd := c2s_packet.Header{
		Cmd:      uint16(c2s_idcmd.Heartbeat),
		ID:       app.pid,
		FlowType: c2s_packet.Request,
	}
	app.pid++

	return c2s_packet.Packet{
		Header: hd,
		Body:   body,
	}
}
