// +build ignore

// Copyright 2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
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
	"github.com/kasworld/genprotocol/example/c2s_connwsgorilla"
	"github.com/kasworld/genprotocol/example/c2s_error"
	"github.com/kasworld/genprotocol/example/c2s_gob"
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

var gMarshalBodyFn func(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error)
var gUnmarshalPacket func(h c2s_packet.Header, bodyData []byte) (interface{}, error)

func main() {
	addr := flag.String("addr", "localhost:8080", "server addr")
	marshaltype := flag.String("marshaltype", "json", "msgp,json,gob")
	nettype := flag.String("nettype", "ws", "tcp,ws websocket")

	flag.Parse()
	fmt.Printf("addr %v \n", *addr)

	switch *marshaltype {
	default:
		fmt.Printf("unsupported marshaltype %v\n", *marshaltype)
		return
	// case "msgp":
	// 	gMarshalBodyFn = c2s_msgp.MarshalBodyFn
	// 	gUnmarshalPacket = c2s_msgp.UnmarshalPacket
	case "json":
		gMarshalBodyFn = c2s_json.MarshalBodyFn
		gUnmarshalPacket = c2s_json.UnmarshalPacket
	case "gob":
		gMarshalBodyFn = c2s_gob.MarshalBodyFn
		gUnmarshalPacket = c2s_gob.UnmarshalPacket
	}
	fmt.Printf("start using marshaltype %v\n", *marshaltype)

	app := NewApp(*addr)
	app.Run(*nettype)
}

type App struct {
	addr string

	c2scWS            *c2s_connwsgorilla.Connection
	c2scTCP           *c2s_conntcp.Connection
	EnqueueSendPacket func(pk c2s_packet.Packet) error

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
	app.EnqueueSendPacket = func(pk c2s_packet.Packet) error {
		fmt.Printf("Too early EnqueueSendPacket call\n")
		return nil
	}
	return app
}

func (app *App) Run(nettype string) {
	ctx, stopFn := context.WithCancel(context.Background())
	app.sendRecvStop = stopFn
	defer app.sendRecvStop()

	switch nettype {
	default:
		fmt.Printf("unsupported nettype %v\n", nettype)
		return
	case "tcp":
		go app.connectTCP(ctx)
	case "ws":
		go app.connectWS(ctx)
	}
	fmt.Printf("start using nettype %v\n", nettype)

	time.Sleep(time.Second)
	app.sendTestPacket()
	app.sendTestPacket()

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (app *App) sendTestPacket() error {
	return app.ReqWithRspFn(
		c2s_idcmd.Heartbeat,
		c2s_obj.ReqHeartbeat_data{time.Now()},
		func(hd c2s_packet.Header, rbody interface{}) error {
			robj, err := gUnmarshalPacket(hd, rbody.([]byte))
			if err != nil {
				return fmt.Errorf("Packet type miss match %v", rbody)
			}
			recvBody, ok := robj.(*c2s_obj.RspHeartbeat_data)
			if !ok {
				return fmt.Errorf("Packet type miss match %v", robj)
			}

			fmt.Printf("ping %v\n", time.Now().Sub(recvBody.Now))
			return nil
		},
	)
}

func (app *App) connectWS(ctx context.Context) {
	app.c2scWS = c2s_connwsgorilla.New(
		readTimeoutSec, writeTimeoutSec,
		gMarshalBodyFn,
		app.handleRecvPacket,
		app.handleSentPacket,
	)
	if err := app.c2scWS.ConnectTo(app.addr); err != nil {
		fmt.Printf("%v\n", err)
		app.sendRecvStop()
		return
	}
	app.EnqueueSendPacket = app.c2scWS.EnqueueSendPacket
	app.c2scWS.Run(ctx)
}

func (app *App) connectTCP(ctx context.Context) {
	app.c2scTCP = c2s_conntcp.New(
		readTimeoutSec, writeTimeoutSec,
		gMarshalBodyFn,
		app.handleRecvPacket,
		app.handleSentPacket,
	)
	if err := app.c2scTCP.ConnectTo(app.addr); err != nil {
		fmt.Printf("%v\n", err)
		app.sendRecvStop()
		return
	}
	app.EnqueueSendPacket = app.c2scTCP.EnqueueSendPacket
	app.c2scTCP.Run(ctx)
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

	if err := app.EnqueueSendPacket(spk); err != nil {
		fmt.Printf("End %v %v %v\n", app, spk, err)
		app.sendRecvStop()
		return fmt.Errorf("Send fail %v %v", app, err)
	}
	return nil
}
