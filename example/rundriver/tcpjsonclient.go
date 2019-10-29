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

	"github.com/kasworld/genprotocol/example/c2s_conntcp"
	"github.com/kasworld/genprotocol/example/c2s_error"
	"github.com/kasworld/genprotocol/example/c2s_idcmd"
	"github.com/kasworld/genprotocol/example/c2s_json"
	"github.com/kasworld/genprotocol/example/c2s_obj"
	"github.com/kasworld/genprotocol/example/c2s_packet"
	"github.com/kasworld/genprotocol/example/c2s_pid2rspfn"
	"github.com/kasworld/genprotocol/example/c2s_statapierror"
	"github.com/kasworld/genprotocol/example/c2s_statcallapi"
	"github.com/kasworld/genprotocol/example/c2s_statnoti"
)

// service const
const (
	// for client
	readTimeoutSec  = 6 * time.Second
	writeTimeoutSec = 3 * time.Second
)

func main() {
	addr := flag.String("addr", "localhost:8081", "server addr")
	flag.Parse()
	fmt.Printf("addr %v \n", *addr)

	app := NewApp(*addr)
	app.sendTestPacket()
	app.sendTestPacket()
	app.Run()
}

type App struct {
	addr         string
	c2sc         *c2s_conntcp.Connection
	sendRecvStop func()
	apistat      *c2s_statcallapi.StatCallAPI
	pid2statobj  *c2s_statcallapi.PacketID2StatObj
	notistat     *c2s_statnoti.StatNotification
	errstat      *c2s_statapierror.StatAPIError
	pid2recv     *c2s_pid2rspfn.PID2RspFn
}

func NewApp(addr string) *App {
	app := &App{
		addr:        addr,
		apistat:     c2s_statcallapi.New(),
		pid2statobj: c2s_statcallapi.NewPacketID2StatObj(),
		notistat:    c2s_statnoti.New(),
		errstat:     c2s_statapierror.New(),
		pid2recv:    c2s_pid2rspfn.New(),
	}
	app.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call\n")
	}
	app.c2sc = c2s_conntcp.New(
		readTimeoutSec, writeTimeoutSec,
		c2s_json.MarshalBodyFn,
		app.handleRecvPacket,
		app.handleSentPacket,
	)
	return app
}

func (app *App) Run() {
	if err := app.c2sc.ConnectTo(app.addr); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	ctx, stopFn := context.WithCancel(context.Background())
	app.sendRecvStop = stopFn
	defer app.sendRecvStop()
	app.c2sc.Run(ctx)
}

func (app *App) handleSentPacket(header c2s_packet.Header) error {
	if err := app.apistat.AfterSendReq(header); err != nil {
		return err
	}
	return nil
}

func (app *App) handleRecvPacket(header c2s_packet.Header, body []byte) error {
	switch header.FlowType {
	default:
		return fmt.Errorf("Invalid packet type %v %v", header, body)
	case c2s_packet.Notification:
		// noti stat
		app.notistat.Add(header)
		//process noti here
		// robj, err := c2s_json.UnmarshalPacket(header, body)

	case c2s_packet.Response:
		// error stat
		app.errstat.Inc(c2s_idcmd.CommandID(header.Cmd), header.ErrorCode)
		// api stat
		if err := app.apistat.AfterRecvRsp(header); err != nil {
			fmt.Printf("%v %v\n", app, err)
			return err
		}
		psobj := app.pid2statobj.Get(header.ID)
		if psobj == nil {
			return fmt.Errorf("no statobj for %v", header.ID)
		}
		psobj.CallServerEnd(header.ErrorCode == c2s_error.None)
		app.pid2statobj.Del(header.ID)

		// process response
		if err := app.pid2recv.HandleRsp(header, body); err != nil {
			return err
		}

		// send new test packet
		go app.sendTestPacket()
	}
	return nil
}

func (app *App) sendTestPacket() error {
	return app.ReqWithRspFn(
		c2s_idcmd.Heartbeat,
		c2s_obj.ReqHeartbeat_data{},
		func(hd c2s_packet.Header, rsp interface{}) error {
			// response here
			// robj, err := c2s_json.UnmarshalPacket(header, body)
			return nil
		},
	)
}

func (app *App) ReqWithRspFn(cmd c2s_idcmd.CommandID, body interface{},
	fn c2s_pid2rspfn.HandleRspFn) error {

	pid := app.pid2recv.NewPID(fn)
	spk := c2s_packet.Packet{
		Header: c2s_packet.Header{
			Cmd:      uint16(cmd),
			ID:       pid,
			FlowType: c2s_packet.Request,
		},
		Body: body,
	}

	// add api stat
	psobj, err := app.apistat.BeforeSendReq(spk.Header)
	if err != nil {
		return nil
	}
	app.pid2statobj.Add(spk.Header.ID, psobj)

	if err := app.c2sc.EnqueueSendPacket(spk); err != nil {
		fmt.Printf("End %v %v %v\n", app, spk, err)
		app.sendRecvStop()
		return fmt.Errorf("Send fail %v %v", app, err)
	}
	return nil
}
