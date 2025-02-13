package sys

import (
	"testing"
)

func TestGenID(t *testing.T) {
	t.Log(GenerateID(UUIDNamespace, "678bb5a27e785308b9e937a3"))
}
