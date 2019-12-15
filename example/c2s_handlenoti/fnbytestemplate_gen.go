// Code generated by "genprotocol -ver=1.0 -prefix=c2s -basedir example -statstype=int"

package c2s_handlenoti

/* bytes base demux fn map template

var DemuxNoti2ByteFnMap = [...]func(me interface{}, hd c2s_packet.Header, rbody []byte) error {
c2s_idnoti.Broadcast : bytesRecvNotiFn_Broadcast,

}

	func bytesRecvNotiFn_Broadcast(me interface{}, hd c2s_packet.Header, rbody []byte) error {
		robj, err := c2s_json.UnmarshalPacket(hd, rbody)
		if err != nil {
			return fmt.Errorf("Packet type miss match %v", rbody)
		}
		recved , ok := robj.(*c2s_obj.NotiBroadcast_data)
		if !ok {
			return fmt.Errorf("packet mismatch %v", robj )
		}
		return fmt.Errorf("Not implemented %v", recved)
	}

*/
