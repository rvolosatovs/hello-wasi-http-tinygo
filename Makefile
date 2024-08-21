PHONY: build serve
build:
	tinygo build -target=wasip2 -wit-package ./wit -wit-world proxy main.go


serve:
	wasmtime serve main.wasm --addr 127.0.0.1:3000