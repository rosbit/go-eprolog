package pl

import (
	"github.com/ichiban/prolog"
	"os"
)

type Prolog struct {
	i *prolog.Interpreter
}

func NewProlog() *Prolog {
	return &Prolog{i: prolog.New(os.Stdin, os.Stderr)}
}

func (p *Prolog) LoadScript(script string) (err error) {
	err = p.i.Exec(script)
	return
}

func (p *Prolog) LoadFile(file string) (err error) {
	err = p.i.Exec(`:- consult(?).`, []string{file})
	return
}
