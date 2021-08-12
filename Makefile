wasm:
	GOOS=js GOARCH=wasm go build -mod vendor -o static/wasm/query.wasm cmd/query/main.go
