package model

// Direction of TCP commnication
type Direction int

const (
	// ClientToRemote is client->proxy->remote
	ClientToRemote Direction = iota
	// RemoteToClient is client<-proxy<-remote
	RemoteToClient
	// SrcToDst is src->dst
	SrcToDst
	// DstToSrc is dst->src
	DstToSrc
	// Unknown direction
	Unknown Direction = 9
)

func (d Direction) String() string {
	switch d {
	case ClientToRemote, SrcToDst:
		return "->"
	case RemoteToClient, DstToSrc:
		return "<-"
	default:
		return "?"
	}
}
