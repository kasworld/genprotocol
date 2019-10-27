# genprotocol - protocol 기반 코드 생성

( goguelike 프로젝트를 하면서 만들어진 )
네트웍 등에서 사용하기 위한 protocol 을 작성하는데 필요한 반복적 이고 기계적인 코드들을 자동으로 생성 해주는 프로그램 입니다. 

example folder 에서 예제를 볼 수 있습니다. 

작성 이유는 하나의 프로젝트에서 여러종의 protocol 을 사용할 일이 생겨서 
반복된 코드를 작성하다 보니 만들게 되었습니다. 

goguelike 를 예로 들면 

game client - tower server 

tower server - ground server 

의 2종의 protocol이 필요하고 사실상 비슷한데 조금 다른 코드들이 서로 구별 되어 쓰입니다. 


## 사용 법 
위 와 같은 경우 

genprotocol -ver=1.2 -prefix=c2t -basedir=. 

genprotocol -ver=1.0 -prefix=t2g -basedir=. 

과 같은 형태로 실행합니다. 

### 인자 설명 

ver : protocol의 version ( protocol 마다 버전이 다를 수 있습니다. )

prefix : 각 protocol을 구별하기 위한 prefix 

basedir : protocol code가 생성될 기본 dir

## example/rundriver 디렉토리 

json을 serializer로 사용하는 
tcp server/client, 
websocket server/client 
예제 입니다. 

## 생성되는 go package (디렉토리)

prefix 는 genprotocol 에 prefix 인자로 준 값 

생성이 끝난 코드들은 import code가 제대로 되어 있지 않으니 
goimports 등으로 정리 해주어야 합니다. 

실행하면 goimports 를 해야할 파일 목록을 찍어 줍니다. 
	
	example.sh 를 실행한 결과 
	goimports -w example/c2s_version/version_gen.go
	goimports -w example/c2s_idcmd/command_gen.go
	goimports -w example/c2s_idnoti/noti_gen.go
	goimports -w example/c2s_error/error_gen.go
	goimports -w example/c2s_packet/packet_gen.go
	goimports -w example/c2s_obj/objtemplate_gen.go
	goimports -w example/c2s_msgp/serialize_gen.go
	goimports -w example/c2s_json/serialize_gen.go
	goimports -w example/c2s_handlersp/fnobjtemplate_gen.go
	goimports -w example/c2s_handlersp/fnbytestemplate_gen.go
	goimports -w example/c2s_handlenoti/fnobjtemplate_gen.go
	goimports -w example/c2s_handlenoti/fnbytestemplate_gen.go
	goimports -w example/c2s_callsendrecv/callsendrecv_gen.go
	goimports -w example/c2s_handlereq/fnobjtemplate_gen.go
	goimports -w example/c2s_handlereq/fnbytestemplate_gen.go
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


prefix_gendata : genprotocol에서 읽어 들이는 파일들 

	각 라인의 첫 단어가 enum 이고 space 로 분리된 뒷 부분은 생성된 코드의 comment가 된다. 
	# 으로 시작하는 라인은 무시(comment취급)
	command.data : request/response packet 용 command id 목록 
	noti.data    : notification packet 용 noti id 목록 
	error.data   : error code 목록 


prefix_obj : protocol struct 들 (packet body)

	생성하는 파일 
	objtemplate_gen.go : 예제 파일 - 참고해서 "_gen"이 없는 파일을 만들것 

prefix_handlereq :  request 를 받아서 api 로 전달 

	json_conn.go 참고, copy해서 필요한 서버 로직을 만듭니다. 
	생성하는 파일 
	demuxreq2api_gen.go : request 를 api로 연결 
	apitemplate_gen.go  : api code template, 참고해서 "_gen"이 없는 파일을 만들것 

prefix_version : protocol version 정보 

	생성하는 파일 
	version_gen.go

prefix_error : protocol error code list, errorcode type은 error를 구현하고 있다.

	생성하는 파일 
	error_gen.go 

prefix_idcmd : protocol command (request,response) list 

	생성하는 파일 
	command_gen.go

prefix_idnoti : protocol notification list 

	생성하는 파일 
	noti_gen.go

prefix_packet : protocol packet( header + body )

	생성하는 파일 
	packet_gen.go

prefix_msgp : messagepack marshal/unmarshal code (https://github.com/tinylib/msgp)

	생성하는 파일 
	serialize_gen.go
 
prefix_json : json marshal/unmarshal code 

	생성하는 파일 
	serialize_gen.go


prefix_loopwsgorilla : go server/client용 gorilla websocket Send/Recv loop ([gorilla](http://www.gorillatoolkit.org/pkg/websocket)) 

	생성하는 파일 
	loopwsgorilla_gen.go

prefix_looptcp : go server/client용 TCP Send/Recv loop

	생성하는 파일 
	looptcp_gen.go

prefix_connwasm : websocket wasm client 용 connection

	생성하는 파일 
	connwasm_gen.go

prefix_connwsgorilla : gorilla websocket  client 용 connection 

	생성하는 파일 
	connwsgorilla_gen.go

prefix_conntcp : tcp client 용 connection 

	생성하는 파일 
	conntcp_gen.go

prefix_handlersp : response 를 받아서 처리 

	생성하는 파일 
	recvrspobjfnmap_gen.go   : 받은 response 처리 

prefix_handlenoti : notification을 받아서 처리 

	생성하는 파일 
	recvnotiobjfnmap_gen.go  : 받은 notification 처리 

prefix_callsendrecv : blocked send/recv : response 를 받을때 까지 wait

	생성하는 파일 
	callsendrecv_gen.go      

prefix_pid2rspfn : callback 형태로 request/response를 처리하기위한 lib(client example참조)

	생성하는 파일 
	pid2rspfn_gen.go      

prefix_statnoti : notification protocol 통계 

	생성하는 파일 
	statnoti_gen.go	

prefix_statcallapi : client api call 통계 

	생성하는 파일 
	statcallapi_gen.go	

prefix_statserveapi : server api 처리 통계 

	생성하는 파일 
	statserveapi_gen.go	

prefix_statapierror : api 결과 error 통계 

	생성하는 파일 
	statapierror_gen.go	
