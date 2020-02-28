package trace

import (
	"fmt"
	"io"
)

// Tracerはコード内での出来事を記録できるオブジェクトを表すインタフェースです．
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

type nilTracer struct{}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

func (t *nilTracer) Trace(a ...interface{}) {}

// OffはTraceメソッドの呼び出しを無視するTracerを返します．
func Off() Tracer {
	return &nilTracer{}
}
