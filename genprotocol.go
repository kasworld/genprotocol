// Copyright 2015,2016,2017,2018,2019,2020 SeukWon Kang (kasworld@gmail.com)
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
	"bytes"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/kasworld/genprotocol/genlib"
)

type GenArgs struct {
	GenComment string
	Version    string
	BaseDir    string
	Prefix     string
	StatsType  string
	Verbose    bool
	CmdIDs     [][]string
	NotiIDs    [][]string
	ErrorIDs   [][]string
}

func loadGenArgs() (GenArgs, error) {
	g_verbose := flag.Bool("verbose", false, "show goimports file")
	g_ver := flag.String("ver", "", "protocol version")
	g_prefix := flag.String("prefix", "", "protocol prefix")
	g_basedir := flag.String("basedir", "", "base directory")
	g_statstype := flag.String("statstype", "", "stats element type, empty not generate")
	flag.Parse()

	if *g_prefix == "" {
		fmt.Println("prefix not set")
	}
	if *g_ver == "" {
		fmt.Println("ver not set")
	}
	if *g_basedir == "" {
		fmt.Println("base dir not set")
	}

	cmddatafile := path.Join(*g_basedir, *g_prefix+"_command.enum")
	cmddata, err := genlib.LoadEnumWithComment(cmddatafile)
	if err != nil {
		fmt.Printf("fail to load %v %v\n", cmddatafile, err)
		return GenArgs{}, err
	}
	notidatafile := path.Join(*g_basedir, *g_prefix+"_noti.enum")
	notidata, err := genlib.LoadEnumWithComment(notidatafile)
	if err != nil {
		fmt.Printf("fail to load %v %v\n", notidatafile, err)
		return GenArgs{}, err
	}
	errordatafile := path.Join(*g_basedir, *g_prefix+"_error.enum")
	errordata, err := genlib.LoadEnumWithComment(errordatafile)
	if err != nil {
		fmt.Printf("fail to load %v %v\n", errordatafile, err)
		return GenArgs{}, err
	}
	genArgs := GenArgs{
		GenComment: genlib.MakeGenComment(),
		Version:    *g_ver,
		BaseDir:    *g_basedir,
		Prefix:     *g_prefix,
		StatsType:  *g_statstype,
		Verbose:    *g_verbose,
		CmdIDs:     cmddata,
		NotiIDs:    notidata,
		ErrorIDs:   errordata,
	}
	return genArgs, nil
}

type MakeDest struct {
	Postfix        string
	Filename       string
	Fn             func(genArgs GenArgs, postfix string) *bytes.Buffer
	OverWriteExist bool
}

func main() {
	genArgs, err := loadGenArgs()
	if err != nil {
		return
	}
	// postfix, filename
	makeDatas := []MakeDest{
		// MakeDest{"_gendata", "",},
		MakeDest{"_version", "version_gen.go", buildVersion, true},
		MakeDest{"_idcmd", "command_gen.go", buildCommandEnum, true},
		MakeDest{"_idnoti", "noti_gen.go", buildNotiEnum, true},
		MakeDest{"_error", "error_gen.go", buildErrorEnum, true},

		MakeDest{"_const", "consttemplate_gen.go", buildConstTemplate, true},
		MakeDest{"_const", "const.go", buildConst, false},

		MakeDest{"_packet", "packet_gen.go", buildPacket, true},

		MakeDest{"_obj", "objtemplate_gen.go", buildObjTemplate, true},
		MakeDest{"_obj", "obj.go", buildObj, false},

		MakeDest{"_msgp", "serialize_gen.go", buildMSGP, true},
		MakeDest{"_json", "serialize_gen.go", buildJSON, true},
		MakeDest{"_gob", "serialize_gen.go", buildGOB, true},

		MakeDest{"_handlersp", "fnobjtemplate_gen.go", buildRecvRspFnObjTemplate, true},
		MakeDest{"_handlersp", "fnobj.go", buildRecvRspFnObj, false},

		MakeDest{"_handlersp", "fnbytestemplate_gen.go", buildRecvRspFnBytesTemplate, true},
		MakeDest{"_handlersp", "fnbytes.go", buildRecvRspFnBytes, false},

		MakeDest{"_handlereq", "fnobjtemplate_gen.go", buildRecvReqFnObjTemplate, true},
		MakeDest{"_handlereq", "fnobj.go", buildRecvReqFnObj, false},

		MakeDest{"_handlereq", "fnbytestemplate_gen.go", buildRecvReqFnBytesAPITemplate, true},
		MakeDest{"_handlereq", "fnbytes.go", buildRecvReqFnBytesAPI, false},

		MakeDest{"_handlenoti", "fnobjtemplate_gen.go", buildRecvNotiFnObjTemplate, true},
		MakeDest{"_handlenoti", "fnobj.go", buildRecvNotiFnObj, false},

		MakeDest{"_handlenoti", "fnbytestemplate_gen.go", buildRecvNotiFnBytesTemplate, true},
		MakeDest{"_handlenoti", "fnbytes.go", buildRecvNotiFnBytes, false},

		MakeDest{"_serveconnbyte", "serveconnbyte_gen.go", buildServeConnByte, true},
		MakeDest{"_connbytemanager", "connbytemanager_gen.go", buildConnByteManager, true},
		MakeDest{"_conntcp", "conntcp_gen.go", buildConnTCP, true},
		MakeDest{"_connwasm", "connwasm_gen.go", buildConnWasm, true},
		MakeDest{"_connwsgorilla", "connwsgorilla_gen.go", buildConnWSGorilla, true},
		MakeDest{"_loopwsgorilla", "loopwsgorilla_gen.go", buildLoopWSGorilla, true},
		MakeDest{"_looptcp", "looptcp_gen.go", buildLoopTCP, true},
		MakeDest{"_pid2rspfn", "pid2rspfn_gen.go", buildPID2RspFn, true},
		MakeDest{"_statnoti", "statnoti_gen.go", buildStatNoti, true},
		MakeDest{"_statcallapi", "statcallapi_gen.go", buildStatCallAPI, true},
		MakeDest{"_statserveapi", "statserveapi_gen.go", buildStatServeAPI, true},
		MakeDest{"_statapierror", "statapierror_gen.go", buildStatAPIError, true},
		MakeDest{"_authorize", "authorize_gen.go", buildAuthorize, true},
	}
	for _, v := range makeDatas {
		os.MkdirAll(path.Join(genArgs.BaseDir, genArgs.Prefix+v.Postfix), os.ModePerm)
		buf := v.Fn(genArgs, v.Postfix)
		filename := path.Join(genArgs.BaseDir, genArgs.Prefix+v.Postfix, v.Filename)
		if !v.OverWriteExist && genlib.IsFileExist(filename) {
			fmt.Printf("skip exist file %v\n", filename)
			continue
		}
		genlib.SaveTo(buf, filename, genArgs.Verbose)
	}

	if genArgs.StatsType != "" {
		dirToMake := [][2]string{
			{"_error", "ErrorCode"},
			{"_idcmd", "CommandID"},
			{"_idnoti", "NotiID"},
		}
		for _, v := range dirToMake {
			packagename := genArgs.Prefix + v[0]
			os.MkdirAll(path.Join(genArgs.BaseDir, packagename+"_stats"), os.ModePerm)
			buf := buildStatsCode(genArgs, packagename, v[1], genArgs.StatsType)
			filename := path.Join(genArgs.BaseDir, packagename+"_stats", packagename+"_stats_gen.go")
			genlib.SaveTo(buf, filename, genArgs.Verbose)
		}
	}
}

func buildVersion(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	const ProtocolVersion = "%[2]s"
	`, genArgs.Prefix+postfix, genArgs.Version)
	return &buf
}

func buildCommandEnum(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import "fmt"
	type CommandID uint16 // use in packet header, DO NOT CHANGE
	const (
	`, genArgs.Prefix+postfix)
	for i, v := range genArgs.CmdIDs {
		if i == 0 {
			fmt.Fprintf(&buf, "%v CommandID = iota // %v \n", v[0], v[1])
		} else {
			fmt.Fprintf(&buf, "%v // %v\n", v[0], v[1])
		}
	}
	fmt.Fprintf(&buf, `
	CommandID_Count int = iota 
	)
	var _CommandID2string = [CommandID_Count][2]string{
	`)
	for _, v := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%v : {\"%v\",\"%v\" },\n", v[0], v[0], v[1])
	}
	fmt.Fprintf(&buf, `
	}
	func (e CommandID) String() string {
		if e >=0 && e < CommandID(CommandID_Count) {
			return _CommandID2string[e][0]
		}
		return fmt.Sprintf("CommandID%%d", uint16(e))
	}

	func (e CommandID) CommentString() string {
		if e >=0 && e < CommandID(CommandID_Count) {
			return _CommandID2string[e][1]
		}
		return ""
	}


	var _string2CommandID = map[string]CommandID{
	`)
	for _, v := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "\"%v\" : %v, \n", v[0], v[0])
	}
	fmt.Fprintf(&buf, `
	}
	func  String2CommandID(s string) (CommandID, bool) {
		v, b :=  _string2CommandID[s]
		return v,b
	}
	`)
	return &buf
}

func buildNotiEnum(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import "fmt"
	type NotiID uint16 // use in packet header, DO NOT CHANGE
	const (
	`, genArgs.Prefix+postfix)
	for i, v := range genArgs.NotiIDs {
		if i == 0 {
			fmt.Fprintf(&buf, "%v NotiID = iota // %v \n", v[0], v[1])
		} else {
			fmt.Fprintf(&buf, "%v // %v\n", v[0], v[1])
		}
	}
	fmt.Fprintf(&buf, `
	NotiID_Count int = iota 
	)
	var _NotiID2string = [NotiID_Count][2]string{
	`)
	for _, v := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, "%v : {\"%v\",\"%v\" },\n", v[0], v[0], v[1])
	}
	fmt.Fprintf(&buf, `
	}
	func (e NotiID) String() string {
		if e >=0 && e < NotiID(NotiID_Count) {
			return _NotiID2string[e][0]
		}
		return fmt.Sprintf("NotiID%%d", uint16(e))
	}
	func (e NotiID) CommentString() string {
		if e >=0 && e < NotiID(NotiID_Count) {
			return _NotiID2string[e][1]
		}
		return ""
	}

	var _string2NotiID = map[string]NotiID{
	`)
	for _, v := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, "\"%v\" : %v, \n", v[0], v[0])
	}
	fmt.Fprintf(&buf, `
	}
	func  String2NotiID(s string) (NotiID, bool) {
		v, b :=  _string2NotiID[s]
		return v,b
	}
	`)
	return &buf
}

func buildErrorEnum(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import "fmt"
	type ErrorCode uint16 // use in packet header, DO NOT CHANGE
	const (
	`, genArgs.Prefix+postfix)
	for i, v := range genArgs.ErrorIDs {
		if i == 0 {
			fmt.Fprintf(&buf, "%v ErrorCode = iota // %v \n", v[0], v[1])
		} else {
			fmt.Fprintf(&buf, "%v // %v\n", v[0], v[1])
		}
	}
	fmt.Fprintf(&buf, `
	ErrorCode_Count int = iota 
	)
	var _ErrorCode2string = [ErrorCode_Count][2]string{
	`)
	for _, v := range genArgs.ErrorIDs {
		fmt.Fprintf(&buf, "%v : {\"%v\",\"%v\" },\n", v[0], v[0], v[1])
	}
	fmt.Fprintf(&buf, `
	}
	func (e ErrorCode) String() string {
		if e >=0 && e < ErrorCode(ErrorCode_Count) {
			return _ErrorCode2string[e][0]
		}
		return fmt.Sprintf("ErrorCode%%d", uint16(e))
	}
	func (e ErrorCode) CommentString() string {
		if e >=0 && e < ErrorCode(ErrorCode_Count) {
			return _ErrorCode2string[e][1]
		}
		return ""
	}

	// implement error interface
	func (e ErrorCode) Error() string {
		return "%[1]s." + e.String()
	}
	var _string2ErrorCode = map[string]ErrorCode{
	`, genArgs.Prefix+postfix)
	for _, v := range genArgs.ErrorIDs {
		fmt.Fprintf(&buf, "\"%v\" : %v, \n", v[0], v[0])
	}
	fmt.Fprintf(&buf, `
	}
	func  String2ErrorCode(s string) (ErrorCode, bool) {
		v, b :=  _string2ErrorCode[s]
		return v,b
	}
	`)
	return &buf
}

func buildConstTemplate(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	/*  copy to no _gen file, and edit it
	const (
		// MaxBodyLen set to max body len, affect send/recv buffer size
		MaxBodyLen = 0xffff

		// ServerAPICallTimeOutDur api call watchdog timer, 0 : no api timeout
		ServerAPICallTimeOutDur = time.Second * 2
	)
	*/
	`, genArgs.Prefix+postfix)
	return &buf
}
func buildConst(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	// edit as you need
	const (
		// MaxBodyLen set to max body len, affect send/recv buffer size
		MaxBodyLen = 0xffff
		// ServerAPICallTimeOutDur api call watchdog timer, 0 : no api timeout
		ServerAPICallTimeOutDur = time.Second * 2
	)
	`, genArgs.Prefix+postfix)
	return &buf
}

func buildPacket(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
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
		// HeaderLen fixed size of header
		HeaderLen = 4 + 4 + 2 + 2 + 1 + 1 + 2

		// MaxPacketLen max total packet size byte of raw packet
		MaxPacketLen = HeaderLen +  %[1]s_const.MaxBodyLen
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
				"Header[%%v:%%v ID:%%v Error:%%v BodyLen:%%v BodyType:%%v Fill:%%v]",
				h.FlowType, h.Cmd, h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
		case invalid:
			return fmt.Sprintf(
				"Header[%%v:%%v ID:%%v Error:%%v BodyLen:%%v BodyType:%%v Fill:%%v]",
				h.FlowType, h.Cmd, h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
		case Request:
			return fmt.Sprintf(
				"Header[%%v:%%v ID:%%v Error:%%v BodyLen:%%v BodyType:%%v Fill:%%v]",
				h.FlowType, %[1]s_idcmd.CommandID(h.Cmd), h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
		case Response:
			return fmt.Sprintf(
				"Header[%%v:%%v ID:%%v Error:%%v BodyLen:%%v BodyType:%%v Fill:%%v]",
				h.FlowType, %[1]s_idcmd.CommandID(h.Cmd), h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
		case Notification:
			return fmt.Sprintf(
				"Header[%%v:%%v ID:%%v Error:%%v BodyLen:%%v BodyType:%%v Fill:%%v]",
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

	func Bytes2HeaderBody(rdata []byte) (Header, []byte, error) {
		if len(rdata) < HeaderLen {
			return Header{}, nil, fmt.Errorf("header not complete")
		}
		header := MakeHeaderFromBytes(rdata)
		if len(rdata) != HeaderLen+int(header.bodyLen) {
			return header, nil, fmt.Errorf("packet not complete")
		}
		return header, rdata[HeaderLen : HeaderLen+int(header.bodyLen)], nil
	}

	func ReadHeaderBody(conn io.Reader) (Header, []byte, error) {
		recvLen := 0
		toRead := HeaderLen
		readBuffer := make([]byte, toRead)
		for recvLen < toRead {
			n, err := conn.Read(readBuffer[recvLen:toRead])
			if err != nil {
				return Header{}, nil, err
			}
			recvLen += n
		}
		header := MakeHeaderFromBytes(readBuffer)
		recvLen = 0
		toRead = int(header.bodyLen)
		readBuffer = make([]byte, toRead)
		for recvLen < toRead {
			n, err := conn.Read(readBuffer[recvLen:toRead])
			if err != nil {
				return header, nil, err
			}
			recvLen += n
		}
		return header, readBuffer, nil
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
		if bodyLen > %[1]s_const.MaxBodyLen {
			return nil,
				fmt.Errorf("fail to serialize large packet %%v, %%v", pk.Header, bodyLen)
		}
		pk.Header.bodyType = bodytype
		pk.Header.bodyLen = uint32(bodyLen)
		pk.Header.toBytesAt(newbuf)
		return newbuf, nil
	}
	`, genArgs.Prefix)
	return &buf
}

func buildObjTemplate(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	/* protocol object template 
	`, genArgs.Prefix+postfix)

	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, `
		// %[2]s %[3]s
		type Req%[2]s_data struct {
			Dummy uint8 // change as you need 
		}
		// %[2]s %[3]s
		type Rsp%[2]s_data struct {
			Dummy uint8 // change as you need 
		}
		`, genArgs.Prefix, f[0], f[1])
	}
	for _, f := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, `
		// %[1]s %[2]s
		type Noti%[1]s_data struct {
			Dummy uint8 // change as you need 
		}
		`, f[0], f[1])
	}
	fmt.Fprintf(&buf, `
	*/`)
	return &buf
}
func buildObj(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	// protocol object
	`, genArgs.Prefix+postfix)

	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, `
		// %[2]s %[3]s
		type Req%[2]s_data struct {
			Dummy uint8 // change as you need 
		}
		// %[2]s %[3]s
		type Rsp%[2]s_data struct {
			Dummy uint8 // change as you need 
		}
		`, genArgs.Prefix, f[0], f[1])
	}
	for _, f := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, `
		// %[1]s %[2]s
		type Noti%[1]s_data struct {
			Dummy uint8 // change as you need 
		}
		`, f[0], f[1])
	}
	return &buf
}

func buildMSGP(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s

	import (
		"github.com/tinylib/msgp/msgp"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	// MarshalBodyFn marshal body and append to oldBufferToAppend
	// and return newbuffer, body type, error
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
	`, genArgs.Prefix)

	const unmarshalMapHeader = `
	var %[2]s = [...]func(h %[1]s_packet.Header,bodyData []byte) (interface{}, error) {
	`
	const unmarshalMapBody = "%[1]s_idcmd:  unmarshal_%[2]s%[3]s, \n"

	// req map
	fmt.Fprintf(&buf, unmarshalMapHeader, genArgs.Prefix, "ReqUnmarshalMap")
	for _, v := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s:  unmarshal_Req%[2]s, \n", genArgs.Prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	// rsp map
	fmt.Fprintf(&buf, unmarshalMapHeader, genArgs.Prefix, "RspUnmarshalMap")
	for _, v := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s:  unmarshal_Rsp%[2]s, \n", genArgs.Prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	fmt.Fprintf(&buf, unmarshalMapHeader, genArgs.Prefix, "NotiUnmarshalMap")
	// noti map
	for _, v := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, "%[1]s_idnoti.%[2]s:  unmarshal_Noti%[2]s, \n", genArgs.Prefix, v[0])
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
	for _, v := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, unmarshalFunc, "Req", v[0], genArgs.Prefix)
		fmt.Fprintf(&buf, unmarshalFunc, "Rsp", v[0], genArgs.Prefix)
	}
	for _, v := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, unmarshalFunc, "Noti", v[0], genArgs.Prefix)
	}
	return &buf
}

func buildJSON(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
		package %[1]s
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	// marshal body and append to oldBufferToAppend
	// and return newbuffer, body type, error
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
	`, genArgs.Prefix)

	const unmarshalMapHeader = `
	var %[2]s = [...]func(h %[1]s_packet.Header,bodyData []byte) (interface{}, error) {
	`
	const unmarshalMapBody = "%[1]s_idcmd:  unmarshal_%[2]s%[3]s, \n"

	// req map
	fmt.Fprintf(&buf, unmarshalMapHeader, genArgs.Prefix, "ReqUnmarshalMap")
	for _, v := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s:  unmarshal_Req%[2]s, \n", genArgs.Prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	// rsp map
	fmt.Fprintf(&buf, unmarshalMapHeader, genArgs.Prefix, "RspUnmarshalMap")
	for _, v := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s:  unmarshal_Rsp%[2]s, \n", genArgs.Prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	fmt.Fprintf(&buf, unmarshalMapHeader, genArgs.Prefix, "NotiUnmarshalMap")
	// noti map
	for _, v := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, "%[1]s_idnoti.%[2]s:  unmarshal_Noti%[2]s, \n", genArgs.Prefix, v[0])
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
	for _, v := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, unmarshalFunc, "Req", v[0], genArgs.Prefix)
		fmt.Fprintf(&buf, unmarshalFunc, "Rsp", v[0], genArgs.Prefix)
	}
	for _, v := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, unmarshalFunc, "Noti", v[0], genArgs.Prefix)
	}
	return &buf
}

func buildGOB(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
		package %[1]s
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	// marshal body and append to oldBufferToAppend
	// and return newbuffer, body type, error
	func MarshalBodyFn(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error) {
		network := bytes.NewBuffer(oldBuffToAppend)
		enc := gob.NewEncoder(network)
		err := enc.Encode(body)
		return network.Bytes(), 0, err
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
	`, genArgs.Prefix)

	const unmarshalMapHeader = `
	var %[2]s = [...]func(h %[1]s_packet.Header,bodyData []byte) (interface{}, error) {
	`
	const unmarshalMapBody = "%[1]s_idcmd:  unmarshal_%[2]s%[3]s, \n"

	// req map
	fmt.Fprintf(&buf, unmarshalMapHeader, genArgs.Prefix, "ReqUnmarshalMap")
	for _, v := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s:  unmarshal_Req%[2]s, \n", genArgs.Prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	// rsp map
	fmt.Fprintf(&buf, unmarshalMapHeader, genArgs.Prefix, "RspUnmarshalMap")
	for _, v := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s:  unmarshal_Rsp%[2]s, \n", genArgs.Prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	fmt.Fprintf(&buf, unmarshalMapHeader, genArgs.Prefix, "NotiUnmarshalMap")
	// noti map
	for _, v := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, "%[1]s_idnoti.%[2]s:  unmarshal_Noti%[2]s, \n", genArgs.Prefix, v[0])
	}
	fmt.Fprintf(&buf, "}\n")

	const unmarshalFunc = `
	func unmarshal_%[1]s%[2]s(h %[3]s_packet.Header,bodyData []byte) (interface{}, error) {
		var args %[3]s_obj.%[1]s%[2]s_data
		network := bytes.NewBuffer(bodyData)
		dec := gob.NewDecoder(network)
		err := dec.Decode(&args)
		return &args, err
	}
	`
	for _, v := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, unmarshalFunc, "Req", v[0], genArgs.Prefix)
		fmt.Fprintf(&buf, unmarshalFunc, "Rsp", v[0], genArgs.Prefix)
	}
	for _, v := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, unmarshalFunc, "Noti", v[0], genArgs.Prefix)
	}
	return &buf
}

func buildRecvRspFnObjTemplate(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	/* obj base demux fn map template 
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf,
		"\nvar DemuxRsp2ObjFnMap = [...]func(me interface{}, hd %[1]s_packet.Header, body interface{}) error {\n",
		genArgs.Prefix)
	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s : objRecvRspFn_%[2]s, // %[2]s %[3]s \n",
			genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, "\n}\n")
	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, `
	// %[2]s %[3]s 
	func objRecvRspFn_%[2]s(me interface{}, hd %[1]s_packet.Header, body interface{}) error {
		robj , ok := body.(*%[1]s_obj.Rsp%[2]s_data)
		if !ok {
			return fmt.Errorf("packet mismatch %%v", body )
		}
		return fmt.Errorf("Not implemented %%v", robj)
	}
	`, genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, `
	*/`)
	return &buf
}
func buildRecvRspFnObj(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	// obj base demux fn map
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf,
		"\nvar DemuxRsp2ObjFnMap = [...]func(me interface{}, hd %[1]s_packet.Header, body interface{}) error {\n",
		genArgs.Prefix)
	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s : objRecvRspFn_%[2]s, // %[2]s %[3]s \n",
			genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, "\n}\n")
	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, `
	// %[2]s %[3]s 
	func objRecvRspFn_%[2]s(me interface{}, hd %[1]s_packet.Header, body interface{}) error {
		robj , ok := body.(*%[1]s_obj.Rsp%[2]s_data)
		if !ok {
			return fmt.Errorf("packet mismatch %%v", body )
		}
		return fmt.Errorf("Not implemented %%v", robj)
	}
	`, genArgs.Prefix, f[0], f[1])
	}
	return &buf
}

func buildRecvRspFnBytesTemplate(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	/* bytes base demux fn map template 
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf,
		"\nvar DemuxRsp2BytesFnMap = [...]func(me interface{}, hd %[1]s_packet.Header, rbody []byte) error {\n",
		genArgs.Prefix)
	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s : bytesRecvRspFn_%[2]s, // %[2]s %[3]s \n",
			genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, "\n}\n")
	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, `
	// %[2]s %[3]s
	func bytesRecvRspFn_%[2]s(me interface{}, hd %[1]s_packet.Header, rbody []byte) error {
		robj, err := %[1]s_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return  fmt.Errorf("Packet type miss match %%v", rbody)
		}
		recved , ok := robj.(*%[1]s_obj.Rsp%[2]s_data)
		if !ok {
			return fmt.Errorf("packet mismatch %%v", robj )
		}
		return fmt.Errorf("Not implemented %%v", recved)
	}
	`, genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, `
	*/`)
	return &buf
}
func buildRecvRspFnBytes(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	// bytes base demux fn map
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf,
		"\nvar DemuxRsp2BytesFnMap = [...]func(me interface{}, hd %[1]s_packet.Header, rbody []byte) error {\n",
		genArgs.Prefix)
	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s : bytesRecvRspFn_%[2]s, // %[2]s %[3]s \n",
			genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, "\n}\n")
	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, `
	// %[2]s %[3]s
	func bytesRecvRspFn_%[2]s(me interface{}, hd %[1]s_packet.Header, rbody []byte) error {
		robj, err := %[1]s_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return  fmt.Errorf("Packet type miss match %%v", rbody)
		}
		recved , ok := robj.(*%[1]s_obj.Rsp%[2]s_data)
		if !ok {
			return fmt.Errorf("packet mismatch %%v", robj )
		}
		return fmt.Errorf("Not implemented %%v", recved)
	}
	`, genArgs.Prefix, f[0], f[1])
	}
	return &buf
}

func buildRecvNotiFnObjTemplate(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	/* obj base demux fn map template 
	`, genArgs.Prefix+postfix)

	fmt.Fprintf(&buf,
		"\nvar DemuxNoti2ObjFnMap = [...]func(me interface{}, hd %[1]s_packet.Header, body interface{}) error {\n",
		genArgs.Prefix)
	for _, f := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, "%[1]s_idnoti.%[2]s : objRecvNotiFn_%[2]s, // %[2]s %[3]s\n",
			genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, "\n}\n")
	for _, f := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, `
	// %[2]s %[3]s
	func objRecvNotiFn_%[2]s(me interface{}, hd %[1]s_packet.Header, body interface{}) error {
		robj , ok := body.(*%[1]s_obj.Noti%[2]s_data)
		if !ok {
			return fmt.Errorf("packet mismatch %%v", body )
		}
		return fmt.Errorf("Not implemented %%v", robj)
	}
	`, genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, `
	*/`)
	return &buf
}
func buildRecvNotiFnObj(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	// obj base demux fn map
	`, genArgs.Prefix+postfix)

	fmt.Fprintf(&buf,
		"\nvar DemuxNoti2ObjFnMap = [...]func(me interface{}, hd %[1]s_packet.Header, body interface{}) error {\n",
		genArgs.Prefix)
	for _, f := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, "%[1]s_idnoti.%[2]s : objRecvNotiFn_%[2]s, // %[2]s %[3]s\n",
			genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, "\n}\n")
	for _, f := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, `
	// %[2]s %[3]s
	func objRecvNotiFn_%[2]s(me interface{}, hd %[1]s_packet.Header, body interface{}) error {
		robj , ok := body.(*%[1]s_obj.Noti%[2]s_data)
		if !ok {
			return fmt.Errorf("packet mismatch %%v", body )
		}
		return fmt.Errorf("Not implemented %%v", robj)
	}
	`, genArgs.Prefix, f[0], f[1])
	}
	return &buf
}

func buildRecvNotiFnBytesTemplate(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	/* bytes base demux fn map template 
	`, genArgs.Prefix+postfix)

	fmt.Fprintf(&buf,
		"\nvar DemuxNoti2ByteFnMap = [...]func(me interface{}, hd %[1]s_packet.Header, rbody []byte) error {\n",
		genArgs.Prefix)
	for _, f := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, "%[1]s_idnoti.%[2]s : bytesRecvNotiFn_%[2]s,// %[2]s %[3]s\n",
			genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, "\n}\n")
	for _, f := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, `
	// %[2]s %[3]s
	func bytesRecvNotiFn_%[2]s(me interface{}, hd %[1]s_packet.Header, rbody []byte) error {
		robj, err := %[1]s_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return fmt.Errorf("Packet type miss match %%v", rbody)
		}
		recved , ok := robj.(*%[1]s_obj.Noti%[2]s_data)
		if !ok {
			return fmt.Errorf("packet mismatch %%v", robj )
		}
		return fmt.Errorf("Not implemented %%v", recved)
	}
	`, genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, `
	*/`)
	return &buf
}
func buildRecvNotiFnBytes(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	// bytes base demux fn map
	`, genArgs.Prefix+postfix)

	fmt.Fprintf(&buf,
		"\nvar DemuxNoti2ByteFnMap = [...]func(me interface{}, hd %[1]s_packet.Header, rbody []byte) error {\n",
		genArgs.Prefix)
	for _, f := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, "%[1]s_idnoti.%[2]s : bytesRecvNotiFn_%[2]s,// %[2]s %[3]s\n",
			genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, "\n}\n")
	for _, f := range genArgs.NotiIDs {
		fmt.Fprintf(&buf, `
	// %[2]s %[3]s
	func bytesRecvNotiFn_%[2]s(me interface{}, hd %[1]s_packet.Header, rbody []byte) error {
		robj, err := %[1]s_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return fmt.Errorf("Packet type miss match %%v", rbody)
		}
		recved , ok := robj.(*%[1]s_obj.Noti%[2]s_data)
		if !ok {
			return fmt.Errorf("packet mismatch %%v", robj )
		}
		return fmt.Errorf("Not implemented %%v", recved)
	}
	`, genArgs.Prefix, f[0], f[1])
	}
	return &buf
}

func buildRecvReqFnObjTemplate(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	`, genArgs.Prefix+postfix)

	fmt.Fprintf(&buf, `
	/* obj base demux fn map template 
	var DemuxReq2ObjAPIFnMap = [...]func(
		me interface{}, hd %[1]s_packet.Header, robj interface{}) (
		%[1]s_packet.Header, interface{}, error){
	`, genArgs.Prefix)
	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s: Req2ObjAPI_%[2]s,// %[2]s %[3]s\n",
			genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, "\n}   // DemuxReq2ObjAPIFnMap\n")

	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, `
	// %[2]s %[3]s
	func Req2ObjAPI_%[2]s(
		me interface{}, hd %[1]s_packet.Header, robj interface{}) (
		%[1]s_packet.Header, interface{},  error) {
		req, ok := robj.(*%[1]s_obj.Req%[2]s_data)
		if !ok {
			return hd, nil, fmt.Errorf("Packet type miss match %%v", robj)
		}
		rhd, rsp, err := objAPIFn_Req%[2]s(me, hd, req)
		return rhd, rsp, err
	}
	// %[2]s %[3]s
	func objAPIFn_Req%[2]s(
		me interface{}, hd %[1]s_packet.Header, robj *%[1]s_obj.Req%[2]s_data) (
		%[1]s_packet.Header, *%[1]s_obj.Rsp%[2]s_data, error) {
		sendHeader := %[1]s_packet.Header{
			ErrorCode : %[1]s_error.None,
		}
		sendBody := &%[1]s_obj.Rsp%[2]s_data{
		}
		return sendHeader, sendBody, nil
	}
		`, genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, `
	*/`)
	return &buf
}
func buildRecvReqFnObj(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	`, genArgs.Prefix+postfix)

	fmt.Fprintf(&buf, `
	// obj base demux fn map
	var DemuxReq2ObjAPIFnMap = [...]func(
		me interface{}, hd %[1]s_packet.Header, robj interface{}) (
		%[1]s_packet.Header, interface{}, error){
	`, genArgs.Prefix)
	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s: Req2ObjAPI_%[2]s,// %[2]s %[3]s\n",
			genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, "\n}   // DemuxReq2ObjAPIFnMap\n")

	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, `
	// %[2]s %[3]s
	func Req2ObjAPI_%[2]s(
		me interface{}, hd %[1]s_packet.Header, robj interface{}) (
		%[1]s_packet.Header, interface{},  error) {
		req, ok := robj.(*%[1]s_obj.Req%[2]s_data)
		if !ok {
			return hd, nil, fmt.Errorf("Packet type miss match %%v", robj)
		}
		rhd, rsp, err := objAPIFn_Req%[2]s(me, hd, req)
		return rhd, rsp, err
	}
	// %[2]s %[3]s
	func objAPIFn_Req%[2]s(
		me interface{}, hd %[1]s_packet.Header, robj *%[1]s_obj.Req%[2]s_data) (
		%[1]s_packet.Header, *%[1]s_obj.Rsp%[2]s_data, error) {
		sendHeader := %[1]s_packet.Header{
			ErrorCode : %[1]s_error.None,
		}
		sendBody := &%[1]s_obj.Rsp%[2]s_data{
		}
		return sendHeader, sendBody, nil
	}
		`, genArgs.Prefix, f[0], f[1])
	}
	return &buf
}

func buildRecvReqFnBytesAPITemplate(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	`, genArgs.Prefix+postfix)

	fmt.Fprintf(&buf, `
	/* bytes base fn map api template , unmarshal in api
	var DemuxReq2BytesAPIFnMap = [...]func(
		me interface{}, hd %[1]s_packet.Header, rbody []byte) (
		%[1]s_packet.Header, interface{}, error){
	`, genArgs.Prefix)
	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s: bytesAPIFn_Req%[2]s,// %[2]s %[3]s\n",
			genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, "\n}   // DemuxReq2BytesAPIFnMap\n")

	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, `
	// %[2]s %[3]s		
	func bytesAPIFn_Req%[2]s(
		me interface{}, hd %[1]s_packet.Header, rbody []byte) (
		%[1]s_packet.Header, interface{}, error) {
		// robj, err := %[1]s_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %%v", rbody)
		// }
		// recvBody, ok := robj.(*%[1]s_obj.Req%[2]s_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %%v", robj)
		// }
		// _ = recvBody
		
		sendHeader := %[1]s_packet.Header{
			ErrorCode : %[1]s_error.None,
		}
		sendBody := &%[1]s_obj.Rsp%[2]s_data{
		}
		return sendHeader, sendBody, nil
	}
		`, genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, `
	*/`)
	return &buf
}
func buildRecvReqFnBytesAPI(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	`, genArgs.Prefix+postfix)

	fmt.Fprintf(&buf, `
	// bytes base fn map api, unmarshal in api
	var DemuxReq2BytesAPIFnMap = [...]func(
		me interface{}, hd %[1]s_packet.Header, rbody []byte) (
		%[1]s_packet.Header, interface{}, error){
	`, genArgs.Prefix)
	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, "%[1]s_idcmd.%[2]s: bytesAPIFn_Req%[2]s,// %[2]s %[3]s\n",
			genArgs.Prefix, f[0], f[1])
	}
	fmt.Fprintf(&buf, "\n}   // DemuxReq2BytesAPIFnMap\n")

	for _, f := range genArgs.CmdIDs {
		fmt.Fprintf(&buf, `
	// %[2]s %[3]s		
	func bytesAPIFn_Req%[2]s(
		me interface{}, hd %[1]s_packet.Header, rbody []byte) (
		%[1]s_packet.Header, interface{}, error) {
		// robj, err := %[1]s_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %%v", rbody)
		// }
		// recvBody, ok := robj.(*%[1]s_obj.Req%[2]s_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %%v", robj)
		// }
		// _ = recvBody
		
		sendHeader := %[1]s_packet.Header{
			ErrorCode : %[1]s_error.None,
		}
		sendBody := &%[1]s_obj.Rsp%[2]s_data{
		}
		return sendHeader, sendBody, nil
	}
		`, genArgs.Prefix, f[0], f[1])
	}
	return &buf
}

func buildServeConnByte(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"context"
		"fmt"
		"net"
		"time"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	type CounterI interface {
		Inc()
	}

	func (scb *ServeConnByte) String() string {
		return fmt.Sprintf("ServeConnByte[SendCh:%%v/%%v]",
			len(scb.sendCh), cap(scb.sendCh))
	}
	
	type ServeConnByte struct {
		connData         interface{} // custom data for this conn
		sendCh           chan %[1]s_packet.Packet
		sendRecvStop     func()
		authorCmdList    *%[1]s_authorize.AuthorizedCmds
		pid2ApiStatObj   *%[1]s_statserveapi.PacketID2StatObj
		apiStat     	 *%[1]s_statserveapi.StatServeAPI
		notiStat         *%[1]s_statnoti.StatNotification
		errorStat        *%[1]s_statapierror.StatAPIError
		sendCounter      CounterI
		recvCounter      CounterI

		demuxReq2BytesAPIFnMap [%[1]s_idcmd.CommandID_Count]func(
			me interface{}, hd %[1]s_packet.Header, rbody []byte) (
			%[1]s_packet.Header, interface{}, error)
	}
	// New with stats local
	func New(
		connData interface{},
		sendBufferSize int,
		authorCmdList *%[1]s_authorize.AuthorizedCmds,
		sendCounter, recvCounter CounterI,
		demuxReq2BytesAPIFnMap [%[1]s_idcmd.CommandID_Count]func(
			me interface{}, hd %[1]s_packet.Header, rbody []byte) (
			%[1]s_packet.Header, interface{}, error),
	) *ServeConnByte {
		scb := &ServeConnByte{
			connData:               connData,
			sendCh:                 make(chan %[1]s_packet.Packet, sendBufferSize),
			pid2ApiStatObj:         %[1]s_statserveapi.NewPacketID2StatObj(),
			apiStat:                %[1]s_statserveapi.New(),
			notiStat:               %[1]s_statnoti.New(),
			errorStat:              %[1]s_statapierror.New(),
			sendCounter:            sendCounter,
			recvCounter:            recvCounter,
			authorCmdList:          authorCmdList,
			demuxReq2BytesAPIFnMap: demuxReq2BytesAPIFnMap,
		}
		scb.sendRecvStop = func() {
			fmt.Printf("Too early sendRecvStop call %%v\n", scb)
		}
		return scb
	}
	// NewWithStats with stats global
	func NewWithStats(
		connData interface{},
		sendBufferSize int,
		authorCmdList    *%[1]s_authorize.AuthorizedCmds,
		sendCounter, recvCounter CounterI,
		apiStat          *%[1]s_statserveapi.StatServeAPI,
		notiStat         *%[1]s_statnoti.StatNotification,
		errorStat        *%[1]s_statapierror.StatAPIError,
		demuxReq2BytesAPIFnMap [%[1]s_idcmd.CommandID_Count]func(
			me interface{}, hd %[1]s_packet.Header, rbody []byte) (
			%[1]s_packet.Header, interface{}, error),
	) *ServeConnByte {
		scb := &ServeConnByte{
			connData:               connData,
			sendCh:                 make(chan %[1]s_packet.Packet, sendBufferSize),
			pid2ApiStatObj:         %[1]s_statserveapi.NewPacketID2StatObj(),
			apiStat:                apiStat,
			notiStat:               notiStat,
			errorStat:              errorStat,
			sendCounter:            sendCounter,
			recvCounter:            recvCounter,
			authorCmdList:          authorCmdList,
			demuxReq2BytesAPIFnMap: demuxReq2BytesAPIFnMap,
		}
		scb.sendRecvStop = func() {
			fmt.Printf("Too early sendRecvStop call %%v\n", scb)
		}
		return scb
	}

	func (scb *ServeConnByte) Disconnect() {
		scb.sendRecvStop()
	}
	func (scb *ServeConnByte) GetConnData() interface{} {
		return scb.connData
	}
	func (scb *ServeConnByte) GetAPIStat() *%[1]s_statserveapi.StatServeAPI {
		return scb.apiStat
	}
	func (scb *ServeConnByte) GetNotiStat() *%[1]s_statnoti.StatNotification {
		return scb.notiStat
	}
	func (scb *ServeConnByte) GetErrorStat() *%[1]s_statapierror.StatAPIError {
		return scb.errorStat
	}
	func (scb *ServeConnByte) GetAuthorCmdList() *%[1]s_authorize.AuthorizedCmds {
		return scb.authorCmdList
	}
	func (scb *ServeConnByte) StartServeWS(
		mainctx context.Context, conn *websocket.Conn,
		readTimeoutSec, writeTimeoutSec time.Duration,
		marshalfn func(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error),
	) error {
		var returnerr error
		sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
		scb.sendRecvStop = sendRecvCancel
		go func() {
			err := %[1]s_loopwsgorilla.RecvLoop(sendRecvCtx, scb.sendRecvStop, conn,
				readTimeoutSec, scb.handleRecvPacket)
			if err != nil {
				returnerr = fmt.Errorf("end RecvLoop %%v", err)
			}
		}()
		go func() {
			err := %[1]s_loopwsgorilla.SendLoop(sendRecvCtx, scb.sendRecvStop, conn,
				writeTimeoutSec, scb.sendCh,
				marshalfn, scb.handleSentPacket)
			if err != nil {
				returnerr = fmt.Errorf("end SendLoop %%v", err)
			}
		}()
	loop:
		for {
			select {
			case <-sendRecvCtx.Done():
				break loop
			}
		}
		return returnerr
	}
	func (scb *ServeConnByte) StartServeTCP(
		mainctx context.Context, conn *net.TCPConn,
		readTimeoutSec, writeTimeoutSec time.Duration,
		marshalfn func(body interface{}, oldBuffToAppend []byte) ([]byte, byte, error),
	) error {
		var returnerr error
		sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
		scb.sendRecvStop = sendRecvCancel
		go func() {
			err := %[1]s_looptcp.RecvLoop(sendRecvCtx, scb.sendRecvStop, conn,
				readTimeoutSec, scb.handleRecvPacket)
			if err != nil {
				returnerr = fmt.Errorf("end RecvLoop %%v", err)
			}
		}()
		go func() {
			err := %[1]s_looptcp.SendLoop(sendRecvCtx, scb.sendRecvStop, conn,
				writeTimeoutSec, scb.sendCh,
				marshalfn, scb.handleSentPacket)
			if err != nil {
				returnerr = fmt.Errorf("end SendLoop %%v", err)
			}
		}()
	loop:
		for {
			select {
			case <-sendRecvCtx.Done():
				break loop
			}
		}
		return returnerr
	}
	func (scb *ServeConnByte) handleSentPacket(header %[1]s_packet.Header) error {
		scb.sendCounter.Inc()
		switch header.FlowType {
		default:
			return fmt.Errorf("invalid packet type %%s %%v", scb, header)
	
		case %[1]s_packet.Request:
			return fmt.Errorf("request packet not supported %%s %%v", scb, header)
	
		case %[1]s_packet.Response:
			statOjb := scb.pid2ApiStatObj.Del(header.ID)
			if statOjb != nil {
				statOjb.AfterSendRsp(header)
			} else {
				return fmt.Errorf("send StatObj not found %%v", header)
			}
		case %[1]s_packet.Notification:
			scb.notiStat.Add(header)
		}
		return nil
	}
	func (scb *ServeConnByte) handleRecvPacket(rheader %[1]s_packet.Header, rbody []byte) error {
		scb.recvCounter.Inc()
		if rheader.FlowType != %[1]s_packet.Request {
			return fmt.Errorf("Unexpected rheader packet type: %%v", rheader)
		}
		if int(rheader.Cmd) >= len(scb.demuxReq2BytesAPIFnMap) {
			return fmt.Errorf("Invalid rheader command %%v", rheader)
		}
		if !scb.authorCmdList.CheckAuth(%[1]s_idcmd.CommandID(rheader.Cmd)) {
			return fmt.Errorf("Not authorized packet %%v", rheader)
		}

		statObj, err := scb.apiStat.AfterRecvReqHeader(rheader)
		if   err != nil {
			return err
		} 
		if err := scb.pid2ApiStatObj.Add(rheader.ID, statObj); err != nil {
			return err
		}
		statObj.BeforeAPICall()


		var sheader %[1]s_packet.Header
		var sbody interface{}
		var apierr error
		if %[1]s_const.ServerAPICallTimeOutDur != 0 {
			// timeout api call
			apiResult := scb.callAPI_timed(rheader, rbody)
			sheader, sbody, apierr = apiResult.header, apiResult.body, apiResult.err
		} else {
			// no timeout api call
			fn := scb.demuxReq2BytesAPIFnMap[rheader.Cmd]
			sheader, sbody, apierr = fn(scb, rheader, rbody)
		}

		statObj.AfterAPICall()

		scb.errorStat.Inc(%[1]s_idcmd.CommandID(rheader.Cmd), sheader.ErrorCode)
		if apierr != nil {
			return apierr
		}
		if sbody == nil {
			return fmt.Errorf("Response body nil")
		}
		sheader.FlowType = %[1]s_packet.Response
		sheader.Cmd = rheader.Cmd
		sheader.ID = rheader.ID
		rpk := %[1]s_packet.Packet{
			Header: sheader,
			Body:   sbody,
		}
		return scb.EnqueueSendPacket(rpk)
	}
	type callAPIResult struct {
		header %[1]s_packet.Header
		body   interface{}
		err    error
	}
	func (scb *ServeConnByte) callAPI_timed(rheader %[1]s_packet.Header, rbody []byte) callAPIResult {
		rtnCh := make(chan callAPIResult, 1)
		go func(rtnCh chan callAPIResult, rheader %[1]s_packet.Header, rbody []byte) {
			fn := scb.demuxReq2BytesAPIFnMap[rheader.Cmd]
			sheader, sbody, apierr := fn(scb, rheader, rbody)
			rtnCh <- callAPIResult{sheader, sbody, apierr}
		}(rtnCh, rheader, rbody)
		timeoutTk := time.NewTicker(%[1]s_const.ServerAPICallTimeOutDur)
		defer timeoutTk.Stop()
		select {
		case apiResult := <-rtnCh:
			return apiResult
		case <-timeoutTk.C:
			return callAPIResult{rheader, nil, fmt.Errorf("APICall Timeout %%v", rheader)}
		}
	}
	func (scb *ServeConnByte) EnqueueSendPacket(pk %[1]s_packet.Packet) error {
		select {
		case scb.sendCh <- pk:
			return nil
		default:
			return fmt.Errorf("Send channel full %%v", scb)
		}
	}
	func (scb *ServeConnByte) SendNotiPacket(
		cmd %[1]s_idnoti.NotiID, body interface{}) error {
		err := scb.EnqueueSendPacket(%[1]s_packet.Packet{
			%[1]s_packet.Header{
				Cmd:      uint16(cmd),
				FlowType: %[1]s_packet.Notification,
			},
			body,
		})
		if err != nil {
			scb.Disconnect()
		}
		return err
	}
	
	`, genArgs.Prefix)
	return &buf
}

func buildConnByteManager(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"fmt"
		"sync"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	type Manager struct {
		mutex   sync.RWMutex
		id2Conn map[string]*%[1]s_serveconnbyte.ServeConnByte
	}
	func New() *Manager {
		rtn := &Manager{
			id2Conn: make(map[string]*%[1]s_serveconnbyte.ServeConnByte),
		}
		return rtn
	}
	func (cm *Manager) Add(id string, c2sc *%[1]s_serveconnbyte.ServeConnByte) error {
		cm.mutex.Lock()
		defer cm.mutex.Unlock()
		if cm.id2Conn[id] != nil {
			return fmt.Errorf("already exist %%v", id)
		}
		cm.id2Conn[id] = c2sc
		return nil
	}
	func (cm *Manager) Del(id string) error {
		cm.mutex.Lock()
		defer cm.mutex.Unlock()
		if cm.id2Conn[id] == nil {
			return fmt.Errorf("not exist %%v", id)
		}
		delete(cm.id2Conn, id)
		return nil
	}
	func (cm *Manager) Get(id string) *%[1]s_serveconnbyte.ServeConnByte {
		cm.mutex.RLock()
		defer cm.mutex.RUnlock()
		return cm.id2Conn[id]
	}
	func (cm *Manager) Len() int {
		return len(cm.id2Conn)
	}
	func (cm *Manager) GetList() []*%[1]s_serveconnbyte.ServeConnByte {
		rtn := make([]*%[1]s_serveconnbyte.ServeConnByte, 0, len(cm.id2Conn))
		cm.mutex.RLock()
		defer cm.mutex.RUnlock()
		for _, v := range cm.id2Conn {
			rtn = append(rtn, v)
		}
		return rtn
	}
	`, genArgs.Prefix)
	return &buf
}

func buildConnTCP(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"context"
		"fmt"
		"net"
		"sync"
		"time"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	type Connection struct {
		conn         *net.TCPConn
		sendCh       chan %[1]s_packet.Packet
		sendRecvStop func()
	
		readTimeoutSec     time.Duration
		writeTimeoutSec    time.Duration
		marshalBodyFn      func(interface{}, []byte) ([]byte, byte, error)
		handleRecvPacketFn func(header %[1]s_packet.Header, body []byte) error
		handleSentPacketFn func(header %[1]s_packet.Header) error
	}
	
	func New(
		readTimeoutSec, writeTimeoutSec time.Duration,
		marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
		handleRecvPacketFn func(header %[1]s_packet.Header, body []byte) error,
		handleSentPacketFn func(header %[1]s_packet.Header) error,
	) *Connection {
		tc := &Connection{
			sendCh:             make(chan %[1]s_packet.Packet, 10),
			readTimeoutSec:     readTimeoutSec,
			writeTimeoutSec:    writeTimeoutSec,
			marshalBodyFn:      marshalBodyFn,
			handleRecvPacketFn: handleRecvPacketFn,
			handleSentPacketFn: handleSentPacketFn,
		}
	
		tc.sendRecvStop = func() {
			fmt.Printf("Too early sendRecvStop call %%v\n", tc)
		}
		return tc
	}
	
	func (tc *Connection) ConnectTo(remoteAddr string) error {
		tcpaddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
		if err != nil {
			return err
		}
		tc.conn, err = net.DialTCP("tcp", nil, tcpaddr)
		if err != nil {
			return err
		}
		return nil
	}
	
	func (tc *Connection) Cleanup() {
		tc.sendRecvStop()
		if tc.conn != nil {
			tc.conn.Close()
		}
	}
	
	func (tc *Connection) Run(mainctx context.Context) error {
		sendRecvCtx, sendRecvCancel := context.WithCancel(mainctx)
		tc.sendRecvStop = sendRecvCancel
		var rtnerr error
		var sendRecvWaitGroup sync.WaitGroup
		sendRecvWaitGroup.Add(2)
		go func() {
			defer sendRecvWaitGroup.Done()
			err := %[1]s_looptcp.RecvLoop(
				sendRecvCtx,
				tc.sendRecvStop,
				tc.conn,
				tc.readTimeoutSec,
				tc.handleRecvPacketFn)
			if err != nil {
				rtnerr = err
			}
		}()
		go func() {
			defer sendRecvWaitGroup.Done()
			err := %[1]s_looptcp.SendLoop(
				sendRecvCtx,
				tc.sendRecvStop,
				tc.conn,
				tc.writeTimeoutSec,
				tc.sendCh,
				tc.marshalBodyFn,
				tc.handleSentPacketFn)
			if err != nil {
				rtnerr = err
			}
		}()
		sendRecvWaitGroup.Wait()
		return rtnerr
	}
	
	func (tc *Connection) EnqueueSendPacket(pk %[1]s_packet.Packet) error {
		select {
		case tc.sendCh <- pk:
			return nil
		default:
			return fmt.Errorf("Send channel full %%v", tc)
		}
	}
	`, genArgs.Prefix)
	return &buf
}

func buildConnWasm(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"context"
		"fmt"
		"sync"
		"syscall/js"
	)
	`, genArgs.Prefix+postfix)
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
		sendBuffer := make([]byte, %[1]s_packet.HeaderLen, %[1]s_packet.MaxPacketLen)
		var err error
	loop:
		for {
			select {
			case <-sendRecvCtx.Done():
				break loop
			case pk := <-wsc.sendCh:
				sendBuffer, err := %[1]s_packet.Packet2Bytes(&pk, wsc.marshalBodyFn, sendBuffer[:%[1]s_packet.HeaderLen])
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
				header, body, lerr := %[1]s_packet.Bytes2HeaderBody(rdata)
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
		`, genArgs.Prefix)
	return &buf
}

func buildConnWSGorilla(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"context"
		"fmt"
		"net/url"
		"sync"
		"time"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	type Connection struct {
		wsConn       *websocket.Conn
		sendRecvStop func()
		sendCh       chan %[1]s_packet.Packet
	
		readTimeoutSec     time.Duration
		writeTimeoutSec    time.Duration
		marshalBodyFn      func(interface{}, []byte) ([]byte, byte, error)
		handleRecvPacketFn func(header %[1]s_packet.Header, body []byte) error
		handleSentPacketFn func(header %[1]s_packet.Header) error
	}
	
	func New(
		readTimeoutSec, writeTimeoutSec time.Duration,
		marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
		handleRecvPacketFn func(header %[1]s_packet.Header, body []byte) error,
		handleSentPacketFn func(header %[1]s_packet.Header) error,
	) *Connection {
		tc := &Connection{
			sendCh:             make(chan %[1]s_packet.Packet, 10),
			readTimeoutSec:     readTimeoutSec,
			writeTimeoutSec:    writeTimeoutSec,
			marshalBodyFn:      marshalBodyFn,
			handleRecvPacketFn: handleRecvPacketFn,
			handleSentPacketFn: handleSentPacketFn,
		}
	
		tc.sendRecvStop = func() {
			fmt.Printf("Too early sendRecvStop call")
		}
		return tc
	}
	
	func (tc *Connection) ConnectTo(connAddr string) error {
		u := url.URL{Scheme: "ws", Host: connAddr, Path: "/ws"}
		var err error
		tc.wsConn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			return err
		}
		return nil
	}
	
	func (tc *Connection) Cleanup() {
		tc.sendRecvStop()
		if tc.wsConn != nil {
			tc.wsConn.Close()
		}
	}
	
	func (tc *Connection) Run(aictx context.Context) error {
		connCtx, ctxCancel := context.WithCancel(aictx)
		tc.sendRecvStop = ctxCancel
		var rtnerr error
		var sendRecvWaitGroup sync.WaitGroup
		sendRecvWaitGroup.Add(2)
		go func() {
			defer sendRecvWaitGroup.Done()
			err := %[1]s_loopwsgorilla.RecvLoop(
				connCtx,
				tc.sendRecvStop,
				tc.wsConn,
				tc.readTimeoutSec,
				tc.handleRecvPacketFn,
			)
			if err != nil {
				rtnerr = err
			}
		}()
		go func() {
			defer sendRecvWaitGroup.Done()
			err := %[1]s_loopwsgorilla.SendLoop(
				connCtx,
				tc.sendRecvStop,
				tc.wsConn,
				tc.writeTimeoutSec,
				tc.sendCh,
				tc.marshalBodyFn,
				tc.handleSentPacketFn,
			)
			if err != nil {
				rtnerr = err
			}
		}()
		sendRecvWaitGroup.Wait()
		return rtnerr
	}
	
	func (tc *Connection) EnqueueSendPacket(pk %[1]s_packet.Packet) error {
		select {
		case tc.sendCh <- pk:
			return nil
		default:
			return fmt.Errorf("Send channel full %%v", tc)
		}
	}
	`, genArgs.Prefix)
	return &buf
}

func buildLoopWSGorilla(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"context"
		"fmt"
		"net"
		"time"
		"github.com/gorilla/websocket"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	func SendControl(
		wsConn *websocket.Conn, mt int, PacketWriteTimeOut time.Duration) error {
	
		return wsConn.WriteControl(mt, []byte{}, time.Now().Add(PacketWriteTimeOut))
	}
	
	func WriteBytes(wsConn *websocket.Conn, sendBuffer []byte) error {
		return wsConn.WriteMessage(websocket.BinaryMessage, sendBuffer)
	}
	
	func SendLoop(sendRecvCtx context.Context, SendRecvStop func(), wsConn *websocket.Conn,
		timeout time.Duration,
		SendCh chan %[1]s_packet.Packet,
		marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
		handleSentPacketFn func(header %[1]s_packet.Header) error,
	) error {

		defer SendRecvStop()
		sendBuffer := make([]byte, %[1]s_packet.HeaderLen)
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
				sendBuffer, err := %[1]s_packet.Packet2Bytes(&pk, marshalBodyFn, sendBuffer[:%[1]s_packet.HeaderLen])
				if err != nil {
					break loop
				}
				if err = WriteBytes(wsConn, sendBuffer); err != nil {
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
		return %[1]s_packet.Bytes2HeaderBody(rdata)
	}
	`, genArgs.Prefix)
	return &buf
}

func buildLoopTCP(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"context"
		"net"
		"time"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
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
		SendCh chan %[1]s_packet.Packet,
		marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
		handleSentPacketFn func(header %[1]s_packet.Header) error,
	) error {

		defer SendRecvStop()
		sendBuffer := make([]byte, %[1]s_packet.HeaderLen, %[1]s_packet.MaxPacketLen)
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
				sendBuffer, err := %[1]s_packet.Packet2Bytes(&pk, marshalBodyFn, sendBuffer[:%[1]s_packet.HeaderLen])
				if err != nil {
					break loop
				}
				if err = WriteBytes(tcpConn, sendBuffer); err != nil {
					break loop
				}
				if err = handleSentPacketFn(pk.Header); err != nil {
					break loop
				}
			}
		}
		return err
	}
	
	func RecvLoop(sendRecvCtx context.Context, SendRecvStop func(), tcpConn *net.TCPConn,
		timeOut time.Duration,
		HandleRecvPacketFn func(header %[1]s_packet.Header, body []byte) error,
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
				header, rbody, err := %[1]s_packet.ReadHeaderBody(tcpConn)
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
	`, genArgs.Prefix)
	return &buf
}

func buildPID2RspFn(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"fmt"
		"sync"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	type HandleRspFn func(%[1]s_packet.Header, interface{}) error
	type PID2RspFn struct {
		mutex      sync.Mutex
		pid2recvfn map[uint32]HandleRspFn
		pid        uint32
	}
	func New() *PID2RspFn {
		rtn := &PID2RspFn{
			pid2recvfn: make(map[uint32]HandleRspFn),
		}
		return rtn
	}
	func (p2r *PID2RspFn) NewPID(fn HandleRspFn) uint32 {
		p2r.mutex.Lock()
		defer p2r.mutex.Unlock()
		p2r.pid++
		p2r.pid2recvfn[p2r.pid] = fn
		return p2r.pid
	}
	func (p2r *PID2RspFn) HandleRsp(header %[1]s_packet.Header, body interface{}) error {
		p2r.mutex.Lock()
		if recvfn, exist := p2r.pid2recvfn[header.ID]; exist {
			delete(p2r.pid2recvfn, header.ID)
			p2r.mutex.Unlock()
			return recvfn(header, body)
		}
		p2r.mutex.Unlock()
		return fmt.Errorf("pid not found")
	}
	`, genArgs.Prefix)
	return &buf
}

func buildStatNoti(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"fmt"
		"net/http"
		"sync"
		"text/template"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	func (ns *StatNotification) String() string {
		return fmt.Sprintf("StatNotification[%%v]", len(ns))
	}
	type StatNotification [%[1]s_idnoti.NotiID_Count]StatRow
	func New() *StatNotification {
		ns := new(StatNotification)
		for i := 0; i < %[1]s_idnoti.NotiID_Count; i++ {
			ns[i].Name = %[1]s_idnoti.NotiID(i).String()
		}
		return ns
	}
	func (ns *StatNotification) Add(hd %[1]s_packet.Header) {
		if int(hd.Cmd) >= %[1]s_idnoti.NotiID_Count {
			return
		}
		ns[hd.Cmd].add(hd)
	}
	func (ns *StatNotification) ToWeb(w http.ResponseWriter, r *http.Request) error {
		tplIndex, err := template.New("index").Parse(%[2]c
	<html><head><title>Notification packet stat Info</title></head><body>
	<table border=1 style="border-collapse:collapse;">%[2]c +
			HTML_tableheader +
			%[2]c{{range $i, $v := .}}%[2]c +
			HTML_row +
			%[2]c{{end}}%[2]c +
			HTML_tableheader +
			%[2]c</table><br/>
	</body></html>%[2]c)
		if err != nil {
			return err
		}
		if err := tplIndex.Execute(w, ns); err != nil {
			return err
		}
		return nil
	}
	const (
		HTML_tableheader = %[2]c<tr>
	<th>Name</th>
	<th>Count</th>
	<th>Total Byte</th>
	<th>Max Byte</th>
	<th>Avg Byte</th>
	</tr>%[2]c
		HTML_row = %[2]c<tr>
	<td>{{$v.Name}}</td>
	<td>{{$v.Count }}</td>
	<td>{{$v.TotalByte }}</td>
	<td>{{$v.MaxByte }}</td>
	<td>{{printf "%%10.3f" $v.Avg }}</td>
	</tr>
	%[2]c
	)
	type StatRow struct {
		mutex     sync.Mutex
		Name      string
		Count     int
		TotalByte int
		MaxByte   int
	}
	func (ps *StatRow) add(hd %[1]s_packet.Header) {
		ps.mutex.Lock()
		ps.Count++
		n := int(hd.BodyLen()) + %[1]s_packet.HeaderLen
		ps.TotalByte += n
		if n > ps.MaxByte {
			ps.MaxByte = n
		}
		ps.mutex.Unlock()
	}
		func (ps *StatRow) Avg() float64 {
		return float64(ps.TotalByte) / float64(ps.Count)
	}
	`, genArgs.Prefix, '`')
	return &buf
}

func buildStatCallAPI(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"fmt"
		"html/template"
		"net/http"
		"sync"
		"time"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	func (cps *StatCallAPI) String() string {
		return fmt.Sprintf("StatCallAPI[%%v]",
			len(cps))
	}
	type StatCallAPI [%[1]s_idcmd.CommandID_Count]StatRow
	func New() *StatCallAPI {
		cps := new(StatCallAPI)
		for i := 0; i < %[1]s_idcmd.CommandID_Count; i++ {
			cps[i].Name = %[1]s_idcmd.CommandID(i).String()
		}
		return cps
	}
	func (cps *StatCallAPI) BeforeSendReq(header %[1]s_packet.Header) (*statObj, error) {
		if int(header.Cmd) >= %[1]s_idcmd.CommandID_Count {
			return nil, fmt.Errorf("CommandID out of range %%v %%v",
				header, %[1]s_idcmd.CommandID_Count)
		}
		return cps[header.Cmd].open(), nil
	}
	func (cps *StatCallAPI) AfterSendReq(header %[1]s_packet.Header) error {
		if int(header.Cmd) >= %[1]s_idcmd.CommandID_Count {
			return fmt.Errorf("CommandID out of range %%v %%v", header, %[1]s_idcmd.CommandID_Count)
		}
		n := int(header.BodyLen()) + %[1]s_packet.HeaderLen
		cps[header.Cmd].addTx(n)
		return nil
	}
	func (cps *StatCallAPI) AfterRecvRsp(header %[1]s_packet.Header) error {
		if int(header.Cmd) >= %[1]s_idcmd.CommandID_Count {
			return fmt.Errorf("CommandID out of range %%v %%v", header, %[1]s_idcmd.CommandID_Count)
		}
		n := int(header.BodyLen()) + %[1]s_packet.HeaderLen
		cps[header.Cmd].addRx(n)
		return nil
	}
	func (ws *StatCallAPI) ToWeb(w http.ResponseWriter, r *http.Request) error {
		tplIndex, err := template.New("index").Parse(%[2]c
	<html><head><title>Call API Stat Info</title></head><body>
	<table border=1 style="border-collapse:collapse;">%[2]c +
			HTML_tableheader +
			%[2]c{{range $i, $v := .}}%[2]c +
			HTML_row +
			%[2]c{{end}}%[2]c +
			HTML_tableheader +
			%[2]c</table><br/>
	</body></html>%[2]c)
		if err != nil {
			return err
		}
		if err := tplIndex.Execute(w, ws); err != nil {
			return err
		}
		return nil
	}
	////////////////////////////////////////////////////////////////////////////////
	type statObj struct {
		StartTime time.Time
		StatRef   *StatRow
	}
	func (so *statObj) CallServerEnd(success bool) {
		so.StatRef.close(success, so.StartTime)
	}
	////////////////////////////////////////////////////////////////////////////////
	type PacketID2StatObj struct {
		mutex sync.RWMutex
		stats map[uint32]*statObj
	}
	func NewPacketID2StatObj() *PacketID2StatObj {
		return &PacketID2StatObj{
			stats: make(map[uint32]*statObj),
		}
	}
	func (som *PacketID2StatObj) Add(pkid uint32, so *statObj) error {
		som.mutex.Lock()
		defer som.mutex.Unlock()
		if _, exist := som.stats[pkid]; exist {
			return fmt.Errorf("pkid exist %%v", pkid)
		}
		som.stats[pkid] = so
		return nil
	}
	func (som *PacketID2StatObj) Del(pkid uint32) *statObj {
		som.mutex.Lock()
		defer som.mutex.Unlock()
		so := som.stats[pkid]
		delete(som.stats, pkid)
		return so
	}
	func (som *PacketID2StatObj) Get(pkid uint32) *statObj {
		som.mutex.RLock()
		defer som.mutex.RUnlock()
		return som.stats[pkid]
	}
	////////////////////////////////////////////////////////////////////////////////
	const (
		HTML_tableheader = %[2]c<tr>
	<th>Name</th>
	<th>Start</th>
	<th>End</th>
	<th>Success</th>
	<th>Running</th>
	<th>Fail</th>
	<th>Avg ms</th>
	<th>TxAvg Byte</th>
	<th>RxAvg Byte</th>
	</tr>%[2]c
		HTML_row = %[2]c<tr>
	<td>{{$v.Name}}</td>
	<td>{{$v.StartCount}}</td>
	<td>{{$v.EndCount}}</td>
	<td>{{$v.SuccessCount}}</td>
	<td>{{$v.RunCount}}</td>
	<td>{{$v.FailCount}}</td>
	<td>{{printf "%%13.6f" $v.Avgms }}</td>
	<td>{{printf "%%10.3f" $v.AvgTx }}</td>
	<td>{{printf "%%10.3f" $v.AvgRx }}</td>
	</tr>
	%[2]c
	)
	type StatRow struct {
		mutex sync.Mutex
		Name  string
		TxCount int
		TxByte  int
		RxCount int
		RxByte  int
		StartCount   int
		EndCount     int
		SuccessCount int
		Sum          time.Duration
	}
	func (sr *StatRow) open() *statObj {
		sr.mutex.Lock()
		defer sr.mutex.Unlock()
		sr.StartCount++
		return &statObj{
			StartTime: time.Now(),
			StatRef:   sr,
		}
	}
	func (sr *StatRow) close(success bool, startTime time.Time) {
		sr.mutex.Lock()
		defer sr.mutex.Unlock()
		sr.EndCount++
		if success {
			sr.SuccessCount++
			sr.Sum += time.Now().Sub(startTime)
		}
	}
	func (sr *StatRow) addTx(n int) {
		sr.mutex.Lock()
		defer sr.mutex.Unlock()
		sr.TxCount++
		sr.TxByte += n
	}
	func (sr *StatRow) addRx(n int) {
		sr.mutex.Lock()
		defer sr.mutex.Unlock()
		sr.RxCount++
		sr.RxByte += n
	}
	func (sr *StatRow) RunCount() int {
		return sr.StartCount - sr.EndCount
	}
	func (sr *StatRow) FailCount() int {
		return sr.EndCount - sr.SuccessCount
	}
	func (sr *StatRow) Avgms() float64 {
		if sr.EndCount != 0 {
			return float64(sr.Sum) / float64(sr.EndCount*1000000)
		}
		return 0.0
	}
	func (sr *StatRow) AvgRx() float64 {
		if sr.EndCount != 0 {
			return float64(sr.RxByte) / float64(sr.RxCount)
		}
		return 0.0
	}
	func (sr *StatRow) AvgTx() float64 {
		if sr.EndCount != 0 {
			return float64(sr.TxByte) / float64(sr.TxCount)
		}
		return 0.0
	}
	`, genArgs.Prefix, '`')
	return &buf
}

func buildStatServeAPI(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"fmt"
		"net/http"
		"sync"
		"text/template"
		"time"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	func (ps *StatServeAPI) String() string {
		return fmt.Sprintf("StatServeAPI[%%v]", len(ps))
	}
	type StatServeAPI [%[1]s_idcmd.CommandID_Count]StatRow
	func New() *StatServeAPI {
		ps := new(StatServeAPI)
		for i := 0; i < %[1]s_idcmd.CommandID_Count; i++ {
			ps[i].Name = %[1]s_idcmd.CommandID(i).String()
		}
		return ps
	}
	func (ps *StatServeAPI) AfterRecvReqHeader(header %[1]s_packet.Header) (*StatObj, error) {
		if int(header.Cmd) >= %[1]s_idcmd.CommandID_Count {
			return nil, fmt.Errorf("CommandID out of range %%v %%v", header, %[1]s_idcmd.CommandID_Count)
		}
		return ps[header.Cmd].open(header), nil
	}
	func (ws *StatServeAPI) ToWeb(w http.ResponseWriter, r *http.Request) error {
		tplIndex, err := template.New("index").Parse(%[2]c
	<html><head><title>Serve API stat Info</title></head><body>
	<table border=1 style="border-collapse:collapse;">%[2]c +
			HTML_tableheader +
			%[2]c{{range $i, $v := .}}%[2]c +
			HTML_row +
			%[2]c{{end}}%[2]c +
			HTML_tableheader +
			%[2]c</table><br/>
	</body></html>%[2]c)
		if err != nil {
			return err
		}
		if err := tplIndex.Execute(w, ws); err != nil {
			return err
		}
		return nil
	}
	////////////////////////////////////////////////////////////////////////////////
	type StatObj struct {
		RecvTime    time.Time
		APICallTime time.Time
		StatRef     *StatRow
	}
	func (sm *StatObj) BeforeAPICall() {
		sm.APICallTime = time.Now().UTC()
		sm.StatRef.afterAuth()
	}
	func (sm *StatObj) AfterAPICall() {
		sm.StatRef.apiEnd(time.Now().UTC().Sub(sm.APICallTime))
	}
	func (sm *StatObj) AfterSendRsp(hd %[1]s_packet.Header) {
		sm.StatRef.afterSend(time.Now().UTC().Sub(sm.RecvTime), hd)
	}
	////////////////////////////////////////////////////////////////////////////////
	type PacketID2StatObj struct {
		mutex sync.RWMutex
		stats map[uint32]*StatObj
	}
	func NewPacketID2StatObj() *PacketID2StatObj {
		return &PacketID2StatObj{
			stats: make(map[uint32]*StatObj),
		}
	}
	func (som *PacketID2StatObj) Add(pkid uint32, so *StatObj) error {
		som.mutex.Lock()
		defer som.mutex.Unlock()
		if _, exist := som.stats[pkid]; exist {
			return fmt.Errorf("pkid exist %%v", pkid)
		}
		som.stats[pkid] = so
		return nil
	}
	func (som *PacketID2StatObj) Del(pkid uint32) *StatObj {
		som.mutex.Lock()
		defer som.mutex.Unlock()
		so := som.stats[pkid]
		delete(som.stats, pkid)
		return so
	}
	func (som *PacketID2StatObj) Get(pkid uint32) *StatObj {
		som.mutex.RLock()
		defer som.mutex.RUnlock()
		return som.stats[pkid]
	}
	////////////////////////////////////////////////////////////////////////////////
	const (
		HTML_tableheader = %[2]c<tr>
	<th>Name</th>
	<th>Recv Count</th>
	<th>Auth Count</th>
	<th>APIEnd Count</th>
	<th>Send Count</th>
	<th>Run Count</th>
	<th>Fail Count</th>
	<th>RecvSend Avg ms</th>
	<th>API Avg ms</th>
	<th>Rx Avg Byte</th>
	<th>Rx Max Byte</th>
	<th>Tx Avg Byte</th>
	<th>Tx Max Byte</th>
	</tr>%[2]c
		HTML_row = %[2]c<tr>
	<td>{{$v.Name}}</td>
	<td>{{$v.RecvCount}}</td>
	<td>{{$v.AuthCount}}</td>
	<td>{{$v.APIEndCount}}</td>
	<td>{{$v.SendCount}}</td>
	<td>{{$v.RunCount}}</td>
	<td>{{$v.FailCount}}</td>
	<td>{{printf "%%13.6f" $v.RSAvgms }}</td>
	<td>{{printf "%%13.6f" $v.APIAvgms }}</td>
	<td>{{printf "%%10.3f" $v.AvgRxByte }}</td>
	<td>{{$v.MaxRecvBytes }}</td>
	<td>{{printf "%%10.3f" $v.AvgTxByte }}</td>
	<td>{{$v.MaxSendBytes }}</td>
	</tr>
	%[2]c
	)
	type StatRow struct {
		mutex sync.Mutex
		Name  string
		RecvCount    int
		MaxRecvBytes int
		RecvBytes    int
		SendCount    int
		MaxSendBytes int
		SendBytes    int
		RecvSendDurSum time.Duration
		AuthCount      int
		APIEndCount    int
		APIDurSum      time.Duration
	}
	func (sr *StatRow) open(hd %[1]s_packet.Header) *StatObj {
		sr.mutex.Lock()
		defer sr.mutex.Unlock()
		sr.RecvCount++
		rxbyte := int(hd.BodyLen()) + %[1]s_packet.HeaderLen
		sr.RecvBytes += rxbyte
		if sr.MaxRecvBytes < rxbyte {
			sr.MaxRecvBytes = rxbyte
		}
		rtn := &StatObj{
			RecvTime: time.Now().UTC(),
			StatRef:  sr,
		}
		return rtn
	}
		func (sr *StatRow) afterAuth() {
		sr.mutex.Lock()
		defer sr.mutex.Unlock()
		sr.AuthCount++
	}
	func (sr *StatRow) apiEnd(diffDur time.Duration) {
		sr.mutex.Lock()
		defer sr.mutex.Unlock()
		sr.APIEndCount++
		sr.APIDurSum += diffDur
	}
	func (sr *StatRow) afterSend(diffDur time.Duration, hd %[1]s_packet.Header) {
		sr.mutex.Lock()
		defer sr.mutex.Unlock()
		sr.SendCount++
		txbyte := int(hd.BodyLen()) + %[1]s_packet.HeaderLen
		sr.SendBytes += txbyte
		if sr.MaxSendBytes < txbyte {
			sr.MaxSendBytes = txbyte
		}
		sr.RecvSendDurSum += diffDur
	}
		////////////////////////////////////////////////////////////////////////////////
	func (sr *StatRow) RunCount() int {
		return sr.AuthCount - sr.APIEndCount
	}
	func (sr *StatRow) FailCount() int {
		return sr.APIEndCount - sr.SendCount
	}
	func (sr *StatRow) RSAvgms() float64 {
		if sr.SendCount == 0 {
			return 0
		}
		return float64(sr.RecvSendDurSum) / float64(sr.SendCount*1000000)
	}
	func (sr *StatRow) APIAvgms() float64 {
		if sr.APIEndCount == 0 {
			return 0
		}
		return float64(sr.APIDurSum) / float64(sr.APIEndCount*1000000)
	}
	func (sr *StatRow) AvgRxByte() float64 {
		if sr.RecvCount == 0 {
			return 0
		}
		return float64(sr.RecvBytes) / float64(sr.RecvCount)
	}
	func (sr *StatRow) AvgTxByte() float64 {
		if sr.SendCount == 0 {
			return 0
		}
		return float64(sr.SendBytes) / float64(sr.SendCount)
	}
	`, genArgs.Prefix, '`')
	return &buf
}

func buildStatAPIError(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"fmt"
		"html/template"
		"net/http"
		"sync"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	func (es *StatAPIError) String() string {
		return fmt.Sprintf(
			"StatAPIError[%%v %%v %%v]",
			len(es.Stat),
			len(es.ECList),
			len(es.CmdList),
		)
	}
	type StatAPIError struct {
		mutex   sync.RWMutex
		Stat    [][]int
		ECList  []string
		CmdList []string
	}
	func New() *StatAPIError {
		es := &StatAPIError{
			Stat: make([][]int, %[1]s_idcmd.CommandID_Count),
		}
		for i, _ := range es.Stat {
			es.Stat[i] = make([]int, %[1]s_error.ErrorCode_Count)
		}
		es.ECList = make([]string, %[1]s_error.ErrorCode_Count)
		for i, _ := range es.ECList {
			es.ECList[i] = fmt.Sprintf("%%s", %[1]s_error.ErrorCode(i).String())
		}
		es.CmdList = make([]string, %[1]s_idcmd.CommandID_Count)
		for i, _ := range es.CmdList {
			es.CmdList[i] = fmt.Sprintf("%%v", %[1]s_idcmd.CommandID(i))
		}
		return es
	}
	func (es *StatAPIError) Inc(cmd %[1]s_idcmd.CommandID, errorcode %[1]s_error.ErrorCode) {
		es.mutex.Lock()
		defer es.mutex.Unlock()
		es.Stat[cmd][errorcode]++
	}
	func (es *StatAPIError) ToWeb(w http.ResponseWriter, r *http.Request) error {
		tplIndex, err := template.New("index").Parse(%[2]c
	<html><head><title>API Error stat Info</title></head><body>
	<table border=1 style="border-collapse:collapse;">
	<tr>
		<td></td>
		{{range $ft, $v := .ECList}}
			<th>{{$v}}</th>
		{{end}}
	</tr>
	{{range $cmd, $w := .Stat}}
		<tr>
			<td>{{index $.CmdList $cmd}}</td>
			{{range $ft, $v := $w}}
				<td>{{$v}}</td>
			{{end}}
		</tr>
	{{end}}
	<tr>
		<td></td>
		{{range $ft, $v := .ECList}}
			<th>{{$v}}</th>
		{{end}}
	</tr>
	</table><br/>
	</body></html>%[2]c)
		if err != nil {
			return err
		}
		if err := tplIndex.Execute(w, es); err != nil {
			return err
		}
		return nil
	}
	`, genArgs.Prefix, '`')
	return &buf
}

func buildAuthorize(genArgs GenArgs, postfix string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s
	import (
		"bytes"
		"fmt"
	)
	`, genArgs.Prefix+postfix)
	fmt.Fprintf(&buf, `
	type AuthorizedCmds [%[1]s_idcmd.CommandID_Count]bool

	func (acidl *AuthorizedCmds) String() string {
		var buff bytes.Buffer
		fmt.Fprintf(&buff, "AuthorizedCmds[")
		for i, v := range acidl {
			if v {
				fmt.Fprintf(&buff, "%%v ", %[1]s_idcmd.CommandID(i))
			}
		}
		fmt.Fprintf(&buff, "]")
		return buff.String()
	}
	
	func NewAllSet() *AuthorizedCmds {
		rtn := new(AuthorizedCmds)
		for i := 0; i < %[1]s_idcmd.CommandID_Count; i++ {
			rtn[i] = true
		}
		return rtn
	}
	
	func NewByCmdIDList(cmdlist []%[1]s_idcmd.CommandID) *AuthorizedCmds {
		rtn := new(AuthorizedCmds)
		for _, id := range cmdlist {
			rtn[id] = true
		}
		return rtn
	}
	
	func (acidl *AuthorizedCmds) Union(src *AuthorizedCmds) *AuthorizedCmds {
		for cmdid, auth := range src {
			if auth {
				acidl[cmdid] = true
			}
		}
		return acidl
	}
	
	func (acidl *AuthorizedCmds) SubIntersection(src *AuthorizedCmds) *AuthorizedCmds {
		for cmdid, auth := range src {
			if auth {
				acidl[cmdid] = false
			}
		}
		return acidl
	}
	
	func (acidl *AuthorizedCmds) Duplicate() *AuthorizedCmds {
		rtn := *acidl
		return &rtn
	}
	
	func (acidl *AuthorizedCmds) CheckAuth(cmdid %[1]s_idcmd.CommandID) bool {
		return acidl[cmdid]
	}
	`, genArgs.Prefix)
	return &buf
}

func buildStatsCode(genArgs GenArgs, pkgname string, typename string, statstype string) *bytes.Buffer {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, genArgs.GenComment)
	fmt.Fprintf(&buf, `
	package %[1]s_stats
	import (
		"bytes"
		"fmt"
		"html/template"
		"net/http"
	)
	`, pkgname, typename)

	fmt.Fprintf(&buf, `
	type %[2]sStat [%[1]s.%[2]s_Count]%[4]s
	func (es %[2]sStat) String() string {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%[2]sStats[")
		for i, v := range es {
			fmt.Fprintf(&buf,
				"%%v:%%v ",
				%[1]s.%[2]s(i), v)
		}
		buf.WriteString("]")
		return buf.String()
	}
	func (es *%[2]sStat) Inc(e %[1]s.%[2]s) {
		es[e]+=1
	}
	func (es *%[2]sStat) Add(e %[1]s.%[2]s, v %[4]s) {
		es[e]+=v
	}
	func (es *%[2]sStat) SetIfGt(e %[1]s.%[2]s, v %[4]s) {
		if es[e] < v {
			es[e]=v
		}
	}
	func (es %[2]sStat) Get(e %[1]s.%[2]s) %[4]s {
		return es[e]
	}
	
	// Iter return true if iter stop, return false if iter all
	// fn return true to stop iter
	func (es %[2]sStat) Iter(fn func(i %[1]s.%[2]s, v %[4]s) bool) bool {
		for i, v := range es {
			if fn(%[1]s.%[2]s(i), v) {
				return true
			}
		}
		return false
	}

	// VectorAdd add element to element
	func (es %[2]sStat) VectorAdd(arg %[2]sStat) %[2]sStat {
		var rtn %[2]sStat
		for i, v := range es {
			rtn[i] = v + arg[i]
		}
		return rtn
	}
	
	// VectorSub sub element to element
	func (es %[2]sStat) VectorSub(arg %[2]sStat) %[2]sStat {
		var rtn %[2]sStat
		for i, v := range es {
			rtn[i] = v - arg[i]
		}
		return rtn
	}


	func (es %[2]sStat) ToWeb(w http.ResponseWriter, r *http.Request) error {
		tplIndex, err := template.New("index").Funcs(IndexFn).Parse(%[3]c
		<html>
		<head>
		<title>%[2]s statistics</title>
		</head>
		<body>
		<table border=1 style="border-collapse:collapse;">%[3]c +
			HTML_tableheader +
			%[3]c{{range $i, $v := .}}%[3]c +
			HTML_row +
			%[3]c{{end}}%[3]c +
			HTML_tableheader +
			%[3]c</table>
	
		<br/>
		</body>
		</html>
		%[3]c)
		if err != nil {
			return err
		}
		if err := tplIndex.Execute(w, es); err != nil {
			return err
		}
		return nil
	}
	
	func Index(i int) string {
		return %[1]s.%[2]s(i).String()
	}
	
	var IndexFn = template.FuncMap{
		"%[2]sIndex": Index,
	}
	
	const (
		HTML_tableheader = %[3]c<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>%[3]c
		HTML_row = %[3]c<tr>
		<td>{{%[2]sIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		%[3]c
	)
	`, pkgname, typename, '`', statstype)

	return &buf
}
