package main

import (
	"fmt"
	"io"
)

type Template struct {
	w   io.Writer
	err error
}

func NewTemplate(w io.Writer) *Template {
	return &Template{w: w}
}

func (t *Template) Write(p []byte) (n int, err error) {
	return t.w.Write(p)
}

func (t *Template) Err() error {
	return t.err
}

func (t *Template) isError() bool {
	return t.err != nil
}

func (t *Template) CheckError(err error) {
	if t.isError() {
		return
	}

	if err != nil {
		t.err = err
	}
}

func (t *Template) T(s string) {
	if t.isError() {
		return
	}

	_, t.err = t.w.Write([]byte(s))
}

func (t *Template) Tfmt(format string, a ...interface{}) {
	if t.isError() {
		return
	}

	_, t.err = fmt.Fprintf(t.w, format, a...)
}

func (t *Template) Tfmtln(format string, a ...interface{}) {
	t.Tfmt(format, a...)
	t.T("\n")
}
