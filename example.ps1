go run genprotocol.go -ver="1.0" -prefix=c2s -basedir=example -statstype=int
# go run genprotocol.go -ver=1.0 -prefix=c2s -basedir example -statstype=int -verbose

goimports -w example

cd example/rundriver 

echo "go build server.go"
go build server.go 
echo "go build client.go"
go build client.go


$env:GOOS="js" 
$env:GOARCH="wasm" 
echo "go build -o clientdata/wasmclient.wasm -ldflags `"-X main.Ver=${BUILD_VER}`" wasmclient.go"
go build -o clientdata/wasmclient.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclient.go
$env:GOOS=""
$env:GOARCH=""

cd ../..

