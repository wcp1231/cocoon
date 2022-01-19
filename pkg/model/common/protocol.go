package common

import "fmt"

type Protocol struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (p *Protocol) String() string {
	return fmt.Sprintf("Proto[%s]", p.Name)
}

var (
	PROTOCOL_UNKNOWN = &Protocol{
		Name:    "Known",
		Version: "0",
	}
	PROTOCOL_HTTP = &Protocol{
		Name:    "HTTP",
		Version: "1.1",
	}
	PROTOCOL_REDIS = &Protocol{
		Name:    "Redis",
		Version: "0",
	}
)
