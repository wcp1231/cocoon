package redis

import "testing"

func TestClassifier_Match(t *testing.T) {
	datas := []string{
		"+PONG\r\n",
		"-ERR\r\n",
		"*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar",
		"$3\r\nfoo",
		":42\r\n",
		"$-1\r\n",
	}
	c := Classifier{}
	for _, data := range datas {
		if !c.Match([]byte(data)) {
			t.Fatalf("Test %v fail", data)
		}
	}
}
