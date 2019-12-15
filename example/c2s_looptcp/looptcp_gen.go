// Code generated by "genprotocol -ver=1.0 -prefix=c2s -basedir example -statstype=int"

package c2s_looptcp

import (
	"context"
	"net"
	"time"

	"github.com/kasworld/genprotocol/example/c2s_const"
	"github.com/kasworld/genprotocol/example/c2s_packet"
)

var bufPool = c2s_packet.NewPool(c2s_const.PacketBufferPoolSize)

func SendPacket(conn *net.TCPConn, buf []byte) error {
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
	SendCh chan c2s_packet.Packet,
	marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
	handleSentPacketFn func(header c2s_packet.Header) error,
) error {

	defer SendRecvStop()
	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			break loop
		case pk := <-SendCh:
			if err = tcpConn.SetWriteDeadline(time.Now().Add(timeOut)); err != nil {
				break loop
			}
			oldbuf := bufPool.Get()
			sendBuffer, err := c2s_packet.Packet2Bytes(&pk, marshalBodyFn, oldbuf)
			if err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			if err = SendPacket(tcpConn, sendBuffer); err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			if err = handleSentPacketFn(pk.Header); err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			bufPool.Put(oldbuf)
		}
	}
	return err
}

func RecvLoop(sendRecvCtx context.Context, SendRecvStop func(), tcpConn *net.TCPConn,
	timeOut time.Duration,
	HandleRecvPacketFn func(header c2s_packet.Header, body []byte) error,
) error {

	defer SendRecvStop()

	pb := c2s_packet.NewRecvPacketBuffer()
	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			return nil

		default:
			if pb.IsPacketComplete() {
				header, rbody, lerr := pb.GetHeaderBody()
				if lerr != nil {
					err = lerr
					break loop
				}
				if err = HandleRecvPacketFn(header, rbody); err != nil {
					break loop
				}
				pb = c2s_packet.NewRecvPacketBuffer()
				if err = tcpConn.SetReadDeadline(time.Now().Add(timeOut)); err != nil {
					break loop
				}
			} else {
				err := pb.Read(tcpConn)
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}
