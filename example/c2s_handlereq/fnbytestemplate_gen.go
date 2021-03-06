// Code generated by "genprotocol.exe -ver=1.0 -prefix=c2s -basedir=example -statstype=int"

package c2s_handlereq

/* bytes base fn map api template , unmarshal in api
	var DemuxReq2BytesAPIFnMap = [...]func(
		me interface{}, hd c2s_packet.Header, rbody []byte) (
		c2s_packet.Header, interface{}, error){
	c2s_idcmd.InvalidCmd: bytesAPIFn_ReqInvalidCmd,// InvalidCmd not used
c2s_idcmd.Login: bytesAPIFn_ReqLogin,// Login make session with nickname and enter stage
c2s_idcmd.Heartbeat: bytesAPIFn_ReqHeartbeat,// Heartbeat prevent connection timeout
c2s_idcmd.Chat: bytesAPIFn_ReqChat,// Chat chat to stage
c2s_idcmd.Act: bytesAPIFn_ReqAct,// Act send user action

}   // DemuxReq2BytesAPIFnMap

	// InvalidCmd not used
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
			ErrorCode : c2s_error.None,
		}
		sendBody := &c2s_obj.RspInvalidCmd_data{
		}
		return sendHeader, sendBody, nil
	}

	// Login make session with nickname and enter stage
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
			ErrorCode : c2s_error.None,
		}
		sendBody := &c2s_obj.RspLogin_data{
		}
		return sendHeader, sendBody, nil
	}

	// Heartbeat prevent connection timeout
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
			ErrorCode : c2s_error.None,
		}
		sendBody := &c2s_obj.RspHeartbeat_data{
		}
		return sendHeader, sendBody, nil
	}

	// Chat chat to stage
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
			ErrorCode : c2s_error.None,
		}
		sendBody := &c2s_obj.RspChat_data{
		}
		return sendHeader, sendBody, nil
	}

	// Act send user action
	func bytesAPIFn_ReqAct(
		me interface{}, hd c2s_packet.Header, rbody []byte) (
		c2s_packet.Header, interface{}, error) {
		// robj, err := c2s_json.UnmarshalPacket(hd, rbody)
		// if err != nil {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
		// }
		// recvBody, ok := robj.(*c2s_obj.ReqAct_data)
		// if !ok {
		// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
		// }
		// _ = recvBody

		sendHeader := c2s_packet.Header{
			ErrorCode : c2s_error.None,
		}
		sendBody := &c2s_obj.RspAct_data{
		}
		return sendHeader, sendBody, nil
	}

*/
