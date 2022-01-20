package kafka

import (
	"fmt"
	"github.com/segmentio/kafka-go/protocol"
	"io"
)

func ReadRequest(r io.Reader) (interface{}, error) {
	apiVersion, correlationID, clientID, msg, err := protocol.ReadRequest(r)
	if err != nil {
		return nil, err
	}
	fmt.Println(apiVersion)
	fmt.Println(correlationID)
	fmt.Println(clientID)
	fmt.Println(msg)
	return nil, nil
}

func WriteRequest(msg interface{}, w io.Writer) {
	protocol.WriteRequest(w, 0, 0, "", nil)
}
