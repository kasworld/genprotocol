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

	"github.com/kasworld/genprotocol/example/c2s_json"
	"github.com/kasworld/genprotocol/example/c2s_serveconnbyte"

	"github.com/kasworld/genprotocol/example/c2s_handlereq"
)

func main() {
	port := flag.String("port", ":8080", "Serve port")
	flag.Parse()

	ctx := context.Background()

	tcpaddr, err := net.ResolveTCPAddr("tcp", *port)
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
	c2sc := c2s_serveconnbyte.New(c2s_json.MarshalBodyFn, c2s_handlereq.DemuxReq2BytesAPIFnMap)
	c2sc.StartServeTCP(ctx, conn)
	conn.Close()
}
