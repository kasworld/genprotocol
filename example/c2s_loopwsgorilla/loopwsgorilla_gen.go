// Code generated by "genprotocol.exe -ver=1.0 -prefix=c2s -basedir=example -statstype=int"

package c2s_loopwsgorilla

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kasworld/genprotocol/example/c2s_packet"
)

func SendControl(
	wsConn *websocket.Conn, mt int, PacketWriteTimeOut time.Duration) error {

	return wsConn.WriteControl(mt, []byte{}, time.Now().Add(PacketWriteTimeOut))
}

func WriteBytes(wsConn *websocket.Conn, sendBuffer []byte) error {
	return wsConn.WriteMessage(websocket.BinaryMessage, sendBuffer)
}

func SendLoop(sendRecvCtx context.Context, SendRecvStop func(), wsConn *websocket.Conn,
	timeout time.Duration,
	SendCh chan *c2s_packet.Packet,
	marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
	handleSentPacketFn func(pk *c2s_packet.Packet) error,
) error {

	defer SendRecvStop()
	sendBuffer := make([]byte, c2s_packet.HeaderLen)
	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			err = SendControl(wsConn, websocket.CloseMessage, timeout)
			break loop
		case pk := <-SendCh:
			sendBuffer, err = c2s_packet.Packet2Bytes(pk, marshalBodyFn, sendBuffer[:c2s_packet.HeaderLen])
			if err != nil {
				break loop
			}
			if err = wsConn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
				break loop
			}
			if err = WriteBytes(wsConn, sendBuffer); err != nil {
				break loop
			}
			if err = handleSentPacketFn(pk); err != nil {
				break loop
			}
		}
	}
	return err
}

func RecvLoop(sendRecvCtx context.Context, SendRecvStop func(), wsConn *websocket.Conn,
	timeout time.Duration,
	HandleRecvPacketFn func(header c2s_packet.Header, body []byte) error) error {

	defer SendRecvStop()
	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			break loop
		default:
			if err = wsConn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
				break loop
			}
			if header, body, lerr := RecvPacket(wsConn); lerr != nil {
				if operr, ok := lerr.(*net.OpError); ok && operr.Timeout() {
					continue
				}
				err = lerr
				break loop
			} else {
				if err = HandleRecvPacketFn(header, body); err != nil {
					break loop
				}
			}
		}
	}
	return err
}

func RecvPacket(wsConn *websocket.Conn) (c2s_packet.Header, []byte, error) {
	mt, rdata, err := wsConn.ReadMessage()
	if err != nil {
		return c2s_packet.Header{}, nil, err
	}
	if mt != websocket.BinaryMessage {
		return c2s_packet.Header{}, nil, fmt.Errorf("message not binary %v", mt)
	}
	return c2s_packet.Bytes2HeaderBody(rdata)
}
