// Copyright (c) 2016 ludovic de luna (github.com/ludovicdeluna)
// License MIT : https://opensource.org/licenses/MIT
//
// The expect Helper - A tiny helper to reuse in your code w/o black box.
// Feel free to update it to your needs
package expect

import "fmt"

var Msg string = "%s :\nGot:  %v\nWant: %v"

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

func It(m string, opts ...interface{}) *Expects {
	return &Expects{title: fmt.Sprintf(m, opts...)}
}

func (e *Expects) For(got interface{}) *Expects {
	e.Got = got
	return e
}

func (e *Expects) It(m string, opts ...interface{}) string {
	e.title = fmt.Sprintf(m, opts...)
	return e.String()
}

func (e *Expects) String() string {
	return fmt.Sprintf(Msg, e.title, e.Got, e.want)
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
