package jarvisbase

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestInitLogger(t *testing.T) {
	_, err := initLogger(zapcore.DebugLevel, true, "../test/testinitlogger")
	if err != nil {
		t.Fatalf("TestBaseFunc TestInitLogger err! %v", err)

		return
	}

	_, err = initLogger(zapcore.DebugLevel, false, "../test/")
	if err != nil {
		t.Fatalf("TestBaseFunc TestInitLogger err! %v", err)

		return
	}

	t.Logf("TestBaseFunc TestInitLogger OK!")
}
