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
	"github.com/kasworld/genprotocol/example/c2s_authorize"
	"github.com/kasworld/genprotocol/example/c2s_json"
	"github.com/kasworld/genprotocol/example/c2s_serveconnbyte"

	"github.com/kasworld/genprotocol/example/c2s_handlereq"
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
	flag.Parse()

	ctx := context.Background()
	go serveTCP(ctx, *tcpport)
	serveHTTP(ctx, *httpport, *httpfolder)
}

func serveHTTP(ctx context.Context, port string, folder string) {
	fmt.Printf("http server dir=%v port=%v , http://localhost%v/\n",
		folder, port, port)
	webMux := http.NewServeMux()
	webMux.Handle("/",
		http.FileServer(http.Dir(folder)),
	)
	webMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebSocketClient(ctx, w, r)
	})
	if err := http.ListenAndServe(port, webMux); err != nil {
		fmt.Println(err.Error())
	}
}

func CheckOrigin(r *http.Request) bool {
	return true
}

func serveWebSocketClient(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: CheckOrigin,
	}
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrade %v\n", err)
		return
	}
	c2sc := c2s_serveconnbyte.New(
		sendBufferSize,
		c2s_authorize.NewAllSet(),
		c2s_handlereq.DemuxReq2BytesAPIFnMap)
	c2sc.StartServeWS(ctx, wsConn,
		readTimeoutSec, writeTimeoutSec, c2s_json.MarshalBodyFn)
	wsConn.Close()
}

func serveTCP(ctx context.Context, port string) {
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
				go serveTCPClient(ctx, conn)
			}
		}
	}
}

func serveTCPClient(ctx context.Context, conn *net.TCPConn) {
	c2sc := c2s_serveconnbyte.New(
		sendBufferSize,
		c2s_authorize.NewAllSet(),
		c2s_handlereq.DemuxReq2BytesAPIFnMap)
	c2sc.StartServeTCP(ctx, conn,
		readTimeoutSec, writeTimeoutSec, c2s_json.MarshalBodyFn)
	conn.Close()
}
