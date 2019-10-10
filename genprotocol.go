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
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func makeGenComment() string {
	return fmt.Sprintf("// Code generated by \"%s %s\"\n",
		filepath.Base(os.Args[0]), strings.Join(os.Args[1:], " "))
}

// loadEnumWithComment load list of enum + comment
func loadEnumWithComment(filename string) ([][]string, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	rtn := make([][]string, 0)
	rd := bufio.NewReader(fd)
	for {
		line, err := rd.ReadString('\n')
		line = strings.TrimSpace(line)
		if len(line) != 0 && line[0] != '#' {
			s2 := strings.SplitN(line, " ", 2)
			if len(s2) == 1 {
				s2 = append(s2, "")
			}
			rtn = append(rtn, s2)
		}
		if err != nil { // eof
			break
		}
	}
	return rtn, nil
}

// saveTo save go source with format, saved file may need goimport
func saveTo(outdata *bytes.Buffer, buferr error, outfilename string) error {
	if buferr != nil {
		fmt.Printf("fail %v %v", outfilename, buferr)
		return buferr
	}
	src, err := format.Source(outdata.Bytes())
	if err != nil {
		fmt.Println(outdata)
		fmt.Printf("fail %v %v", outfilename, err)
		return err
	}
	if werr := ioutil.WriteFile(outfilename, src, 0644); werr != nil {
		fmt.Printf("fail %v %v", outfilename, werr)
		return werr
	}
	fmt.Printf("goimports -w %v\n", outfilename)
	return nil
}

var (
	ver     = flag.String("ver", "", "protocol version")
	prefix  = flag.String("prefix", "", "protocol prefix")
	basedir = flag.String("basedir", "", "base directory")
)

func main() {
	flag.Parse()

	if *prefix == "" {
		fmt.Println("prefix not set")
	}
	if *ver == "" {
		fmt.Println("ver not set")
	}
	if *basedir == "" {
		fmt.Println("base dir not set")
	}

	dirToMake := []string{
		"_client",
		"_server",
		"_msgp",
		"_json",
		"_error",
		"_gendata",
		"_idcmd",
		"_idnoti",
		"_obj",
		"_packet",
		"_version",
		"_wasmconn",
		"_wsgorilla",
	}

	cmddatafile := path.Join(*basedir, *prefix+"_gendata", "command.data")
	cmddata, err := loadEnumWithComment(cmddatafile)
	if err != nil {
		fmt.Printf("fail to load %v %v\n", cmddatafile, err)
		return
	}
	notidatafile := path.Join(*basedir, *prefix+"_gendata", "noti.data")
	notidata, err := loadEnumWithComment(notidatafile)
	if err != nil {
		fmt.Printf("fail to load %v %v\n", notidatafile, err)
		return
	}
	errordatafile := path.Join(*basedir, *prefix+"_gendata", "error.data")
	errordata, err := loadEnumWithComment(errordatafile)
	if err != nil {
		fmt.Printf("fail to load %v %v\n", errordatafile, err)
		return
	}

	for _, v := range dirToMake {
		os.MkdirAll(path.Join(*basedir, *prefix+v), os.ModePerm)
	}

	buf, err := buildVersion(*prefix, *ver)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_version", "version_gen.go"))

	buf, err = buildDataCode(*prefix+"_idcmd", "CommandID", cmddata)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_idcmd", "command_gen.go"))

	buf, err = buildDataCode(*prefix+"_idnoti", "NotiID", notidata)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_idnoti", "noti_gen.go"))

	buf, err = buildDataCode(*prefix+"_error", "ErrorCode", errordata)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_error", "error_gen.go"))

	buf, err = buildPacket(*prefix)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_packet", "packet_gen.go"))

	buf, err = buildObjTemplate(*prefix, cmddata, notidata)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_obj", "objtemplate_gen.go"))

	buf, err = buildMSGP(*prefix, cmddata, notidata)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_msgp", "serialize_gen.go"))

	buf, err = buildJSON(*prefix, cmddata, notidata)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_json", "serialize_gen.go"))

	buf, err = buildRecvRspFnMap(*prefix, cmddata, notidata)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_client", "recvrspobjfnmap_gen.go"))

	buf, err = buildRecvNotiFnMap(*prefix, cmddata, notidata)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_client", "recvnotiobjfnmap_gen.go"))

	buf, err = buildClientSendRecv(*prefix, cmddata, notidata)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_client", "callsendrecv_gen.go"))

	buf, err = buildDemuxReq2API(*prefix, cmddata, notidata)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_server", "demuxreq2api_gen.go"))

	buf, err = buildAPITemplate(*prefix, cmddata, notidata)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_server", "apitemplate_gen.go"))

	buf, err = buildWasmConn(*prefix)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_wasmconn", "wasmconn_gen.go"))

	buf, err = buildWSGorilla(*prefix)
	saveTo(buf, err, path.Join(*basedir, *prefix+"_wsgorilla", "wsgorilla_gen.go"))
}

func buildDataCode(pkgname string, enumtype string, data [][]string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
		package %[1]s
		import "fmt"
		type %v uint16 // use in packet header, DO NOT CHANGE
	`, pkgname, enumtype)

	fmt.Fprintf(&buf, "const (\n")
	for i, v := range data {
		if i == 0 {
			fmt.Fprintf(&buf, "%v %v = iota // %v \n", v[0], enumtype, v[1])
		} else {
			fmt.Fprintf(&buf, "%v // %v\n", v[0], v[1])
		}
	}
	fmt.Fprintf(&buf, `
	%[1]s_Count int = iota 
	)
	var _%[1]s_str = map[%[1]s]string{
	`, enumtype)

	for _, v := range data {
		fmt.Fprintf(&buf, "%v : \"%v\", \n", v[0], v[0])
	}
	fmt.Fprintf(&buf, "\n}\n")
	fmt.Fprintf(&buf, `
	func (e %[1]s) String() string {
		if s, exist := _%[1]s_str[e]; exist {
			return s
		}
		return fmt.Sprintf("%[1]s%%d", uint16(e))
	}
		`, enumtype)

	return &buf, nil
}

func buildVersion(prefix, ver string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
	package %[1]s_version
	const ProtocolVersion = "%[2]s"
	`, prefix, ver)
	return &buf, nil
}

func buildMSGP(prefix string, cmddata, notidata [][]string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
		package %[1]s_msgp

		// MarshalBodyFn marshal body and append to oldBufferToAppend
		// and return newbuffer, body type(marshaltype,compress,encryption), error
		func MarshalBodyFn(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error) {
			newBuffer, err := body.(msgp.Marshaler).MarshalMsg(oldBuffToAppend)
			return newBuffer, 0, err
		}

		func UnmarshalPacket(h %[1]s_packet.Header,  bodyData []byte) (interface{}, error) {
			switch h.FlowType {
			case %[1]s_packet.Request:
				if int(h.Cmd) >= len(ReqUnmarshalMap) {
					return nil, fmt.Errorf("unknown request command: %%v %%v", 
					h.FlowType, %[1]s_idcmd.CommandID(h.Cmd))
				}
				return ReqUnmarshalMap[h.Cmd](h,bodyData)

			case %[1]s_packet.Response:
				if int(h.Cmd) >= len(RspUnmarshalMap) {
					return nil, fmt.Errorf("unknown response command: %%v %%v", 
					h.FlowType, %[1]s_idcmd.CommandID(h.Cmd))
				}
				return RspUnmarshalMap[h.Cmd](h,bodyData)

			case %[1]s_packet.Notification:
				if int(h.Cmd) >= len(NotiUnmarshalMap) {
					return nil, fmt.Errorf("unknown notification command: %%v %%v", 
					h.FlowType, %[1]s_idcmd.CommandID(h.Cmd))
				}
				return NotiUnmarshalMap[h.Cmd](h,bodyData)
			}
			return nil, fmt.Errorf("unknown packet FlowType %%v",h.FlowType)
		}
		`, prefix, ver)

	const unmarshalMapHeader = `
	var %[2]s = [...]func(h %[1]s_packet.Header,bodyData []byte) (interface{}, error) {
	`
	const unmarshalMapBody = "%[1]s_idcmd:  unmarshal_%[2]s%[3]s, \n"

	// req map
	fmt.Fprintf(&buf, unmarshalMapHeader, prefix, "ReqUnmarshalMap")
	for _, v := range cmddata {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s:  unmarshal_Req%[2]s, \n", prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	// rsp map
	fmt.Fprintf(&buf, unmarshalMapHeader, prefix, "RspUnmarshalMap")
	for _, v := range cmddata {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s:  unmarshal_Rsp%[2]s, \n", prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	fmt.Fprintf(&buf, unmarshalMapHeader, prefix, "NotiUnmarshalMap")
	// noti map
	for _, v := range notidata {
		fmt.Fprintf(&buf, "%[1]s_idnoti.%[2]s:  unmarshal_Noti%[2]s, \n", prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	const unmarshalFunc = `
	func unmarshal_%[1]s%[2]s(h %[3]s_packet.Header,bodyData []byte) (interface{}, error) {
		var args %[3]s_obj.%[1]s%[2]s_data
		if _, err := args.UnmarshalMsg(bodyData); err != nil {
			return nil, err
		}
		return &args, nil
	}
	`
	for _, v := range cmddata {
		fmt.Fprintf(&buf, unmarshalFunc, "Req", v[0], prefix)
		fmt.Fprintf(&buf, unmarshalFunc, "Rsp", v[0], prefix)
	}
	for _, v := range notidata {
		fmt.Fprintf(&buf, unmarshalFunc, "Noti", v[0], prefix)
	}
	return &buf, nil
}

func buildJSON(prefix string, cmddata, notidata [][]string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
		package %[1]s_json

		// marshal body and append to oldBufferToAppend
		// and return newbuffer, body type(marshaltype,compress,encryption), error
		func MarshalBodyFn(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error) {
			var newBuffer []byte
			mdata, err := json.Marshal(body)
			if err == nil {
				newBuffer = append(oldBuffToAppend, mdata...)
			}
			return newBuffer, 0, err
		}
		
		func UnmarshalPacket(h %[1]s_packet.Header,  bodyData []byte) (interface{}, error) {
			switch h.FlowType {
			case %[1]s_packet.Request:
				if int(h.Cmd) >= len(ReqUnmarshalMap) {
					return nil, fmt.Errorf("unknown request command: %%v %%v", 
					h.FlowType, %[1]s_idcmd.CommandID(h.Cmd))
				}
				return ReqUnmarshalMap[h.Cmd](h,bodyData)

			case %[1]s_packet.Response:
				if int(h.Cmd) >= len(RspUnmarshalMap) {
					return nil, fmt.Errorf("unknown response command: %%v %%v", 
					h.FlowType, %[1]s_idcmd.CommandID(h.Cmd))
				}
				return RspUnmarshalMap[h.Cmd](h,bodyData)

			case %[1]s_packet.Notification:
				if int(h.Cmd) >= len(NotiUnmarshalMap) {
					return nil, fmt.Errorf("unknown notification command: %%v %%v", 
					h.FlowType, %[1]s_idcmd.CommandID(h.Cmd))
				}
				return NotiUnmarshalMap[h.Cmd](h,bodyData)
			}
			return nil, fmt.Errorf("unknown packet FlowType %%v",h.FlowType)
		}
		`, prefix, ver)

	const unmarshalMapHeader = `
	var %[2]s = [...]func(h %[1]s_packet.Header,bodyData []byte) (interface{}, error) {
	`
	const unmarshalMapBody = "%[1]s_idcmd:  unmarshal_%[2]s%[3]s, \n"

	// req map
	fmt.Fprintf(&buf, unmarshalMapHeader, prefix, "ReqUnmarshalMap")
	for _, v := range cmddata {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s:  unmarshal_Req%[2]s, \n", prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	// rsp map
	fmt.Fprintf(&buf, unmarshalMapHeader, prefix, "RspUnmarshalMap")
	for _, v := range cmddata {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s:  unmarshal_Rsp%[2]s, \n", prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	fmt.Fprintf(&buf, unmarshalMapHeader, prefix, "NotiUnmarshalMap")
	// noti map
	for _, v := range notidata {
		fmt.Fprintf(&buf, "%[1]s_idnoti.%[2]s:  unmarshal_Noti%[2]s, \n", prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	const unmarshalFunc = `
	func unmarshal_%[1]s%[2]s(h %[3]s_packet.Header,bodyData []byte) (interface{}, error) {
		var args %[3]s_obj.%[1]s%[2]s_data
		if err := json.Unmarshal(bodyData, &args) ; err != nil {
			return nil, err
		}
		return &args, nil
	}
	`
	for _, v := range cmddata {
		fmt.Fprintf(&buf, unmarshalFunc, "Req", v[0], prefix)
		fmt.Fprintf(&buf, unmarshalFunc, "Rsp", v[0], prefix)
	}
	for _, v := range notidata {
		fmt.Fprintf(&buf, unmarshalFunc, "Noti", v[0], prefix)
	}
	return &buf, nil
}

func buildRecvRspFnMap(prefix string, cmddata, notidata [][]string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
		package %[1]s_client
		import (
			"fmt"
		)
		`, prefix)

	fmt.Fprintf(&buf,
		"\nvar RecvRspObjFnMap = [...]func(me interface{}, hd %[1]s_packet.Header, body interface{}) error {\n",
		prefix)
	for _, f := range cmddata {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s : defaultRecvObjFn_%[2]s,\n", prefix, f[0])
	}
	fmt.Fprintf(&buf, "\n}\n")
	for _, f := range cmddata {
		fmt.Fprintf(&buf, `
			func defaultRecvObjFn_%[2]s(me interface{}, hd %[1]s_packet.Header, body interface{}) error {
				robj , ok := body.(*%[1]s_obj.Rsp%[2]s_data)
				if !ok {
					return fmt.Errorf("packet mismatch %%v", body )
				}
				return fmt.Errorf("Not implemented %%v", robj)
			}
			`, prefix, f[0])
	}
	return &buf, nil
}

func buildRecvNotiFnMap(prefix string, cmddata, notidata [][]string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
		package %[1]s_client
		import (
			"fmt"
		)
		`, prefix)

	fmt.Fprintf(&buf,
		"\nvar RecvNotiObjFnMap = [...]func(me interface{}, hd %[1]s_packet.Header, body interface{}) error {\n",
		prefix)
	for _, f := range notidata {
		fmt.Fprintf(&buf, "%[1]s_idnoti.%[2]s : defaultRecvObjNotiFn_%[2]s,\n", prefix, f[0])
	}
	fmt.Fprintf(&buf, "\n}\n")
	for _, f := range notidata {
		fmt.Fprintf(&buf, `
			func defaultRecvObjNotiFn_%[2]s(me interface{}, hd %[1]s_packet.Header, body interface{}) error {
				robj , ok := body.(*%[1]s_obj.Noti%[2]s_data)
				if !ok {
					return fmt.Errorf("packet mismatch %%v", body )
				}
				return fmt.Errorf("Not implemented %%v", robj)
			}
			`, prefix, f[0])
	}
	return &buf, nil
}

func buildClientSendRecv(prefix string, cmddata, notidata [][]string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
		package %[1]s_client

		type C2SConnectI interface {
			SendRecv(
				cmd %[1]s_idcmd.CommandID, body interface{}) (
				%[1]s_packet.Header, interface{}, error)
		
			CheckAPI(hd %[1]s_packet.Header) error
	
		}
		
		`, prefix)

	for _, f := range cmddata {
		fmt.Fprintf(&buf, `
			func Call_%[2]s(c2sc C2SConnectI,arg *%[1]s_obj.Req%[2]s_data) (*%[1]s_obj.Rsp%[2]s_data, error) {
				if arg == nil {
					arg = &%[1]s_obj.Req%[2]s_data{}
				}
				hd, rsp, err := c2sc.SendRecv(
					%[1]s_idcmd.%[2]s,
					arg)
				if err != nil {
					return nil, err
				}
				robj := rsp.(*%[1]s_obj.Rsp%[2]s_data)
				return robj, c2sc.CheckAPI(hd)
			}
			`, prefix, f[0])
	}
	return &buf, nil
}

func buildPacket(prefix string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
	package %[1]s_packet

	type FlowType byte // packet flow type
	
	const (
		invalid FlowType = iota // make uninitalized packet error
		Request // Request for request packet (response packet expected)
		Response // Response is reply of request packet
		Notification // Notification is just send and forget packet
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
		return fmt.Sprintf("FlowType%%d", byte(e))
	}

	///////////////////////////////////////////////////////////////////////////////

	const (
		// MaxBodyLen set to max body len, affect send/recv buffer size
		MaxBodyLen = 0xfffff

		// HeaderLen fixed size of header
		HeaderLen = 4 + 4 + 2 + 2 + 1 + 1 + 2

		// MaxPacketLen max total packet size byte of raw packet
		MaxPacketLen = HeaderLen + MaxBodyLen
	)

	func (pk Packet) String() string {
		return fmt.Sprintf("Packet[%%v %%+v]", pk.Header, pk.Body)
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
				"Header[%%v:%%v ID:%%v Error:%%v bodyLen:%%v Compress:%%v Fill:%%v]",
				h.FlowType, h.Cmd, h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
		case invalid:
			return fmt.Sprintf(
				"Header[%%v:%%v ID:%%v Error:%%v bodyLen:%%v Compress:%%v Fill:%%v]",
				h.FlowType, h.Cmd, h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
		case Request:
			return fmt.Sprintf(
				"Header[%%v:%%v ID:%%v Error:%%v bodyLen:%%v Compress:%%v Fill:%%v]",
				h.FlowType, %[1]s_idcmd.CommandID(h.Cmd), h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
		case Response:
			return fmt.Sprintf(
				"Header[%%v:%%v ID:%%v Error:%%v bodyLen:%%v Compress:%%v Fill:%%v]",
				h.FlowType, %[1]s_idcmd.CommandID(h.Cmd), h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
		case Notification:
			return fmt.Sprintf(
				"Header[%%v:%%v ID:%%v Error:%%v bodyLen:%%v Compress:%%v Fill:%%v]",
				h.FlowType, %[1]s_idnoti.NotiID(h.Cmd), h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
		}
	}

	// Header is fixed size header of packet
	type Header struct {
		bodyLen   uint32              // set at marshal(Packet2Bytes)
		ID        uint32              // sender set, unique id per packet (wrap around reuse)
		Cmd       uint16              // sender set, application demux received packet
		ErrorCode %[1]s_error.ErrorCode // sender set, Response error
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
		h.ErrorCode = %[1]s_error.ErrorCode(binary.LittleEndian.Uint16(buf[10:12]))
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

	// func NewRecvPacketBuffer() *RecvPacketBuffer {
	// 	pb := &RecvPacketBuffer{
	// 		RecvBuffer: make([]byte, MaxPacketLen),
	// 		RecvLen:    0,
	// 	}
	// 	return pb
	// }

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
	func Packet2Bytes(pk *Packet, marshalBodyFn func(interface{}, []byte) ([]byte, byte, error)) ([]byte, error) {
		newbuf, bodytype, err := marshalBodyFn(pk.Body, make([]byte, HeaderLen, MaxPacketLen))
		if err != nil {
			return nil, err
		}
		bodyLen := len(newbuf) - HeaderLen
		if bodyLen > MaxBodyLen {
			return nil,
				fmt.Errorf("fail to serialize large packet %%v, %%v", pk.Header, bodyLen)
		}
		pk.Header.bodyType = bodytype
		pk.Header.bodyLen = uint32(bodyLen)
		pk.Header.toBytesAt(newbuf)
		return newbuf, nil
	}
	`, prefix)
	return &buf, nil
}

func buildObjTemplate(prefix string, cmddata, notidata [][]string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
		package %[1]s_obj
		/* protocol template 
		`, prefix)

	for _, f := range cmddata {
		fmt.Fprintf(&buf, `
		type Req%[2]s_data struct {
			Dummy uint8
		}
		type Rsp%[2]s_data struct {
			Dummy uint8
		}
		`, prefix, f[0])
	}
	for _, f := range notidata {
		fmt.Fprintf(&buf, `
		type Noti%[1]s_data struct {
			Dummy uint8
		}
		`, f[0])
	}
	fmt.Fprintf(&buf, `
	*/`)
	return &buf, nil
}

func buildDemuxReq2API(prefix string, cmddata, notidata [][]string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
		package %[1]s_server
		import (
			"fmt"
		)`, prefix)

	fmt.Fprintf(&buf, `
		var DemuxReq2APIFnMap = [...]func(
		c *ServeClientConn, hd %[1]s_packet.Header, robj interface{}) (
		%[1]s_packet.Header, interface{}, error){
		`, prefix)
	for _, f := range cmddata {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s: Req2API_%[2]s,\n", prefix, f[0])
	}
	fmt.Fprintf(&buf, "\n}   // DemuxReq2APIFnMap\n")

	for _, f := range cmddata {
		fmt.Fprintf(&buf, `
		func Req2API_%[2]s(
			c *ServeClientConn, hd %[1]s_packet.Header, robj interface{}) (
			%[1]s_packet.Header, interface{},  error) {
		req, ok := robj.(*%[1]s_obj.Req%[2]s_data)
		if !ok {
			return hd, nil, fmt.Errorf("Packet type miss match %%v", robj)
		}
		rhd, rsp, err := apifn_Req%[2]s(c, hd, req)
		return rhd, rsp, err
		}
		`, prefix, f[0])
	}
	return &buf, nil
}

func buildAPITemplate(prefix string, cmddata, notidata [][]string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
		package %[1]s_server
		/* api template 
		`, prefix)

	for _, f := range cmddata {
		fmt.Fprintf(&buf, `
		func apifn_Req%[2]s(
			c2sc *ServeClientConn, hd %[1]s_packet.Header, robj *%[1]s_obj.Req%[2]s_data) (
			%[1]s_packet.Header, *%[1]s_obj.Rsp%[2]s_data, error) {
			rhd := %[1]s_packet.Header{
				ErrorCode : %[1]s_error.None,
			}
			spacket := &%[1]s_obj.Rsp%[2]s_data{
			}
			return rhd, spacket, nil
		}
		`, prefix, f[0])
	}
	fmt.Fprintf(&buf, `
	*/`)
	return &buf, nil
}

func buildWasmConn(prefix string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
		package %[1]s_wasmconn

		import (
			"context"
			"fmt"
			"sync"
			"syscall/js"
		)

		`, prefix)

	fmt.Fprintf(&buf, `
	type Connection struct {
		remoteAddr   string
		conn         js.Value
		SendRecvStop func()
		sendCh       chan %[1]s_packet.Packet
	
		marshalBodyFn      func(interface{}, []byte) ([]byte, byte, error)
		handleRecvPacketFn func(header %[1]s_packet.Header, body []byte) error
		handleSentPacketFn func(header %[1]s_packet.Header) error
	}
	
	func (wsc *Connection) String() string {
		return fmt.Sprintf("Connection[%%v SendCh:%%v]",
			wsc.remoteAddr, len(wsc.sendCh))
	}
	
	func New(
		connAddr string,
		marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
		handleRecvPacketFn func(header %[1]s_packet.Header, body []byte) error,
		handleSentPacketFn func(header %[1]s_packet.Header) error,
	) *Connection {
		wsc := &Connection{
			remoteAddr:         connAddr,
			sendCh:             make(chan %[1]s_packet.Packet, 10),
			marshalBodyFn:      marshalBodyFn,
			handleRecvPacketFn: handleRecvPacketFn,
			handleSentPacketFn: handleSentPacketFn,
		}
		wsc.SendRecvStop = func() {
			JsLogErrorf("Too early SendRecvStop call %%v", wsc)
		}
		return wsc
	}
	
	func (wsc *Connection) Connect(ctx context.Context, wg *sync.WaitGroup) error {
		connCtx, ctxCancel := context.WithCancel(ctx)
		wsc.SendRecvStop = ctxCancel
	
		wsc.conn = js.Global().Get("WebSocket").New(wsc.remoteAddr)
		if !wsc.conn.Truthy() {
			err := fmt.Errorf("fail to connect %%v", wsc.remoteAddr)
			JsLogErrorf("%%v", err)
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
				var sendBuffer []byte
				sendBuffer, err = %[1]s_packet.Packet2Bytes(&pk, wsc.marshalBodyFn)
				if err != nil {
					break loop
				}
				if err = wsc.sendPacket(sendBuffer); err != nil {
					break loop
				}
				if err = wsc.handleSentPacketFn(pk.Header); err != nil {
					break loop
				}
			}
		}
		JsLogErrorf("end SendLoop %%v\n", err)
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
				rPk := %[1]s_packet.NewRecvPacketBufferByData(rdata)
				header, body, lerr := rPk.GetHeaderBody()
				if lerr != nil {
					JsLogError(lerr.Error())
					wsc.SendRecvStop()
					return nil
				} else {
					if err := wsc.handleRecvPacketFn(header, body); err != nil {
						JsLogErrorf("%%v", err)
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
	
	func (wsc *Connection) EnqueueSendPacket(pk %[1]s_packet.Packet) error {
		select {
		case wsc.sendCh <- pk:
			return nil
		default:
			return fmt.Errorf("Send channel full %%v", wsc)
		}
	}
	
	/////////
	
	func JsLogError(v ...interface{}) {
		js.Global().Get("console").Call("error", v...)
	}
	
	func JsLogErrorf(format string, v ...interface{}) {
		js.Global().Get("console").Call("error", fmt.Sprintf(format, v...))
	}
		`, prefix)
	return &buf, nil
}

func buildWSGorilla(prefix string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, makeGenComment())
	fmt.Fprintf(&buf, `
		package %[1]s_wsgorilla

		import (
			"context"
			"fmt"
			"net"
			"time"
		
			"github.com/gorilla/websocket"
		)

		`, prefix)

	fmt.Fprintf(&buf, `
	func SendControl(
		wsConn *websocket.Conn, mt int, PacketWriteTimeOut time.Duration) error {
	
		return wsConn.WriteControl(mt, []byte{}, time.Now().Add(PacketWriteTimeOut))
	}
	
	func SendPacket(wsConn *websocket.Conn, sendBuffer []byte) error {
		return wsConn.WriteMessage(websocket.BinaryMessage, sendBuffer)
	}
	
	func SendLoop(sendRecvCtx context.Context, SendRecvStop func(), wsConn *websocket.Conn,
		timeout time.Duration,
		SendCh chan %[1]s_packet.Packet,
		marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
		handleSentPacketFn func(header %[1]s_packet.Header) error,
	) error {
	
		defer SendRecvStop()
		var err error
	loop:
		for {
			select {
			case <-sendRecvCtx.Done():
				err = SendControl(wsConn, websocket.CloseMessage, timeout)
				break loop
			case pk := <-SendCh:
				if err = wsConn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
					break loop
				}
				sendBuffer, err := %[1]s_packet.Packet2Bytes(&pk, marshalBodyFn)
				if err != nil {
					break loop
				}
				if err = SendPacket(wsConn, sendBuffer); err != nil {
					break loop
				}
				if err = handleSentPacketFn(pk.Header); err != nil {
					break loop
				}
			}
		}
		return err
	}
	
	func RecvLoop(sendRecvCtx context.Context, SendRecvStop func(), wsConn *websocket.Conn,
		timeout time.Duration,
		HandleRecvPacketFn func(header %[1]s_packet.Header, body []byte) error) error {
	
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
	
	func RecvPacket(wsConn *websocket.Conn) (%[1]s_packet.Header, []byte, error) {
		mt, rdata, err := wsConn.ReadMessage()
		if err != nil {
			return %[1]s_packet.Header{}, nil, err
		}
		if mt != websocket.BinaryMessage {
			return %[1]s_packet.Header{}, nil, fmt.Errorf("message not binary %%v", mt)
		}
		return %[1]s_packet.NewRecvPacketBufferByData(rdata).GetHeaderBody()
	}
	`, prefix)
	return &buf, nil
}
