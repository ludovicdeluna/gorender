package expect

import "fmt"

var msg string = "%s :\nGot:  %v\nWant: %v"

type Expects struct {
	got   interface{}
	want  interface{}
	title string
  fail  bool
}

func For(got interface{}) *Expects {
	expects := Expects{got: got}
	return &expects
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
  e.fail = e.got != e.want
  return e
}

func (e *Expects) HasError(want interface{}) *Expects {
	e.want = want
  switch t := e.got.(type) {
  case error:
  	e.fail = t.Error() != e.want
  default:
  	e.fail = true
  }
  return e
}

func (e *Expects) Fail() bool {
  return e.fail
}
