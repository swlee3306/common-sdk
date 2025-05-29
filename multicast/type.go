package multicast

import "encoding/json"

type HostInfoReceiver struct {
	Version      string   `json:"version"`
	BuildDate    string   `json:"buildDate"`
	Revision     string   `json:"revision"`
	Hostname     string   `json:"hostName"`
	IPs          []string `json:"ips"`
	Endpoint     string   `json:"endpoint"`
	EndpointPort int      `json:"endpointPort"`
}

type GenericMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
