package c2s_obj

import "time"

// ReqInvalidCmd_data not used
type ReqInvalidCmd_data struct {
	Dummy uint8
}
type RspInvalidCmd_data struct {
	Dummy uint8
}

// ReqLogin_data make session with nickname and enter stage
type ReqLogin_data struct {
	Dummy uint8
}
type RspLogin_data struct {
	Dummy uint8
}

// ReqHeartbeat_data prevent connection timeout
type ReqHeartbeat_data struct {
	Now time.Time
}
type RspHeartbeat_data struct {
	Now time.Time
}

// ReqChat_data chat to stage
type ReqChat_data struct {
	Msg string
}
type RspChat_data struct {
	Msg string
}

// ReqAct_data send user action
type ReqAct_data struct {
	Dummy uint8
}
type RspAct_data struct {
	Dummy uint8
}

// NotiBroadcast_data send message to all connection
type NotiBroadcast_data struct {
	Dummy uint8
}
