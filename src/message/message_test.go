package message

import (
	"testing"
)

func TestMessageParser(t *testing.T) {
	message := ParseRawMessage("<0>'e' 2010<\\0>\n")
	if message == nil || message.Name != "0" || message.Content != "'e' 2010" {
		t.FailNow()
	}
}
