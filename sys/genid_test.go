package sys

import (
	"testing"
)

func TestGenID(t *testing.T) {
	t.Log(GenerateID(UUIDNamespace, "678bb5a27e785308b9e937a3"))
}

func TestGenV5(t *testing.T) {
	t.Log(GenerateV5("ac35628dce0e5a34ab6854368a6973f7", "32796c8170b050aa915dc45acb9c23c8", "3df0bd87e067452aa952d3915e319461"))
}
