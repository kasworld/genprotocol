// Code generated by "genprotocol -ver=1.0 -prefix=c2s -basedir example -statstype=int"

package c2s_statapierror

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/kasworld/genprotocol/example/c2s_error"
	"github.com/kasworld/genprotocol/example/c2s_idcmd"
)

func (es *StatAPIError) String() string {
	return fmt.Sprintf(
		"StatAPIError[%v %v %v]",
		len(es.Stat),
		len(es.ECList),
		len(es.CmdList),
	)
}

type StatAPIError struct {
	mutex   sync.RWMutex
	Stat    [][]int
	ECList  []string
	CmdList []string
}

func New() *StatAPIError {
	es := &StatAPIError{
		Stat: make([][]int, c2s_idcmd.CommandID_Count),
	}
	for i, _ := range es.Stat {
		es.Stat[i] = make([]int, c2s_error.ErrorCode_Count)
	}
	es.ECList = make([]string, c2s_error.ErrorCode_Count)
	for i, _ := range es.ECList {
		es.ECList[i] = fmt.Sprintf("%s", c2s_error.ErrorCode(i).String())
	}
	es.CmdList = make([]string, c2s_idcmd.CommandID_Count)
	for i, _ := range es.CmdList {
		es.CmdList[i] = fmt.Sprintf("%v", c2s_idcmd.CommandID(i))
	}
	return es
}
func (es *StatAPIError) Inc(cmd c2s_idcmd.CommandID, errorcode c2s_error.ErrorCode) {
	es.mutex.Lock()
	defer es.mutex.Unlock()
	es.Stat[cmd][errorcode]++
}
func (es *StatAPIError) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Parse(`
	<html><head><title>API Error stat Info</title></head><body>
	<table border=1 style="border-collapse:collapse;">
	<tr>
		<td></td>
		{{range $ft, $v := .ECList}}
			<th>{{$v}}</th>
		{{end}}
	</tr>
	{{range $cmd, $w := .Stat}}
		<tr>
			<td>{{index $.CmdList $cmd}}</td>
			{{range $ft, $v := $w}}
				<td>{{$v}}</td>
			{{end}}
		</tr>
	{{end}}
	<tr>
		<td></td>
		{{range $ft, $v := .ECList}}
			<th>{{$v}}</th>
		{{end}}
	</tr>
	</table><br/>
	</body></html>`)
	if err != nil {
		return err
	}
	if err := tplIndex.Execute(w, es); err != nil {
		return err
	}
	return nil
}