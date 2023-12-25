package testutil

import (
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

func MustUnmarshalTextproto[T proto.Message](t *testing.T, data string, m T) T {
	if err := prototext.Unmarshal([]byte(data), m); err != nil {
		t.Errorf("Unmarshal(%q) failed: %v", data, err)
	}

	return m
}
