package c2s_handlereq

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kasworld/genprotocol/example/c2s_error"
	"github.com/kasworld/genprotocol/example/c2s_idcmd"
	"github.com/kasworld/genprotocol/example/c2s_json"
	"github.com/kasworld/genprotocol/example/c2s_looptcp"
	"github.com/kasworld/genprotocol/example/c2s_loopwsgorilla"
	"github.com/kasworld/genprotocol/example/c2s_packet"
	"github.com/kasworld/genprotocol/example/c2s_statapierror"
	"github.com/kasworld/genprotocol/example/c2s_statnoti"
	"github.com/kasworld/genprotocol/example/c2s_statserveapi"
)

// service const
const (
	SendBufferSize = 10

	PacketReadTimeoutSec  = 6 * time.Second
	PacketWriteTimeoutSec = 3 * time.Second
)

func (c2sc *ServeClientConn) String() string {
	return fmt.Sprintf("ServeClientConn[SendCh:%v]", len(c2sc.sendCh))
}

type ServeClientConn struct {
	sendCh       chan c2s_packet.Packet
	sendRecvStop func()
	tid2StatObj  *c2s_statserveapi.PacketID2StatObj
	protocolStat *c2s_statserveapi.StatServeAPI
	notiStat     *c2s_statnoti.StatNotification
	errorStat    *c2s_statapierror.StatAPIError
}

func NewServeClientConn() *ServeClientConn {
	c2sc := &ServeClientConn{
		sendCh:       make(chan c2s_packet.Packet, SendBufferSize),
		tid2StatObj:  c2s_statserveapi.NewPacketID2StatObj(),
		protocolStat: c2s_statserveapi.New(),
		notiStat:     c2s_statnoti.New(),
		errorStat:    c2s_statapierror.New(),
	}
	c2sc.sendRecvStop = func() {
		fmt.Printf("Too early sendRecvStop call %v\n", c2sc)
	}
	return c2sc
}

func (c2sc *ServeClientConn) GetProtocolStat() *c2s_statserveapi.StatServeAPI {
	return c2sc.protocolStat
}
func (c2sc *ServeClientConn) GetNotiStat() *c2s_statnoti.StatNotification {
	return c2sc.notiStat
}
func (c2sc *ServeClientConn) GetErrorStat() *c2s_statapierror.StatAPIError {
	return c2sc.errorStat
}

func (c2sc *ServeClientConn) StartServeWS(mainctx context.Context, conn *websocket.Conn) {
	sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
	c2sc.sendRecvStop = sendRecvCancel
	go func() {
		err := c2s_loopwsgorilla.RecvLoop(sendRecvCtx, c2sc.sendRecvStop, conn,
			PacketReadTimeoutSec, c2sc.HandleRecvPacket)
		if err != nil {
			fmt.Printf("end RecvLoop %v\n", err)
		}
	}()
	go func() {
		err := c2s_loopwsgorilla.SendLoop(sendRecvCtx, c2sc.sendRecvStop, conn,
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

func (c2sc *ServeClientConn) StartServeTCP(mainctx context.Context, conn *net.TCPConn) {
	sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
	c2sc.sendRecvStop = sendRecvCancel
	go func() {
		err := c2s_looptcp.RecvLoop(sendRecvCtx, c2sc.sendRecvStop, conn,
			PacketReadTimeoutSec, c2sc.HandleRecvPacket)
		if err != nil {
			fmt.Printf("end RecvLoop %v\n", err)
		}
	}()
	go func() {
		err := c2s_looptcp.SendLoop(sendRecvCtx, c2sc.sendRecvStop, conn,
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
	switch header.FlowType {
	default:
		panic(fmt.Sprintf("invalid packet type %s %v", c2sc, header))

	case c2s_packet.Request:
		panic(fmt.Sprintf("request packet not supported %s %v", c2sc, header))

	case c2s_packet.Response:
		statOjb := c2sc.tid2StatObj.Del(header.ID)
		if statOjb != nil {
			statOjb.AfterSendRsp(header)
		} else {
			panic(fmt.Sprintf("send StatObj not found %v", header))
		}
	case c2s_packet.Notification:
		c2sc.GetNotiStat().Add(header)
	}
	return nil
}

func (c2sc *ServeClientConn) HandleRecvPacket(header c2s_packet.Header, rbody []byte) error {
	// robj, err := c2s_json.UnmarshalPacket(header, rbody)
	// if err != nil {
	// 	return err
	// }
	if header.FlowType != c2s_packet.Request {
		return fmt.Errorf("Unexpected header packet type: %v", header)
	}
	if int(header.Cmd) >= len(DemuxReq2BytesAPIFnMap) {
		return fmt.Errorf("Invalid header command %v", header)
	}

	statObj, err := c2sc.GetProtocolStat().AfterRecvReqHeader(header)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		if err := c2sc.tid2StatObj.Add(header.ID, statObj); err != nil {
			fmt.Printf("%v\n", err)
			return err
		}
	}

	sObj := c2sc.tid2StatObj.Get(header.ID)
	if sObj == nil {
		return fmt.Errorf("protocol stat obj nil %v, maybe pkid duplicate?", header.ID)
	}
	sObj.BeforeAPICall()
	response, errcode, apierr := DemuxReq2BytesAPIFnMap[header.Cmd](c2sc, header, rbody)
	c2sc.GetErrorStat().Inc(c2s_idcmd.CommandID(header.Cmd), response.ErrorCode)
	sObj.AfterAPICall()

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
