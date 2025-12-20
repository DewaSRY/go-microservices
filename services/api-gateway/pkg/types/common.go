package types

import "encoding/json"

type WsMessageData struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
