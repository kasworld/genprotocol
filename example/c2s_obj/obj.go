package c2s_obj

import "time"

type ReqInvalidCmd_data struct {
	Dummy uint8
}
type RspInvalidCmd_data struct {
	Dummy uint8
}

type ReqLogin_data struct {
	Dummy uint8
}
type RspLogin_data struct {
	Dummy uint8
}

type ReqHeartbeat_data struct {
	Now time.Time
}
type RspHeartbeat_data struct {
	Now time.Time
}

type ReqChat_data struct {
	Msg string
}
type RspChat_data struct {
	Msg string
}

type NotiBroadcast_data struct {
	Dummy uint8
}
