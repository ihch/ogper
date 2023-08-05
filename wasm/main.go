package main

import (
	"encoding/json"
	"syscall/js"

	Ogper "github.com/ihch/ogper"
)

func getOGPForWasm(this js.Value, args []js.Value) interface{} {
	url := args[0].String()

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			ogp, err := Ogper.GetOGP(url)
			if err != nil {
				reject.Invoke(err)
			}

			responseJson, err := json.Marshal(ogp)
			if err != nil {
				reject.Invoke(err)
			}

			resolve.Invoke(string(responseJson))
		}()

		return nil
	})

	return js.Global().Get("Promise").New(handler)
}

func main() {
	c := make(chan struct{})
	js.Global().Set("ogper", js.ValueOf(map[string]any{
		"getOGP": js.FuncOf(getOGPForWasm),
	}))

	// wasmがJavaScriptからのリクエストに応答できるように、チャネルでプログラムを動作させ続ける
	<-c
}
