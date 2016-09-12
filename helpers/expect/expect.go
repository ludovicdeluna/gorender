package expect

import "fmt"

var msg string = "%s :\nGot:  %v\nWant: %v"

type Expects struct {
	Got   interface{}
	want  interface{}
	title string
	Fail  bool
}

func New() *Expects {
	return &Expects{}
}

func For(got interface{}) *Expects {
	return &Expects{Got: got}
}

func (e *Expects) For(got interface{}) *Expects {
	e.Got = got
	return e
}

func (e *Expects) String() string {
	return fmt.Sprintf(msg, e.title, e.Got, e.want)
}

func (e *Expects) It(m string) string {
	e.title = m
	return e.String()
}

func (e *Expects) Title(m string) *Expects {
	e.title = m
	return e
}

func (e *Expects) Equals(want interface{}) *Expects {
	e.want = want
	e.Fail = e.Got != e.want
	return e
}

func (e *Expects) HasError(want interface{}) *Expects {
	e.want = want
	switch t := e.Got.(type) {
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
