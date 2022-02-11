package common

type Protocol struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Pass    bool
	Dump    bool
	Mock    bool
}

func (p *Protocol) String() string {
	return p.Name
}

var (
	PROTOCOL_UNKNOWN = &Protocol{
		Name:    "Unknown",
		Version: "0",
		Pass:    true,
		Dump:    false,
		Mock:    false,
	}
	PROTOCOL_NOT_SUPPORTED = &Protocol{
		Name:    "Not Supported",
		Version: "0",
		Pass:    true,
		Dump:    false,
		Mock:    false,
	}
	PROTOCOL_HTTP = &Protocol{
		Name:    "HTTP",
		Version: "1.1",
		Pass:    false,
		Dump:    true,
		Mock:    true,
	}
	PROTOCOL_REDIS = &Protocol{
		Name:    "Redis",
		Version: "0",
		Pass:    false,
		Dump:    true,
		Mock:    true,
	}
	PROTOCOL_MONGO = &Protocol{
		Name:    "Mongo",
		Version: "0",
		Pass:    false,
		Dump:    true,
		Mock:    true,
	}
	PROTOCOL_DUBBO = &Protocol{
		Name:    "Dubbo",
		Version: "0",
		Pass:    false,
		Dump:    true,
		Mock:    true,
	}
	PROTOCOL_MYSQL = &Protocol{
		Name:    "Mysql",
		Version: "0",
		Pass:    false,
		Dump:    true,
		Mock:    true,
	}
)
