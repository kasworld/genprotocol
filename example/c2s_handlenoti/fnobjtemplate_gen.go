// Code generated by "genprotocol.exe -ver=1.0 -prefix=c2s -basedir=example -statstype=int"

package c2s_handlenoti

/* obj base demux fn map template

var DemuxNoti2ObjFnMap = [...]func(me interface{}, hd c2s_packet.Header, body interface{}) error {
c2s_idnoti.Broadcast : objRecvNotiFn_Broadcast, // Broadcast send message to all connection

}

	// Broadcast send message to all connection
	func objRecvNotiFn_Broadcast(me interface{}, hd c2s_packet.Header, body interface{}) error {
		robj , ok := body.(*c2s_obj.NotiBroadcast_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", body )
		}
		return fmt.Errorf("Not implemented %v", robj)
	}

*/
