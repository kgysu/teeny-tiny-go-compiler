package comp

import (
	"fmt"
	"os"
)

type Emitter struct {
	fullPath string
	header   string
	code     string
}

func NewEmitter(p string) *Emitter {
	return &Emitter{
		fullPath: p,
	}
}

func (e *Emitter) emit(code string) {
	e.code += code
}

func (e *Emitter) emitLine(code string) {
	e.code += code + "\n"
}

func (e *Emitter) headerLine(code string) {
	e.header += code + "\n"
}

func (e *Emitter) WriteFile() {
	f, err := os.OpenFile(e.fullPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintln(f, e.header+e.code)
	if err != nil {
		panic(err)
	}
}
