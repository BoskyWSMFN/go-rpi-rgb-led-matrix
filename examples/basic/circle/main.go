package main

import (
	"flag"
	"github.com/BoskyWSMFN/go-rpi-rgb-led-matrix/pkg/canvas"
	"github.com/BoskyWSMFN/go-rpi-rgb-led-matrix/pkg/matrix"
	"image/color"
	"os"
	"os/signal"
	"syscall"
)

var (
	rows                   = flag.Int("led-rows", 32, "number of rows supported")
	cols                   = flag.Int("led-cols", 64, "number of columns supported")
	parallel               = flag.Int("led-parallel", 1, "number of daisy-chained panels")
	chain                  = flag.Int("led-chain", 2, "number of displays daisy-chained")
	brightness             = flag.Int("brightness", 100, "brightness (0-100)")
	hardwareMapping        = flag.String("led-gpio-mapping", "adafruit-hat-pwm", "Name of GPIO mapping used.")
	showRefresh            = flag.Bool("led-show-refresh", false, "Show refresh rate.")
	inverseColors          = flag.Bool("led-inverse", false, "Switch if your matrix has inverse colors on.")
	disableHardwarePulsing = flag.Bool("led-no-hardware-pulse", false, "Don't use hardware pin-pulse generation.")
	pixelMapping           = flag.String("led-pixel-mapper", "U-mapper", "Pixel mapping from api")
	slowGPIO               = flag.Int("led-slowdown-gpio", 5, "GPIO slowdown for pwm rpi hats")
)

func main() {
	config := &matrix.DefaultConfig
	config.Rows = *rows
	config.Cols = *cols
	config.Parallel = *parallel
	config.ChainLength = *chain
	config.Brightness = *brightness
	config.HardwareMapping = *hardwareMapping
	config.ShowRefreshRate = *showRefresh
	config.InverseColors = *inverseColors
	config.DisableHardwarePulsing = *disableHardwarePulsing
	config.PixelMapping = *pixelMapping
	config.SlowdownGPIO = *slowGPIO

	m, err := matrix.NewRGBLedMatrix(config)
	fatal(err)

	tk := canvas.NewToolKit(m)
	defer tk.Close()

	geomX, geomY := m.Geometry()
	centerX := geomX / 2
	centerY := geomY / 2
	radius := 30

	tk.DrawCircle(centerX, centerY, radius, color.RGBA{G: 255, A: 255})

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	// Main routine sleeps until signal is received or global context is done.
	case <-osSignal:
		return
	}
}

func init() {
	flag.Parse()
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
