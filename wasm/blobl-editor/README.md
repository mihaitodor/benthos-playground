# WebAssembly (WASM) bloblang editor

This is a WASM implementation of `benthos blobl server`.

## Build WASM module

```shell
> GOOS=js GOARCH=wasm go build -o ./static/blobl.wasm ./cmd/wasm
> cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./static
```

## Run the webserver

```shell
> go run ./cmd/webserver
```

Open the following page in a browser: http://localhost:3000/index.html
