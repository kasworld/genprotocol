// Code generated by "genprotocol -ver=1.0 -prefix=c2s -basedir example -statstype=int"

package c2s_handlersp

/* obj base demux fn map template

var DemuxRsp2ObjFnMap = [...]func(me interface{}, hd c2s_packet.Header, body interface{}) error {
c2s_idcmd.InvalidCmd : objRecvRspFn_InvalidCmd,
c2s_idcmd.Login : objRecvRspFn_Login,
c2s_idcmd.Heartbeat : objRecvRspFn_Heartbeat,
c2s_idcmd.Chat : objRecvRspFn_Chat,

}

	func objRecvRspFn_InvalidCmd(me interface{}, hd c2s_packet.Header, body interface{}) error {
		robj , ok := body.(*c2s_obj.RspInvalidCmd_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

	func objRecvRspFn_Login(me interface{}, hd c2s_packet.Header, body interface{}) error {
		robj , ok := body.(*c2s_obj.RspLogin_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

	func objRecvRspFn_Heartbeat(me interface{}, hd c2s_packet.Header, body interface{}) error {
		robj , ok := body.(*c2s_obj.RspHeartbeat_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

	func objRecvRspFn_Chat(me interface{}, hd c2s_packet.Header, body interface{}) error {
		robj , ok := body.(*c2s_obj.RspChat_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

*/