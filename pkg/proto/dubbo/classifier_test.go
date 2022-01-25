package dubbo

import (
	"bufio"
	"bytes"
	"os"
	"testing"
)

func TestClassifier_Match(t *testing.T) {
	c := Classifier{}
	testDubboClassifier(t, c, "./test/dubbo_heart_request")
	testDubboClassifier(t, c, "./test/dubbo_heart_response")
}

func testDubboClassifier(t *testing.T, c Classifier, filename string) {
	reqBytes, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("read file %s failed. %v\n", filename, err)
		return
	}
	r := bufio.NewReader(bytes.NewReader(reqBytes))
	if !c.Match(r) {
		t.Fatal("Test match request fail")
	}
}
