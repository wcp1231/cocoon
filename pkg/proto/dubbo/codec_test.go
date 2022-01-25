package dubbo

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

func TestDubboCodec_Decode(t *testing.T) {
	testDubboCodec_ReadPacket(t, "./test/dubbo_heart_request", 1)
	testDubboCodec_ReadPacket(t, "./test/dubbo_heart_response", 1)
}

func testDubboCodec_ReadPacket(t *testing.T, filename string, count int) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("read file %s failed. %v\n", filename, err)
		return
	}
	r := bufio.NewReader(bytes.NewReader(bs))
	encoded := []byte{}
	for i := 0; i < count; i++ {
		msg, err := ReadPacket(r)
		if err != nil {
			t.Fatalf("read packet failed. %v\n", err)
		}
		encoded = append(encoded, *msg.Raw...)
	}

	res := bytes.Compare(bs, encoded)
	if res != 0 {
		t.Fatal("Encode packet failed")
	}
}
