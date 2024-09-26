//go:generate go run github.com/bytecodealliance/wasm-tools-go/cmd/wit-bindgen-go generate --world proxy ./wit

package main

import (
	"github.com/bytecodealliance/wasm-tools-go/cm"
	incominghandler "github.com/mossaka/wasi-http-with-wasm-tools-go/wasi/http/incoming-handler"
	"github.com/mossaka/wasi-http-with-wasm-tools-go/wasi/http/types"
)

func init() {
	incominghandler.Exports.Handle = func(request types.IncomingRequest, responseOut types.ResponseOutparam) {
		headers := types.NewFields()

		httpResponse := types.NewOutgoingResponse(headers)
		httpResponse.SetStatusCode(200)

		bodyResult := httpResponse.Body()
		if bodyResult.IsErr() {
			handleError(responseOut, "Failed to retrieve body")
			return
		}
		body := bodyResult.OK()

		streamResult := body.Write()
		if streamResult.IsErr() {
			handleError(responseOut, "Failed to write to stream")
			return
		}
		outputStream := streamResult.OK()

		writerResult := outputStream.BlockingWriteAndFlush(cm.ToList([]uint8("Hello from Go!\n")))
		if writerResult.IsErr() {
			handleError(responseOut, "Failed to write and flush")
			return
		}
		outputStream.ResourceDrop()

		types.OutgoingBodyFinish(*body, cm.None[types.Fields]())

		result := cm.OK[cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]](httpResponse)
		types.ResponseOutparamSet(responseOut, result)
	}
}

func handleError(responseOut types.ResponseOutparam, message string) {
	httpResponse := types.NewOutgoingResponse(types.NewFields())
	httpResponse.SetStatusCode(500)
	bodyResult := httpResponse.Body()
	if bodyResult.IsErr() {
		// If handling error fails, there's not much we can do.
		return
	}
	body := bodyResult.OK()

	streamResult := body.Write()
	if streamResult.IsErr() {
		// If handling error fails, there's not much we can do.
		return
	}
	outputStream := streamResult.OK()
	outputStream.BlockingWriteAndFlush(cm.ToList([]uint8(message)))
	outputStream.ResourceDrop()

	types.OutgoingBodyFinish(*body, cm.None[types.Fields]())

	result := cm.OK[cm.Result[types.ErrorCodeShape, types.OutgoingResponse, types.ErrorCode]](httpResponse)
	types.ResponseOutparamSet(responseOut, result)
}

func main() {}
