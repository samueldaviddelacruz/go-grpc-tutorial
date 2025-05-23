package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtobufToJSON converts protocol buffer message to JSON string
func ProtobufToJSON(message proto.Message) (string, error) {
	marshaler := protojson.MarshalOptions{
		Indent:            "  ",
		EmitDefaultValues: true,
		UseProtoNames:     false,
	}
	result, err := marshaler.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// JSONToProtobufMessage converts JSON string to protocol buffer message
func JSONToProtobufMessage(data string, message proto.Message) error {
	return protojson.Unmarshal([]byte(data), message)
}
