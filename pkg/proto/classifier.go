package proto

import (
	"bufio"
	"cocoon/pkg/model/api"
	"cocoon/pkg/model/common"
	"cocoon/pkg/proto/http"
	"cocoon/pkg/proto/mongo"
	"cocoon/pkg/proto/redis"
)

type ProtoClassifier struct {
	size        int
	classifiers []api.ProtoClassifier
}

func NewProtoClassifier() *ProtoClassifier {
	var classifiers []api.ProtoClassifier
	classifiers = append(classifiers, &redis.Classifier{})
	classifiers = append(classifiers, &http.Classifier{})
	classifiers = append(classifiers, &mongo.Classifier{})
	classifiers = append(classifiers, &DefaultClassifier{})
	return &ProtoClassifier{
		size:        len(classifiers),
		classifiers: classifiers,
	}
}

func (c *ProtoClassifier) Classify(r *bufio.Reader) *common.Protocol {
	for _, c := range c.classifiers {
		if c.Match(r) {
			return c.Protocol()
		}
	}
	return c.DefaultProto()
}

func (c *ProtoClassifier) DefaultProto() *common.Protocol {
	return c.classifiers[c.size-1].Protocol()
}