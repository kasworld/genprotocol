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
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kasworld/actpersec"
	"github.com/kasworld/genprotocol/example/c2s_authorize"
	"github.com/kasworld/genprotocol/example/c2s_error"
	"github.com/kasworld/genprotocol/example/c2s_gob"
	"github.com/kasworld/genprotocol/example/c2s_idcmd"
	"github.com/kasworld/genprotocol/example/c2s_json"
	"github.com/kasworld/genprotocol/example/c2s_obj"
	"github.com/kasworld/genprotocol/example/c2s_packet"
	"github.com/kasworld/genprotocol/example/c2s_serveconnbyte"
	"github.com/kasworld/genprotocol/example/c2s_statapierror"
	"github.com/kasworld/genprotocol/example/c2s_statnoti"
	"github.com/kasworld/genprotocol/example/c2s_statserveapi"
)

// service const
const (
	sendBufferSize  = 10
	readTimeoutSec  = 6 * time.Second
	writeTimeoutSec = 3 * time.Second
)

func main() {
	httpport := flag.String("httpport", ":8080", "Serve httpport")
	httpfolder := flag.String("httpdir", "www", "Serve http Dir")
	tcpport := flag.String("tcpport", ":8081", "Serve tcpport")
	marshaltype := flag.String("marshaltype", "json", "msgp,json,gob")
	flag.Parse()

	svr := NewServer(*marshaltype)
	svr.Run(*tcpport, *httpport, *httpfolder)
}

type Server struct {
	sendRecvStop func()

	SendStat *actpersec.ActPerSec `prettystring:"simple"`
	RecvStat *actpersec.ActPerSec `prettystring:"simple"`

	apiStat                *c2s_statserveapi.StatServeAPI
	notiStat               *c2s_statnoti.StatNotification
	errStat                *c2s_statapierror.StatAPIError
	marshalBodyFn          func(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error)
	unmarshalPacketFn      func(h c2s_packet.Header, bodyData []byte) (interface{}, error)
	DemuxReq2BytesAPIFnMap [c2s_idcmd.CommandID_Count]func(
		me interface{}, hd c2s_packet.Header, rbody []byte) (
		c2s_packet.Header, interface{}, error)
}

func NewServer(marshaltype string) *Server {
	svr := &Server{
		SendStat: actpersec.New(),
		RecvStat: actpersec.New(),

		apiStat:  c2s_statserveapi.New(),
		notiStat: c2s_statnoti.New(),
		errStat:  c2s_statapierror.New(),
	}
	svr.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call\n")
	}

	switch marshaltype {
	default:
		fmt.Printf("unsupported marshaltype %v\n", marshaltype)
		return nil
	// case "msgp":
	// 	gMarshalBodyFn = c2s_msgp.MarshalBodyFn
	// 	gUnmarshalPacket = c2s_msgp.UnmarshalPacket
	case "json":
		svr.marshalBodyFn = c2s_json.MarshalBodyFn
		svr.unmarshalPacketFn = c2s_json.UnmarshalPacket
	case "gob":
		svr.marshalBodyFn = c2s_gob.MarshalBodyFn
		svr.unmarshalPacketFn = c2s_gob.UnmarshalPacket
	}
	fmt.Printf("start using marshaltype %v\n", marshaltype)
	svr.DemuxReq2BytesAPIFnMap = [...]func(
		me interface{}, hd c2s_packet.Header, rbody []byte) (
		c2s_packet.Header, interface{}, error){
		c2s_idcmd.InvalidCmd: svr.bytesAPIFn_ReqInvalidCmd,
		c2s_idcmd.Login:      svr.bytesAPIFn_ReqLogin,
		c2s_idcmd.Heartbeat:  svr.bytesAPIFn_ReqHeartbeat,
		c2s_idcmd.Chat:       svr.bytesAPIFn_ReqChat,
	}
	return svr
}

func (svr *Server) Run(tcpport string, httpport string, httpfolder string) {
	ctx, stopFn := context.WithCancel(context.Background())
	svr.sendRecvStop = stopFn
	defer svr.sendRecvStop()

	go svr.serveTCP(ctx, tcpport)
	go svr.serveHTTP(ctx, httpport, httpfolder)

	timerInfoTk := time.NewTicker(1 * time.Second)
	defer timerInfoTk.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timerInfoTk.C:
			svr.SendStat.UpdateLap()
			svr.RecvStat.UpdateLap()
		}
	}
}

func (svr *Server) serveHTTP(ctx context.Context, port string, folder string) {
	fmt.Printf("http server dir=%v port=%v , http://localhost%v/\n",
		folder, port, port)
	webMux := http.NewServeMux()
	webMux.Handle("/",
		http.FileServer(http.Dir(folder)),
	)
	webMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		svr.serveWebSocketClient(ctx, w, r)
	})
	if err := http.ListenAndServe(port, webMux); err != nil {
		fmt.Println(err.Error())
	}
}

func CheckOrigin(r *http.Request) bool {
	return true
}

func (svr *Server) serveWebSocketClient(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: CheckOrigin,
	}
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrade %v\n", err)
		return
	}
	c2sc := c2s_serveconnbyte.NewWithStats(
		nil,
		sendBufferSize,
		c2s_authorize.NewAllSet(),
		svr.SendStat, svr.RecvStat,
		svr.apiStat,
		svr.notiStat,
		svr.errStat,
		svr.DemuxReq2BytesAPIFnMap)
	c2sc.StartServeWS(ctx, wsConn,
		readTimeoutSec, writeTimeoutSec, svr.marshalBodyFn)
	wsConn.Close()
}

func (svr *Server) serveTCP(ctx context.Context, port string) {
	fmt.Printf("tcp server port=%v\n", port)
	tcpaddr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}
	defer listener.Close()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			listener.SetDeadline(time.Now().Add(time.Duration(1 * time.Second)))
			conn, err := listener.AcceptTCP()
			if err != nil {
				operr, ok := err.(*net.OpError)
				if ok && operr.Timeout() {
					continue
				}
				fmt.Printf("error %#v\n", err)
			} else {
				go svr.serveTCPClient(ctx, conn)
			}
		}
	}
}

func (svr *Server) serveTCPClient(ctx context.Context, conn *net.TCPConn) {
	c2sc := c2s_serveconnbyte.NewWithStats(
		nil,
		sendBufferSize,
		c2s_authorize.NewAllSet(),
		svr.SendStat, svr.RecvStat,
		svr.apiStat,
		svr.notiStat,
		svr.errStat,
		svr.DemuxReq2BytesAPIFnMap)
	c2sc.StartServeTCP(ctx, conn,
		readTimeoutSec, writeTimeoutSec, svr.marshalBodyFn)
	conn.Close()
}

///////////////////////////////////////////////////////////////

func (svr *Server) bytesAPIFn_ReqInvalidCmd(
	me interface{}, hd c2s_packet.Header, rbody []byte) (
	c2s_packet.Header, interface{}, error) {
	// robj, err := gUnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*c2s_obj.ReqInvalidCmd_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	hd.ErrorCode = c2s_error.None
	sendBody := &c2s_obj.RspInvalidCmd_data{}
	return hd, sendBody, nil
}

func (svr *Server) bytesAPIFn_ReqLogin(
	me interface{}, hd c2s_packet.Header, rbody []byte) (
	c2s_packet.Header, interface{}, error) {
	// robj, err := gUnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*c2s_obj.ReqLogin_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	hd.ErrorCode = c2s_error.None
	sendBody := &c2s_obj.RspLogin_data{}
	return hd, sendBody, nil
}

func (svr *Server) bytesAPIFn_ReqHeartbeat(
	me interface{}, hd c2s_packet.Header, rbody []byte) (
	c2s_packet.Header, interface{}, error) {
	robj, err := svr.unmarshalPacketFn(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*c2s_obj.ReqHeartbeat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	hd.ErrorCode = c2s_error.None
	sendBody := &c2s_obj.RspHeartbeat_data{recvBody.Now}
	return hd, sendBody, nil
}

func (svr *Server) bytesAPIFn_ReqChat(
	me interface{}, hd c2s_packet.Header, rbody []byte) (
	c2s_packet.Header, interface{}, error) {
	robj, err := svr.unmarshalPacketFn(hd, rbody)
	if err != nil {
		return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	}
	recvBody, ok := robj.(*c2s_obj.ReqChat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	_ = recvBody

	hd.ErrorCode = c2s_error.None
	sendBody := &c2s_obj.RspChat_data{recvBody.Msg}
	return hd, sendBody, nil
}
