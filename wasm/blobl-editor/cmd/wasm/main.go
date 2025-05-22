package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/redpanda-data/benthos/v4/public/bloblang"
	"github.com/redpanda-data/benthos/v4/public/service"

	_ "github.com/redpanda-data/benthos/v4/public/components/pure"
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

	result, err := service.NewMessage([]byte(args[1].String())).BloblangQuery(mapping)
	if err != nil {
		return fmt.Sprintf("Failed to execute mapping: %s", err)
	}

	var message any
	var metadata map[string]any
	// The result is nil if the mapping assigns `deleted()` to `root`
	if result != nil {
		message, err = result.AsStructured()
		if err != nil {
			res, err := result.AsBytes()
			if err != nil {
				return fmt.Errorf("failed to extract message: %s", err)
			}
			message = string(res)
		}

		if err = result.MetaWalkMut(func(key string, value any) error {
			if metadata == nil {
				metadata = make(map[string]any)
			}
			metadata[key] = value
			return nil
		}); err != nil {
			return fmt.Errorf("failed to extract metadata: %s", err)
		}
	}

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
