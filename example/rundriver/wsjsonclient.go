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
	"github.com/kasworld/genprotocol/example/c2s_statapierror"
	"github.com/kasworld/genprotocol/example/c2s_statcallapi"
	"github.com/kasworld/genprotocol/example/c2s_statnoti"
)

// service const
const (
	// for client
	ReadTimeoutSec  = 6 * time.Second
	WriteTimeoutSec = 3 * time.Second
)

func main() {
	addr := flag.String("addr", "localhost:8080", "server addr")
	flag.Parse()
	fmt.Printf("addr %v \n", *addr)

	app := NewApp(*addr)
	app.Run()
}

type App struct {
	pid      uint32
	addr     string
	c2sc     *c2s_connwsgorilla.Connection
	DoClose  func()
	apistat  *c2s_statcallapi.StatCallAPI
	notistat *c2s_statnoti.StatNotification
	errstat  *c2s_statapierror.StatAPIError
}

func NewApp(addr string) *App {
	app := &App{
		addr:     addr,
		apistat:  c2s_statcallapi.New(),
		notistat: c2s_statnoti.New(),
		errstat:  c2s_statapierror.New(),
	}
	app.DoClose = func() {
		fmt.Printf("Too early DoClose call\n")
	}
	return app
}

func (app *App) Run() {
	app.c2sc = c2s_connwsgorilla.New(
		ReadTimeoutSec, WriteTimeoutSec,
		c2s_json.MarshalBodyFn,
		app.HandleRecvPacket,
		app.handleSentPacket,
	)
	if err := app.c2sc.ConnectTo(app.addr); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	app.c2sc.EnqueueSendPacket(app.makePacket())
	app.c2sc.EnqueueSendPacket(app.makePacket())
	ctx, closeCtx := context.WithCancel(context.Background())
	app.DoClose = closeCtx
	defer app.DoClose()
	app.c2sc.Run(ctx)
}

func (app *App) handleSentPacket(header c2s_packet.Header) error {
	return nil
}

func (app *App) HandleRecvPacket(header c2s_packet.Header, body []byte) error {
	// robj, err := c2s_json.UnmarshalPacket(header, body)

	switch header.FlowType {
	default:
		return fmt.Errorf("Invalid packet type %v %v", header, body)
	case c2s_packet.Response:
		app.errstat.Inc(c2s_idcmd.CommandID(header.Cmd), header.ErrorCode)

		app.c2sc.EnqueueSendPacket(app.makePacket())

	case c2s_packet.Notification:
		app.notistat.Add(header)
	}
	return nil
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
