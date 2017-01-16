package trace

import (
	"fmt"
	"io"
)

// Tracer はコード内での出来事を記録できるオブジェクトを表すインターフェース
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off は Trace メソッドの呼び出しを無視する Tracer を返す
func Off() Tracer {
	return &nilTracer{}
}
