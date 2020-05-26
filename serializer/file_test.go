package serializer_test

import (
	"testing"

	"github.com/idirall22/grpc/pb"
	"google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/require"

	"github.com/idirall22/grpc/sample"
	"github.com/idirall22/grpc/serializer"
)

func TestWriteProtobufToBinaryFile(t *testing.T) {
	t.Parallel()
	binaryFile := "../tmp/file.bin"
	jsonFile := "../tmp/file.json"

	laptop := sample.NewLaptop()
	err := serializer.WriteProtobufToBinaryFile(laptop, binaryFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err = serializer.ReadProtobufFromBinary(binaryFile, laptop2)
	require.NoError(t, err)
	require.True(t, proto.Equal(laptop, laptop2))

	err = serializer.WriteProtobufToJSONFile(laptop, jsonFile)
	require.NoError(t, err)
}
