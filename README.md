# genprotocol - protocol 기반 코드 생성

( goguelike 프로젝트를 하면서 만들어진 )
네트웍 등에서 사용하기 위한 protocol 을 작성하는데 필요한 반복적 이고 기계적인 코드들을 자동으로 생성 해주는 프로그램 입니다. 

https://github.com/kasworld/wasmwebsocket 을 보면 예제를 볼 수 있습니다. 

작성 이유는 하나의 프로젝트에서 여러종의 protocol 을 사용할 일이 생겨서 
반복된 코드를 작성하다 보니 만들게 되었습니다. 

goguelike 를 예로 들면 

game client - tower server 

tower server - ground server 

의 2종의 protocol이 필요하고 사실상 비슷한데 조금 다른 코드들이 서로 구별 되어 쓰입니다. 


## 사용 법 
위 와 같은 경우 

genprotocol -ver=1.0 -prefix=c2t -basedir=. 

과 같은 형태로 실행합니다. 

### 인자 설명 

ver : protocol의 version ( protocol 마다 버전이 다를 수 있습니다. )

prefix : 각 protocol을 구별하기 위한 prefix 

basedir : protocol code가 생성될 기본 dir


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
	goimports -w example/c2s_client/recvrspobjfnmap_gen.go
	goimports -w example/c2s_client/recvnotiobjfnmap_gen.go
	goimports -w example/c2s_client/callsendrecv_gen.go
	goimports -w example/c2s_server/demuxreq2api_gen.go
	goimports -w example/c2s_server/apitemplate_gen.go


prefix_gendata : genprotocol에서 읽어 들이는 파일들 

	command.data : request/response packet 용 command id 목록 
	noti.data    : notification packet 용 noti id 목록 
	error.data   : error code 목록 


prefix_client : client (접속을 시도 하는 쪽)에서 사용 

	생성하는 파일 
	recvrspobjfnmap_gen.go   
	recvnotiobjfnmap_gen.go
	callsendrecv_gen.go

prefix_server : server ( 접속을 받아 주는 쪽) 

	생성하는 파일 
	demuxreq2api_gen.go
	apitemplate_gen.go

prefix_msgp : messagepack marshal/unmarshal code (https://github.com/tinylib/msgp)

	생성하는 파일 
	serialize_gen.go
 
prefix_json : json marshal/unmarshal code 

	생성하는 파일 
	serialize_gen.go

prefix_error : protocol error 

	생성하는 파일 
	error_gen.go 


prefix_idcmd : protocol command (request,response) 

	생성하는 파일 
	command_gen.go

prefix_idnoti : protocol notification 

	생성하는 파일 
	noti_gen.go

prefix_obj : protocol struct 들 (packet body)

	생성하는 파일 
	objtemplate_gen.go : 예제 파일 - 참고해서 "_gen"이 없는 파일을 만들것 

prefix_packet : protocol packet( header + body )

	생성하는 파일 
	packet_gen.go

prefix_version : protocol version 정보 

	생성하는 파일 
	version_gen.go

