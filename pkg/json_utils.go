package pkg

import (
	"encoding/json"
	"fmt"
	"io"
)

type JsonUtils struct{}

func (*JsonUtils) Encode(destination io.Writer, dataToEncode any) error {
	// Encoding the json
	err := json.NewEncoder(destination).Encode(dataToEncode)
	// Returning error if any
	return fmt.Errorf("failed to encode json: %w", err)
}

func (*JsonUtils) Decode(destination any, dataToDecode io.Reader) error {
	// Decoding the json
	err := json.NewDecoder(dataToDecode).Decode(destination)
	// Returning error if any
	return fmt.Errorf("failed to decode json: %w", err)
}