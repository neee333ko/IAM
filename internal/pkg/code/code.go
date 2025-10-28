package code

import (
	"github.com/neee333ko/errors"
	"github.com/novalagung/gubrak"
)

var ref string = "http://reference.com"

type MyCoder struct {
	C    int
	Http int
	Msg  string
	Ref  string
}

func (c *MyCoder) Code() int {
	return c.C
}

func (c *MyCoder) HttpStatus() int {
	return c.Http
}

func (c *MyCoder) Message() string {
	return c.Msg
}

func (c *MyCoder) Reference() string {
	return c.Ref
}

func register(code int, httpStatus int, msg string) {
	ok, _ := gubrak.Includes([]int{200, 400, 401, 403, 404, 500}, httpStatus)

	if !ok {
		panic("invalid http status")
	}

	c := &MyCoder{
		C:    code,
		Http: httpStatus,
		Msg:  msg,
		Ref:  ref,
	}

	errors.MustRegister(c)
}
