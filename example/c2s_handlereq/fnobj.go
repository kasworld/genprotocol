// Code generated by "genprotocol.exe -ver=1.0 -prefix=c2s -basedir=example -statstype=int"

package c2s_handlereq

import (
	"fmt"

	"github.com/kasworld/genprotocol/example/c2s_error"
	"github.com/kasworld/genprotocol/example/c2s_idcmd"
	"github.com/kasworld/genprotocol/example/c2s_obj"
	"github.com/kasworld/genprotocol/example/c2s_packet"
)

// obj base demux fn map
var DemuxReq2ObjAPIFnMap = [...]func(
	me interface{}, hd c2s_packet.Header, robj interface{}) (
	c2s_packet.Header, interface{}, error){
	c2s_idcmd.InvalidCmd: Req2ObjAPI_InvalidCmd, // InvalidCmd not used
	c2s_idcmd.Login:      Req2ObjAPI_Login,      // Login make session with nickname and enter stage
	c2s_idcmd.Heartbeat:  Req2ObjAPI_Heartbeat,  // Heartbeat prevent connection timeout
	c2s_idcmd.Chat:       Req2ObjAPI_Chat,       // Chat chat to stage
	c2s_idcmd.Act:        Req2ObjAPI_Act,        // Act send user action

} // DemuxReq2ObjAPIFnMap

// InvalidCmd not used
func Req2ObjAPI_InvalidCmd(
	me interface{}, hd c2s_packet.Header, robj interface{}) (
	c2s_packet.Header, interface{}, error) {
	req, ok := robj.(*c2s_obj.ReqInvalidCmd_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	rhd, rsp, err := objAPIFn_ReqInvalidCmd(me, hd, req)
	return rhd, rsp, err
}

// InvalidCmd not used
func objAPIFn_ReqInvalidCmd(
	me interface{}, hd c2s_packet.Header, robj *c2s_obj.ReqInvalidCmd_data) (
	c2s_packet.Header, *c2s_obj.RspInvalidCmd_data, error) {
	sendHeader := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	sendBody := &c2s_obj.RspInvalidCmd_data{}
	return sendHeader, sendBody, nil
}

// Login make session with nickname and enter stage
func Req2ObjAPI_Login(
	me interface{}, hd c2s_packet.Header, robj interface{}) (
	c2s_packet.Header, interface{}, error) {
	req, ok := robj.(*c2s_obj.ReqLogin_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	rhd, rsp, err := objAPIFn_ReqLogin(me, hd, req)
	return rhd, rsp, err
}

// Login make session with nickname and enter stage
func objAPIFn_ReqLogin(
	me interface{}, hd c2s_packet.Header, robj *c2s_obj.ReqLogin_data) (
	c2s_packet.Header, *c2s_obj.RspLogin_data, error) {
	sendHeader := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	sendBody := &c2s_obj.RspLogin_data{}
	return sendHeader, sendBody, nil
}

// Heartbeat prevent connection timeout
func Req2ObjAPI_Heartbeat(
	me interface{}, hd c2s_packet.Header, robj interface{}) (
	c2s_packet.Header, interface{}, error) {
	req, ok := robj.(*c2s_obj.ReqHeartbeat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	rhd, rsp, err := objAPIFn_ReqHeartbeat(me, hd, req)
	return rhd, rsp, err
}

// Heartbeat prevent connection timeout
func objAPIFn_ReqHeartbeat(
	me interface{}, hd c2s_packet.Header, robj *c2s_obj.ReqHeartbeat_data) (
	c2s_packet.Header, *c2s_obj.RspHeartbeat_data, error) {
	sendHeader := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	sendBody := &c2s_obj.RspHeartbeat_data{}
	return sendHeader, sendBody, nil
}

// Chat chat to stage
func Req2ObjAPI_Chat(
	me interface{}, hd c2s_packet.Header, robj interface{}) (
	c2s_packet.Header, interface{}, error) {
	req, ok := robj.(*c2s_obj.ReqChat_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	rhd, rsp, err := objAPIFn_ReqChat(me, hd, req)
	return rhd, rsp, err
}

// Chat chat to stage
func objAPIFn_ReqChat(
	me interface{}, hd c2s_packet.Header, robj *c2s_obj.ReqChat_data) (
	c2s_packet.Header, *c2s_obj.RspChat_data, error) {
	sendHeader := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	sendBody := &c2s_obj.RspChat_data{}
	return sendHeader, sendBody, nil
}

// Act send user action
func Req2ObjAPI_Act(
	me interface{}, hd c2s_packet.Header, robj interface{}) (
	c2s_packet.Header, interface{}, error) {
	req, ok := robj.(*c2s_obj.ReqAct_data)
	if !ok {
		return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	}
	rhd, rsp, err := objAPIFn_ReqAct(me, hd, req)
	return rhd, rsp, err
}

// Act send user action
func objAPIFn_ReqAct(
	me interface{}, hd c2s_packet.Header, robj *c2s_obj.ReqAct_data) (
	c2s_packet.Header, *c2s_obj.RspAct_data, error) {
	sendHeader := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	sendBody := &c2s_obj.RspAct_data{}
	return sendHeader, sendBody, nil
}