// Code generated by "genprotocol -ver=1.0 -prefix=c2s -basedir example -statstype=int"

package c2s_connwasm

import (
	"context"
	"fmt"
	"sync"
	"syscall/js"

	"github.com/kasworld/genprotocol/example/c2s_const"
	"github.com/kasworld/genprotocol/example/c2s_packet"
)

var bufPool = c2s_packet.NewPool(c2s_const.PacketBufferPoolSize)

type Connection struct {
	remoteAddr   string
	conn         js.Value
	SendRecvStop func()
	sendCh       chan c2s_packet.Packet

	marshalBodyFn      func(interface{}, []byte) ([]byte, byte, error)
	handleRecvPacketFn func(header c2s_packet.Header, body []byte) error
	handleSentPacketFn func(header c2s_packet.Header) error
}

func (wsc *Connection) String() string {
	return fmt.Sprintf("Connection[%v SendCh:%v]",
		wsc.remoteAddr, len(wsc.sendCh))
}

func New(
	connAddr string,
	marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
	handleRecvPacketFn func(header c2s_packet.Header, body []byte) error,
	handleSentPacketFn func(header c2s_packet.Header) error,
) *Connection {
	wsc := &Connection{
		remoteAddr:         connAddr,
		sendCh:             make(chan c2s_packet.Packet, 10),
		marshalBodyFn:      marshalBodyFn,
		handleRecvPacketFn: handleRecvPacketFn,
		handleSentPacketFn: handleSentPacketFn,
	}
	wsc.SendRecvStop = func() {
		JsLogErrorf("Too early SendRecvStop call %v", wsc)
	}
	return wsc
}

func (wsc *Connection) Connect(ctx context.Context, wg *sync.WaitGroup) error {
	connCtx, ctxCancel := context.WithCancel(ctx)
	wsc.SendRecvStop = ctxCancel

	wsc.conn = js.Global().Get("WebSocket").New(wsc.remoteAddr)
	if !wsc.conn.Truthy() {
		err := fmt.Errorf("fail to connect %v", wsc.remoteAddr)
		JsLogErrorf("%v", err)
		return err
	}
	wsc.conn.Call("addEventListener", "open", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			wsc.conn.Call("addEventListener", "message", js.FuncOf(wsc.handleWebsocketMessage))
			go wsc.sendLoop(connCtx)
			wg.Done()
			return nil
		}))
	wsc.conn.Call("addEventListener", "close", js.FuncOf(wsc.wsClosed))
	wsc.conn.Call("addEventListener", "error", js.FuncOf(wsc.wsError))
	return nil
}

func (wsc *Connection) wsClosed(this js.Value, args []js.Value) interface{} {
	wsc.SendRecvStop()
	JsLogError("ws closed")
	return nil
}

func (wsc *Connection) wsError(this js.Value, args []js.Value) interface{} {
	wsc.SendRecvStop()
	JsLogError(this, args)
	return nil
}

func (wsc *Connection) sendLoop(sendRecvCtx context.Context) {
	defer wsc.SendRecvStop()
	var err error
loop:
	for {
		select {
		case <-sendRecvCtx.Done():
			break loop
		case pk := <-wsc.sendCh:
			oldbuf := bufPool.Get()
			sendBuffer, err := c2s_packet.Packet2Bytes(&pk, wsc.marshalBodyFn, oldbuf)
			if err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			if err = wsc.sendPacket(sendBuffer); err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			if err = wsc.handleSentPacketFn(pk.Header); err != nil {
				bufPool.Put(oldbuf)
				break loop
			}
			bufPool.Put(oldbuf)
		}
	}
	JsLogErrorf("end SendLoop %v\n", err)
	return
}

func (wsc *Connection) sendPacket(sendBuffer []byte) error {
	sendData := js.Global().Get("Uint8Array").New(len(sendBuffer))
	js.CopyBytesToJS(sendData, sendBuffer)
	wsc.conn.Call("send", sendData)
	return nil
}

func (wsc *Connection) handleWebsocketMessage(this js.Value, args []js.Value) interface{} {
	data := args[0].Get("data") // blob
	aBuff := data.Call("arrayBuffer")
	aBuff.Call("then",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			rdata := ArrayBufferToSlice(args[0])
			rPk := c2s_packet.NewRecvPacketBufferByData(rdata)
			header, body, lerr := rPk.GetHeaderBody()
			if lerr != nil {
				JsLogError(lerr.Error())
				wsc.SendRecvStop()
				return nil
			} else {
				if err := wsc.handleRecvPacketFn(header, body); err != nil {
					JsLogErrorf("%v", err)
					wsc.SendRecvStop()
					return nil
				}
			}
			return nil
		}))

	return nil
}

func Uint8ArrayToSlice(value js.Value) []byte {
	s := make([]byte, value.Get("byteLength").Int())
	js.CopyBytesToGo(s, value)
	return s
}

func ArrayBufferToSlice(value js.Value) []byte {
	return Uint8ArrayToSlice(js.Global().Get("Uint8Array").New(value))
}

func (wsc *Connection) EnqueueSendPacket(pk c2s_packet.Packet) error {
	select {
	case wsc.sendCh <- pk:
		return nil
	default:
		return fmt.Errorf("Send channel full %v", wsc)
	}
}

/////////

func JsLogError(v ...interface{}) {
	js.Global().Get("console").Call("error", v...)
}

func JsLogErrorf(format string, v ...interface{}) {
	js.Global().Get("console").Call("error", fmt.Sprintf(format, v...))
}