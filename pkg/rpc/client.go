package rpc

import (
	"encoding/gob"
	"github.com/BoskyWSMFN/go-rpi-rgb-led-matrix/tools"
	"image/color"
	"net/rpc"
)

func init() {
	gob.Register(color.RGBA{})
}

// RGBLedMatrix matrix representation for ws281x
type Client struct {
	network string
	addr    string
	client  *rpc.Client
	leds    []color.Color
}

// NewRGBLedMatrix returns a new matrix using the given size and config
func NewClient(network, addr string, size int) (tools.Matrix, error) {
	client, err := rpc.DialHTTP(network, addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		network: network,
		addr:    addr,
		client:  client,
		leds:    make([]color.Color, size),
	}, nil
}

// Geometry returns the width and the height of the matrix
func (c *Client) Geometry() (width, height int) {
	var reply *GeometryReply
	err := c.client.Call("Matrix.Geometry", &GeometryArgs{}, &reply)
	if err != nil {
		panic(err)
	}

	return reply.Width, reply.Height
}

func (c *Client) Apply(leds []color.Color) error {
	defer func() { c.leds = make([]color.Color, 2048) }()

	var reply *ApplyReply
	return c.client.Call("Matrix.Apply", &ApplyArgs{Colors: leds}, &reply)
}

// Render update the display with the data from the LED buffer
func (c *Client) Render() error {
	return c.Apply(c.leds)
}

// At return an Color which allows access to the LED display data as
// if it were a sequence of 24-bit RGB values.
func (c *Client) At(position int) color.Color {
	if c.leds[position] == nil {
		return color.Black
	}

	return c.leds[position]
}

// Set set LED at position x,y to the provided 24-bit color value.
func (c *Client) Set(position int, col color.Color) {
	c.leds[position] = color.RGBAModel.Convert(col)
}

// Close finalizes the ws281x interface
func (c *Client) Close() error {
	return c.Apply(make([]color.Color, 2048))
}

func (c *Client) SetBrightness(b uint8) error {
	return c.client.Call("Matrix.Geometry", &SetBrightnessArgs{B: b}, &SetBrightnessReply{})
}
