// Code generated by "genprotocol -ver=1.0 -prefix=c2s -basedir example -statstype=int"

package c2s_idcmd

import "fmt"

type CommandID uint16 // use in packet header, DO NOT CHANGE
const (
	InvalidCmd CommandID = iota //
	Login                       //
	Heartbeat                   //
	Chat                        //

	CommandID_Count int = iota
)

var _CommandID2string = [CommandID_Count]string{
	InvalidCmd: "InvalidCmd",
	Login:      "Login",
	Heartbeat:  "Heartbeat",
	Chat:       "Chat",
}

func (e CommandID) String() string {
	if e >= 0 && e < CommandID(CommandID_Count) {
		return _CommandID2string[e]
	}
	return fmt.Sprintf("CommandID%d", uint16(e))
}

var _string2CommandID = map[string]CommandID{
	"InvalidCmd": InvalidCmd,
	"Login":      Login,
	"Heartbeat":  Heartbeat,
	"Chat":       Chat,
}

func String2CommandID(s string) (CommandID, bool) {
	v, b := _string2CommandID[s]
	return v, b
}
