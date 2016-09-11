package expect

import "fmt"

var msg string = "%s :\nGot:  %v\nWant: %v"

type Expects struct {
	got   interface{}
	want  interface{}
	title string
	Fail  bool
}

func New() *Expects {
	return &Expects{}
}

func For(got interface{}) *Expects {
	return &Expects{got: got}
}

func (e *Expects) For(got interface{}) *Expects {
	e.got = got
	return e
}

func (e *Expects) String() string {
	return fmt.Sprintf(msg, e.title, e.got, e.want)
}

func (e *Expects) It(m string) string {
	e.title = m
	return e.String()
}

func (e *Expects) Equals(want interface{}) *Expects {
	e.want = want
	e.Fail = e.got != e.want
	return e
}

func (e *Expects) HasError(want interface{}) *Expects {
	e.want = want
	switch t := e.got.(type) {
	case error:
		e.Fail = t.Error() != e.want
	default:
		e.Fail = true
	}
	return e
}

func (e *Expects) HasNoError() *Expects {
	return e.Equals(error(nil))
}
