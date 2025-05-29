package multicast

import "encoding/json"

type HostInfoReceiver struct {
	Hostname string   `json:"hostname"`
	IPs      []string `json:"ips"`
}

type GenericMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
