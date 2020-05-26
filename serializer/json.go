package serializer

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// protobufToJSON protobuf to json
func protobufToJSON(message proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		Indent:       "  ",
		OrigName:     true,
		EmitDefaults: true,
	}
	return marshaler.MarshalToString(message)
}
