package c2s_handlereq

import (
	"github.com/kasworld/genprotocol/example/c2s_error"
	"github.com/kasworld/genprotocol/example/c2s_idcmd"
	"github.com/kasworld/genprotocol/example/c2s_obj"
	"github.com/kasworld/genprotocol/example/c2s_packet"
)

var DemuxReq2BytesAPIFnMap = [...]func(
	me interface{}, hd c2s_packet.Header, rbody []byte) (
	c2s_packet.Header, interface{}, error){
	c2s_idcmd.InvalidCmd: bytesAPIFn_ReqInvalidCmd,
	c2s_idcmd.Login:      bytesAPIFn_ReqLogin,
	c2s_idcmd.Heartbeat:  bytesAPIFn_ReqHeartbeat,
	c2s_idcmd.Chat:       bytesAPIFn_ReqChat,
} // DemuxReq2BytesAPIFnMap

func bytesAPIFn_ReqInvalidCmd(
	me interface{}, hd c2s_packet.Header, rbody []byte) (
	c2s_packet.Header, interface{}, error) {
	// robj, err := c2s_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*c2s_obj.ReqInvalidCmd_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	sendBody := &c2s_obj.RspInvalidCmd_data{}
	return sendHeader, sendBody, nil
}

func bytesAPIFn_ReqLogin(
	me interface{}, hd c2s_packet.Header, rbody []byte) (
	c2s_packet.Header, interface{}, error) {
	// robj, err := c2s_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*c2s_obj.ReqLogin_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	sendBody := &c2s_obj.RspLogin_data{}
	return sendHeader, sendBody, nil
}

func bytesAPIFn_ReqHeartbeat(
	me interface{}, hd c2s_packet.Header, rbody []byte) (
	c2s_packet.Header, interface{}, error) {
	// robj, err := c2s_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*c2s_obj.ReqHeartbeat_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	sendBody := &c2s_obj.RspHeartbeat_data{}
	return sendHeader, sendBody, nil
}

func bytesAPIFn_ReqChat(
	me interface{}, hd c2s_packet.Header, rbody []byte) (
	c2s_packet.Header, interface{}, error) {
	// robj, err := c2s_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*c2s_obj.ReqChat_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	sendBody := &c2s_obj.RspChat_data{}
	return sendHeader, sendBody, nil
}
