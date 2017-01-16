package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("New からの返り値が nil です")
	} else {
		tracer.Trace("こんにちは，trace パッケージ")
		if buf.String() != "こんにちは，trace パッケージ\n" {
			t.Errorf("'%s' という誤った文字列が出力されました", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("データ")
}
