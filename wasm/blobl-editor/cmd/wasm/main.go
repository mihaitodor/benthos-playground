package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/redpanda-data/benthos/v4/public/bloblang"
	"github.com/redpanda-data/benthos/v4/public/service"
)

func main() {
	js.Global().Set("blobl", js.FuncOf(blobl))

	// Wait for a signal to shut down
	select {}
}

func blobl(_ js.Value, args []js.Value) any {
	if len(args) != 2 {
		return fmt.Sprintf("Expected two arguments, received %d instead", len(args))
	}

	mapping, err := bloblang.NewEnvironment().Parse(args[0].String())
	if err != nil {
		return fmt.Sprintf("Failed to parse mapping: %s", err)
	}

	msg, err := service.NewMessage([]byte(args[1].String())).BloblangQuery(mapping)
	if err != nil {
		return fmt.Sprintf("Failed to execute mapping: %s", err)
	}

	message, err := msg.AsStructured()
	if err != nil {
		return fmt.Sprintf("Failed to marshal message: %s", err)
	}

	var metadata map[string]any
	msg.MetaWalkMut(func(key string, value any) error {
		if metadata == nil {
			metadata = make(map[string]any)
		}
		metadata[key] = value
		return nil
	})

	var output []byte
	if output, err = json.MarshalIndent(struct {
		Msg  any            `json:"msg"`
		Meta map[string]any `json:"meta,omitempty"`
	}{
		Msg:  message,
		Meta: metadata,
	}, "", "  "); err != nil {
		return fmt.Sprintf("Failed to marshal output: %s", err)
	}

	return string(output)
}
