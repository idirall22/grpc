package serializer

import (
	"fmt"
	"io/ioutil"

	// "google.golang.org/protobuf/proto"
	"github.com/golang/protobuf/proto"
)

// WriteProtobufToBinaryFile write protobuf to a binary file
func WriteProtobufToBinaryFile(message proto.Message, fileName string) error {
	data, err := proto.Marshal(message)

	if err != nil {
		return fmt.Errorf("Could not marshal proto message to binary: %w", err)
	}

	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("Could not write message data to file: %w", err)
	}

	return nil
}

// WriteProtobufToJSONFile write protobuf to json file
func WriteProtobufToJSONFile(message proto.Message, fileName string) error {
	data, err := protobufToJSON(message)
	if err != nil {
		return fmt.Errorf("Could not marshal message to json %w", err)
	}

	err = ioutil.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("Could not write json file %w", err)
	}

	return nil
}

// ReadProtobufFromBinary read protobuf from abinary file
func ReadProtobufFromBinary(fileName string, message proto.Message) error {
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return fmt.Errorf("Could not read binary file %w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("Could not unmarshal data to message %w", err)
	}
	return nil
}
