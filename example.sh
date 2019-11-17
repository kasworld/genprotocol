go run genprotocol.go -ver=1.0 -prefix=c2s -basedir example -statstype=int

goimports -w example/c2s_version/version_gen.go
goimports -w example/c2s_idcmd/command_gen.go
goimports -w example/c2s_idnoti/noti_gen.go
goimports -w example/c2s_error/error_gen.go
goimports -w example/c2s_const/const_gen.go
goimports -w example/c2s_packet/packet_gen.go
goimports -w example/c2s_obj/objtemplate_gen.go
goimports -w example/c2s_msgp/serialize_gen.go
goimports -w example/c2s_json/serialize_gen.go
goimports -w example/c2s_handlersp/fnobjtemplate_gen.go
goimports -w example/c2s_handlersp/fnbytestemplate_gen.go
goimports -w example/c2s_handlenoti/fnobjtemplate_gen.go
goimports -w example/c2s_handlenoti/fnbytestemplate_gen.go
goimports -w example/c2s_handlereq/fnobjtemplate_gen.go
goimports -w example/c2s_handlereq/fnbytestemplate_gen.go
goimports -w example/c2s_serveconnbyte/serveconnbyte_gen.go
goimports -w example/c2s_conntcp/conntcp_gen.go
goimports -w example/c2s_connwasm/connwasm_gen.go
goimports -w example/c2s_connwsgorilla/connwsgorilla_gen.go
goimports -w example/c2s_loopwsgorilla/loopwsgorilla_gen.go
goimports -w example/c2s_looptcp/looptcp_gen.go
goimports -w example/c2s_pid2rspfn/pid2rspfn_gen.go
goimports -w example/c2s_statnoti/statnoti_gen.go
goimports -w example/c2s_statcallapi/statcallapi_gen.go
goimports -w example/c2s_statserveapi/statserveapi_gen.go
goimports -w example/c2s_statapierror/statapierror_gen.go
goimports -w example/c2s_authorize/authorize_gen.go

goimports -w example/c2s_error_stats/c2s_error_stats_gen.go
goimports -w example/c2s_idcmd_stats/c2s_idcmd_stats_gen.go
goimports -w example/c2s_idnoti_stats/c2s_idnoti_stats_gen.go


cd example/rundriver 

echo "build wasm client"
GOOS=js GOARCH=wasm go build -o www/wasmclient.wasm wasmclient.go

# echo "build go client"
# go build -o bin/goclient goclient.go

