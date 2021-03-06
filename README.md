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

자동 생성된 파일은 모두 _gen.go 로 끝나는 이름을 가집니다.
 

### 인자 설명 

ver : protocol의 version ( protocol 마다 버전이 다를 수 있습니다. )

prefix : 각 protocol을 구별하기 위한 prefix 

basedir : protocol code가 생성될 기본 폴더, 여기서 prefix_*.enum 을 읽습니다. 

## example/rundriver 디렉토리 

json을 serializer로 사용하는 
tcp server/client, 
websocket server/client 
예제 입니다. 

## genprotocol에서 읽어 들이는 파일들 

	인자로준 basedir 폴더내에서 찾습니다. 
	
	각 라인의 첫 단어가 enum 이고 space 로 분리된 뒷 부분은 생성된 코드의 comment가 된다. 
	# 으로 시작하는 라인은 무시(comment취급)
	prefix_command.enum : request/response packet 용 command id 목록 
	prefix_noti.enum    : notification packet 용 noti id 목록 
	prefix_error.enum   : error code 목록 


## 생성되는 go package (디렉토리)

prefix 는 genprotocol 에 prefix 인자로 준 값 

생성이 끝난 코드들은 import code가 제대로 되어 있지 않으니 
goimports 등으로 정리 해주어야 합니다. 

-verbose 인자를 주고 실행하면 goimports 를 해야할 파일 목록을 찍어 줍니다. 
	
	example.sh 를 실행한 결과 
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
	goimports -w example/c2s_connbytemanager/connbytemanager_gen.go
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


	인자로 -statstype=elementtype 이 주어진 경우 추가로 생성
	goimports -w example/c2s_error_stats/c2s_error_stats_gen.go
	goimports -w example/c2s_idcmd_stats/c2s_idcmd_stats_gen.go
	goimports -w example/c2s_idnoti_stats/c2s_idnoti_stats_gen.go



prefix_obj : protocol struct 들 (packet body)

	생성하는 파일 
	objtemplate_gen.go : 예제 파일 - 참고해서 "_gen"이 없는 파일을 만들것 

prefix_handlereq :   commandid -> api function map, 과 예제 api function들 

	생성하는 파일 
	fnbytestemplate_gen : packetbody 가 []byte 인 형태로 api로 demux, unmarshal을 api 쪽에서 해야 함 
	fntemplate_gen  : packetbody 가 interface{} 인 형태로 api로 demux, unmarshal 후에 demux map을 호출 해야 함. 

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

prefix_const : 각종 상수들 , 설정의 변경이 가능하도록 모아둠. 

	생성하는 파일 
	const_gen.go : 복사해서 _gen 이 없는 파일을 만들고 comment를 풀어 사용할것.

prefix_packet : protocol packet( header + body )

	생성하는 파일 
	packet_gen.go

prefix_msgp : messagepack marshal/unmarshal code (https://github.com/tinylib/msgp)

	생성하는 파일 
	serialize_gen.go
 
prefix_json : json marshal/unmarshal code 

	생성하는 파일 
	serialize_gen.go

prefix_gob : gob marshal/unmarshal code 

	생성하는 파일 
	serialize_gen.go


prefix_loopwsgorilla : go server/client용 gorilla websocket Send/Recv loop ([gorilla](http://www.gorillatoolkit.org/pkg/websocket)) 

	생성하는 파일 
	loopwsgorilla_gen.go

prefix_looptcp : go server/client용 TCP Send/Recv loop

	생성하는 파일 
	looptcp_gen.go

prefix_serveconnbyte : server 용 connection api 처리 (tcp, websocket) packet body []byte 형태 

	생성하는 파일 
	serveconnbyte_gen.go

prefix_connbytemanager : server 용 connection manager, server에 연결된 connection(prefix_serveconnbyte) 을 관리한다. 

	생성하는 파일 
	connbytemanager_gen.go

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

prefix_authorize : client api call의 권한 관리 (commandid 기준)

	생성하는 파일 
	authorize_gen.go	

## 인자로 -statstype=elementtype 이 주어진 경우 추가로 생성하는 패키지 

prefix_error_stats : elementtype을 구성요소로한 simple errorcode 통계 (genenum에서 가져옴)

	생성하는 파일 
	prefix_error_stats_gen.go	

prefix_idcmd_stats : elementtype을 구성요소로한 simple commandid 통계 (genenum에서 가져옴)

	생성하는 파일 
	prefix_idcmd_stats_gen.go	

prefix_idnoti_stats : elementtype을 구성요소로한 simple notiid 통계 (genenum에서 가져옴)

	생성하는 파일 
	prefix_idnoti_stats_gen.go	


## websocket 을 사용하려면 

	go get github.com/gorilla/websocket
