genprotocol -ver=1.0 -prefix=c2s -basedir example

goimports -w example/c2s_version/version_gen.go
goimports -w example/c2s_idcmd/command_gen.go
goimports -w example/c2s_idnoti/noti_gen.go
goimports -w example/c2s_error/error_gen.go
goimports -w example/c2s_packet/packet_gen.go
goimports -w example/c2s_obj/objtemplate_gen.go
goimports -w example/c2s_msgp/serialize_gen.go
goimports -w example/c2s_json/serialize_gen.go
goimports -w example/c2s_handlersp/recvrspobjfnmap_gen.go
goimports -w example/c2s_handlenoti/recvnotiobjfnmap_gen.go
goimports -w example/c2s_callsendrecv/callsendrecv_gen.go
goimports -w example/c2s_handlereq/recvreqobjfnmap_gen.go
goimports -w example/c2s_handlereq/apitemplate_gen.go
goimports -w example/c2s_conntcp/conntcp_gen.go
goimports -w example/c2s_connwasm/connwasm_gen.go
goimports -w example/c2s_connwsgorilla/connwsgorilla_gen.go
goimports -w example/c2s_loopwsgorilla/loopwsgorilla_gen.go
goimports -w example/c2s_looptcp/looptcp_gen.go
goimports -w example/c2s_statnoti/statnoti_gen.go
goimports -w example/c2s_statcallapi/statcallapi_gen.go
goimports -w example/c2s_statserveapi/statserveapi_gen.go
goimports -w example/c2s_statapierror/statapierror_gen.go
