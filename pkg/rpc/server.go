package rpc

import (
	"fmt"
	"image/color"
	"net"
	"net/http"
	"net/rpc"

	"github.com/pkg/errors"
	"go-rpi-rgb-led-matrix/tools"
)

type Matrix struct {
	m tools.Matrix
}

type GeometryArgs struct{}
type GeometryReply struct{ Width, Height int }

func (m *Matrix) Geometry(_ *GeometryArgs, reply *GeometryReply) error {
	w, h := m.m.Geometry()
	reply.Width = w
	reply.Height = h

	return nil
}

type ApplyArgs struct{ Colors []color.Color }
type ApplyReply struct{}

func (m *Matrix) Apply(args *ApplyArgs, _ *ApplyReply) error {
	return m.m.Apply(args.Colors)
}

type CloseArgs struct{}
type CloseReply struct{}

func (m *Matrix) Close(_ *CloseArgs, _ *CloseReply) error {
	return m.m.Close()
}

type SetBrightnessArgs struct{ B uint8 }
type SetBrightnessReply struct{}

func (m *Matrix) SetBrightness(a *SetBrightnessArgs, _ *SetBrightnessReply) {
	m.m.SetBrightness(a.B)
}

func Serve(m tools.Matrix) error {
	err := rpc.Register(&Matrix{m})
	if err != nil {
		return errors.Wrap(err, "register error:")
	}

	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		return errors.Wrap(err, "listen error:")
	}

	fmt.Println(l)
	err = http.Serve(l, nil)

	return errors.Wrap(err, "serve error:")
}
