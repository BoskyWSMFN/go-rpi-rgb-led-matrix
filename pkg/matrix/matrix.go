package matrix

// TODO refactor following C code

/*

#cgo CFLAGS: -std=c99 -I${SRCDIR}/include -DSHOW_REFRESH_RATE
#cgo amd64 LDFLAGS: -lrgbmatrix -L${SRCDIR}/lib/amd64 -lstdc++ -lm
#cgo arm LDFLAGS: -lrgbmatrix -L${SRCDIR}/lib/arm -lstdc++ -lm
#cgo arm64 LDFLAGS: -lrgbmatrix -L${SRCDIR}/lib/arm64 -lstdc++ -lm
#include <led-matrix-c.h>
#include <stdio.h>

void led_matrix_swap(struct RGBLedMatrix *matrix, struct LedCanvas *offscreen_canvas,
                     int width, int height, const uint32_t pixels[]) {


  int i, x, y;
  uint32_t color;
  for (x = 0; x < width; ++x) {
    for (y = 0; y < height; ++y) {
      i = x + (y * width);
      color = pixels[i];

      led_canvas_set_pixel(offscreen_canvas, x, y,
        (color >> 16) & 255, (color >> 8) & 255, color & 255);
    }
  }

  offscreen_canvas = led_matrix_swap_on_vsync(matrix, offscreen_canvas);
}

void set_show_refresh_rate(struct RGBLedMatrixOptions *o, int show_refresh_rate) {
  o->show_refresh_rate = show_refresh_rate != 0 ? 1 : 0;
}

void set_disable_hardware_pulsing(struct RGBLedMatrixOptions *o, int disable_hardware_pulsing) {
  o->disable_hardware_pulsing = disable_hardware_pulsing != 0 ? 1 : 0;
}

void set_inverse_colors(struct RGBLedMatrixOptions *o, int inverse_colors) {
  o->inverse_colors = inverse_colors != 0 ? 1 : 0;
}

void set_daemon(struct RGBLedRuntimeOptions *rt, int daemon) {
  rt->daemon = daemon != 0 ? 1 : 0;
}

void set_drop_privileges(struct RGBLedRuntimeOptions *rt, int drop) {
  rt->drop_privileges = drop != 0 ? 1 : 0;
}
*/
import "C" //nolint: typecheck
import (
	"fmt"
	"github.com/BoskyWSMFN/go-rpi-rgb-led-matrix/pkg/canvas"
	"image/color"
	"os"
	"unsafe"

	"github.com/BoskyWSMFN/go-rpi-rgb-led-matrix/internal/emulator"
)

// DefaultConfig default WS281x configuration
var DefaultConfig = HardwareConfig{
	Rows:              32,
	Cols:              32,
	ChainLength:       1,
	Parallel:          1,
	PWMBits:           11,
	PWMLSBNanoseconds: 130,
	Brightness:        100,
	ScanMode:          Progressive,
}

// HardwareConfig rgb-led-matrix configuration
type HardwareConfig struct {
	// Rows the number of rows supported by the display, so 32 or 16.
	Rows int
	// Cols the number of columns supported by the display, so 32 or 64 .
	Cols int
	// ChainLengthis the number of displays daisy-chained together
	// (output of one connected to input of next).
	ChainLength int
	// Parallel is the number of parallel chains connected to the Pi; in old Pis
	// with 26 GPIO pins, that is 1, in newer Pis with 40 interfaces pins, that
	// can also be 2 or 3. The effective number of pixels in vertical direction is
	// then thus rows * parallel.
	Parallel int
	// Set PWM bits used for output. Default is 11, but if you only deal with
	// limited comic-colors, 1 might be sufficient. Lower require less CPU and
	// increases refresh-rate.
	PWMBits int
	// Change the base time-unit for the on-time in the lowest significant bit in
	// nanoseconds.  Higher numbers provide better quality (more accurate color,
	// less ghosting), but have a negative impact on the frame rate.
	PWMLSBNanoseconds int // the DMA channel to use
	// Brightness is the initial brightness of the panel in percent. Valid range
	// is 1..100
	Brightness int
	// ScanMode progressive or interlaced
	ScanMode ScanMode // strip color layout
	// Disable the PWM hardware subsystem to create pulses. Typically, you don't
	// want to disable hardware pulsing, this is mostly for debugging and figuring
	// out if there is interference with the sound system.
	// This won't do anything if output enable is not connected to GPIO 18 in
	// non-standard wirings.
	DisableHardwarePulsing bool

	ShowRefreshRate bool
	InverseColors   bool

	// Name of GPIO mapping used
	HardwareMapping string

	Multiplexing int
	PixelMapping string

	SlowdownGPIO   int
	Daemon         bool
	DropPrivileges bool
	DoGPIOInit     bool
}

func (c *HardwareConfig) geometry() (width, height int) {
	return c.Cols * c.ChainLength, c.Rows * c.Parallel
	//return c.Cols, c.Rows
}

func (c *HardwareConfig) toC() (*C.struct_RGBLedMatrixOptions, *C.struct_RGBLedRuntimeOptions) {
	o := &C.struct_RGBLedMatrixOptions{}
	o.rows = C.int(c.Rows)
	o.cols = C.int(c.Cols)
	o.chain_length = C.int(c.ChainLength)
	o.parallel = C.int(c.Parallel)
	o.pwm_bits = C.int(c.PWMBits)
	o.pwm_lsb_nanoseconds = C.int(c.PWMLSBNanoseconds)
	o.brightness = C.int(c.Brightness)
	o.scan_mode = C.int(c.ScanMode)
	o.hardware_mapping = C.CString(c.HardwareMapping)
	o.pixel_mapper_config = C.CString(c.PixelMapping)
	o.multiplexing = C.int(c.Multiplexing)

	if c.ShowRefreshRate { // TODO also refactor. Change to bools.
		C.set_show_refresh_rate(o, C.int(1))
	} else {
		C.set_show_refresh_rate(o, C.int(0))
	}

	if c.DisableHardwarePulsing { // TODO same here
		C.set_disable_hardware_pulsing(o, C.int(1))
	} else {
		C.set_disable_hardware_pulsing(o, C.int(0))
	}

	if c.InverseColors { // TODO same here
		C.set_inverse_colors(o, C.int(1))
	} else {
		C.set_inverse_colors(o, C.int(0))
	}

	rt := &C.struct_RGBLedRuntimeOptions{}
	rt.gpio_slowdown = C.int(c.SlowdownGPIO)
	rt.do_gpio_init = C.bool(c.DoGPIOInit)

	if c.Daemon {
		C.set_daemon(rt, C.int(1))
	} else {
		C.set_daemon(rt, C.int(0))
	}

	if c.DropPrivileges {
		C.set_drop_privileges(rt, C.int(1))
	} else {
		C.set_drop_privileges(rt, C.int(0))
	}

	return (*C.struct_RGBLedMatrixOptions)(unsafe.Pointer(o)),
		(*C.struct_RGBLedRuntimeOptions)(unsafe.Pointer(rt))
}

type ScanMode int8

const (
	Progressive ScanMode = 0
	Interlaced  ScanMode = 1
)

// RGBLedMatrix matrix representation for ws281x
type RGBLedMatrix struct {
	Config *HardwareConfig

	height int
	width  int
	matrix *C.struct_RGBLedMatrix
	buffer *C.struct_LedCanvas
	leds   []C.uint32_t
}

const EmulatorENV = "MATRIX_EMULATOR"

// NewRGBLedMatrix returns a new matrix using the given size and config
func NewRGBLedMatrix(config *HardwareConfig) (c canvas.Matrix, err error) {
	defer func() {
		if r := recover(); r != nil {
			_, ok := r.(error)
			if !ok {
				err = fmt.Errorf("error creating matrix: %v", r)
			}
		}
	}()

	if isMatrixEmulator() {
		return buildMatrixEmulator(config), nil
	}

	m := C.led_matrix_create_from_options_and_rt_options(config.toC()) // TODO C exceptions handling
	b := C.led_matrix_create_offscreen_canvas(m)                       // TODO same here

	cW, cH := C.int(0), C.int(0)

	C.led_canvas_get_size(b, (*C.int)(unsafe.Pointer(&cW)), (*C.int)(unsafe.Pointer(&cH)))

	w, h := int(cW), int(cH)

	c = &RGBLedMatrix{
		Config: config,
		width:  w, height: h,
		matrix: m,
		buffer: b,
		leds:   make([]C.uint32_t, w*h),
	}
	if m == nil {
		return nil, fmt.Errorf("unable to allocate memory")
	}

	return c, nil
}

func isMatrixEmulator() bool {
	return os.Getenv(EmulatorENV) == "1"
}

func buildMatrixEmulator(config *HardwareConfig) canvas.Matrix {
	w, h := config.geometry()
	return emulator.NewEmulator(w, h, emulator.DefaultPixelPitch, true)
}

// Geometry returns the width and the height of the matrix
func (c *RGBLedMatrix) Geometry() (width, height int) {
	return c.width, c.height
}

// Apply set all the pixels to the values contained in leds
func (c *RGBLedMatrix) Apply(leds []color.Color) error {
	for position, l := range leds {
		c.Set(position, l)
	}

	return c.Render()
}

// Render update the display with the data from the LED buffer
func (c *RGBLedMatrix) Render() error {
	w, h := c.Geometry()

	C.led_matrix_swap(
		c.matrix,
		c.buffer,
		C.int(w), C.int(h),
		(*C.uint32_t)(unsafe.Pointer(&c.leds[0])),
	)

	c.leds = make([]C.uint32_t, w*h)
	return nil
}

// At return an Color which allows access to the LED display data as
// if it were a sequence of 24-bit RGB values.
func (c *RGBLedMatrix) At(position int) color.Color {
	return uint32ToColor(c.leds[position])
}

// Set set LED at position x,y to the provided 24-bit color value.
func (c *RGBLedMatrix) Set(position int, color color.Color) {
	c.leds[position] = C.uint32_t(colorToUint32(color))
}

// DrawCircle - requires rendering before it is called.
func (c *RGBLedMatrix) DrawCircle(x, y, radius int, col color.Color) {
	r, g, b, _ := col.RGBA()

	C.draw_circle(c.buffer, C.int(x), C.int(y), C.int(radius), C.uint8_t(r), C.uint8_t(g), C.uint8_t(b))
}

// SetBrightness - sets brightness in real time. Issue with network activity still remains.
func (c *RGBLedMatrix) SetBrightness(b uint8) error {
	if b > 100 || b == 0 {
		return fmt.Errorf("must be from 1 to 100")
	}

	C.led_matrix_set_brightness(c.matrix, C.uint8_t(b))

	return nil
}

// Close finalizes the ws281x interface
func (c *RGBLedMatrix) Close() error {
	C.led_matrix_delete(c.matrix)

	return nil
}

func colorToUint32(c color.Color) uint32 {
	if c == nil {
		return 0
	}

	// A color's RGBA method returns values in the range [0, 65535]
	red, green, blue, _ := c.RGBA()

	return (red>>8)<<16 | (green>>8)<<8 | blue>>8
}

func uint32ToColor(u C.uint32_t) color.Color {
	return color.RGBA{
		R: uint8(u>>16) & 255,
		G: uint8(u>>8) & 255,
		B: uint8(u>>0) & 255,
	}
}
