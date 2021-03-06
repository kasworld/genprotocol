// Code generated by "genprotocol.exe -ver=1.0 -prefix=c2s -basedir=example -statstype=int"

package c2s_pid2rspfn

import (
	"fmt"
	"sync"

	"github.com/kasworld/genprotocol/example/c2s_packet"
)

type HandleRspFn func(c2s_packet.Header, interface{}) error
type PID2RspFn struct {
	mutex      sync.Mutex
	pid2recvfn map[uint32]HandleRspFn
	pid        uint32
}

func New() *PID2RspFn {
	rtn := &PID2RspFn{
		pid2recvfn: make(map[uint32]HandleRspFn),
	}
	return rtn
}
func (p2r *PID2RspFn) NewPID(fn HandleRspFn) uint32 {
	p2r.mutex.Lock()
	defer p2r.mutex.Unlock()
	p2r.pid++
	p2r.pid2recvfn[p2r.pid] = fn
	return p2r.pid
}
func (p2r *PID2RspFn) HandleRsp(header c2s_packet.Header, body interface{}) error {
	p2r.mutex.Lock()
	if recvfn, exist := p2r.pid2recvfn[header.ID]; exist {
		delete(p2r.pid2recvfn, header.ID)
		p2r.mutex.Unlock()
		return recvfn(header, body)
	}
	p2r.mutex.Unlock()
	return fmt.Errorf("pid not found")
}
