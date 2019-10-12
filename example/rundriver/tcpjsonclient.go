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
	"github.com/kasworld/genprotocol/example/c2s_statcallapi"
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

	app := NewApp()
	if err := app.c2sc.ConnectTo(*addr); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	app.sendTestPacket()
	app.sendTestPacket()
	app.Run()
}

type App struct {
	pid          uint32
	c2sc         *c2s_conntcp.Connection
	sendRecvStop func()
	statCallAPI  *c2s_statcallapi.StatCallAPI
	pid2statobj  *c2s_statcallapi.PacketID2StatObj
}

func NewApp() *App {
	app := &App{
		statCallAPI: c2s_statcallapi.New(),
		pid2statobj: c2s_statcallapi.NewPacketID2StatObj(),
	}
	app.c2sc = c2s_conntcp.New(
		PacketReadTimeoutSec, PacketWriteTimeoutSec,
		c2s_json.MarshalBodyFn,
		app.HandleRecvPacket,
		app.handleSentPacket,
	)
	return app
}

func (app *App) Run() {
	ctx := context.Background()
	sendRecvCtx, sendRecvCancel := context.WithCancel(ctx)
	app.sendRecvStop = sendRecvCancel
	app.c2sc.Run(sendRecvCtx)
}

func (app *App) handleSentPacket(header c2s_packet.Header) error {
	if err := app.statCallAPI.AfterSendReq(header); err != nil {
		return err
	}

	return nil
}

func (app *App) HandleRecvPacket(header c2s_packet.Header, body []byte) error {
	robj, err := c2s_json.UnmarshalPacket(header, body)
	_ = robj
	switch header.FlowType {
	case c2s_packet.Response:
		if err := app.statCallAPI.AfterRecvRsp(header); err != nil {
			fmt.Printf("%v %v\n", app, err)
			app.sendRecvStop()
			return err
		}
		psobj := app.pid2statobj.Get(header.ID)
		if psobj == nil {
			app.sendRecvStop()
			return fmt.Errorf("no statobj for %v", header.ID)
		}
		psobj.CallServerEnd(header.ErrorCode == c2s_error.None)
		app.pid2statobj.Del(header.ID)
	}
	if err == nil {
		if err = app.sendTestPacket(); err != nil {
			return err
		}
	}
	return err
}

func (app *App) sendTestPacket() error {
	spk := app.makePacket()
	psobj, err := app.statCallAPI.BeforeSendReq(spk.Header)
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
