genprotocol -ver=1.0 -prefix=c2s -basedir example


goimports -w example/c2s_version/version_gen.go
goimports -w example/c2s_idcmd/command_gen.go
goimports -w example/c2s_idnoti/noti_gen.go
goimports -w example/c2s_error/error_gen.go
goimports -w example/c2s_packet/packet_gen.go
goimports -w example/c2s_obj/objtemplate_gen.go
goimports -w example/c2s_msgp/serialize_gen.go
goimports -w example/c2s_json/serialize_gen.go
goimports -w example/c2s_client/recvrspobjfnmap_gen.go
goimports -w example/c2s_client/recvnotiobjfnmap_gen.go
goimports -w example/c2s_client/callsendrecv_gen.go
goimports -w example/c2s_server/demuxreq2api_gen.go
goimports -w example/c2s_server/apitemplate_gen.go
goimports -w example/c2s_wasmconn/wasmconn_gen.go
goimports -w example/c2s_loopwsgorilla/loopwsgorilla_gen.go
goimports -w example/c2s_looptcp/looptcp_gen.go
