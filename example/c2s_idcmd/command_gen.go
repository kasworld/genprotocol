// Code generated by "genprotocol.exe -ver=1.0 -prefix=c2s -basedir=example -statstype=int"

package c2s_idcmd

import "fmt"

type CommandID uint16 // use in packet header, DO NOT CHANGE
const (
	InvalidCmd CommandID = iota // not used
	Login                       // make session with nickname and enter stage
	Heartbeat                   // prevent connection timeout
	Chat                        // chat to stage
	Act                         // send user action

	CommandID_Count int = iota
)

var _CommandID2string = [CommandID_Count][2]string{
	InvalidCmd: {"InvalidCmd", "not used"},
	Login:      {"Login", "make session with nickname and enter stage"},
	Heartbeat:  {"Heartbeat", "prevent connection timeout"},
	Chat:       {"Chat", "chat to stage"},
	Act:        {"Act", "send user action"},
}

func (e CommandID) String() string {
	if e >= 0 && e < CommandID(CommandID_Count) {
		return _CommandID2string[e][0]
	}
	return fmt.Sprintf("CommandID%d", uint16(e))
}

func (e CommandID) CommentString() string {
	if e >= 0 && e < CommandID(CommandID_Count) {
		return _CommandID2string[e][1]
	}
	return ""
}

var _string2CommandID = map[string]CommandID{
	"InvalidCmd": InvalidCmd,
	"Login":      Login,
	"Heartbeat":  Heartbeat,
	"Chat":       Chat,
	"Act":        Act,
}

func String2CommandID(s string) (CommandID, bool) {
	v, b := _string2CommandID[s]
	return v, b
}
