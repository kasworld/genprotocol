package c2s_server

import (
	"github.com/kasworld/genprotocol/example/c2s_error"
	"github.com/kasworld/genprotocol/example/c2s_obj"
	"github.com/kasworld/genprotocol/example/c2s_packet"
)

func apifn_ReqInvalidCmd(
	c2sc *ServeClientConn, hd c2s_packet.Header, robj *c2s_obj.ReqInvalidCmd_data) (
	c2s_packet.Header, *c2s_obj.RspInvalidCmd_data, error) {
	rhd := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	spacket := &c2s_obj.RspInvalidCmd_data{}
	return rhd, spacket, nil
}

func apifn_ReqLogin(
	c2sc *ServeClientConn, hd c2s_packet.Header, robj *c2s_obj.ReqLogin_data) (
	c2s_packet.Header, *c2s_obj.RspLogin_data, error) {
	rhd := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	spacket := &c2s_obj.RspLogin_data{}
	return rhd, spacket, nil
}

func apifn_ReqHeartbeat(
	c2sc *ServeClientConn, hd c2s_packet.Header, robj *c2s_obj.ReqHeartbeat_data) (
	c2s_packet.Header, *c2s_obj.RspHeartbeat_data, error) {
	rhd := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	spacket := &c2s_obj.RspHeartbeat_data{}
	return rhd, spacket, nil
}

func apifn_ReqChat(
	c2sc *ServeClientConn, hd c2s_packet.Header, robj *c2s_obj.ReqChat_data) (
	c2s_packet.Header, *c2s_obj.RspChat_data, error) {
	rhd := c2s_packet.Header{
		ErrorCode: c2s_error.None,
	}
	spacket := &c2s_obj.RspChat_data{}
	return rhd, spacket, nil
}
