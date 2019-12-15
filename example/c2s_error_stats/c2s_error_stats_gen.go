// Code generated by "genprotocol -ver=1.0 -prefix=c2s -basedir example -statstype=int"

package c2s_error_stats

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/genprotocol/example/c2s_error"
)

type ErrorCodeStat [c2s_error.ErrorCode_Count]int

func (es *ErrorCodeStat) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "ErrorCodeStats[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			c2s_error.ErrorCode(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *ErrorCodeStat) Inc(e c2s_error.ErrorCode) {
	es[e] += 1
}
func (es *ErrorCodeStat) Add(e c2s_error.ErrorCode, v int) {
	es[e] += v
}
func (es *ErrorCodeStat) SetIfGt(e c2s_error.ErrorCode, v int) {
	if es[e] < v {
		es[e] = v
	}
}
func (es *ErrorCodeStat) Get(e c2s_error.ErrorCode) int {
	return es[e]
}

func (es *ErrorCodeStat) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>ErrorCode statistics</title>
		</head>
		<body>
		<table border=1 style="border-collapse:collapse;">` +
		HTML_tableheader +
		`{{range $i, $v := .}}` +
		HTML_row +
		`{{end}}` +
		HTML_tableheader +
		`</table>
	
		<br/>
		</body>
		</html>
		`)
	if err != nil {
		return err
	}
	if err := tplIndex.Execute(w, es); err != nil {
		return err
	}
	return nil
}

func Index(i int) string {
	return c2s_error.ErrorCode(i).String()
}

var IndexFn = template.FuncMap{
	"ErrorCodeIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{ErrorCodeIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)
