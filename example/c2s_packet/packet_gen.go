// Code generated by "genprotocol -ver=1.0 -prefix=c2s -basedir example -statstype=int"

package c2s_packet

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/kasworld/genprotocol/example/c2s_const"
	"github.com/kasworld/genprotocol/example/c2s_error"
	"github.com/kasworld/genprotocol/example/c2s_idcmd"
	"github.com/kasworld/genprotocol/example/c2s_idnoti"
)

type FlowType byte // packet flow type

const (
	invalid      FlowType = iota // make uninitalized packet error
	Request                      // Request for request packet (response packet expected)
	Response                     // Response is reply of request packet
	Notification                 // Notification is just send and forget packet
)

var _FlowType_str = map[FlowType]string{
	invalid:      "invalid",
	Request:      "Request",
	Response:     "Response",
	Notification: "Notification",
}

func (e FlowType) String() string {
	if s, exist := _FlowType_str[e]; exist {
		return s
	}
	return fmt.Sprintf("FlowType%d", byte(e))
}

///////////////////////////////////////////////////////////////////////////////

const (
	// HeaderLen fixed size of header
	HeaderLen = 4 + 4 + 2 + 2 + 1 + 1 + 2

	// MaxPacketLen max total packet size byte of raw packet
	MaxPacketLen = HeaderLen + c2s_const.MaxBodyLen
)

func (pk Packet) String() string {
	return fmt.Sprintf("Packet[%v %+v]", pk.Header, pk.Body)
}

// Packet is header + body as object (not byte list)
type Packet struct {
	Header Header
	Body   interface{}
}

func (h Header) String() string {
	switch h.FlowType {
	default:
		return fmt.Sprintf(
			"Header[%v:%v ID:%v Error:%v BodyLen:%v BodyType:%v Fill:%v]",
			h.FlowType, h.Cmd, h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
	case invalid:
		return fmt.Sprintf(
			"Header[%v:%v ID:%v Error:%v BodyLen:%v BodyType:%v Fill:%v]",
			h.FlowType, h.Cmd, h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
	case Request:
		return fmt.Sprintf(
			"Header[%v:%v ID:%v Error:%v BodyLen:%v BodyType:%v Fill:%v]",
			h.FlowType, c2s_idcmd.CommandID(h.Cmd), h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
	case Response:
		return fmt.Sprintf(
			"Header[%v:%v ID:%v Error:%v BodyLen:%v BodyType:%v Fill:%v]",
			h.FlowType, c2s_idcmd.CommandID(h.Cmd), h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
	case Notification:
		return fmt.Sprintf(
			"Header[%v:%v ID:%v Error:%v BodyLen:%v BodyType:%v Fill:%v]",
			h.FlowType, c2s_idnoti.NotiID(h.Cmd), h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
	}
}

// Header is fixed size header of packet
type Header struct {
	bodyLen   uint32              // set at marshal(Packet2Bytes)
	ID        uint32              // sender set, unique id per packet (wrap around reuse)
	Cmd       uint16              // sender set, application demux received packet
	ErrorCode c2s_error.ErrorCode // sender set, Response error
	FlowType  FlowType            // sender set, flow control, Request, Response, Notification
	bodyType  byte                // set at marshal(Packet2Bytes), body compress, marshal type
	Fill      uint16              // sender set, any data
}

// MakeHeaderFromBytes unmarshal header from bytelist
func MakeHeaderFromBytes(buf []byte) Header {
	var h Header
	h.bodyLen = binary.LittleEndian.Uint32(buf[0:4])
	h.ID = binary.LittleEndian.Uint32(buf[4:8])
	h.Cmd = binary.LittleEndian.Uint16(buf[8:10])
	h.ErrorCode = c2s_error.ErrorCode(binary.LittleEndian.Uint16(buf[10:12]))
	h.FlowType = FlowType(buf[12])
	h.bodyType = buf[13]
	h.Fill = binary.LittleEndian.Uint16(buf[14:16])
	return h
}

func (h Header) toBytesAt(buf []byte) {
	binary.LittleEndian.PutUint32(buf[0:4], h.bodyLen)
	binary.LittleEndian.PutUint32(buf[4:8], h.ID)
	binary.LittleEndian.PutUint16(buf[8:10], h.Cmd)
	binary.LittleEndian.PutUint16(buf[10:12], uint16(h.ErrorCode))
	buf[12] = byte(h.FlowType)
	buf[13] = h.bodyType
	binary.LittleEndian.PutUint16(buf[14:16], h.Fill)
}

// ToByteList marshal header to bytelist
func (h Header) ToByteList() []byte {
	buf := make([]byte, HeaderLen)
	h.toBytesAt(buf)
	return buf
}

// GetBodyLenFromHeaderBytes return packet body len from bytelist of header
func GetBodyLenFromHeaderBytes(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf[0:4])
}

// BodyLen return bodylen field
func (h *Header) BodyLen() uint32 {
	return h.bodyLen
}

// BodyType return bodyType field
func (h *Header) BodyType() byte {
	return h.bodyType
}

///////////////////////////////////////////////////////////////////////////////

// func NewSendPacketBuffer() []byte {
// 	return make([]byte, MaxPacketLen)
// }

func NewRecvPacketBuffer() *RecvPacketBuffer {
	pb := &RecvPacketBuffer{
		RecvBuffer: make([]byte, MaxPacketLen),
		RecvLen:    0,
	}
	return pb
}

// RecvPacketBuffer used for packet receive
type RecvPacketBuffer struct {
	RecvBuffer []byte
	RecvLen    int
}

// NewRecvPacketBufferByData make RecvPacketBuffer by exist data
func NewRecvPacketBufferByData(rdata []byte) *RecvPacketBuffer {
	pb := &RecvPacketBuffer{
		RecvBuffer: rdata,
		RecvLen:    len(rdata),
	}
	return pb
}

// GetHeader make header and return
// if data is insufficent, return empty header
func (pb *RecvPacketBuffer) GetHeader() Header {
	if !pb.IsHeaderComplete() {
		return Header{}
	}
	header := MakeHeaderFromBytes(pb.RecvBuffer)
	return header
}

// GetBodyBytes return body ready to unmarshal.
func (pb *RecvPacketBuffer) GetBodyBytes() ([]byte, error) {
	if !pb.IsPacketComplete() {
		return nil, fmt.Errorf("packet not complete")
	}
	header := pb.GetHeader()
	body := pb.RecvBuffer[HeaderLen : HeaderLen+int(header.bodyLen)]
	return body, nil
}

// GetHeaderBody return header and Body as bytelist
// application need demux by header.FlowType, header.Cmd
// unmarshal body with header.bodyType
// and check header.ID(if response packet)
func (pb *RecvPacketBuffer) GetHeaderBody() (Header, []byte, error) {
	if !pb.IsPacketComplete() {
		return Header{}, nil, fmt.Errorf("packet not complete")
	}
	header := pb.GetHeader()
	body, err := pb.GetBodyBytes()
	return header, body, err
}

// IsHeaderComplete check recv data is sufficient for header
func (pb *RecvPacketBuffer) IsHeaderComplete() bool {
	return pb.RecvLen >= HeaderLen
}

// IsPacketComplete check recv data is sufficient for packet
func (pb *RecvPacketBuffer) IsPacketComplete() bool {
	if !pb.IsHeaderComplete() {
		return false
	}
	bodylen := GetBodyLenFromHeaderBytes(pb.RecvBuffer)
	return pb.RecvLen == HeaderLen+int(bodylen)
}

// NeedRecvLen return need data len for header or body
func (pb *RecvPacketBuffer) NeedRecvLen() int {
	if !pb.IsHeaderComplete() {
		return HeaderLen
	}
	bodylen := GetBodyLenFromHeaderBytes(pb.RecvBuffer)
	return HeaderLen + int(bodylen)
}

// Read use for partial recv like tcp read
func (pb *RecvPacketBuffer) Read(conn io.Reader) error {
	toRead := pb.NeedRecvLen()
	if toRead > MaxPacketLen {
		return fmt.Errorf("packet size over MaxPacketLen")
	}
	for pb.RecvLen < toRead {
		n, err := conn.Read(pb.RecvBuffer[pb.RecvLen:toRead])
		if err != nil {
			return err
		}
		pb.RecvLen += n
	}
	return nil
}

// Packet2Bytes make packet to bytelist
// marshalBodyFn append marshaled(+compress) body to buffer and return total buffer, bodyType, error
// set Packet.Header.bodyLen, Packet.Header.bodyType
// return bytelist, error
func Packet2Bytes(pk *Packet,
	marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
	oldbuf []byte,
) ([]byte, error) {
	newbuf, bodytype, err := marshalBodyFn(pk.Body, oldbuf)
	if err != nil {
		return nil, err
	}
	bodyLen := len(newbuf) - HeaderLen
	if bodyLen > c2s_const.MaxBodyLen {
		return nil,
			fmt.Errorf("fail to serialize large packet %v, %v", pk.Header, bodyLen)
	}
	pk.Header.bodyType = bodytype
	pk.Header.bodyLen = uint32(bodyLen)
	pk.Header.toBytesAt(newbuf)
	return newbuf, nil
}

///////////////////////////////////////////////////////////////////////////////
type Buffer []byte

type Pool struct {
	mutex    sync.Mutex
	buffPool []Buffer
	count    int
}

func NewPool(count int) *Pool {
	return &Pool{
		buffPool: make([]Buffer, 0, count),
		count:    count,
	}
}

func (p *Pool) String() string {
	return fmt.Sprintf("PacketPool[%v %v/%v]",
		len(p.buffPool), p.count,
	)
}

func (p *Pool) Get() Buffer {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var rtn Buffer
	if l := len(p.buffPool); l > 0 {
		rtn = p.buffPool[l-1]
		p.buffPool = p.buffPool[:l-1]
	} else {
		rtn = make(Buffer, HeaderLen, MaxPacketLen)
	}
	return rtn
}

func (p *Pool) Put(pb Buffer) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if len(p.buffPool) < p.count {
		p.buffPool = append(p.buffPool, pb[:HeaderLen])
	}
}