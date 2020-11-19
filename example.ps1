go run genprotocol.go -ver="1.0" -prefix=c2s -basedir=example -statstype=int
# go run genprotocol.go -ver=1.0 -prefix=c2s -basedir example -statstype=int -verbose

cd example
goimports -w .
cd ..



cd example/rundriver 

# echo "build wasm client"
# GOOS=js GOARCH=wasm go build -o www/wasmclient.wasm wasmclient.go

go build server.go 
go build client.go


$env:GOOS="js" 
$env:GOARCH="wasm" 
Write-Output "go build -o clientdata/wasmclient.wasm -ldflags `"-X main.Ver=${BUILD_VER}`" wasmclient.go"
go build -o clientdata/wasmclient.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclient.go
$env:GOOS=""
$env:GOARCH=""

cd ../..


# echo "build go client"
# go build -o bin/goclient goclient.go

