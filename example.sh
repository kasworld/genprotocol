go run genprotocol.go -ver=1.0 -prefix=c2s -basedir example -statstype=int
# go run genprotocol.go -ver=1.0 -prefix=c2s -basedir example -statstype=int -verbose

cd example
goimports -w .
cd ..

cd example/rundriver 

echo "build wasm client"
GOOS=js GOARCH=wasm go build -o www/wasmclient.wasm wasmclient.go

# echo "build go client"
# go build -o bin/goclient goclient.go

