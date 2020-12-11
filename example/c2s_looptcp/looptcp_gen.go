// Code generated by "genprotocol.exe -ver=1.0 -prefix=c2s -basedir=example -statstype=int"

package c2s_looptcp

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/kasworld/genprotocol/example/c2s_packet"
)

func WriteBytes(conn io.Writer, buf []byte) error {
	toWrite := len(buf)
	for l := 0; l < toWrite; {
		n, err := conn.Write(buf[l:toWrite])
		if err != nil {
			return err
		}
		l += n
	}
	return nil
}

func SendLoop(sendRecvCtx context.Context, SendRecvStop func(), tcpConn *net.TCPConn,
	timeOut time.Duration,
	SendCh chan *c2s_packet.Packet,
	marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
	handleSentPacketFn func(pk *c2s_packet.Packet) error,
) error {

	defer SendRecvStop()
	sendBuffer := make([]byte, c2s_packet.HeaderLen, c2s_packet.MaxPacketLen)
	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			break loop
		case pk := <-SendCh:
			sendBuffer, err = c2s_packet.Packet2Bytes(pk, marshalBodyFn, sendBuffer[:c2s_packet.HeaderLen])
			if err != nil {
				break loop
			}
			if err = tcpConn.SetWriteDeadline(time.Now().Add(timeOut)); err != nil {
				break loop
			}
			if err = WriteBytes(tcpConn, sendBuffer); err != nil {
				break loop
			}
			if err = handleSentPacketFn(pk); err != nil {
				break loop
			}
		}
	}
	return err
}

func RecvLoop(sendRecvCtx context.Context, SendRecvStop func(), tcpConn *net.TCPConn,
	timeOut time.Duration,
	HandleRecvPacketFn func(header c2s_packet.Header, body []byte) error,
) error {

	defer SendRecvStop()

	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			return nil

		default:
			if err = tcpConn.SetReadDeadline(time.Now().Add(timeOut)); err != nil {
				break loop
			}
			header, rbody, err := c2s_packet.ReadHeaderBody(tcpConn)
			if err != nil {
				return err
			}
			if err = HandleRecvPacketFn(header, rbody); err != nil {
				break loop
			}
		}
	}
	return err
}
